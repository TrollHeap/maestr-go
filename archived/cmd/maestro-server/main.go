package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"

	"maestro/internal/api"
	"maestro/internal/domain"
	"maestro/internal/storage"
)

func main() {
	// Configuration
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

	// Create data directory
	if err := os.MkdirAll(*dataDir, 0o755); err != nil {
		log.Fatal(err)
	}

	// CORRECTION: NewJSONStore retourne (store, error)
	store, err := storage.NewJSONStore(*dataDir)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	// Initialize domain services
	scheduler := domain.NewScheduler()
	recommender := domain.NewRecommender(scheduler)
	planner := domain.NewPlanner() // ✅ AJOUTÉ: Création du planner

	// ✅ CORRECTION: Passer planner au lieu de streak
	handler := api.NewExerciseHandler(store, scheduler, recommender, planner)

	// Setup router
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Serve static files
	staticDir := filepath.Join(".", "public")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(staticDir)))

	// Start server
	addr := fmt.Sprintf(":%s", *port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
