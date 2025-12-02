package main

import (
	"log"
	"net/http"
	"os"

	"maestro/internal/config"
	"maestro/internal/store"
)

func main() {
	// === BANNER ===
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🚀 Maestro Go v2.0 - Low-Power Learning")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// === DATABASE INIT ===
	dbPath := getEnv("DB_PATH", "data/maestro.db")
	log.Printf("📦 Connexion DB: %s", dbPath)

	if err := store.InitDB(dbPath); err != nil {
		log.Fatalf("❌ Erreur init DB: %v", err)
	}
	defer func() {
		log.Println("🔒 Fermeture DB...")
		store.CloseDB()
	}()

	log.Println("✅ DB initialisée")

	// === ROUTES ===
	log.Println("🔧 Configuration routes...")
	mux := config.Routes()

	// === SERVER START ===
	port := getEnv("PORT", "8080")
	addr := ":" + port

	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("✅ Serveur démarré sur http://localhost:%s", port)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("❌ Erreur serveur: %v", err)
	}
}

// getEnv récupère variable d'environnement avec fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
