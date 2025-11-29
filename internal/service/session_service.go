package service

import (
	"database/sql"
	"fmt"
	"time"

	"maestro/internal/domain/session" // ✅ NOUVEAU
	"maestro/internal/models"
	"maestro/internal/store"
)

type SessionService struct{}

func NewSessionService() *SessionService {
	return &SessionService{}
}

// StartSession : Crée une session adaptative
func (s *SessionService) StartSession(
	energy models.EnergyLevel,
	exerciseIDs []int, // ✅ Reçoit directement les IDs limités
) (int64, *models.AdaptiveSession, error) {
	// 1. Récupère config depuis domain
	config := session.GetConfig(energy)

	// 2. Charge exercices complets depuis store
	exercises := make([]models.Exercise, 0, len(exerciseIDs))
	for _, id := range exerciseIDs {
		ex, err := store.FindExercise(id)
		if err != nil {
			return 0, nil, fmt.Errorf("find exercise %d: %w", id, err)
		}
		exercises = append(exercises, *ex)
	}

	// 3. Trie par priorité (domain)
	exercises = session.SortByPriority(exercises)

	// 4. Build session model
	sessionModel := models.AdaptiveSession{
		Mode:          config.Mode,
		EnergyLevel:   energy,
		EstimatedTime: config.Duration,
		Exercises:     exerciseIDs, // Garde les IDs uniquement
		BreakSchedule: config.BreakSchedule,
		StartedAt:     time.Now(),
		CurrentIndex:  0,
	}

	// 5. Stocke dans SQLite
	sessionID, err := store.StartSession(energy, exercises)
	if err != nil {
		return 0, nil, fmt.Errorf("start session: %w", err)
	}

	sessionModel.ID = sessionID

	return sessionID, &sessionModel, nil
}

// CompleteExercise : Marque un exercice comme complété dans la session
func (s *SessionService) CompleteExercise(sessionID int64, exerciseID int, quality int) error {
	if err := store.CompleteSessionExercise(sessionID, exerciseID, quality); err != nil {
		return fmt.Errorf("complete exercise %d in session %d: %w", exerciseID, sessionID, err)
	}
	return nil
}

// EndSession : Termine une session
func (s *SessionService) EndSession(sessionID int64) error {
	if err := store.EndSession(sessionID); err != nil {
		return fmt.Errorf("end session %d: %w", sessionID, err)
	}
	return nil
}

// GetActiveSession : Session en cours (retourne ID ou 0)
func (s *SessionService) GetActiveSession() (int64, error) {
	sessionID, err := store.GetActiveSession()
	if err != nil {
		return 0, fmt.Errorf("get active session: %w", err)
	}
	return sessionID, nil
}

// ClearAllSessions : Ferme toutes les sessions actives
func (s *SessionService) ClearAllSessions() error {
	sessionID, err := s.GetActiveSession()
	if err != nil || sessionID == 0 {
		return err
	}
	return s.EndSession(sessionID)
}

// GetSessionResult : Récupère le résultat d'une session terminée
func (s *SessionService) GetSessionResult(sessionID int64) (*models.SessionResult, error) {
	result, err := store.GetSessionResult(sessionID)
	if err != nil {
		return nil, fmt.Errorf("get session result %d: %w", sessionID, err)
	}
	return result, nil
}

// GetNextExercise : Prochain exercice dans la session
func (s *SessionService) GetNextExercise(sessionID int64) (*models.Exercise, error) {
	exerciseID, err := store.GetNextSessionExercise(sessionID)
	if err == sql.ErrNoRows || exerciseID == 0 {
		return nil, nil // Plus d'exercices
	}
	if err != nil {
		return nil, fmt.Errorf("get next exercise in session %d: %w", sessionID, err)
	}

	// Charge l'exercice complet
	ex, err := store.FindExercise(exerciseID)
	if err != nil {
		return nil, fmt.Errorf("find exercise %d: %w", exerciseID, err)
	}

	return ex, nil
}

// StopSession : Alias pour EndSession (compatibilité handlers)
func (s *SessionService) StopSession(sessionID int64) error {
	return s.EndSession(sessionID)
}
