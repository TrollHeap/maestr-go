package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"maestro/internal/domain/session" // ✅ NOUVEAU
	"maestro/internal/models"
	"maestro/internal/store"
)

// ============================================
// SESSION BUILDER (Choix énergie)
// ============================================

// HandleSessionBuilder : Page de sélection d'énergie
func HandleSessionBuilder(w http.ResponseWriter, r *http.Request) {
	// ✅ Utilise domain configs au lieu de models
	configs := []session.Config{
		session.GetConfig(models.EnergyLow),
		session.GetConfig(models.EnergyMedium),
		session.GetConfig(models.EnergyHigh),
	}

	data := map[string]any{
		"Configs": configs,
	}

	if err := Tmpl.ExecuteTemplate(w, "session-builder", data); err != nil {
		log.Printf("❌ Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ============================================
// SESSION START (Démarrage)
// ============================================

func HandleStartSession(w http.ResponseWriter, r *http.Request) {
	// 1. PARSE energy
	energyStr := r.URL.Query().Get("energy")
	energy, err := strconv.Atoi(energyStr)
	if err != nil || energy < 1 || energy > 3 {
		energy = 2 // Default medium
	}

	energyLevel := models.EnergyLevel(energy)
	log.Printf("🔍 START SESSION: energy=%d", energy)

	// 2. RÉCUPÈRE EXERCICES DISPONIBLES
	report, exercises, err := store.GetTodayReport()
	if err != nil {
		log.Printf("❌ GetTodayReport failed: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("🔍 [SESSION] Disponibles: %d dus + %d nouveaux = %d total",
		report.TodayDue, report.TodayNew, len(exercises))

	// 3. AUCUN EXERCICE ? Affiche rapport
	if len(exercises) == 0 {
		renderNoExercises(w, report)
		return
	}

	// 4. ✅ APPLIQUE LIMITE ÉNERGIE (domain)
	exerciseIDs := make([]int, len(exercises))
	for i, ex := range exercises {
		exerciseIDs[i] = ex.ID
	}

	limitedIDs := session.LimitExercises(exerciseIDs, energyLevel)

	log.Printf("🔍 [SESSION] Limité à %d exercices (max=%d pour energy=%d)",
		len(limitedIDs),
		session.GetMaxExercises(energyLevel),
		energy,
	)

	// 5. CRÉE SESSION (via service)
	sessionID, sessionData, err := sessionService.StartSession(energyLevel, limitedIDs)
	if err != nil {
		log.Printf("❌ StartSession failed: %v", err)
		http.Error(w, "Erreur création session", http.StatusInternalServerError)
		return
	}

	// 6. REDIRIGE vers premier exercice
	if len(sessionData.Exercises) == 0 {
		log.Printf("⚠️ Session created but no exercises")
		renderNoExercises(w, report)
		return
	}

	firstExerciseID := sessionData.Exercises[0]
	redirectURL := fmt.Sprintf("/exercise/%d?from=session&session=%d",
		firstExerciseID, sessionID)

	log.Printf("🚀 Session %d started → exo #%d", sessionID, firstExerciseID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// ============================================
// SESSION COMPLETE (Résultats)
// ============================================

// HandleSessionComplete : Page de complétion
func HandleSessionComplete(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.URL.Query().Get("id")

	// Si pas d'ID fourni, essaie de récupérer la dernière session active
	if sessionIDStr == "" {
		sessionID, err := sessionService.GetActiveSession()
		if err != nil || sessionID == 0 {
			log.Printf("⚠️ No active session found")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		sessionIDStr = fmt.Sprintf("%d", sessionID)
	}

	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("❌ Invalid session ID: %s", sessionIDStr)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Récupère résultat depuis service
	result, err := sessionService.GetSessionResult(sessionID)
	if err != nil {
		log.Printf("❌ GetSessionResult failed: %v", err)

		// Fallback : affiche page vide
		renderEmptySessionComplete(w)
		return
	}

	// Affiche résultats
	data := map[string]any{
		"SessionID":      sessionID,
		"CompletedCount": result.CompletedCount,
		"Duration":       result.Duration.Round(time.Second),
		"CompletedAt":    result.CompletedAt.Format("15:04"),
		"ExerciseCount":  len(result.Exercises),
		"ExerciseIDs":    result.Exercises,
	}

	if err := Tmpl.ExecuteTemplate(w, "session-complete", data); err != nil {
		log.Printf("❌ Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ============================================
// SESSION STOP (Arrêt manuel)
// ============================================

// HandleStopSession : Arrête une session en cours
func HandleStopSession(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.PathValue("id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("❌ Invalid session ID: %s", sessionIDStr)
		http.Error(w, "ID de session invalide", http.StatusBadRequest)
		return
	}

	// Termine la session
	if err := sessionService.StopSession(sessionID); err != nil {
		log.Printf("❌ StopSession failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Session %d stopped manually", sessionID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ============================================
// HELPERS (Render Templates)
// ============================================

// renderNoExercises : Affiche message "aucun exercice disponible"
func renderNoExercises(w http.ResponseWriter, report models.SessionReport) {
	data := map[string]interface{}{
		"Message":         "🎉 Aucun exercice à réviser aujourd'hui !",
		"Report":          report,
		"TodayDue":        report.TodayDue,
		"TodayNew":        report.TodayNew,
		"NextReviewDate":  report.NextReviewDate,
		"UpcomingReviews": report.UpcomingReviews,
	}

	if err := Tmpl.ExecuteTemplate(w, "no-exercises-today", data); err != nil {
		log.Printf("❌ Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderEmptySessionComplete : Affiche page session vide (fallback)
func renderEmptySessionComplete(w http.ResponseWriter) {
	data := map[string]any{
		"SessionID":      0,
		"CompletedCount": 0,
		"Duration":       0,
		"CompletedAt":    time.Now().Format("15:04"),
		"ExerciseCount":  0,
		"ExerciseIDs":    []int{},
	}

	if err := Tmpl.ExecuteTemplate(w, "session-complete", data); err != nil {
		log.Printf("❌ Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
