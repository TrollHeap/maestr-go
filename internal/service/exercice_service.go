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

func NewExerciseService() *ExerciseService {
	return &ExerciseService{}
}

// ReviewExercise : Applique SRS + Log historique
func (s *ExerciseService) ReviewExercise(
	exerciseID int,
	quality srs.ReviewQuality,
) (*models.Exercise, error) {
	// 1. Charge depuis SQLite
	ex, err := store.FindExercise(exerciseID)
	if err != nil || ex == nil {
		return nil, fmt.Errorf("exercice %d introuvable", exerciseID)
	}

	// 2. Calcul SRS (pur, sans effet de bord)
	result := srs.CalculateNextReview(
		quality,
		ex.IntervalDays,
		ex.EaseFactor,
		ex.Repetitions,
	)

	// 3. Applique le résultat
	now := time.Now()
	ex.LastReviewed = &now
	ex.IntervalDays = result.IntervalDays
	ex.EaseFactor = result.EaseFactor
	ex.Repetitions = result.Repetitions
	ex.NextReviewAt = result.NextReview

	// 4. Sauve dans SQLite
	if err := store.SaveExercise(ex); err != nil {
		return nil, fmt.Errorf("erreur sauvegarde: %w", err)
	}

	// 5. ✅ NOUVEAU : Log dans l'historique
	if err := store.LogProgress(exerciseID, int(quality), ex); err != nil {
		// Non-bloquant, log l'erreur
		fmt.Printf("Avertissement: échec log progress: %v\n", err)
	}

	return ex, nil
}

// ToggleExerciseDone : Toggle statut TODO/DONE
func (s *ExerciseService) ToggleExerciseDone(exerciseID int) (*models.Exercise, error) {
	ex, err := store.FindExercise(exerciseID)
	if err != nil || ex == nil {
		return nil, fmt.Errorf("exercice %d introuvable", exerciseID)
	}

	// Toggle
	ex.Done = !ex.Done

	// Si DONE → complète toutes les étapes
	if ex.Done {
		ex.CompletedSteps = []int{}
		for i := range ex.Steps {
			ex.CompletedSteps = append(ex.CompletedSteps, i)
		}
	}

	// Sauve
	if err := store.SaveExercise(ex); err != nil {
		return nil, fmt.Errorf("erreur sauvegarde: %w", err)
	}

	return ex, nil
}

// ToggleExerciseStep : Toggle une étape individuelle
func (s *ExerciseService) ToggleExerciseStep(exerciseID, step int) (*models.Exercise, error) {
	ex, err := store.FindExercise(exerciseID)
	if err != nil || ex == nil {
		return nil, fmt.Errorf("exercice %d introuvable", exerciseID)
	}

	// Toggle étape
	found := false
	for i, s := range ex.CompletedSteps {
		if s == step {
			// Retire l'étape
			ex.CompletedSteps = append(ex.CompletedSteps[:i], ex.CompletedSteps[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		// Ajoute l'étape
		ex.CompletedSteps = append(ex.CompletedSteps, step)
	}

	// Sauve
	if err := store.SaveExercise(ex); err != nil {
		return nil, fmt.Errorf("erreur sauvegarde: %w", err)
	}

	return ex, nil
}

// GetExerciseWithMarkdown : Récupère exercice complet
func (s *ExerciseService) GetExerciseWithMarkdown(exerciseID int) (*models.Exercise, error) {
	ex, err := store.FindExercise(exerciseID)
	if err != nil || ex == nil {
		return nil, fmt.Errorf("exercice %d introuvable", exerciseID)
	}
	return ex, nil
}

// GetFilteredExercises : Liste filtrée
func (s *ExerciseService) GetFilteredExercises(
	filter models.ExerciseFilter,
) ([]models.Exercise, error) {
	return store.GetFiltered(filter)
}

// GetAllExercises : Tous les exercices
func (s *ExerciseService) GetAllExercises() ([]models.Exercise, error) {
	return store.GetFiltered(models.ExerciseFilter{})
}

// GetExerciseStats : Stats par vue
func (s *ExerciseService) GetExerciseStats() map[string]int {
	return map[string]int{
		"urgent":   store.CountByView("urgent"),
		"today":    store.CountByView("today"),
		"upcoming": store.CountByView("upcoming"),
		"active":   store.CountByView("active"),
		"new":      store.CountByView("new"),
	}
}

// ✅ NOUVEAU : Historique d'un exercice
func (s *ExerciseService) GetExerciseHistory(
	exerciseID int,
	limit int,
) ([]map[string]any, error) {
	return store.GetProgressHistory(exerciseID, limit)
}
