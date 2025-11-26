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
	mux.HandleFunc("GET /exercises", handlers.HandleExercisesPage) // ← PAGE complète
	mux.HandleFunc("GET /stats", handlers.HandleStatsPage)
	mux.HandleFunc("GET /planner", handlers.HandlePlannerPage)

	// --- GROUPE 2 : FRAGMENTS HTMX (GET) ---
	mux.HandleFunc("GET /exercises/list", handlers.HandleListExercice)  // ← FRAGMENT filtres
	mux.HandleFunc("GET /exercise/{id}", handlers.HandleDetailExercice) // ← FRAGMENT détail

	mux.HandleFunc("GET /stats/metrics", handlers.HandleStatsMetrics)
	mux.HandleFunc("GET /stats/domains", handlers.HandleStatsDomains)
	mux.HandleFunc("GET /stats/difficulties", handlers.HandleStatsDifficulties)

	mux.HandleFunc("GET /planner/day", handlers.HandlePlannerDay)
	mux.HandleFunc("GET /planner/week", handlers.HandlePlannerWeek)
	mux.HandleFunc("GET /planner/month", handlers.HandlePlannerMonth)

	// --- GROUPE 3 : ACTIONS (POST) ---
	mux.HandleFunc("POST /toggle-status", handlers.HandleToggleStatus)
	mux.HandleFunc("POST /toggle-done", handlers.HandleToggleDone)
	mux.HandleFunc("POST /exercise/{id}/toggle-step", handlers.HandleToggleStep)
	mux.HandleFunc("POST /exercise/{id}/review", handlers.HandleReview)

	// --- GROUPE 4 : ASSETS ---
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	return mux
}
