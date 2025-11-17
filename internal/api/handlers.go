package api

import (
	"encoding/json"
	"maestro/internal/domain"
	"maestro/internal/models"
	"maestro/internal/storage"
	"net/http"
)

// ExerciseHandler gère les endpoints pour les exercices
type ExerciseHandler struct {
	store       storage.Store
	scheduler   *domain.Scheduler
	recommender *domain.Recommender
}

// NewExerciseHandler crée une nouvelle instance ExerciseHandler
func NewExerciseHandler(store storage.Store, scheduler *domain.Scheduler, recommender *domain.Recommender) *ExerciseHandler {
	return &ExerciseHandler{
		store:       store,
		scheduler:   scheduler,
		recommender: recommender,
	}
}

// GetExercises retourne tous les exercices
func (h *ExerciseHandler) GetExercises(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "Failed to load exercises", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(exercises)
}

// GetRecommended retourne les exercices recommandés
func (h *ExerciseHandler) GetRecommended(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "Failed to load exercises", http.StatusInternalServerError)
		return
	}

	// Obtenir les 3 recommandés
	recommended := h.recommender.GetNextExercises(exercises, 3)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(recommended)
}

// RateExercise met à jour la note d'un exercice
func (h *ExerciseHandler) RateExercise(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var input models.ReviewInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Valider le rating
	if input.Rating < 1 || input.Rating > 4 {
		http.Error(w, "Rating must be 1-4", http.StatusBadRequest)
		return
	}

	// Charger l'exercice
	exercise, err := h.store.GetByID(ctx, input.ExerciseID)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	// Appliquer l'algorithme SM-2
	h.scheduler.ReviewExercise(exercise, input.Rating)

	// Persister
	if err := h.store.Update(ctx, exercise); err != nil {
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}

	// Message d'encouragement
	messages := map[int]string{
		1: "😅 Pas grave, tu vas y arriver!",
		2: "💪 Bonne tentative, continue!",
		3: "👏 Nickel, bien joué!",
		4: "🔥 Excellent! Parfaitement maîtrisé!",
	}

	// Répondre
	response := models.ReviewResponse{
		Exercise:     exercise,
		NextReviewIn: exercise.IntervalDays,
		Message:      messages[input.Rating],
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

// GetStats retourne les statistiques
func (h *ExerciseHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "Failed to load exercises", http.StatusInternalServerError)
		return
	}

	stats := h.recommender.CalculateStats(exercises)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(stats)
}

// HealthCheck endpoint pour vérifier que l'API est alive
func (h *ExerciseHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Maestro Backend is running",
	})
}
