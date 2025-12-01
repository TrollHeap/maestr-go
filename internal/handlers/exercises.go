package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"maestro/internal/domain/exercise"
	"maestro/internal/domain/srs"
	"maestro/internal/models"
	"maestro/internal/service"
	"maestro/internal/store"
)

// ============================================
// SERVICES GLOBAUX
// ============================================

var (
	exerciseService *service.ExerciseService
	sessionService  *service.SessionService
)

func init() {
	exerciseService = service.NewExerciseService()
	sessionService = service.NewSessionService()
}

// ============================================
// 1️⃣ PAGE PRINCIPALE EXERCICES
// ============================================

func HandleExercisesPage(w http.ResponseWriter, r *http.Request) {
	// 1. Récupère données
	allExercises, err := exerciseService.GetAllExercises()
	if err != nil {
		log.Printf("❌ GetAllExercises error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	stats := exerciseService.GetExerciseStats()

	// 2. Structure données
	data := map[string]any{
		"Exercises":     allExercises,
		"UrgentCount":   stats["urgent"],
		"TodayCount":    stats["today"],
		"UpcomingCount": stats["upcoming"],
		"ActiveCount":   stats["active"],
		"NewCount":      stats["new"],
	}

	// 3. ✅ Render avec helper (gère errors automatiquement)
	RenderTemplateOrError(w, "exercise-list-page", data)
}

// ============================================
// 2️⃣ FRAGMENT HTMX : Liste filtrée
// ============================================

func HandleListExercice(w http.ResponseWriter, r *http.Request) {
	// 1. Parse query params
	view := r.URL.Query().Get("view")
	domain := r.URL.Query().Get("domain")
	difficulty, _ := strconv.Atoi(r.URL.Query().Get("difficulty"))

	// 2. Construit filtre
	filter := models.ExerciseFilter{
		View:       view,
		Domain:     domain,
		Difficulty: difficulty,
	}

	// 3. Récupère exercices filtrés
	filteredList, err := exerciseService.GetFilteredExercises(filter)
	if err != nil {
		log.Printf("❌ GetFilteredExercises error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// 4. ✅ Render fragment
	RenderTemplateOrError(w, "exercise-list", filteredList)
}

// ============================================
// 3️⃣ PAGE DÉTAIL EXERCICE
// ============================================

func HandleDetailExercice(w http.ResponseWriter, r *http.Request) {
	// 1. Parse params
	id, _ := strconv.Atoi(r.PathValue("id"))
	fromSession := r.URL.Query().Get("from") == "session"
	sessionIDStr := r.URL.Query().Get("session")

	log.Printf("🔍 DetailExercice: id=%d, fromSession=%v, sessionID=%s",
		id, fromSession, sessionIDStr)

	// 2. Validation ID
	if err := exercise.ValidateID(id); err != nil {
		log.Printf("❌ Invalid ID: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 3. Récupère exercice
	ex, err := exerciseService.GetExerciseWithMarkdown(id)
	if err != nil {
		log.Printf("❌ Exercice #%d non trouvé: %v", id, err)
		http.NotFound(w, r)
		return
	}

	// 4. Structure données (struct typé)
	data := struct {
		Exercise    *models.Exercise
		FromSession bool
		SessionID   string
	}{
		Exercise:    ex,
		FromSession: fromSession,
		SessionID:   sessionIDStr,
	}

	// 5. ✅ Render template
	RenderTemplateOrError(w, "exercise-detail-page", data)
}

// ============================================
// 4️⃣ ACTION : Toggle Done
// ============================================

func HandleToggleDone(w http.ResponseWriter, r *http.Request) {
	// 1. Parse params
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	fromSession := r.URL.Query().Get("from") == "session"
	sessionIDStr := r.URL.Query().Get("session")

	log.Printf("🔄 ToggleDone: id=%d, fromSession=%v, sessionID=%s",
		id, fromSession, sessionIDStr)

	// 2. Validation
	if err := exercise.ValidateID(id); err != nil {
		log.Printf("❌ Validation ID failed: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 3. Toggle via service
	ex, err := exerciseService.ToggleExerciseDone(id)
	if err != nil {
		log.Printf("❌ ToggleExerciseDone error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Exercice #%d toggled: Done=%v", id, ex.Done)

	// 4. MODE SESSION : Gestion flow
	if fromSession && ex.Done && sessionIDStr != "" {
		sessionID, _ := strconv.ParseInt(sessionIDStr, 10, 64)

		// a) Marque exercice complété
		if err := sessionService.CompleteExercise(sessionID, id, 3); err != nil {
			log.Printf("❌ CompleteExercise error: %v", err)
		}

		// b) Récupère prochain exercice
		nextEx, err := sessionService.GetNextExercise(sessionID)
		if err != nil {
			log.Printf("❌ GetNextExercise error: %v", err)
		}

		if nextEx != nil {
			// → Redirige vers exercice suivant
			redirectURL := fmt.Sprintf("/exercise/%d?from=session&session=%d",
				nextEx.ID, sessionID)
			log.Printf("➡️ Redirect to: %s", redirectURL)

			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		} else {
			// → Session terminée
			log.Println("✅ Session complete, no more exercises")

			if err := sessionService.EndSession(sessionID); err != nil {
				log.Printf("❌ EndSession error: %v", err)
			}

			http.Redirect(w, r, fmt.Sprintf("/session/complete?id=%d", sessionID),
				http.StatusSeeOther)
			return
		}
	}

	// 5. MODE NORMAL : Render fragment status
	RenderTemplateOrError(w, "status-indicator", ex)
}

// ============================================
// 5️⃣ ACTION : Toggle Step
// ============================================

func HandleToggleStep(w http.ResponseWriter, r *http.Request) {
	// 1. Parse params
	id, _ := strconv.Atoi(r.PathValue("id"))
	step, _ := strconv.Atoi(r.URL.Query().Get("step"))

	log.Printf("🔄 ToggleStep: id=%d, step=%d", id, step)

	// 2. Validation ID
	if err := exercise.ValidateID(id); err != nil {
		log.Printf("❌ Invalid ID: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 3. Récupère exercice
	ex, err := exerciseService.GetExerciseWithMarkdown(id)
	if err != nil {
		log.Printf("❌ Exercise #%d not found: %v", id, err)
		http.NotFound(w, r)
		return
	}

	// 4. Validation step
	if err := exercise.ValidateStep(step, len(ex.Steps)); err != nil {
		log.Printf("❌ Invalid step: %v", err)
		http.Error(w, "Step invalide", http.StatusBadRequest)
		return
	}

	// 5. Toggle step
	ex, err = exerciseService.ToggleExerciseStep(id, step)
	if err != nil {
		log.Printf("❌ ToggleExerciseStep error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Step #%d toggled", step)

	// 6. ✅ Render fragment (HTMX swap)
	RenderTemplateOrError(w, "steps-exo", ex)
}

// ============================================
// 6️⃣ ACTION : Review (SRS)
// ============================================

func HandleReview(w http.ResponseWriter, r *http.Request) {
	// 1. Parse params
	id, _ := strconv.Atoi(r.PathValue("id"))
	quality, _ := strconv.Atoi(r.URL.Query().Get("quality"))
	fromSession := r.URL.Query().Get("from") == "session"
	sessionIDStr := r.URL.Query().Get("session")

	log.Printf("🔍 [Review] id=%d, quality=%d, fromSession=%v, sessionID=%s",
		id, quality, fromSession, sessionIDStr)

	// 2. Validation ID
	if err := exercise.ValidateID(id); err != nil {
		log.Printf("❌ Validation ID failed: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 3. Validation quality (0-5)
	if err := exercise.ValidateQuality(quality); err != nil {
		log.Printf("❌ Validation Quality failed: %v", err)
		http.Error(w, "Quality invalide", http.StatusBadRequest)
		return
	}

	// 4. Applique algorithme SRS
	ex, err := exerciseService.ReviewExercise(id, srs.ReviewQuality(quality))
	if err != nil {
		log.Printf("❌ ReviewExercise error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Review applied: ease=%.2f, nextReview=%s",
		ex.EaseFactor, ex.NextReviewAt.Format("2006-01-02"))

	// 5. Marque DONE si quality >= 1
	if quality >= 1 {
		ex.Done = true
		if err := store.SaveExercise(ex); err != nil {
			log.Printf("❌ SaveExercise error: %v", err)
			http.Error(w, "Erreur sauvegarde", http.StatusInternalServerError)
			return
		}
		log.Printf("✅ Exercise marked DONE")
	}

	// 6. MODE SESSION : Flow exercice suivant
	if fromSession && sessionIDStr != "" {
		sessionID, _ := strconv.ParseInt(sessionIDStr, 10, 64)

		log.Printf("🔄 Session mode: sessionID=%d", sessionID)

		// a) Enregistre dans session
		if err := sessionService.CompleteExercise(sessionID, id, quality); err != nil {
			log.Printf("❌ CompleteExercise error: %v", err)
		} else {
			log.Printf("✅ Exercise completed in session")
		}

		// b) Prochain exercice
		nextEx, err := sessionService.GetNextExercise(sessionID)
		if err != nil {
			log.Printf("❌ GetNextExercise error: %v", err)
		}

		if nextEx != nil {
			// → Redirection HTMX vers exercice suivant
			redirectURL := fmt.Sprintf("/exercise/%d?from=session&session=%d",
				nextEx.ID, sessionID)
			log.Printf("➡️ HX-Redirect to: %s", redirectURL)

			w.Header().Set("HX-Redirect", redirectURL)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			// → Session terminée
			log.Println("✅ Session complete, no more exercises")

			if err := sessionService.EndSession(sessionID); err != nil {
				log.Printf("❌ EndSession error: %v", err)
			}

			w.Header().Set("HX-Redirect", fmt.Sprintf("/session/complete?id=%d", sessionID))
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// 7. MODE LIBRE : Recharge détail exercice
	log.Println("🔄 Free mode, reload detail")

	data := struct {
		Exercise    *models.Exercise
		FromSession bool
		SessionID   string
	}{
		Exercise:    ex,
		FromSession: fromSession,
		SessionID:   sessionIDStr,
	}

	// ✅ Render template
	RenderTemplateOrError(w, "exercise-detail", data)
}

// ============================================
// 7️⃣ ACTION : Prochain exercice (Mode libre)
// ============================================

func HandleNextExercise(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 HandleNextExercise: Free mode")

	// 1. Récupère rapport + exercices disponibles
	report, exercises, err := store.GetTodayReport()
	if err != nil {
		log.Printf("❌ GetTodayReport error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// 2. Aucun exercice disponible → Affiche page "No exercises"
	if len(exercises) == 0 {
		log.Println("ℹ️ No exercises available, showing report")

		data := map[string]any{
			"Message":         "🎉 Aucun exercice à réviser aujourd'hui !",
			"Report":          report,
			"TodayDue":        report.TodayDue,
			"TodayNew":        report.TodayNew,
			"NextReviewDate":  report.NextReviewDate,
			"UpcomingReviews": report.UpcomingReviews,
		}

		// ✅ Render template
		RenderTemplateOrError(w, "no-exercises-today", data)
		return
	}

	// 3. Exercice(s) disponible(s) → Redirige vers le plus urgent
	nextExercise := exercises[0]
	redirectURL := fmt.Sprintf("/exercise/%d", nextExercise.ID)

	log.Printf("➡️ Redirect to exercise #%d: %s", nextExercise.ID, nextExercise.Title)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
