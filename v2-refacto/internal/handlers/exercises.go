package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"maestro/v2-refacto/internal/models"
	"maestro/v2-refacto/internal/store"
	"maestro/v2-refacto/internal/validator"
)

// Vue : Page complète (affiche toute la structure HTML)
// Vue : Page complète exercices
func HandleExercisesPage(w http.ResponseWriter, r *http.Request) {
	allExercises := store.GetAll()

	data := map[string]any{
		"Exercises":     allExercises,
		"UrgentCount":   store.CountByView("urgent"),
		"TodayCount":    store.CountByView("today"),
		"UpcomingCount": store.CountByView("upcoming"),
		"ActiveCount":   store.CountByView("active"),
		"NewCount":      store.CountByView("new"),
	}

	if err := Tmpl.ExecuteTemplate(w, "exercise-list-page", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleListExercice reste inchangé (fragment pour filtres)

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

	filteredList := store.GetFiltered(filter)
	Tmpl.ExecuteTemplate(w, "exercise-list", filteredList)
}

// Vue : Détail (Fragment)
func HandleDetailExercice(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	ex := store.FindExercise(id)

	if ex == nil {
		http.NotFound(w, r)
		return
	}

	// ✅ Renvoie la PAGE COMPLÈTE au lieu du fragment
	if err := Tmpl.ExecuteTemplate(w, "exercise-detail-page", ex); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Action : Toggle Done
// Action : Toggle Done
func HandleToggleDone(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	ex := store.FindExercise(id)
	if ex == nil {
		http.NotFound(w, r)
		return
	}

	// Logique de transition
	if ex.Done {
		// Done → WIP (garde les CompletedSteps)
		ex.Done = false
	} else if len(ex.CompletedSteps) > 0 {
		// WIP → TODO (reset les étapes)
		ex.CompletedSteps = []int{}
	} else {
		// TODO → Done
		ex.Done = true
		// Optionnel : marque toutes les étapes comme complétées
		for i := range ex.Steps {
			ex.CompletedSteps = append(ex.CompletedSteps, i)
		}
	}

	store.Save()
	Tmpl.ExecuteTemplate(w, "exo-card", ex)
}

// Action : Toggle Step
func HandleToggleStep(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	step, _ := strconv.Atoi(r.URL.Query().Get("step"))

	// 2. SAS DE SÉCURITÉ (Nouveau)
	if err := validator.ValidateID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ex := store.FindExercise(id)
	if ex == nil {
		http.NotFound(w, r)
		return
	}

	// 1. Toggle de l'étape
	found := false
	for i, s := range ex.CompletedSteps {
		if s == step {
			ex.CompletedSteps = append(ex.CompletedSteps[:i], ex.CompletedSteps[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		ex.CompletedSteps = append(ex.CompletedSteps, step)
	}

	// ✅ SAUVEGARDE (Crucial)
	if err := store.Save(); err != nil {
		http.Error(w, "Erreur sauvegarde", http.StatusInternalServerError)
		return
	}
	Tmpl.ExecuteTemplate(w, "exercise-detail", *ex)
}

// Cycle: TODO → WIP → DONE → TODO
func HandleToggleStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	ex := store.FindExercise(id)
	if ex == nil {
		http.NotFound(w, r)
		return
	}

	if ex.Done {
		// DONE → TODO (reset)
		ex.Done = false
		ex.CompletedSteps = []int{}
	} else if len(ex.CompletedSteps) > 0 {
		// WIP → DONE
		ex.Done = true
		// Optionnel : marquer toutes les étapes
		ex.CompletedSteps = []int{}
		for i := range ex.Steps {
			ex.CompletedSteps = append(ex.CompletedSteps, i)
		}
	} else {
		// TODO → WIP (marque première étape)
		ex.CompletedSteps = append(ex.CompletedSteps, 0)
	}

	store.Save()
	Tmpl.ExecuteTemplate(w, "exo-card", ex)
}

// Action : Enregistrer une révision
func HandleReview(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	quality, _ := strconv.Atoi(r.URL.Query().Get("quality")) // 0=Oublié, 1=Dur, 3=Bien, 5=Facile

	ex := store.FindExercise(id)
	if ex == nil {
		http.NotFound(w, r)
		return
	}

	// Calcul SM-2 adapté
	now := time.Now()
	ex.LastReviewed = &now

	switch quality {
	case 0: // ❌ OUBLIÉ (Again)
		// Reset complet : retour à l'apprentissage actif
		ex.IntervalDays = 0                              // Révision dans 10 minutes (même session)
		ex.Repetitions = 0                               // Reset compteur
		ex.EaseFactor = math.Max(1.3, ex.EaseFactor-0.3) // Forte pénalité
		// Prochaine révision dans 10 minutes
		ex.NextReviewAt = now.Add(10 * time.Minute)

	case 1: // 😓 DUR
		ex.IntervalDays = 1
		ex.Repetitions++ // On compte quand même la répétition
		ex.EaseFactor = math.Max(1.3, ex.EaseFactor-0.2)
		ex.NextReviewAt = now.AddDate(0, 0, 1)

	case 3: // 🙂 BIEN
		if ex.IntervalDays == 0 {
			ex.IntervalDays = 1 // Première révision réussie
		} else {
			ex.IntervalDays = ex.IntervalDays * 2
		}
		ex.Repetitions++
		ex.NextReviewAt = now.AddDate(0, 0, ex.IntervalDays)

	case 5: // 😎 FACILE
		if ex.IntervalDays == 0 {
			ex.IntervalDays = 4 // Saute directement à 4 jours
		} else {
			ex.IntervalDays = ex.IntervalDays * 3
		}
		ex.Repetitions++
		ex.EaseFactor = math.Min(2.5, ex.EaseFactor+0.1)
		ex.NextReviewAt = now.AddDate(0, 0, ex.IntervalDays)
	}

	// Sauvegarde
	if err := store.Save(); err != nil {
		http.Error(w, "Erreur sauvegarde", http.StatusInternalServerError)
		return
	}

	Tmpl.ExecuteTemplate(w, "exercise-detail", ex)
}
