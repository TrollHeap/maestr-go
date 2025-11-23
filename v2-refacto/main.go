package main

import (
	"log"
	"net/http"

	"maestro/v2-refacto/internal/handlers"
)

func main() {
	// A. Init Templates
	handlers.InitTemplates()

	// B. Routeur V2 (Récupère le routeur)
	mux := handlers.Routes() // <-- Capture le retour

	// C. Lancement (Passe le routeur)
	log.Println("Serveur sur http://localhost:8080")

	// ❌ http.ListenAndServe(":8080") // ERREUR : Il manque le mux !
	// ✅ CORRECTION :
	log.Fatal(http.ListenAndServe(":8080", mux))
}
