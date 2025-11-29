package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"maestro/internal/models"
	"maestro/internal/store"
)

// ============================================
// SESSION BUILDER (Choix énergie)
// ============================================

// HandleSessionBuilder : Page de sélection d'énergie
func HandleSessionBuilder(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Configs": models.SessionConfigs,
	}

	if err := Tmpl.ExecuteTemplate(w, "session-builder", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ============================================
// SESSION START (Démarrage)
// ============================================

func HandleStartSession(w http.ResponseWriter, r *http.Request) {
	energyStr := r.URL.Query().Get("energy")
	energy, err := strconv.Atoi(energyStr)
	if err != nil || energy < 1 || energy > 3 {
		energy = 2 // Default medium
	}

	log.Printf("🔍 START SESSION: energy=%d", energy)

	// Récupère exercices DISPONIBLES
	report, _, err := store.GetTodayReport()
	if err != nil {
		log.Printf("❌ Erreur GetTodayReport: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// ✅ LOG DÉTAILLÉ
	log.Printf("🔍 SESSION disponible: %d dus + %d nouveaux = %d total",
		report.TodayDue, report.TodayNew, report.TotalAvailable)

	// AUCUN exercice ? Affiche rapport
	if report.TotalAvailable == 0 {
		data := map[string]interface{}{
			"Message":         "🎉 Aucun exercice à réviser aujourd'hui !",
			"Report":          report,
			"TodayDue":        report.TodayDue,
			"TodayNew":        report.TodayNew,
			"NextReviewDate":  report.NextReviewDate,
			"UpcomingReviews": report.UpcomingReviews,
		}
		Tmpl.ExecuteTemplate(w, "no-exercises-today", data)
		return
	}

	// CRÉE SESSION
	sessionID, session, err := sessionService.StartSession(models.EnergyLevel(energy))
	if err != nil {
		if noExErr, ok := err.(*models.NoExercisesTodayError); ok {
			data := map[string]interface{}{
				"Message":         "🎉 Aucun exercice à réviser aujourd'hui !",
				"Report":          noExErr.Report,
				"TodayDue":        noExErr.Report.TodayDue,
				"TodayNew":        noExErr.Report.TodayNew,
				"NextReviewDate":  noExErr.Report.NextReviewDate,
				"UpcomingReviews": noExErr.Report.UpcomingReviews,
			}
			Tmpl.ExecuteTemplate(w, "no-exercises-today", data)
			return
		}
		log.Printf("❌ Erreur StartSession: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirige vers premier exercice
	firstExercise := session.Exercises[0]
	redirectURL := fmt.Sprintf("/exercise/%d?from=session&session=%d",
		firstExercise.ID, sessionID)
	log.Printf("🚀 Session %d démarrée → exo #%d '%s'",
		sessionID, firstExercise.ID, firstExercise.Title)
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
			log.Printf("Aucune session active trouvée")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		sessionIDStr = fmt.Sprintf("%d", sessionID)
	}

	sessionID, _ := strconv.ParseInt(sessionIDStr, 10, 64)

	// Récupère résultat depuis SQLite
	result, err := sessionService.GetSessionResult(sessionID)
	if err != nil {
		log.Printf("Erreur récupération résultat: %v", err)

		// Fallback : affiche page vide
		data := map[string]any{
			"CompletedCount": 0,
			"Duration":       0,
			"CompletedAt":    time.Now().Format("15:04"),
			"ExerciseCount":  0,
		}
		Tmpl.ExecuteTemplate(w, "session-complete", data)
		return
	}

	// Affiche résultats
	data := map[string]any{
		"CompletedCount": result.CompletedCount,
		"Duration":       result.Duration.Round(time.Second),
		"CompletedAt":    result.CompletedAt.Format("15:04"),
		"ExerciseCount":  len(result.Exercises),
	}

	if err := Tmpl.ExecuteTemplate(w, "session-complete", data); err != nil {
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
		http.Error(w, "ID de session invalide", http.StatusBadRequest)
		return
	}

	// Termine la session
	if err := sessionService.StopSession(sessionID); err != nil {
		log.Printf("Erreur stop session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Session %d arrêtée", sessionID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Dans internal/service/session.go (ou équivalent)
