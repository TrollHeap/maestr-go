package handlers

import (
	"log"
	"net/http"
)

// Initialiser les templates au démarrage (dans main.go ou init)
// Handler simple pour afficher une page
func HandlePico(w http.ResponseWriter, r *http.Request) {
	// Définir le type de contenu
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Exécuter le template
	err := Tmpl.ExecuteTemplate(w, "stats.html", nil)
	if err != nil {
		http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
		log.Printf("Erreur ExecuteTemplate: %v", err)
	}
}
