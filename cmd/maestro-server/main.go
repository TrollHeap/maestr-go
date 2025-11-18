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
	streak := domain.NewStreakManager()
	planner := domain.NewPlanner()

	// Initialize handlers
	exerciseHandler := api.NewExerciseHandler(store, scheduler, recommender, streak)
	plannerHandler := api.NewPlannerHandler(planner, store)

	// Create router mux
	router := mux.NewRouter()

	// Register all API routes
	exerciseHandler.RegisterRoutes(router)
	plannerHandler.RegisterRoutes(router)

	// Serve frontend static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	// Server info
	fmt.Printf("🚀 Maestro Backend listening on http://localhost:%s\n", *port)
	fmt.Printf("📂 Data directory: %s\n", *dataDir)
	fmt.Printf("📄 Exercises file: %s\n\n", storeFilepath)

	// Start server
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
