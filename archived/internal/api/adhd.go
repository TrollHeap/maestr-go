package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// GetQuickWins retourne les exercices "quick-win" (5-15 min)
func (h *ExerciseHandler) GetQuickWins(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "Failed to load exercises", http.StatusInternalServerError)
		return
	}

	// Filtrer les quick-wins (non complétés, pas supprimés, durée courte)
	quickWins := []any{}
	for _, ex := range exercises {
		if !ex.Completed && !ex.Deleted && len(ex.Steps) <= 3 {
			quickWins = append(quickWins, map[string]any{
				"id":          ex.ID,
				"title":       ex.Title,
				"description": ex.Description,
				"steps":       ex.Steps,
				"difficulty":  ex.Difficulty,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quickWins)
}

// GetVisualExercises retourne les exercices visuels
func (h *ExerciseHandler) GetVisualExercises(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "Failed to load exercises", http.StatusInternalServerError)
		return
	}

	// Filtrer les exercices avec du contenu visuel
	visual := []any{}
	for _, ex := range exercises {
		if !ex.Deleted {
			visual = append(visual, map[string]any{
				"id":          ex.ID,
				"title":       ex.Title,
				"description": ex.Description,
				"content":     ex.Content,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(visual)
}

// StartFocusSession démarre une session de focus
func (h *ExerciseHandler) StartFocusSession(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ExerciseIDs []string `json:"exercise_ids"`
		Duration    int      `json:"duration"` // minutes
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Créer session de focus
	session := map[string]any{
		"exercise_ids": input.ExerciseIDs,
		"duration":     input.Duration,
		"started_at":   time.Now(),
		"status":       "active",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

// TakeBreak enregistre une pause
func (h *ExerciseHandler) TakeBreak(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Duration int `json:"duration"` // minutes
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	breakInfo := map[string]any{
		"duration":   input.Duration,
		"started_at": time.Now(),
		"type":       "break",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breakInfo)
}

// SkipExercise enregistre un skip (ADHD pattern)
func (h *ExerciseHandler) SkipExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	exercise, err := h.store.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	// ✅ Incrémenter skip count
	exercise.SkippedCount++
	now := time.Now()
	exercise.LastSkipped = &now // ✅ Utilise le champ qui existe maintenant
	exercise.UpdatedAt = now

	if err := h.store.Update(ctx, exercise); err != nil {
		http.Error(w, "Failed to update exercise", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercise)
}
