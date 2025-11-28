package main

import (
	"log"
	"net/http"

	"maestro/internal/config"
	"maestro/internal/handlers"
	"maestro/internal/store"
)

func main() {
	// 1. INIT TEMPLATES
	if err := handlers.InitTemplates(); err != nil {
		log.Fatalf("❌ Erreur chargement templates: %v", err)
	}
	log.Println("✅ Templates chargés")

	// 2. CHARGE LES DONNÉES EXERCICES
	if err := store.Load(); err != nil {
		log.Fatalf("❌ Erreur chargement données: %v", err)
	}
	log.Println("✅ Données exercices chargées")

	// 3. ✅ CHARGE LES SESSIONS SAUVEGARDÉES
	if err := store.LoadSessions(); err != nil {
		// Non fatal : fichier peut ne pas exister au premier lancement
		log.Printf("⚠️  Sessions non chargées (normal au 1er lancement): %v", err)
	} else {
		log.Println("✅ Sessions chargées")
	}

	// 4. INITIALISE AVEC DONNÉES PAR DÉFAUT SI VIDE
	if err := store.InitDefaultExercises(); err != nil {
		log.Fatalf("❌ Erreur initialisation: %v", err)
	}
	log.Println("✅ Exercices initialisés")

	// 5. ROUTEUR
	mux := config.Routes()
	log.Println("✅ Routes configurées")

	// 6. LANCEMENT SERVEUR
	log.Println("🚀 Serveur sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
