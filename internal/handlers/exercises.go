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

	if err := validator.ValidateID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validator.ValidateQuality(quality); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Appel au service SRS
	ex, err := exerciseService.ReviewExercise(id, srs.ReviewQuality(quality))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marque DONE si qualité >= 1 (tous sauf "Oublié")
	if quality >= 1 {
		ex.Done = true
		if err := store.Save(); err != nil {
			http.Error(w, "Erreur sauvegarde", http.StatusInternalServerError)
			return
		}
	}
	// Si quality = 0 : Done reste false, mais on continue en session

	// 🎯 MODE SESSION : Passer au suivant dans TOUS les cas
	if fromSession {
		activeSession := sessionService.GetActiveSession()
		if activeSession != nil {
			activeSession.MarkCompleted(ex.ID)
			nextEx := activeSession.NextExercise()

			if nextEx != nil {
				// Redirection vers exercice suivant
				redirectURL := fmt.Sprintf("/exercise/%d?from=session", nextEx.ID)
				w.Header().Set("HX-Redirect", redirectURL)
				w.WriteHeader(http.StatusOK)
				return
			} else {
				// ✅ Session terminée : stocke le résultat
				result := &models.SessionResult{
					CompletedCount: len(activeSession.CompletedIDs),
					Duration:       time.Since(activeSession.StartedAt),
					CompletedAt:    time.Now(),
					Exercises:      activeSession.CompletedIDs,
				}

				store.StoreSessionResult(result)
				sessionService.ClearAllSessions()

				// Redirect vers page de complétion
				w.Header().Set("HX-Redirect", "/session/complete")
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	// 🎯 MODE LIBRE : Recharge la page de détail
	view := struct {
		Exercise    *models.Exercise
		FromSession bool
	}{
		Exercise:    ex,
		FromSession: fromSession,
	}
	Tmpl.ExecuteTemplate(w, "exercise-detail", view)
}

func HandleSessionNext(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sid")
	currentID, _ := strconv.Atoi(r.URL.Query().Get("current"))

	activeSession := sessionService.GetActiveSession()
	if activeSession == nil || activeSession.ID != sessionID {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Marque l'exercice actuel comme complété
	activeSession.MarkCompleted(currentID)

	// Récupère le suivant
	nextEx := activeSession.NextExercise()

	if nextEx != nil {
		// Redirige vers exercice suivant
		http.Redirect(w, r,
			fmt.Sprintf("/exercise/%d?from=session&sid=%s", nextEx.ID, sessionID),
			http.StatusSeeOther)
	} else {
		// Session terminée
		sessionService.ClearAllSessions()
		data := map[string]any{
			"CompletedCount": len(activeSession.CompletedIDs),
			"Duration":       time.Since(activeSession.StartedAt).Round(time.Minute),
		}
		Tmpl.ExecuteTemplate(w, "session-complete", data)
	}
}

// HandleNextExercise retourne le prochain exercice à réviser
func HandleNextExercise(w http.ResponseWriter, r *http.Request) {
	fromSession := r.URL.Query().Get("from") == "session"

	// Récupère les IDs de la session active (si applicable)
	var sessionExercises []int
	if fromSession {
		activeSession := sessionService.GetActiveSession()
		if activeSession != nil {
			for _, ex := range activeSession.Session.Exercises {
				sessionExercises = append(sessionExercises, ex.ID)
			}
		}
	}

	// 1️⃣ Récupère le prochain exercice
	nextExercise, err := store.GetNextDueExercise(fromSession, sessionExercises)
	if err != nil {
		log.Println("Erreur GetNextDueExercise:", err)
		http.Error(w, "Erreur serveur", 500)
		return
	}

	// 2️⃣ Aucun exercice disponible
	if nextExercise == nil {
		// Template de fin de session
		data := map[string]any{
			"Message": "🎉 Plus d'exercices à réviser !",
		}
		Tmpl.ExecuteTemplate(w, "no-more-exercises", data)
		return
	}

	// 3️⃣ Construit les données template
	view := models.ExerciseView{
		Exercise:    nextExercise,
		FromSession: fromSession,
	}

	// 4️⃣ Exécute le template
	if err := Tmpl.ExecuteTemplate(w, "exercise-detail", view); err != nil {
		log.Println("Erreur template next exercise:", err)
		http.Error(w, "Erreur serveur", 500)
	}
}
