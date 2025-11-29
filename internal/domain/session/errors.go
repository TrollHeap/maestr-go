// internal/domain/session/errors.go
package session

import (
	"fmt"

	"maestro/internal/models"
)

// ============================================
// ERREURS MÉTIER
// ============================================

// NoExercisesAvailableError : Aucun exercice disponible
type NoExercisesAvailableError struct {
	Report models.SessionReport
}

func (e *NoExercisesAvailableError) Error() string {
	if e.Report.NextReviewDate.IsZero() {
		return "Aucun exercice disponible aujourd'hui. Aucune révision programmée."
	}
	return fmt.Sprintf(
		"Aucun exercice disponible aujourd'hui. Prochaine révision : %s (%d à venir)",
		e.Report.NextReviewDate.Format("2006-01-02"),
		e.Report.TodayDue,
	)
}

// InvalidEnergyError : Niveau d'énergie invalide
type InvalidEnergyError struct {
	Level int
}

func (e *InvalidEnergyError) Error() string {
	return fmt.Sprintf("invalid energy level %d (must be 1-3)", e.Level)
}

// SessionNotFoundError : Session introuvable
type SessionNotFoundError struct {
	SessionID int64
}

func (e *SessionNotFoundError) Error() string {
	return fmt.Sprintf("session %d not found", e.SessionID)
}
