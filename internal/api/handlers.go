package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"maestro/internal/domain"
	"maestro/internal/models"
	"maestro/internal/storage"
)

// ExerciseHandler gère les endpoints pour les exercices
type ExerciseHandler struct {
	store       storage.Store
	scheduler   *domain.Scheduler
	recommender *domain.Recommender
	streak      *domain.StreakManager
}

// NewExerciseHandler crée une nouvelle instance ExerciseHandler
func NewExerciseHandler(
	store storage.Store,
	scheduler *domain.Scheduler,
	recommender *domain.Recommender,
	streak *domain.StreakManager,
) *ExerciseHandler {
	return &ExerciseHandler{
		store:       store,
		scheduler:   scheduler,
		recommender: recommender,
		streak:      streak,
	}
}

// GetExercises retourne tous les exercices avec pagination
func (h *ExerciseHandler) GetExercises(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "Failed to load exercises", http.StatusInternalServerError)
		return
	}

	// Pagination
	page := 1
	pageSize := 10

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	// Calculer pagination
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(exercises) {
		start = len(exercises)
	}
	if end > len(exercises) {
		end = len(exercises)
	}

	paginatedExercises := exercises[start:end]

	// Ajouter les dates de révision
	reviewDates := domain.CalculateNextReviewDates(exercises)

	response := map[string]interface{}{
		"total":        len(exercises),
		"page":         page,
		"page_size":    pageSize,
		"exercises":    paginatedExercises,
		"review_dates": reviewDates,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
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

	// Ajouter les dates de révision
	reviewDates := domain.CalculateNextReviewDates(exercises)

	response := map[string]interface{}{
		"recommended":  recommended,
		"review_dates": reviewDates,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
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

	// Mettre à jour le streak
	h.streak.UpdateStreak(time.Now())

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

	// Ajouter les infos de streak
	streakStats := h.streak.GetStats()

	response := map[string]interface{}{
		"total_completed": stats.TotalCompleted,
		"total_reviews":   stats.TotalReviews,
		"domain_stats":    stats.DomainStats,
		"streak": map[string]interface{}{
			"current":      streakStats.CurrentStreak,
			"display":      streakStats.StreakDisplay,
			"last_session": streakStats.LastSessionAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
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
