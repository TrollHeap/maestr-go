package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"maestro/internal/domain"
	"maestro/internal/models"
	"maestro/internal/storage"
)

// ExerciseHandler gère les requêtes HTTP pour les exercices
type ExerciseHandler struct {
	store       storage.Store
	scheduler   *domain.Scheduler
	recommender *domain.Recommender
	planner     *domain.Planner
	// ❌ SUPPRIMÉ: adhdManager *domain.ADHDManager
}

// NewExerciseHandler crée un nouveau handler
func NewExerciseHandler(
	store storage.Store,
	scheduler *domain.Scheduler,
	recommender *domain.Recommender,
	planner *domain.Planner,
	// ❌ SUPPRIMÉ: adhdManager *domain.ADHDManager,
) *ExerciseHandler {
	return &ExerciseHandler{
		store:       store,
		scheduler:   scheduler,
		recommender: recommender,
		planner:     planner,
		// ❌ SUPPRIMÉ: adhdManager: adhdManager,
	}
}

// ... reste du code identique

// RegisterRoutes enregistre toutes les routes
func (h *ExerciseHandler) RegisterRoutes(router *mux.Router) {
	// Exercises CRUD
	router.HandleFunc("/api/exercises", h.GetExercises).Methods("GET")
	router.HandleFunc("/api/exercises", h.CreateExercise).Methods("POST")
	router.HandleFunc("/api/exercises/{id}", h.GetExercise).Methods("GET")
	router.HandleFunc("/api/exercises/{id}", h.UpdateExercise).Methods("PUT")
	router.HandleFunc("/api/exercises/{id}", h.DeleteExercise).Methods("DELETE")

	// Exercise actions
	router.HandleFunc("/api/exercises/{id}/steps", h.UpdateExerciseSteps).Methods("PUT")
	router.HandleFunc("/api/exercises/{id}/completion", h.ToggleExerciseCompletion).Methods("PUT")
	router.HandleFunc("/api/exercises/{id}/review", h.ReviewExercise).Methods("POST")

	// Recommendations & Stats
	router.HandleFunc("/api/recommendations", h.GetRecommendations).Methods("GET")
	router.HandleFunc("/api/stats", h.GetStats).Methods("GET")

	// ADHD Features (depuis adhd.go)
	router.HandleFunc("/api/adhd/quick-wins", h.GetQuickWins).Methods("GET")
	router.HandleFunc("/api/adhd/visual", h.GetVisualExercises).Methods("GET")
	router.HandleFunc("/api/adhd/focus-session", h.StartFocusSession).Methods("POST")
	router.HandleFunc("/api/adhd/break", h.TakeBreak).Methods("POST")
	router.HandleFunc("/api/exercises/{id}/skip", h.SkipExercise).Methods("POST")

	// Planner (depuis planner_handlers.go)
	router.HandleFunc("/api/planner/today", h.GetTodayPlan).Methods("GET")
	router.HandleFunc("/api/planner/week", h.GetWeekPlan).Methods("GET")
	router.HandleFunc("/api/planner/stats", h.GetPlannerStats).Methods("GET")
	router.HandleFunc("/api/planner/sessions", h.CreatePlannerSession).Methods("POST")
	router.HandleFunc("/api/planner/sessions/{id}", h.UpdatePlannerSession).Methods("PUT")
	router.HandleFunc("/api/planner/sessions/{id}", h.DeletePlannerSession).Methods("DELETE")
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
	json.NewEncoder(w).Encode(exercises)
}

// GetExercise retourne un exercice spécifique
func (h *ExerciseHandler) GetExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	exercise, err := h.store.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercise)
}

