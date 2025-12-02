package main

import (
	"log"
	"net/http"
	"go-retro-terminal/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Routes - Pages (templates)
	router.HandleFunc("/", handlers.HandleIndex).Methods("GET")
	router.HandleFunc("/stats", handlers.HandleStats).Methods("GET")
	router.HandleFunc("/terminal", handlers.HandleTerminal).Methods("GET")

	// Routes - API (HTMX / JSON)
	router.HandleFunc("/api/stats", handlers.HandleStatsAPI).Methods("GET")
	router.HandleFunc("/api/terminal/execute", handlers.HandleTerminalExecute).Methods("POST")
	router.HandleFunc("/api/menu", handlers.HandleMenu).Methods("GET")

	// Assets statiques
	router.PathPrefix("/assets/").Handler(
		http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))),
	)

	log.Println("🟢 Retro Terminal started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
