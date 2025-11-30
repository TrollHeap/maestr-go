package handlers

import (
	"context"
	"net/http"

	"go-retro-terminal/internal/models"
	"go-retro-terminal/views"
)

// Page: Dashboard principal
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.DashboardPage().Render(context.Background(), w)
}

// Page: Stats détaillées
func HandleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	stats := &models.StatsData{
		CPU:     45.2,
		Memory:  62.1,
		Uptime:  "127h 42m",
		Status:  "OPERATIONAL",
		Version: "1.0.0",
	}

	views.StatsPage(stats).Render(context.Background(), w)
}

// Page: Terminal interactif
func HandleTerminal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.TerminalPage().Render(context.Background(), w)
}

// API: Menu items (pour HTMX)
func HandleMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.MenuItems().Render(context.Background(), w)
}

// API: Stats JSON (pour fetch JS)
func HandleStatsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats := `{
		"cpu": 45.2,
		"memory": 62.1,
		"uptime": "127h 42m",
		"status": "OPERATIONAL"
	}`

	w.Write([]byte(stats))
}

// API: Execute terminal command (HTMX POST)
func HandleTerminalExecute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	command := r.FormValue("command")

	// Simule l'exécution d'une commande
	var output string
	switch command {
	case "help":
		output = `Available commands:
- help       Show this help
- status     System status
- clear      Clear screen
- date       Show current date
- about      About this app`
	case "status":
		output = `System Status: OPERATIONAL
CPU: 45.2%
Memory: 62.1%
Uptime: 127h 42m`
	case "date":
		output = "Sunday, November 30, 2025, 4:53 AM CET"
	case "about":
		output = `Go Retro Terminal v1.0.0
Built with: Go, templ, HTMX, CSS
Design: Retro CRT aesthetic`
	case "clear":
		w.Write(
			[]byte(`<script>document.getElementById('terminal-output').innerHTML = '';</script>`),
		)
		return
	default:
		output = "Command not found: " + command + "\nType 'help' for available commands"
	}

	// Rendu du composant ligne de terminal
	views.TerminalLine(command, output).Render(context.Background(), w)
}
