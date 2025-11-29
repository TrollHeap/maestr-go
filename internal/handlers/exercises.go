package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"maestro/internal/models"
	"maestro/internal/service"
	"maestro/internal/srs"
	"maestro/internal/store"
	"maestro/internal/validator"
)

// Services globaux
var (
	exerciseService *service.ExerciseService
	sessionService  *service.SessionService
)

func init() {
	exerciseService = service.NewExerciseService()
	sessionService = service.NewSessionService()
}

// ============================================
// PAGES COMPLÈTES
// ============================================

// HandleExercisesPage : Page principale exercices
func HandleExercisesPage(w http.ResponseWriter, r *http.Request) {
	allExercises, err := exerciseService.GetAllExercises()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats := exerciseService.GetExerciseStats()

	data := map[string]any{
		"Exercises":     allExercises,
		"UrgentCount":   stats["urgent"],
		"TodayCount":    stats["today"],
		"UpcomingCount": stats["upcoming"],
		"ActiveCount":   stats["active"],
		"NewCount":      stats["new"],
	}

	if err := Tmpl.ExecuteTemplate(w, "exercise-list-page", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ============================================
// FRAGMENTS HTMX
// ============================================

// HandleListExercice : Liste filtrée (Fragment)
func HandleListExercice(w http.ResponseWriter, r *http.Request) {
	view := r.URL.Query().Get("view")
	domain := r.URL.Query().Get("domain")
	difficulty, _ := strconv.Atoi(r.URL.Query().Get("difficulty"))

	filter := models.ExerciseFilter{
		View:       view,
		Domain:     domain,
		Difficulty: difficulty,
	}

	filteredList, err := exerciseService.GetFilteredExercises(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Tmpl.ExecuteTemplate(w, "exercise-list", filteredList)
}

// HandleDetailExercice : Détail d'un exercice (Fragment)
func HandleDetailExercice(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	sessionIDStr := r.URL.Query().Get("session")

	if err := validator.ValidateID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ex, err := exerciseService.GetExerciseWithMarkdown(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// ✅ AJOUTE SessionID dans les données
	view := struct {
		Exercise    *models.Exercise
		FromSession bool
		SessionID   string // ← NOUVEAU
	}{
		Exercise:    ex,
		FromSession: r.URL.Query().Get("from") == "session",
		SessionID:   sessionIDStr, // ← NOUVEAU
	}

	if err := Tmpl.ExecuteTemplate(w, "exercise-detail-page", view); err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ============================================
// ACTIONS (POST)
// ============================================

// HandleToggleDone : Toggle TODO ↔ DONE
func HandleToggleDone(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	fromSession := r.URL.Query().Get("from") == "session"
	sessionIDStr := r.URL.Query().Get("session")

	if err := validator.ValidateID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Toggle via service
	ex, err := exerciseService.ToggleExerciseDone(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Gestion session si applicable
	if fromSession && ex.Done && sessionIDStr != "" {
		sessionID, _ := strconv.ParseInt(sessionIDStr, 10, 64)

		// Marque exercice complété dans la session
		if err := sessionService.CompleteExercise(sessionID, id, 3); err != nil {
			log.Printf("Erreur complete exercise: %v", err)
		}

		// Récupère prochain exercice
		nextEx, err := sessionService.GetNextExercise(sessionID)
		if err != nil {
			log.Printf("Erreur next exercise: %v", err)
		}

		if nextEx != nil {
			// Redirige vers suivant
			redirectURL := fmt.Sprintf("/exercise/%d?from=session&session=%d", nextEx.ID, sessionID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		} else {
			// Session terminée
			if err := sessionService.EndSession(sessionID); err != nil {
				log.Printf("Erreur end session: %v", err)
			}

			// Redirige vers page de complétion
			http.Redirect(w, r, fmt.Sprintf("/session/complete?id=%d", sessionID), http.StatusSeeOther)
			return
		}
	}

	// Mode normal : renvoie juste le bouton status
	Tmpl.ExecuteTemplate(w, "status-indicator", ex)
}

// HandleToggleStep : Toggle une étape individuelle
func HandleToggleStep(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	step, _ := strconv.Atoi(r.URL.Query().Get("step"))

	if err := validator.ValidateID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ex, err := exerciseService.GetExerciseWithMarkdown(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := validator.ValidateStep(step, len(ex.Steps)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ex, err = exerciseService.ToggleExerciseStep(id, step)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Tmpl.ExecuteTemplate(w, "steps-exo", ex)
}

// HandleReview : Enregistre une révision SRS
func HandleReview(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	quality, _ := strconv.Atoi(r.URL.Query().Get("quality"))
	fromSession := r.URL.Query().Get("from") == "session"
	sessionIDStr := r.URL.Query().Get("session")

	// ✅ AJOUTE CES LOGS
	log.Printf("🔍 [HandleReview] exID=%d, quality=%d, from=%v, session=%s",
		id, quality, fromSession, sessionIDStr)

	if err := validator.ValidateID(id); err != nil {
		log.Printf("❌ Validation ID échouée: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.ValidateQuality(quality); err != nil {
		log.Printf("❌ Validation Quality échouée: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Applique SRS
	ex, err := exerciseService.ReviewExercise(id, srs.ReviewQuality(quality))
	if err != nil {
		log.Printf("❌ Erreur ReviewExercise: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Review appliquée: newEaseFactor=%.2f, nextReview=%s",
		ex.EaseFactor, ex.NextReviewAt.Format("2006-01-02"))

	// Marque DONE si qualité >= 1
	if quality >= 1 {
		ex.Done = true
		if err := store.SaveExercise(ex); err != nil {
			log.Printf("❌ Erreur SaveExercise: %v", err)
			http.Error(w, "Erreur sauvegarde", http.StatusInternalServerError)
			return
		}
		log.Printf("✅ Exercice marqué DONE")
	}

	// MODE SESSION : Passer au suivant
	if fromSession && sessionIDStr != "" {
		sessionID, _ := strconv.ParseInt(sessionIDStr, 10, 64)

		log.Printf("🔄 Mode session: sessionID=%d", sessionID)

		// Enregistre dans session
		if err := sessionService.CompleteExercise(sessionID, id, quality); err != nil {
			log.Printf("❌ Erreur CompleteExercise: %v", err)
		} else {
			log.Printf("✅ Exercice complété dans session")
		}

		// Prochain exercice
		nextEx, err := sessionService.GetNextExercise(sessionID)
		if err != nil {
			log.Printf("❌ Erreur GetNextExercise: %v", err)
		}

		if nextEx != nil {
			// Redirection HTMX
			redirectURL := fmt.Sprintf("/exercise/%d?from=session&session=%d", nextEx.ID, sessionID)
			log.Printf("➡️ Redirection vers: %s", redirectURL)

			w.Header().Set("HX-Redirect", redirectURL)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			// Session terminée
			log.Println("✅ Session terminée, aucun exercice suivant")

			if err := sessionService.EndSession(sessionID); err != nil {
				log.Printf("❌ Erreur EndSession: %v", err)
			}

			w.Header().Set("HX-Redirect", fmt.Sprintf("/session/complete?id=%d", sessionID))
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// MODE LIBRE : Recharge détail
	log.Println("🔄 Mode libre, recharge détail")
	view := struct {
		Exercise    *models.Exercise
		FromSession bool
		SessionID   string
	}{
		Exercise:    ex,
		FromSession: fromSession,
		SessionID:   sessionIDStr,
	}
	Tmpl.ExecuteTemplate(w, "exercise-detail", view)
}

// HandleNextExercise : Prochain exercice à réviser
func HandleNextExercise(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 HandleNextExercise: Mode libre")

	// 1. Récupère rapport complet + exercices disponibles aujourd'hui
	report, exercises, err := store.GetTodayReport()
	if err != nil {
		log.Printf("❌ Erreur GetTodayReport: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// 2. Aucun exercice disponible aujourd'hui → Affiche rapport
	if len(exercises) == 0 {
		log.Println("ℹ️ Aucun exercice disponible, affichage du rapport")

		data := map[string]interface{}{
			"Message":         "🎉 Aucun exercice à réviser aujourd'hui !",
			"Report":          report,
			"TodayDue":        report.TodayDue,
			"TodayNew":        report.TodayNew,
			"NextReviewDate":  report.NextReviewDate,
			"UpcomingReviews": report.UpcomingReviews,
		}

		if err := Tmpl.ExecuteTemplate(w, "no-exercises-today", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// 3. Exercice(s) disponible(s) → Redirige vers le premier (plus urgent)
	nextExercise := exercises[0]
	redirectURL := fmt.Sprintf("/exercise/%d", nextExercise.ID)

	log.Printf("➡️ Redirection vers exercice #%d: %s", nextExercise.ID, nextExercise.Title)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
