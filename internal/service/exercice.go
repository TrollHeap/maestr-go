package service

import (
	"fmt"
	"strings"
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

// ============================================
// CREATE & UPDATE (User Management)
// ============================================

// CreateExercise : Crée un nouvel exercice
func (s *ExerciseService) CreateExercise(ex *models.Exercise) error {
	// 1. Validation domaine (business rules)
	if err := exercise.ValidateExerciseInput(ex.Title, ex.Difficulty, ex.Domain); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// 2. Validation contenu (sécurité)
	if len(ex.Content) > 50000 {
		return fmt.Errorf("content too large: max 50KB")
	}

	// 3. Trim whitespace
	ex.Title = strings.TrimSpace(ex.Title)
	ex.Description = strings.TrimSpace(ex.Description)
	ex.Mnemonic = strings.TrimSpace(ex.Mnemonic)

	// 4. Defaults SRS (domain rules)
	ex.EaseFactor = 2.5
	ex.IntervalDays = 0
	ex.Repetitions = 0
	ex.Done = false
	ex.NextReviewAt = time.Now() // Disponible immédiatement
	ex.CompletedSteps = []int{}  // Aucune étape complétée

	// 5. Defaults visuals si vide
	if ex.ConceptualVisuals == nil {
		ex.ConceptualVisuals = []models.VisualAid{}
	}

	// 6. Insert DB (store layer)
	if err := store.CreateExercise(ex); err != nil {
		return fmt.Errorf("create exercise in store: %w", err)
	}

	return nil
}

// UpdateExercise : Met à jour le contenu d'un exercice existant
func (s *ExerciseService) UpdateExercise(ex *models.Exercise) error {
	// 1. Validation ID
	if err := exercise.ValidateID(ex.ID); err != nil {
		return fmt.Errorf("invalid exercise ID: %w", err)
	}

	// 2. Validation domaine (business rules)
	if err := exercise.ValidateExerciseInput(ex.Title, ex.Difficulty, ex.Domain); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// 3. Validation contenu
	if len(ex.Content) > 50000 {
		return fmt.Errorf("content too large: max 50KB")
	}

	// 4. Trim whitespace
	ex.Title = strings.TrimSpace(ex.Title)
	ex.Description = strings.TrimSpace(ex.Description)
	ex.Mnemonic = strings.TrimSpace(ex.Mnemonic)

	// 5. Vérifie existence + récupère données SRS à préserver
	existing, err := store.FindExercise(ex.ID)
	if err != nil {
		return fmt.Errorf("find existing exercise: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("exercise %d not found", ex.ID)
	}

	// 6. ⚠️ PRÉSERVE les données SRS (pas de reset !)
	ex.EaseFactor = existing.EaseFactor
	ex.IntervalDays = existing.IntervalDays
	ex.Repetitions = existing.Repetitions
	ex.Done = existing.Done
	ex.NextReviewAt = existing.NextReviewAt
	ex.LastReviewed = existing.LastReviewed
	ex.SkippedCount = existing.SkippedCount
	ex.LastSkipped = existing.LastSkipped

	// 7. ⚠️ PRÉSERVE completed_steps (progression utilisateur)
	ex.CompletedSteps = existing.CompletedSteps

	// 8. Update DB (store layer - UPDATE contenu uniquement)
	if err := store.UpdateExercise(ex); err != nil {
		return fmt.Errorf("update exercise in store: %w", err)
	}

	return nil
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

// DeleteExercise : Soft delete d'un exercice (préserve l'historique)
func (s *ExerciseService) DeleteExercise(id int) error {
	// 1. Validation ID (règle métier)
	if err := exercise.ValidateID(id); err != nil {
		return fmt.Errorf("invalid exercise ID: %w", err)
	}

	// 2. Vérifie existence (optionnel mais plus propre pour message d'erreur)
	ex, err := store.FindExercise(id)
	if err != nil {
		return fmt.Errorf("find exercise before delete: %w", err)
	}
	if ex == nil {
		return fmt.Errorf("exercise %d not found", id)
	}

	// 3. Soft delete via store
	if err := store.DeleteExercise(id); err != nil {
		return fmt.Errorf("delete exercise in store: %w", err)
	}

	return nil
}

func (s *ExerciseService) RestoreExercise(id int) error {
	if err := exercise.ValidateID(id); err != nil {
		return fmt.Errorf("invalid exercise ID: %w", err)
	}
	if err := store.RestoreExercise(id); err != nil {
		return fmt.Errorf("restore exercise in store: %w", err)
	}
	return nil
}
