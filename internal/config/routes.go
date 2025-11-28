package config

import (
	"net/http"

	"maestro/internal/handlers"
)

// Routes définit et RETOURNE le routeur configuré
func Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// --- GROUPE 1 : PAGES COMPLÈTES (GET) ---
	mux.HandleFunc("/", handlers.HandleDashboard)

	// --- GROUPE 2 : FRAGMENTS HTMX (GET) ---
	mux.HandleFunc("GET /exercises", handlers.HandleExercisesPage)      // ← PAGE complète
	mux.HandleFunc("GET /exercises/list", handlers.HandleListExercice)  // ← FRAGMENT filtres
	mux.HandleFunc("GET /exercise/{id}", handlers.HandleDetailExercice) // ← FRAGMENT détail
	mux.HandleFunc("GET /exercise/next", handlers.HandleNextExercise)
	mux.HandleFunc("POST /exercise/{id}/toggle-step", handlers.HandleToggleStep)
	mux.HandleFunc("POST /exercise/{id}/review", handlers.HandleReview)

	// internal/handlers/routes.go (AJOUTER)

	// Routes sessions - Version corrigée
	mux.HandleFunc("GET /session/builder", handlers.HandleSessionBuilder)
	mux.HandleFunc("GET /session/start", handlers.HandleStartSession) // ← Corrigé
	mux.HandleFunc("GET /session/{id}", handlers.HandleCurrentSession)
	mux.HandleFunc("GET /session/complete", handlers.HandleSessionComplete)
	mux.HandleFunc("POST /session/{id}/stop", handlers.HandleStopSession)

	mux.HandleFunc("GET /stats", handlers.HandleStatsPage)
	mux.HandleFunc("GET /stats/metrics", handlers.HandleStatsMetrics)
	mux.HandleFunc("GET /stats/domains", handlers.HandleStatsDomains)
	mux.HandleFunc("GET /stats/difficulties", handlers.HandleStatsDifficulties)

	mux.HandleFunc("GET /planner", handlers.HandlePlannerPage)
	mux.HandleFunc("GET /planner/day", handlers.HandlePlannerDay)     // ← AJOUTE
	mux.HandleFunc("GET /planner/week", handlers.HandlePlannerWeek)   // ← AJOUTE
	mux.HandleFunc("GET /planner/month", handlers.HandlePlannerMonth) // ← AJOUTE

	// --- GROUPE 3 : ACTIONS (POST) ---
	mux.HandleFunc("POST /toggle-done", handlers.HandleToggleDone)

	// --- GROUPE 4 : ASSETS ---
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	return mux
}
