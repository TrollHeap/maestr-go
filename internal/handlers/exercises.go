package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"maestro/internal/models"
	"maestro/internal/service"
	"maestro/internal/srs"
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

// Vue : Page complète exercices
func HandleExercisesPage(w http.ResponseWriter, r *http.Request) {
	allExercises := exerciseService.GetAllExercises()
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

// Vue : Liste seule (Fragment)
func HandleListExercice(w http.ResponseWriter, r *http.Request) {
	view := r.URL.Query().Get("view")
	domain := r.URL.Query().Get("domain")
	difficulty, _ := strconv.Atoi(r.URL.Query().Get("difficulty"))

	filter := models.ExerciseFilter{
		View:       view,
		Domain:     domain,
		Difficulty: difficulty,
	}

	filteredList := exerciseService.GetFilteredExercises(filter)
	Tmpl.ExecuteTemplate(w, "exercise-list", filteredList)
}

// Vue : Détail (Fragment)
func HandleDetailExercice(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := validator.ValidateID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ex, err := exerciseService.GetExerciseWithMarkdown(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	view := struct {
		Exercise    *models.Exercise
		FromSession bool
	}{
		Exercise:    ex,
		FromSession: r.URL.Query().Get("from") == "session",
	}

	if err := Tmpl.ExecuteTemplate(w, "exercise-detail-page", view); err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Action : Toggle Done (TODO ↔ DONE avec gestion session)
func HandleToggleDone(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	fromSession := r.URL.Query().Get("from") == "session"

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
	if fromSession && ex.Done {
		activeSession := sessionService.GetActiveSession()
		if activeSession != nil {
			activeSession.MarkCompleted(ex.ID)
			nextEx := activeSession.NextExercise()

			if nextEx != nil {
				// Passer au suivant
				http.Redirect(w, r,
					fmt.Sprintf("/exercise/%d?from=session", nextEx.ID),
					http.StatusSeeOther)
				return
			} else {
				// Session terminée
				sessionService.ClearAllSessions()
				data := map[string]any{
					"CompletedCount": len(activeSession.CompletedIDs),
					"Duration":       time.Since(activeSession.StartedAt).Round(time.Minute),
				}
				Tmpl.ExecuteTemplate(w, "session-complete", data)
				return
			}
		}
	}

	// Mode normal : renvoie juste le bouton status
	Tmpl.ExecuteTemplate(w, "status-indicator", ex)
}

// Action : Toggle Step
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

// Action : Enregistrer une révision
func HandleReview(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	quality, _ := strconv.Atoi(r.URL.Query().Get("quality"))
	fromSession := r.URL.Query().Get("from") == "session"

	// Validation
	if err := validator.ValidateID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validator.ValidateQuality(quality); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Appel SRS
	ex, err := exerciseService.ReviewExercise(id, srs.ReviewQuality(quality))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 🎯 MODE SESSION : Passer au suivant
	if fromSession {
		activeSession := sessionService.GetActiveSession()
		if activeSession != nil {
			nextEx := activeSession.NextExercise()

			if nextEx != nil {
				// Redirection vers exercice suivant
				http.Redirect(w, r,
					fmt.Sprintf("/exercise/%d?from=session", nextEx.ID),
					http.StatusSeeOther)
				return
			} else {
				// Session terminée
				sessionService.ClearAllSessions()
				data := map[string]any{
					"CompletedCount": len(activeSession.CompletedIDs),
					"Duration":       time.Since(activeSession.StartedAt).Round(time.Minute),
				}
				Tmpl.ExecuteTemplate(w, "session-complete", data)
				return
			}
		}
	}

	// 🎯 MODE LIBRE : Affiche confirmation et mise à jour
	Tmpl.ExecuteTemplate(w, "review-panel", ex)
}
