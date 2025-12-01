package main

import (
	"log"
	"net/http"
	"os"

	"maestro/internal/config"
	"maestro/internal/handlers"
	"maestro/internal/store"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🚀 Maestro Go v2.0 - Low-Power Learning")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// DB
	dbPath := getEnv("DB_PATH", "data/maestro.db")
	log.Printf("📦 Initialisation DB: %s", dbPath)
	if err := store.InitDB(dbPath); err != nil {
		log.Fatalf("❌ Erreur init DB: %v", err)
	}
	defer func() {
		log.Println("🔒 Fermeture DB...")
		store.CloseDB()
	}()
	log.Println("✅ DB initialisée")

	// ✅ Templates (auto-validation)
	log.Println("📄 Chargement templates...")
	if err := handlers.InitTemplates(); err != nil {
		log.Fatalf("❌ Erreur init templates: %v", err)
	}

	// ✅ Liste tous les templates chargés (auto-découverte)
	loadedTemplates := handlers.ListTemplates()
	if len(loadedTemplates) == 0 {
		log.Fatal("❌ Aucun template chargé")
	}

	log.Printf("✅ %d templates chargés:", len(loadedTemplates))
	for _, tmpl := range loadedTemplates {
		log.Printf("   • %s", tmpl)
	}

	// ✅ Validation minimale : au moins templates de base
	minimumRequired := []string{
		"dashboard.html",          // ✅ Nom fichier (pas "dashboard")
		"exercise-list-page.html", // ✅ Nom fichier
		"base",                    // ✅ Layout ({{ define "base" }})
	}
	for _, tmpl := range minimumRequired {
		if !handlers.HasTemplate(tmpl) {
			log.Fatalf("❌ Template critique manquant: %s", tmpl)
		}
	}
	log.Println("✅ Templates critiques validés")

	// Dossiers
	log.Println("📁 Validation dossiers...")
	requiredDirs := []string{"data", "templates"}
	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Fatalf("❌ Dossier manquant: %s/", dir)
		}
	}
	log.Printf("✅ %d dossiers validés", len(requiredDirs))

	// Routes + serveur
	port := getEnv("PORT", "8080")
	mux := config.Routes()

	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("✅ Serveur prêt")
	log.Printf("✅ http://localhost:%s", port)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
