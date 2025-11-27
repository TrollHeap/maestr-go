package service

import (
	"fmt"
	"time"

	"maestro/internal/models"
	"maestro/internal/srs"
	"maestro/internal/store"
)

// ExerciseService gère la logique métier des exercices
type ExerciseService struct{}

// NewExerciseService crée une nouvelle instance
func NewExerciseService() *ExerciseService {
	return &ExerciseService{}
}

// ReviewExercise enregistre une révision et applique l'algorithme SRS
func (s *ExerciseService) ReviewExercise(
	exerciseID int,
	quality srs.ReviewQuality,
) (*models.Exercise, error) {
	ex := store.FindExercise(exerciseID)
	if ex == nil {
		return nil, fmt.Errorf("exercice %d introuvable", exerciseID)
	}

	// Appel de l'algorithme SRS (calcul pur)
	result := srs.CalculateNextReview(
		quality,
		ex.IntervalDays,
		ex.EaseFactor,
		ex.Repetitions,
	)

	// Application du résultat
	now := time.Now()
	ex.LastReviewed = &now
	ex.IntervalDays = result.IntervalDays
	ex.EaseFactor = result.EaseFactor
	ex.Repetitions = result.Repetitions
	ex.NextReviewAt = result.NextReview

	if err := store.Save(); err != nil {
		return nil, fmt.Errorf("erreur sauvegarde: %w", err)
	}

	return ex, nil
}

// ToggleExerciseDone bascule entre TODO et DONE (logique simplifiée)
func (s *ExerciseService) ToggleExerciseDone(exerciseID int) (*models.Exercise, error) {
	ex := store.FindExercise(exerciseID)
	if ex == nil {
		return nil, fmt.Errorf("exercice %d introuvable", exerciseID)
	}

	// Toggle simple : inverse le statut
	ex.Done = !ex.Done

	// Si on marque DONE : complète toutes les étapes
	if ex.Done {
		ex.CompletedSteps = []int{}
		for i := range ex.Steps {
			ex.CompletedSteps = append(ex.CompletedSteps, i)
		}
	}
	// Si on marque TODO : garde les étapes complétées (permet de reprendre)

	if err := store.Save(); err != nil {
		return nil, fmt.Errorf("erreur sauvegarde: %w", err)
	}

	return ex, nil
}

// ToggleExerciseStep toggle une étape individuelle
func (s *ExerciseService) ToggleExerciseStep(exerciseID, step int) (*models.Exercise, error) {
	ex := store.FindExercise(exerciseID)
	if ex == nil {
		return nil, fmt.Errorf("exercice %d introuvable", exerciseID)
	}

	// Toggle de l'étape
	found := false
	for i, s := range ex.CompletedSteps {
		if s == step {
			ex.CompletedSteps = append(ex.CompletedSteps[:i], ex.CompletedSteps[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		ex.CompletedSteps = append(ex.CompletedSteps, step)
	}

	if err := store.Save(); err != nil {
		return nil, fmt.Errorf("erreur sauvegarde: %w", err)
	}

	return ex, nil
}

// GetExerciseWithMarkdown récupère un exercice
func (s *ExerciseService) GetExerciseWithMarkdown(exerciseID int) (*models.Exercise, error) {
	ex := store.FindExercise(exerciseID)
	if ex == nil {
		return nil, fmt.Errorf("exercice %d introuvable", exerciseID)
	}
	return ex, nil
}

// GetFilteredExercises récupère les exercices filtrés
func (s *ExerciseService) GetFilteredExercises(filter models.ExerciseFilter) []models.Exercise {
	return store.GetFiltered(filter)
}

// GetAllExercises récupère tous les exercices
func (s *ExerciseService) GetAllExercises() []models.Exercise {
	return store.GetAll()
}

// GetExerciseStats récupère les statistiques par vue
func (s *ExerciseService) GetExerciseStats() map[string]int {
	return map[string]int{
		"urgent":   store.CountByView("urgent"),
		"today":    store.CountByView("today"),
		"upcoming": store.CountByView("upcoming"),
		"active":   store.CountByView("active"),
		"new":      store.CountByView("new"),
	}
}
