package main

import (
	"log"
	"net/http"

	"maestro/v2-refacto/internal/handlers"
)

// --- 2. MAIN (Le Chef d'Orchestre) ---
func main() {
	// A. Init Templates
	handlers.InitTemplates()

	// B. Routeur V2
	mux := http.NewServeMux()

	// --- GROUPE 1 : VUES (GET) ---
	mux.HandleFunc("GET /", handlers.HandleDashboard)
	mux.HandleFunc("GET /exercises", handlers.HandleListExercice)
	mux.HandleFunc("GET /exercise/{id}", handlers.HandleDetailExercice)

	// --- GROUPE 2 : ACTIONS (POST) ---
	mux.HandleFunc("POST /toggle-done", handlers.HandleToggleDone)
	mux.HandleFunc("POST /exercise/{id}/toggle-step", handlers.HandleToggleStep)

	// --- GROUPE 3 : ASSETS ---
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// C. Lancement
	log.Println("Serveur sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// --- 3. HANDLERS (Les Spécialistes) ---
