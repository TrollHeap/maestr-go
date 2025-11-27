// internal/handlers/sessions.go
package handlers

import (
	"net/http"
	"strconv"

	"maestro/internal/models"
	"maestro/internal/store"
)

// Page sélection énergie
func HandleSessionBuilder(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Configs": models.SessionConfigs,
	}

	if err := Tmpl.ExecuteTemplate(w, "session-builder", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Démarre une session
func HandleStartSession(w http.ResponseWriter, r *http.Request) {
	energyStr := r.URL.Query().Get("energy") // "1", "2", ou "3"
	energy, _ := strconv.Atoi(energyStr)

	if energy < 1 || energy > 3 {
		http.Error(w, "Invalid energy level", http.StatusBadRequest)
		return
	}

	session := store.BuildAdaptiveSession(models.EnergyLevel(energy))

	data := map[string]any{
		"Session": session,
	}

	if err := Tmpl.ExecuteTemplate(w, "session-current", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
