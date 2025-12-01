package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"maestro/internal/domain/session"
	"maestro/internal/models"
	"maestro/internal/store"
)

// ============================================
// 1️⃣ SESSION BUILDER (Choix énergie)
// ============================================

// HandleSessionBuilder : Page de sélection d'énergie
func HandleSessionBuilder(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 SessionBuilder: displaying energy selection")

	// 1. Récupère configs énergie (domain)
	configs := []session.Config{
		session.GetConfig(models.EnergyLow),
		session.GetConfig(models.EnergyMedium),
		session.GetConfig(models.EnergyHigh),
	}

	// 2. Structure données
	data := map[string]any{
		"Configs": configs,
	}

	log.Printf("✅ Builder: %d energy configs available", len(configs))

	// 3. ✅ Render template
	RenderTemplateOrError(w, "session-builder", data)
}

// ============================================
// 2️⃣ SESSION START (Démarrage)
// ============================================

func HandleStartSession(w http.ResponseWriter, r *http.Request) {
	// 1. Parse energy level
	energyStr := r.URL.Query().Get("energy")
	energy, err := strconv.Atoi(energyStr)
	if err != nil || energy < 1 || energy > 3 {
		log.Printf("⚠️ Invalid energy '%s', defaulting to 2 (medium)", energyStr)
		energy = 2 // Default medium
	}

	energyLevel := models.EnergyLevel(energy)
	log.Printf("🔍 StartSession: energy=%d (%s)", energy, energyLevel)

	// 2. Récupère exercices disponibles
	report, exercises, err := store.GetTodayReport()
	if err != nil {
		log.Printf("❌ GetTodayReport failed: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("📊 Available exercises: due=%d, new=%d, total=%d",
		report.TodayDue, report.TodayNew, len(exercises))

	// 3. Aucun exercice disponible
	if len(exercises) == 0 {
		log.Println("ℹ️ No exercises available, showing report")
		renderNoExercises(w, report)
		return
	}

	// 4. Applique limite énergie (domain)
	exerciseIDs := make([]int, len(exercises))
	for i, ex := range exercises {
		exerciseIDs[i] = ex.ID
	}

	limitedIDs := session.LimitExercises(exerciseIDs, energyLevel)
	maxExercises := session.GetMaxExercises(energyLevel)

	log.Printf("⚡ Energy limit applied: %d/%d exercises (max=%d for level %d)",
		len(limitedIDs), len(exerciseIDs), maxExercises, energy)

	// 5. Crée session (via service)
	sessionID, sessionData, err := sessionService.StartSession(energyLevel, limitedIDs)
	if err != nil {
		log.Printf("❌ StartSession failed: %v", err)
		http.Error(w, "Erreur création session", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Session #%d created: %d exercises queued", sessionID, len(sessionData.Exercises))

	// 6. Vérifie qu'il y a des exercices
	if len(sessionData.Exercises) == 0 {
		log.Printf("⚠️ Session #%d created but no exercises", sessionID)
		renderNoExercises(w, report)
		return
	}

	// 7. Redirige vers premier exercice
	firstExerciseID := sessionData.Exercises[0]
	redirectURL := fmt.Sprintf("/exercise/%d?from=session&session=%d",
		firstExerciseID, sessionID)

	log.Printf("🚀 Session #%d started → redirecting to exercise #%d", sessionID, firstExerciseID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// ============================================
// 3️⃣ SESSION COMPLETE (Résultats)
// ============================================

// HandleSessionComplete : Page de complétion
func HandleSessionComplete(w http.ResponseWriter, r *http.Request) {
	// 1. Parse session ID
	sessionIDStr := r.URL.Query().Get("id")

	// Fallback : récupère dernière session active
	if sessionIDStr == "" {
		log.Println("⚠️ No session ID provided, looking for active session")

		sessionID, err := sessionService.GetActiveSession()
		if err != nil || sessionID == 0 {
			log.Printf("❌ No active session found: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		sessionIDStr = fmt.Sprintf("%d", sessionID)
		log.Printf("✅ Found active session: #%s", sessionIDStr)
	}

	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("❌ Invalid session ID '%s': %v", sessionIDStr, err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Printf("🔍 SessionComplete: sessionID=%d", sessionID)

	// 2. Récupère résultat session
	result, err := sessionService.GetSessionResult(sessionID)
	if err != nil {
		log.Printf("❌ GetSessionResult failed for #%d: %v", sessionID, err)

		// Fallback : affiche page vide
		renderEmptySessionComplete(w, sessionID)
		return
	}

	log.Printf("✅ Session #%d results: completed=%d, duration=%s, exercises=%d",
		sessionID, result.CompletedCount, result.Duration.Round(time.Second), len(result.Exercises))

	// 3. Structure données
	data := map[string]any{
		"SessionID":      sessionID,
		"CompletedCount": result.CompletedCount,
		"Duration":       result.Duration.Round(time.Second),
		"DurationMin":    int(result.Duration.Minutes()),
		"CompletedAt":    result.CompletedAt.Format("15:04"),
		"ExerciseCount":  len(result.Exercises),
		"ExerciseIDs":    result.Exercises,
	}

	// 4. ✅ Render template
	RenderTemplateOrError(w, "session-complete", data)
}

// ============================================
// 4️⃣ SESSION STOP (Arrêt manuel)
// ============================================

// HandleStopSession : Arrête une session en cours
func HandleStopSession(w http.ResponseWriter, r *http.Request) {
	// 1. Parse session ID
	sessionIDStr := r.PathValue("id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("❌ Invalid session ID '%s': %v", sessionIDStr, err)
		http.Error(w, "ID de session invalide", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 StopSession: sessionID=%d", sessionID)

	// 2. Termine session
	if err := sessionService.StopSession(sessionID); err != nil {
		log.Printf("❌ StopSession failed for #%d: %v", sessionID, err)
		http.Error(w, "Erreur arrêt session", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Session #%d stopped manually", sessionID)

	// 3. Redirige vers dashboard
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ============================================
// HELPERS (Render Templates)
// ============================================

// renderNoExercises : Affiche message "aucun exercice disponible"
func renderNoExercises(w http.ResponseWriter, report models.SessionReport) {
	log.Println("📄 Rendering no-exercises-today template")

	data := map[string]any{
		"Message":         "🎉 Aucun exercice à réviser aujourd'hui !",
		"Report":          report,
		"TodayDue":        report.TodayDue,
		"TodayNew":        report.TodayNew,
		"NextReviewDate":  report.NextReviewDate,
		"UpcomingReviews": report.UpcomingReviews,
	}

	// ✅ Render avec helper
	RenderTemplateOrError(w, "no-exercises-today", data)
}

// renderEmptySessionComplete : Affiche page session vide (fallback)
func renderEmptySessionComplete(w http.ResponseWriter, sessionID int64) {
	log.Printf("⚠️ Rendering empty session-complete for #%d", sessionID)

	data := map[string]any{
		"SessionID":      sessionID,
		"CompletedCount": 0,
		"Duration":       time.Duration(0),
		"DurationMin":    0,
		"CompletedAt":    time.Now().Format("15:04"),
		"ExerciseCount":  0,
		"ExerciseIDs":    []int{},
	}

	// ✅ Render avec helper
	RenderTemplateOrError(w, "session-complete", data)
}
