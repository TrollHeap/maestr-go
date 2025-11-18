package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"maestro/internal/domain"
	"maestro/internal/models"
)

// ============= EXERCISE ACTIONS HANDLERS =============

// UncompleteExercise - Toggle uncomplete
func (h *ExerciseHandler) UncompleteExercise(w http.ResponseWriter, r *http.Request) {
	exerciseID := r.URL.Query().Get("id")
	if exerciseID == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	exercise, err := h.store.GetByID(ctx, exerciseID)
	if err != nil {
		http.Error(w, "exercise not found", http.StatusNotFound)
		return
	}

	// Reset to initial state
	domain.UncompleteExercise(exercise)

	if err := h.store.Update(ctx, exercise); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"exercise": exercise,
		"message":  "Exercise marked as incomplete",
	})
}

// DeleteExercise - Soft delete
func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	exerciseID := r.URL.Query().Get("id")
	if exerciseID == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	exercise, err := h.store.GetByID(ctx, exerciseID)
	if err != nil {
		http.Error(w, "exercise not found", http.StatusNotFound)
		return
	}

	// Soft delete
	now := time.Now()
	exercise.Deleted = true
	exercise.DeletedAt = &now

	if err := h.store.Update(ctx, exercise); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Exercise deleted",
	})
}

// SkipExercise - Mark for later review
func (h *ExerciseHandler) SkipExercise(w http.ResponseWriter, r *http.Request) {
	exerciseID := r.URL.Query().Get("id")
	if exerciseID == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	exercise, err := h.store.GetByID(ctx, exerciseID)
	if err != nil {
		http.Error(w, "exercise not found", http.StatusNotFound)
		return
	}

	// Mark as skipped (no penalty!)
	domain.SkipExercise(exercise)

	if err := h.store.Update(ctx, exercise); err != nil {
		http.Error(w, "skip failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"exercise": exercise,
		"message":  "Exercise skipped - review later",
	})
}

// ToggleStep - Toggle step completion
func (h *ExerciseHandler) ToggleStep(w http.ResponseWriter, r *http.Request) {
	exerciseID := r.URL.Query().Get("id")
	stepStr := r.URL.Query().Get("step")

	if exerciseID == "" || stepStr == "" {
		http.Error(w, "missing id or step parameter", http.StatusBadRequest)
		return
	}

	stepIndex, err := strconv.Atoi(stepStr)
	if err != nil {
		http.Error(w, "invalid step index", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	exercise, err := h.store.GetByID(ctx, exerciseID)
	if err != nil {
		http.Error(w, "exercise not found", http.StatusNotFound)
		return
	}

	if stepIndex < 0 || stepIndex >= len(exercise.Steps) {
		http.Error(w, "step index out of range", http.StatusBadRequest)
		return
	}

	// Toggle the step
	domain.ToggleStep(exercise, stepIndex)

	if err := h.store.Update(ctx, exercise); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	// Check if all steps done - emit reward
	reward := ""
	rewardEngine := domain.NewRewardEngine()
	if len(exercise.CompletedSteps) == len(exercise.Steps) {
		reward = rewardEngine.GetMessage(domain.EventExerciseCompleted)
	} else {
		reward = rewardEngine.GetMessage(domain.EventStepCompleted)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"exercise":        exercise,
		"completed_steps": exercise.CompletedSteps,
		"total_steps":     len(exercise.Steps),
		"reward":          reward,
	})
}

// ResetAll - Nuclear reset (with confirmation)
func (h *ExerciseHandler) ResetAll(w http.ResponseWriter, r *http.Request) {
	confirm := r.URL.Query().Get("confirm")
	if confirm != "YES" {
		http.Error(w, "confirm=YES required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "load failed", http.StatusInternalServerError)
		return
	}

	// Reset all
	for i := range exercises {
		domain.UncompleteExercise(&exercises[i])
		exercises[i].SkippedCount = 0
		exercises[i].Deleted = false
		exercises[i].DeletedAt = nil
	}

	if err := h.store.Save(ctx, exercises); err != nil {
		http.Error(w, "save failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "All exercises reset",
		"total":   len(exercises),
	})
}

// GetNextExercise - ADHD-friendly recommendation
func (h *ExerciseHandler) GetNextExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "load failed", http.StatusInternalServerError)
		return
	}

	// Filter non-deleted
	active := []models.Exercise{}
	for _, ex := range exercises {
		if !ex.Deleted {
			active = append(active, ex)
		}
	}

	recommender := domain.NewADHDRecommender()
	next := recommender.GetNextExercise(active)

	if next == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"exercise": nil,
			"message":  "🎉 All done! Great work!",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"exercise": next,
		"message":  "Next up:",
	})
}

// GetSessionStats - Stats for current session
func (h *ExerciseHandler) GetSessionStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "load failed", http.StatusInternalServerError)
		return
	}

	// Count what matters for ADHD
	toDoCount := 0
	overdue := 0

	for _, ex := range exercises {
		if ex.Deleted {
			continue
		}

		if !ex.Completed {
			toDoCount++
			// Check if overdue
			if ex.LastReviewed != nil {
				nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
				if time.Now().After(nextReview) {
					overdue++
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"to_do":     toDoCount,
		"overdue":   overdue,
		"total":     len(exercises),
		"completed": len(exercises) - toDoCount,
		"goal":      3, // ADHD goal: do 3-5 per session
	})
}

// RegisterRoutes registers the new ADHD routes
func (h *ExerciseHandler) RegisterADHDRoutes(mux *http.ServeMux) {
	// Existing
	// (Keep your existing routes)

	// NEW - ADHD features
	mux.HandleFunc("/api/exercise/uncomplete", h.UncompleteExercise)
	mux.HandleFunc("/api/exercise/delete", h.DeleteExercise)
	mux.HandleFunc("/api/exercise/skip", h.SkipExercise)
	mux.HandleFunc("/api/exercise/toggle-step", h.ToggleStep)
	mux.HandleFunc("/api/exercise/next", h.GetNextExercise)
	mux.HandleFunc("/api/session/stats", h.GetSessionStats)
	mux.HandleFunc("/api/reset-all", h.ResetAll)
}
