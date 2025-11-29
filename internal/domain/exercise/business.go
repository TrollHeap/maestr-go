package exercise

import (
	"slices"

	"maestro/internal/models"
)

// MarkAsDone : Règle métier "marquer exercice DONE"
func MarkAsDone(ex *models.Exercise) {
	ex.Done = true

	// Règle : DONE = toutes les étapes complétées
	ex.CompletedSteps = make([]int, len(ex.Steps))
	for i := range ex.Steps {
		ex.CompletedSteps[i] = i
	}
}

// MarkAsNotDone : Règle métier "marquer exercice TODO"
func MarkAsNotDone(ex *models.Exercise) {
	ex.Done = false
	// Garde les étapes complétées (ne pas reset)
}

// ToggleStep : Règle métier "toggle une étape"
func ToggleStep(ex *models.Exercise, stepIndex int) error {
	if stepIndex < 0 || stepIndex >= len(ex.Steps) {
		return ErrInvalidStep
	}

	// Cherche si l'étape est déjà complétée
	for i, s := range ex.CompletedSteps {
		if s == stepIndex {
			// Retire l'étape
			ex.CompletedSteps = append(
				ex.CompletedSteps[:i],
				ex.CompletedSteps[i+1:]...,
			)
			return nil
		}
	}

	// Ajoute l'étape
	ex.CompletedSteps = append(ex.CompletedSteps, stepIndex)
	return nil
}

// ShouldMarkDone : Règle "quand marquer DONE après review ?"
func ShouldMarkDone(quality int) bool {
	return quality >= 1 // Hard/Good/Easy → DONE
}

// IsStepCompleted : Vérifie si une étape est complétée
func IsStepCompleted(ex *models.Exercise, stepIndex int) bool {
	return slices.Contains(ex.CompletedSteps, stepIndex)
}

// CompletionRate : Taux de complétion des étapes
func CompletionRate(ex *models.Exercise) float64 {
	if len(ex.Steps) == 0 {
		return 1.0
	}
	return float64(len(ex.CompletedSteps)) / float64(len(ex.Steps))
}
