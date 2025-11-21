package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv" // ← Ajouté
	"strings"
	"time"
)

type Exercise struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Domain         string    `json:"domain"`
	Difficulty     int       `json:"difficulty"`
	Steps          []string  `json:"steps"`
	CompletedSteps []int     `json:"completed_steps"` // ✅ NOUVEAU : indices des étapes validées (ex: [0, 2] = étapes 0 et 2 faites)
	Content        string    `json:"content"`
	Done           bool      `json:"done"`
	CreatedAt      time.Time `json:"created_at"`
}

// AllStepsCompleted vérifie si toutes les étapes sont validées
func (e *Exercise) AllStepsCompleted() bool {
	if len(e.Steps) == 0 {
		return false // Pas d'étapes = pas validable
	}
	return len(e.CompletedSteps) == len(e.Steps)
}

func main() {
	tmpl := template.New("").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"lower": strings.ToLower, // ← AJOUTE CETTE LIGNE
	})
	tmpl = template.Must(tmpl.ParseGlob("internal/templates/**/*.html"))

	exercises := []Exercise{
		{
			ID:          1,
			Title:       "Tri rapide (Quicksort)",
			Description: "Implémenter l'algorithme de tri rapide en Go",
			Domain:      "Algorithmes",
			Difficulty:  4,
			Steps: []string{
				"Comprendre le principe du pivot",
				"Implémenter la partition",
				"Récursion gauche et droite",
				"Tester avec des cas limites",
			},
			CompletedSteps: []int{}, // Aucune étape faite pour l'instant
			Content:        "# Quicksort...",
			Done:           false,
			CreatedAt:      time.Now(),
		},
		{
			ID:          2,
			Title:       "Filtrage de slice",
			Description: "Créer une fonction générique de filtrage",
			Domain:      "Go",
			Difficulty:  2,
			Steps: []string{
				"Définir la signature avec generics",
				"Implémenter la boucle de filtrage",
				"Écrire les tests unitaires",
			},
			CompletedSteps: []int{}, // Aucune étape faite pour l'instant
			Content:        "# Filtrage avec Go Generics\n\nDepuis Go 1.18...",
			Done:           false,
			CreatedAt:      time.Now(),
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"Exercises": exercises,
		}
		tmpl.ExecuteTemplate(w, "base", data)
	})

	http.HandleFunc("/exercises", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "exercise-list", exercises)
	})

	http.HandleFunc("/toggle-done", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest) // 400 Bad Request
			return                                                    // On arrête l'exécution ici
		}

		// 2. Parcours & Modification (Syntaxe Struct)
		// Attention : on utilise l'index 'i' pour modifier l'original dans le slice
		for i := range exercises {
			if exercises[i].ID == id {
				exercises[i].Done = !exercises[i].Done // Simple bascule true/false

				// Rendu du fragment HTML mis à jour
				tmpl.ExecuteTemplate(w, "exercise-card", exercises[i])
				return
			}
		}
		// Si on arrive ici, c'est que l'ID n'existe pas
		http.Error(w, "Exercise not found", http.StatusNotFound)
	})

	http.HandleFunc("/exercise/", func(w http.ResponseWriter, r *http.Request) {
		// Pattern : /exercise/1/toggle-step?step=0
		path := r.URL.Path[len("/exercise/"):]

		// Séparer ID et action
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			// Cas simple : /exercise/1 (affichage detail)
			id, err := strconv.Atoi(path)
			if err != nil {
				http.Error(w, "Invalid exercise ID", http.StatusBadRequest)
				return
			}

			for _, ex := range exercises {
				if ex.ID == id {
					tmpl.ExecuteTemplate(w, "exercise-detail", ex)
					return
				}
			}
			http.Error(w, "Exercise not found", http.StatusNotFound)
			return
		}

		// Cas toggle-step : /exercise/1/toggle-step?step=0
		if parts[1] == "toggle-step" {
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				http.Error(w, "Invalid exercise ID", http.StatusBadRequest)
				return
			}

			stepStr := r.URL.Query().Get("step")
			step, err := strconv.Atoi(stepStr)
			if err != nil {
				http.Error(w, "Invalid step index", http.StatusBadRequest)
				return
			}

			// Trouver l'exercice
			var exerciseIndex int = -1
			for i := range exercises {
				if exercises[i].ID == id {
					exerciseIndex = i
					break
				}
			}

			if exerciseIndex == -1 {
				http.Error(w, "Exercise not found", http.StatusNotFound)
				return
			}

			// Toggle de l'étape
			ex := &exercises[exerciseIndex]

			// Vérifier si l'étape est déjà complétée
			stepCompleted := false
			completedIndex := -1
			for i, s := range ex.CompletedSteps {
				if s == step {
					stepCompleted = true
					completedIndex = i
					break
				}
			}

			if stepCompleted {
				// Retirer l'étape (décoche)
				ex.CompletedSteps = append(
					ex.CompletedSteps[:completedIndex],
					ex.CompletedSteps[completedIndex+1:]...)
				ex.Done = false // Si on décoche une étape, ce n'est plus "Done"
			} else {
				// Ajouter l'étape (coche)
				ex.CompletedSteps = append(ex.CompletedSteps, step)
			}

			// Renvoyer DEUX fragments :
			// 1. Le détail mis à jour (target principal)
			tmpl.ExecuteTemplate(w, "exercise-detail", *ex)

			// 2. La carte mise à jour (out-of-band)
			w.Write([]byte("\n"))
			tmpl.ExecuteTemplate(w, "exercise-card-oob", *ex)
			return
		}

		http.Error(w, "Not found", http.StatusNotFound)
	})
	// ← Déplacé AVANT ListenAndServe
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Println("Serveur sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
