package handlers

import (
	"log"
	"net/http"

	"maestro/v2-refacto/internal/store"
)

// internal/handlers/exercises.go
func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 HandleDashboard appelé - Path: %s", r.URL.Path)

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := map[string]any{
		"Exercises": store.GetAll(),
	}

	log.Println("✅ Exécution template 'dashboard'")
	if err := Tmpl.ExecuteTemplate(w, "dashboard", data); err != nil {
		log.Printf("❌ Erreur template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
