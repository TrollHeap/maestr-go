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

	"github.com/russross/blackfriday/v2"
)

// Service global (à initialiser dans main.go)
var exerciseService *service.ExerciseService

func init() {
	exerciseService = service.NewExerciseService()
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

	// Validation
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

// Action : Toggle Done
func HandleToggleDone(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	// Validation
	if err := validator.ValidateID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Appel au service
	ex, err := exerciseService.ToggleExerciseDone(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Tmpl.ExecuteTemplate(w, "exo-card", ex)
}

// Action : Toggle Step
func HandleToggleStep(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	step, _ := strconv.Atoi(r.URL.Query().Get("step"))

	// Validation
	if err := validator.ValidateID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Récupère l'exercice pour valider le step
	ex, err := exerciseService.GetExerciseWithMarkdown(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Valide le step par rapport au nombre d'étapes
	if err := validator.ValidateStep(step, len(ex.Steps)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Appel au service
	ex, err = exerciseService.ToggleExerciseStep(id, step)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Tmpl.ExecuteTemplate(w, "exercise-detail", *ex)
}

// Cycle: TODO → WIP → DONE → TODO
func HandleToggleStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	fromSession := r.URL.Query().Get("from") == "session"

	// Validation
	if err := validator.ValidateID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Appel au service
	ex, err := exerciseService.ToggleExerciseStatus(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Gestion session (logique coordination)
	if fromSession {
		activeSession := models.GetActiveSession()

		if activeSession != nil && ex.Done {
			activeSession.MarkCompleted(ex.ID)
			nextEx := activeSession.NextExercise()

			if nextEx != nil {
				http.Redirect(
					w,
					r,
					fmt.Sprintf("/exercise/%d?from=session", nextEx.ID),
					http.StatusSeeOther,
				)
				return
			} else {
				// Session terminée
				models.ClearActiveSession()

				data := map[string]any{
					"CompletedCount": len(activeSession.CompletedIDs),
					"Duration":       time.Since(activeSession.StartedAt).Round(time.Minute),
				}

				Tmpl.ExecuteTemplate(w, "session-complete", data)
				return
			}
		}
	}

	// Mode libre : conversion markdown
	htmlContent := blackfriday.Run([]byte(ex.Content))
	ex.Content = string(htmlContent)

	Tmpl.ExecuteTemplate(w, "exercise-detail", ex)
}

// Action : Enregistrer une révision
func HandleReview(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	quality, _ := strconv.Atoi(r.URL.Query().Get("quality"))

	// Validation
	if err := validator.ValidateID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validator.ValidateQuality(quality); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Appel au service avec conversion vers ReviewQuality
	ex, err := exerciseService.ReviewExercise(id, srs.ReviewQuality(quality))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Tmpl.ExecuteTemplate(w, "exercise-detail", ex)
}
