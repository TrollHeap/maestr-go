package handlers

import (
	"net/http"
)

// Routes définit et RETOURNE le routeur configuré
func Routes() *http.ServeMux {
	// B. Routeur V2
	mux := http.NewServeMux()

	// --- GROUPE 1 : VUES (GET) ---
	mux.HandleFunc("GET /", HandleDashboard)
	mux.HandleFunc("GET /exercises", HandleListExercice)
	mux.HandleFunc("GET /exercise/{id}", HandleDetailExercice)

	// --- GROUPE 2 : ACTIONS (POST) ---
	mux.HandleFunc("POST /toggle-done", HandleToggleDone)
	mux.HandleFunc("POST /exercise/{id}/toggle-step", HandleToggleStep)

	// --- GROUPE 3 : ASSETS ---
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	return mux
}