// CreateExercise crée un nouvel exercice
func (h *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Domain      string   `json:"domain"`
		Difficulty  int      `json:"difficulty"`
		Steps       []string `json:"steps"`
		Content     string   `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	now := time.Now()
	exercise := models.Exercise{
		ID:             uuid.New().String(),
		Title:          input.Title,
		Description:    input.Description,
		Domain:         input.Domain,
		Difficulty:     input.Difficulty,
		Steps:          input.Steps,
		Content:        input.Content,
		Completed:      false,
		CompletedSteps: []int{},
		EaseFactor:     2.5,
		IntervalDays:   0,
		Repetitions:    0,
		SkippedCount:   0,
		Deleted:        false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "Failed to load exercises", http.StatusInternalServerError)
		return
	}

	exercises = append(exercises, exercise)

	if err := h.store.Save(ctx, exercises); err != nil {
		http.Error(w, "Failed to save exercise", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exercise)
}

// UpdateExercise met à jour un exercice
func (h *ExerciseHandler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Domain      string   `json:"domain"`
		Difficulty  int      `json:"difficulty"`
		Steps       []string `json:"steps"`
		Content     string   `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exercise, err := h.store.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	exercise.Title = input.Title
	exercise.Description = input.Description
	exercise.Domain = input.Domain
	exercise.Difficulty = input.Difficulty
	exercise.Steps = input.Steps
	exercise.Content = input.Content
	exercise.UpdatedAt = time.Now()

	if err := h.store.Update(ctx, exercise); err != nil {
		http.Error(w, "Failed to update exercise", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercise)
}

// DeleteExercise supprime (soft delete) un exercice
func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	exercise, err := h.store.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	// Soft delete
	now := time.Now()
	exercise.Deleted = true
	exercise.DeletedAt = &now
	exercise.UpdatedAt = now

	if err := h.store.Update(ctx, exercise); err != nil {
		http.Error(w, "Failed to delete exercise", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateExerciseSteps met à jour les steps complétés
func (h *ExerciseHandler) UpdateExerciseSteps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	var input struct {
		CompletedSteps []int `json:"completed_steps"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exercise, err := h.store.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	exercise.CompletedSteps = input.CompletedSteps
	exercise.UpdatedAt = time.Now()

	// Auto-compléter si tous les steps sont faits
	if len(input.CompletedSteps) == len(exercise.Steps) {
		exercise.Completed = true
	}

	if err := h.store.Update(ctx, exercise); err != nil {
		http.Error(w, "Failed to update steps", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercise)
}

// ToggleExerciseCompletion bascule le statut de complétion
func (h *ExerciseHandler) ToggleExerciseCompletion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	var input struct {
		Completed bool `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exercise, err := h.store.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	exercise.Completed = input.Completed
	exercise.UpdatedAt = time.Now()

	if input.Completed && len(exercise.CompletedSteps) == 0 {
		exercise.CompletedSteps = make([]int, len(exercise.Steps))
		for i := range exercise.CompletedSteps {
			exercise.CompletedSteps[i] = i
		}
	}

	if err := h.store.Update(ctx, exercise); err != nil {
		http.Error(w, "Failed to update completion", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercise)
}

// ReviewExercise enregistre une révision
// ReviewExercise enregistre une révision
func (h *ExerciseHandler) ReviewExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	var input models.ReviewInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exercise, err := h.store.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	// Appliquer l'algorithme SM-2
	h.scheduler.ReviewExercise(exercise, input.Rating)

	if err := h.store.Update(ctx, exercise); err != nil {
		http.Error(w, "Failed to update exercise", http.StatusInternalServerError)
		return
	}

	nextReview := h.scheduler.GetNextReviewDate(exercise)

	// ✅ CORRECTION: Déréférencer le pointeur et enlever Message
	response := models.ReviewResponse{
		Exercise:   *exercise, // ✅ Déréférencer le pointeur
		NextReview: nextReview,
		// ❌ Pas de Message dans ReviewResponse
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetRecommendations retourne les exercices recommandés
func (h *ExerciseHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "Failed to load exercises", http.StatusInternalServerError)
		return
	}

	recommended := h.recommender.GetNextExercises(exercises, 5)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommended)
}

// GetStats retourne les statistiques globales
func (h *ExerciseHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exercises, err := h.store.Load(ctx)
	if err != nil {
		http.Error(w, "Failed to load exercises", http.StatusInternalServerError)
		return
	}

	stats := h.recommender.CalculateStats(exercises)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ✅ CORRECTION 2: GetTodayPlan avec GetTodayPlan()
// GetTodayPlan retourne le plan du jour
func (h *ExerciseHandler) GetTodayPlan(w http.ResponseWriter, r *http.Request) {
	plan := h.planner.GetToday() // ✅ Utilise GetToday() au lieu de GetTodayPlan()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

// GetWeekPlan retourne le plan de la semaine
func (h *ExerciseHandler) GetWeekPlan(w http.ResponseWriter, r *http.Request) {
	plan := h.planner.GetWeek() // ✅ Utilise GetWeek() au lieu de GetWeekPlan()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

// GetPlannerStats retourne les statistiques du planner
func (h *ExerciseHandler) GetPlannerStats(w http.ResponseWriter, r *http.Request) {
	stats := h.planner.GetStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// CreatePlannerSession crée une nouvelle session planifiée
func (h *ExerciseHandler) CreatePlannerSession(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Date        string   `json:"date"`
		TimeSlot    string   `json:"time_slot"`
		ExerciseIDs []string `json:"exercise_ids"`
		Duration    int      `json:"duration"`
		Status      string   `json:"status"`
		Notes       string   `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validation
	if input.Date == "" {
		http.Error(w, "Date is required", http.StatusBadRequest)
		return
	}
	if input.TimeSlot == "" {
		http.Error(w, "Time slot is required", http.StatusBadRequest)
		return
	}
	if len(input.ExerciseIDs) == 0 {
		http.Error(w, "At least one exercise is required", http.StatusBadRequest)
		return
	}

	// Parser la date
	dateTime, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid date format: %v", err), http.StatusBadRequest)
		return
	}

	// Créer la session
	session := models.PlannedSession{
		ID:          uuid.New().String(),
		Date:        dateTime,
		TimeSlot:    input.TimeSlot,
		ExerciseIDs: input.ExerciseIDs,
		Duration:    input.Duration,
		Status:      "planned", // Toujours "planned" à la création
		Notes:       input.Notes,
	}

	h.planner.AddSession(session) // ✅ CORRECTION 4: Méthode existe dans planner.go

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func (h *ExerciseHandler) UpdatePlannerSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var session models.PlannedSession
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session.ID = id

	// ✅ CORRECTION: Appeler UpdateSession avec session complet
	if err := h.planner.UpdateSession(session); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update session: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// DeletePlannerSession supprime une session
func (h *ExerciseHandler) DeletePlannerSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.planner.DeleteSession(id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete session: %v", err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
