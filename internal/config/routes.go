package config

import (
	"net/http"

	"maestro/internal/handlers"
)

// Routes définit et RETOURNE le routeur configuré
func Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// ============================================
	// GROUPE 1 : PAGES COMPLÈTES (GET)
	// ============================================
	mux.HandleFunc("GET /", handlers.HandleDashboard)
	mux.HandleFunc("GET /exercises", handlers.HandleExercisesPage)
	mux.HandleFunc("GET /planner", handlers.HandlePlannerPage)

	// ============================================
	// GROUPE 2 : EXERCICES - FRAGMENTS HTMX
	// ============================================
	mux.HandleFunc("GET /exercises/list", handlers.HandleListExercice)  // Liste filtrée
	mux.HandleFunc("GET /exercise/{id}", handlers.HandleDetailExercice) // Détail
	mux.HandleFunc("GET /exercise/next", handlers.HandleNextExercise)   // Prochain à réviser

	// Actions exercices (POST)
	mux.HandleFunc("POST /exercise/{id}/toggle-step", handlers.HandleToggleStep)
	mux.HandleFunc("POST /exercise/{id}/review", handlers.HandleReview)
	mux.HandleFunc("POST /toggle-done", handlers.HandleToggleDone)

	// ============================================
	// GROUPE 3 : SESSIONS
	// ============================================
	mux.HandleFunc("GET /session/builder", handlers.HandleSessionBuilder)
	mux.HandleFunc("GET /session/start", handlers.HandleStartSession)       // Démarre session
	mux.HandleFunc("GET /session/complete", handlers.HandleSessionComplete) // Page fin
	mux.HandleFunc("POST /session/{id}/stop", handlers.HandleStopSession)   // Arrête session

	// ============================================
	// GROUPE 4 : PLANNER (Calendrier)
	// ============================================
	mux.HandleFunc("GET /planner/day", handlers.HandlePlannerDay)
	mux.HandleFunc("GET /planner/week", handlers.HandlePlannerWeek)
	mux.HandleFunc("GET /planner/month", handlers.HandlePlannerMonth)

	// ============================================
	// GROUPE 5 : ASSETS STATIQUES
	// ============================================
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	return mux
}
