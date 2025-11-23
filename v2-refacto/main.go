package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Variable globale pour la démo (En vrai : Database struct)
var exercises = []Exercise{
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

// Variable globale pour les templates (Initialisée dans le main)
var tmpl *template.Template

// --- 2. MAIN (Le Chef d'Orchestre) ---
func main() {
	// A. Init Templates
	initTemplates()

	// B. Routeur V2
	mux := http.NewServeMux()

	// --- GROUPE 1 : VUES (GET) ---
	mux.HandleFunc("GET /", HandleDashboard)
	mux.HandleFunc("GET /exercises", HandleListExercice)
	mux.HandleFunc("GET /exercise/{id}", HandleDetailExercice)

	// --- GROUPE 2 : ACTIONS (POST) ---
	mux.HandleFunc("POST /toggle-done", HandleToggleDone)
	mux.HandleFunc("POST /exercise/{id}/toggle-step", HandleToggleStep)

	// --- GROUPE 3 : ASSETS ---
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// C. Lancement
	log.Println("Serveur sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// --- 3. HANDLERS (Les Spécialistes) ---

func initTemplates() {
	tmpl = template.New("").Funcs(template.FuncMap{
		"add":   func(a, b int) int { return a + b },
		"lower": strings.ToLower,
	})
	tmpl = template.Must(tmpl.ParseGlob("templates/**/*.html"))
}

// Vue : Dashboard complet
func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := map[string]any{"Exercises": exercises}
	tmpl.ExecuteTemplate(w, "base", data)
}

// Vue : Liste seule (Fragment)
func HandleListExercice(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "exercise-list", exercises)
}

// Vue : Détail (Fragment)
func HandleDetailExercice(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	if ex := findExercise(id); ex != nil {
		tmpl.ExecuteTemplate(w, "exercise-detail", ex)
	} else {
		http.NotFound(w, r)
	}
}

// Action : Toggle Done
func HandleToggleDone(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if ex := findExercise(id); ex != nil {
		ex.Done = !ex.Done
		tmpl.ExecuteTemplate(w, "exercise-card", ex)
	}
}

// Action : Toggle Step
func HandleToggleStep(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	step, _ := strconv.Atoi(r.URL.Query().Get("step"))

	ex := findExercise(id)
	if ex == nil {
		http.NotFound(w, r)
		return
	}

	// Logique Métier (devrait être dans un Service, mais ok pour Handler ici)
	found := false
	for i, s := range ex.CompletedSteps {
		if s == step {
			ex.CompletedSteps = append(ex.CompletedSteps[:i], ex.CompletedSteps[i+1:]...)
			ex.Done = false
			found = true
			break
		}
	}
	if !found {
		ex.CompletedSteps = append(ex.CompletedSteps, step)
	}

	// Réponse Double (OOB)
	tmpl.ExecuteTemplate(w, "exercise-detail", *ex)
	w.Write([]byte("\n"))
	tmpl.ExecuteTemplate(w, "exercise-card-oob", *ex)
}

// --- 4. HELPERS (Outils internes) ---

// findExercise retourne un pointeur pour pouvoir modifier l'original
func findExercise(id int) *Exercise {
	for i := range exercises {
		if exercises[i].ID == id {
			return &exercises[i]
		}
	}
	return nil
}
