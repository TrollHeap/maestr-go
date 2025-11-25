package handlers

import (
	"net/http"
)

// Routes définit et RETOURNE le routeur configuré
func Routes() *http.ServeMux {
	// B. Routeur V2
	mux := http.NewServeMux()

	// --- GROUPE 1 : VUES (GET) ---
	mux.HandleFunc("/", HandleDashboard) // ✅ "/" = dashboard
	// Stats routes
	mux.HandleFunc("GET /stats", HandleStatsPage)
	mux.HandleFunc("GET /stats/metrics", HandleStatsMetrics)
	mux.HandleFunc("GET /stats/domains", HandleStatsDomains)
	mux.HandleFunc("GET /stats/difficulties", HandleStatsDifficulties)
	mux.HandleFunc("GET /planner", HandlePlannerPage)
	mux.HandleFunc("GET /planner/day", HandlePlannerDay)
	mux.HandleFunc("GET /planner/week", HandlePlannerWeek)
	mux.HandleFunc("GET /planner/month", HandlePlannerMonth)

	mux.HandleFunc("GET /exercises", HandleListExercice)
	mux.HandleFunc("GET /exercise/{id}", HandleDetailExercice)

	// --- GROUPE 2 : ACTIONS (POST) ---
	mux.HandleFunc("POST /toggle-done", HandleToggleDone)
	mux.HandleFunc("POST /exercise/{id}/toggle-step", HandleToggleStep)
	mux.HandleFunc("POST /exercise/{id}/review", HandleReview)

	// --- GROUPE 3 : ASSETS ---
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	return mux
}
