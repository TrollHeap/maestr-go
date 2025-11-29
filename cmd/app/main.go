package main

import (
	"log"
	"net/http"

	"maestro/internal/config"
	"maestro/internal/handlers"
	"maestro/internal/store"
)

func main() {
	log.Println("🚀 Démarrage Maestro Go v2...")

	// ✅ NOUVEAU : Init SQLite
	if err := store.InitDB("data/maestro.db"); err != nil {
		log.Fatal("Erreur init DB:", err)
	}
	defer store.CloseDB()

	// Init templates
	if err := handlers.InitTemplates(); err != nil {
		log.Fatal("Erreur templates:", err)
	}

	// Routes
	mux := config.Routes()

	log.Println("✅ Serveur démarré sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
