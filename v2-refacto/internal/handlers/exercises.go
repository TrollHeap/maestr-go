package handlers

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"maestro/v2-refacto/internal/models"
	"maestro/v2-refacto/internal/store"
	"maestro/v2-refacto/internal/validator"
)

var Tmpl *template.Template

func InitTemplates() {
	Tmpl = template.New("").Funcs(template.FuncMap{
		"add":   func(a, b int) int { return a + b },
		"lower": strings.ToLower,
	})
	Tmpl = template.Must(Tmpl.ParseGlob("templates/**/*.html"))
}

// Vue : Dashboard complet
func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := map[string]any{"Exercises": store.GetAll()}
	Tmpl.ExecuteTemplate(w, "base", data)
}

// Vue : Liste seule (Fragment)
// Vue : Liste seule (Fragment)
func HandleListExercice(w http.ResponseWriter, r *http.Request) {
	// 1. Lire TOUS les paramètres de filtre
	status := r.URL.Query().Get("status") // "todo", "done" ou vide
	domain := r.URL.Query().Get("domain") // "Algorithmes", "Go", etc. ou vide

	// 2. Construire le filtre composite
	filter := models.ExerciseFilter{
		Status: status,
		Domain: domain,
		// Difficulty: 0 (pas encore utilisé, mais prêt pour Phase 2)
	}

	// 3. Appeler le store avec le filtre complet
	filteredList := store.GetFiltered(filter)

	// 4. Renvoyer le fragment
	Tmpl.ExecuteTemplate(w, "exercise-list", filteredList)
}

// Vue : Détail (Fragment)
func HandleDetailExercice(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	if ex := store.FindExercise(id); ex != nil {
		Tmpl.ExecuteTemplate(w, "exercise-detail", ex)
	} else {
		http.NotFound(w, r)
	}
}

// Action : Toggle Done
func HandleToggleDone(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if ex := store.FindExercise(id); ex != nil {
		ex.Done = !ex.Done // Simple bascule true/false
		store.Save()       // <= Save appelé à chaque modification
		Tmpl.ExecuteTemplate(w, "exercise-card", ex)
	}
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
