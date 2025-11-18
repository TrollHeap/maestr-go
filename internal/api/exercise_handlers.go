package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// UpdateExerciseSteps met à jour les steps complétés d'un exercice
func (h *ExerciseHandler) UpdateExerciseSteps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	exerciseID := vars["id"]

	var req struct {
		CompletedSteps []int `json:"completed_steps"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Récupérer l'exercice
	exercise, err := h.store.GetExercise(exerciseID)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	// Mettre à jour les steps complétés
	exercise.CompletedSteps = req.CompletedSteps

	// Sauvegarder
	if err := h.store.UpdateExercise(exercise); err != nil {
		http.Error(w, "Failed to update exercise", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"exercise": exercise,
	})
}

// ToggleCompletion marque un exercice comme complété ou non
func (h *ExerciseHandler) ToggleCompletion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	exerciseID := vars["id"]

	var req struct {
		Completed bool `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exercise, err := h.store.GetExercise(exerciseID)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	exercise.Completed = req.Completed

	if err := h.store.UpdateExercise(exercise); err != nil {
		http.Error(w, "Failed to update exercise", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"exercise": exercise,
	})
}

// RegisterRoutes enregistre toutes les routes
func (h *ExerciseHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/exercises", h.GetExercises).Methods("GET")
	router.HandleFunc("/api/rate", h.RateExercise).Methods("POST")
	router.HandleFunc("/api/exercises/{id}/steps", h.UpdateExerciseSteps).Methods("PUT")
	router.HandleFunc("/api/stats", h.GetStats).Methods("GET")

	// AJOUTER CETTE LIGNE
	router.HandleFunc("/api/exercises/{id}/completion", h.ToggleCompletion).Methods("PUT")
}
