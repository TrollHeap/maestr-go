package main

import (
	"log"
	"net/http"

	"maestro/v2-refacto/internal/handlers"
	"maestro/v2-refacto/internal/store"
)

func main() {
	// A. Init Templates
	handlers.InitTemplates()

	// 2. Charge les données
	if err := store.Load(); err != nil {
		log.Fatalf("Erreur chargement données: %v", err)
	}

	// 3. Initialise avec des données par défaut si vide
	if err := store.InitDefaultExercises(); err != nil {
		log.Fatalf("Erreur initialisation: %v", err)
	}

	// B. Routeur V2 (Récupère le routeur)
	mux := handlers.Routes() // <-- Capture le retour

	// C. Lancement (Passe le routeur)
	log.Println("Serveur sur http://localhost:8080")

	// ❌ http.ListenAndServe(":8080") // ERREUR : Il manque le mux !
	// ✅ CORRECTION :
	log.Fatal(http.ListenAndServe(":8080", mux))
}
