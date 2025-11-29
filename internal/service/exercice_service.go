package service

import (
	"fmt"
	"time"

	"maestro/internal/domain/exercise"
	"maestro/internal/domain/srs"
	"maestro/internal/models"
	"maestro/internal/store"
)

type ExerciseService struct{}

func NewExerciseService() *ExerciseService {
	return &ExerciseService{}
}

// ReviewExercise : Applique SRS + Log historique
func (s *ExerciseService) ReviewExercise(
	exerciseID int,
	quality srs.ReviewQuality,
) (*models.Exercise, error) {
	// 1. Récupère depuis store
	ex, err := store.FindExercise(exerciseID)
	if err != nil || ex == nil {
		return nil, fmt.Errorf("review exercise %d: %w", exerciseID, err)
	}

	// 2. Applique SRS (domain)
	result := srs.CalculateNextReview(
		quality,
		ex.IntervalDays,
		ex.EaseFactor,
		ex.Repetitions,
	)

	// 3. Met à jour modèle
	now := time.Now()
	ex.LastReviewed = &now
	ex.IntervalDays = result.IntervalDays
	ex.EaseFactor = result.EaseFactor
	ex.Repetitions = result.Repetitions
	ex.NextReviewAt = result.NextReview

	// 4. Applique règle métier "mark done" (domain)
	if exercise.ShouldMarkDone(int(quality)) {
		exercise.MarkAsDone(ex) // ✅ Délègue à domain
	}

	// 5. Sauvegarde
	if err := store.SaveExercise(ex); err != nil {
		return nil, fmt.Errorf("save reviewed exercise %d: %w", exerciseID, err)
	}

	// 6. Log historique (non-bloquant)
	if err := store.LogProgress(exerciseID, int(quality), ex); err != nil {
		fmt.Printf("⚠️ Log progress failed: %v\n", err)
	}

	return ex, nil
}

// ToggleExerciseDone : Toggle statut TODO/DONE
func (s *ExerciseService) ToggleExerciseDone(exerciseID int) (*models.Exercise, error) {
	ex, err := store.FindExercise(exerciseID)
	if err != nil || ex == nil {
		return nil, fmt.Errorf("toggle done %d: %w", exerciseID, err)
	}

	// ✅ Applique règle métier (domain)
	if ex.Done {
		exercise.MarkAsNotDone(ex)
	} else {
		exercise.MarkAsDone(ex)
	}

	// Sauvegarde
	if err := store.SaveExercise(ex); err != nil {
		return nil, fmt.Errorf("save toggled exercise %d: %w", exerciseID, err)
	}

	return ex, nil
}

// ToggleExerciseStep : Toggle une étape individuelle
func (s *ExerciseService) ToggleExerciseStep(exerciseID, stepIndex int) (*models.Exercise, error) {
	ex, err := store.FindExercise(exerciseID)
	if err != nil || ex == nil {
		return nil, fmt.Errorf("toggle step exercise %d: %w", exerciseID, err)
	}

	// ✅ Applique règle métier (domain)
	if err := exercise.ToggleStep(ex, stepIndex); err != nil {
		return nil, fmt.Errorf("toggle step %d on exercise %d: %w", stepIndex, exerciseID, err)
	}

	// Sauvegarde
	if err := store.SaveExercise(ex); err != nil {
		return nil, fmt.Errorf("save stepped exercise %d: %w", exerciseID, err)
	}

	return ex, nil
}

// GetExerciseWithMarkdown : Récupère exercice complet
func (s *ExerciseService) GetExerciseWithMarkdown(exerciseID int) (*models.Exercise, error) {
	ex, err := store.FindExercise(exerciseID)
	if err != nil || ex == nil {
		return nil, fmt.Errorf("get exercise %d: %w", exerciseID, err)
	}
	return ex, nil
}

// GetFilteredExercises : Liste filtrée (délègue à store)
func (s *ExerciseService) GetFilteredExercises(
	filter models.ExerciseFilter,
) ([]models.Exercise, error) {
	exercises, err := store.GetFiltered(filter)
	if err != nil {
		return nil, fmt.Errorf("get filtered exercises: %w", err)
	}
	return exercises, nil
}

// GetAllExercises : Tous les exercices
func (s *ExerciseService) GetAllExercises() ([]models.Exercise, error) {
	return s.GetFilteredExercises(models.ExerciseFilter{})
}

// GetExerciseStats : Stats par vue (délègue à store)
func (s *ExerciseService) GetExerciseStats() map[string]int {
	return map[string]int{
		"urgent":   store.CountByView("urgent"),
		"today":    store.CountByView("today"),
		"upcoming": store.CountByView("upcoming"),
		"active":   store.CountByView("active"),
		"new":      store.CountByView("new"),
	}
}

// GetExerciseHistory : Historique d'un exercice
func (s *ExerciseService) GetExerciseHistory(
	exerciseID int,
	limit int,
) ([]map[string]any, error) {
	history, err := store.GetProgressHistory(exerciseID, limit)
	if err != nil {
		return nil, fmt.Errorf("get history for exercise %d: %w", exerciseID, err)
	}
	return history, nil
}
