package main

import (
	"log"
	"net/http"

	"maestro/internal/config"
	"maestro/internal/handlers"
	"maestro/internal/store"
)

func main() {
	// 1. INIT TEMPLATES (AVEC GESTION D'ERREUR !)
	if err := handlers.InitTemplates(); err != nil {
		log.Fatalf("❌ Erreur chargement templates: %v", err)
	}
	log.Println("✅ Templates chargés")

	// 2. CHARGE LES DONNÉES
	if err := store.Load(); err != nil {
		log.Fatalf("❌ Erreur chargement données: %v", err)
	}
	log.Println("✅ Données chargées")

	// 3. INITIALISE AVEC DONNÉES PAR DÉFAUT SI VIDE
	if err := store.InitDefaultExercises(); err != nil {
		log.Fatalf("❌ Erreur initialisation: %v", err)
	}
	log.Println("✅ Exercices initialisés")

	// 4. ROUTEUR
	mux := config.Routes()
	log.Println("✅ Routes configurées")

	// 5. LANCEMENT SERVEUR
	log.Println("🚀 Serveur sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
