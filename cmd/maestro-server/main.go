package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"maestro/internal/api"
	"maestro/internal/domain"
	"maestro/internal/storage"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	dataDir := flag.String("data-dir", "", "Data directory")
	flag.Parse()

	// Setup data directory
	if *dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		*dataDir = filepath.Join(home, ".maestro")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(*dataDir, 0o755); err != nil {
		log.Fatal(err)
	}

	// Initialize store
	storeFilepath := filepath.Join(*dataDir, "exercises.json")
	store := storage.NewJSONStore(storeFilepath)

	// Initialize domain logic
	scheduler := domain.NewScheduler()
	recommender := domain.NewRecommender(scheduler)

	// ✅ NOUVEAU: Initialiser StreakManager
	streak := domain.NewStreakManager()

	// ✅ MODIFIÉ: Ajouter streak au handler
	exerciseHandler := api.NewExerciseHandler(store, scheduler, recommender, streak)

	// Setup routes
	http.HandleFunc("/api/health", exerciseHandler.HealthCheck)
	http.HandleFunc("/api/exercises", exerciseHandler.GetExercises)
	http.HandleFunc("/api/recommended", exerciseHandler.GetRecommended)
	http.HandleFunc("/api/rate", exerciseHandler.RateExercise)
	http.HandleFunc("/api/stats", exerciseHandler.GetStats)

	// Serve frontend (static files)
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	fmt.Printf("🎯 Maestro Backend listening on http://localhost:%s\n", *port)
	fmt.Printf("📁 Data directory: %s\n", *dataDir)
	fmt.Printf("📄 Exercises file: %s\n", storeFilepath)
	fmt.Printf("\n✨ Endpoints:\n")
	fmt.Printf("  GET  http://localhost:%s/api/health\n", *port)
	fmt.Printf("  GET  http://localhost:%s/api/exercises\n", *port)
	fmt.Printf("  GET  http://localhost:%s/api/recommended\n", *port)
	fmt.Printf("  POST http://localhost:%s/api/rate\n", *port)
	fmt.Printf("  GET  http://localhost:%s/api/stats\n", *port)
	fmt.Printf("\n🌐 Web UI: http://localhost:%s\n", *port)

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal(err)
	}
}
