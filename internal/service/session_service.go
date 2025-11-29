package service

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"maestro/internal/models"
	"maestro/internal/store"
)

type SessionService struct{}

func NewSessionService() *SessionService {
	return &SessionService{}
}

// StartSession : Crée une session adaptative
// StartSession : Crée une session adaptative
// StartSession : Crée une session adaptative
func (s *SessionService) StartSession(
	energy models.EnergyLevel,
) (int64, *models.AdaptiveSession, error) {
	config := models.SessionConfigs[energy]

	// 1. Récupère rapport + exercices disponibles aujourd'hui
	report, exercises, err := store.GetTodayReport()
	if err != nil {
		return 0, nil, err
	}

	// 2. Si aucun exercice disponible aujourd'hui
	if len(exercises) == 0 {
		return 0, nil, &models.NoExercisesTodayError{
			Report: report,
		}
	}

	// Tri par priorité
	exercises = sortByPriority(exercises)

	// Limite au nombre d'exercices configuré
	if len(exercises) > config.ExerciseCount {
		exercises = exercises[:config.ExerciseCount]
	}

	// 3. Build session
	session := models.AdaptiveSession{
		Mode:          config.Mode,
		EnergyLevel:   energy,
		EstimatedTime: config.Duration,
		Exercises:     exercises,
		BreakSchedule: config.Breaks,
		StartedAt:     time.Now(),
		CurrentIndex:  0,
	}

	// 4. Stocke dans SQLite
	sessionID, err := store.StartSession(energy, exercises)
	if err != nil {
		return 0, nil, err
	}

	return sessionID, &session, nil
}

// ============================================
// HELPERS
// ============================================
// CompleteExercise : Marque un exercice comme complété dans la session
func (s *SessionService) CompleteExercise(sessionID int64, exerciseID int, quality int) error {
	return store.CompleteSessionExercise(sessionID, exerciseID, quality)
}

// EndSession : Termine une session
func (s *SessionService) EndSession(sessionID int64) error {
	return store.EndSession(sessionID)
}

// GetActiveSession : Session en cours (retourne ID ou 0)
func (s *SessionService) GetActiveSession() (int64, error) {
	return store.GetActiveSession()
}

// ClearAllSessions : Ferme toutes les sessions actives
func (s *SessionService) ClearAllSessions() error {
	sessionID, err := store.GetActiveSession()
	if err != nil || sessionID == 0 {
		return err
	}
	return store.EndSession(sessionID)
}

// GetSessionResult : Récupère le résultat d'une session terminée
func (s *SessionService) GetSessionResult(sessionID int64) (*models.SessionResult, error) {
	query := `SELECT 
        completed_count, 
        duration_min, 
        ended_at
    FROM sessions WHERE id = ?`

	var completedCount, durationMin int
	var endedAt int64

	db := store.GetDB()
	err := db.QueryRow(query, sessionID).Scan(
		&completedCount, &durationMin, &endedAt,
	)
	if err != nil {
		return nil, err
	}

	// Récupère les IDs des exercices complétés
	exerciseQuery := `SELECT exercise_id FROM session_exercises 
                      WHERE session_id = ? AND completed = 1
                      ORDER BY position`

	rows, _ := db.Query(exerciseQuery, sessionID)
	defer rows.Close()

	var exerciseIDs []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		exerciseIDs = append(exerciseIDs, id)
	}

	return &models.SessionResult{
		CompletedCount: completedCount,
		Duration:       time.Duration(durationMin) * time.Minute,
		CompletedAt:    time.Unix(endedAt, 0),
		Exercises:      exerciseIDs,
	}, nil
}

// GetNextExercise : Prochain exercice dans la session
func (s *SessionService) GetNextExercise(sessionID int64) (*models.Exercise, error) {
	// Récupère l'ID du prochain exercice non complété
	query := `SELECT se.exercise_id
              FROM session_exercises se
              WHERE se.session_id = ? AND se.completed = 0
              ORDER BY se.position ASC
              LIMIT 1`

	var exerciseID int
	db := store.GetDB()
	err := db.QueryRow(query, sessionID).Scan(&exerciseID)

	if err == sql.ErrNoRows {
		return nil, nil // Plus d'exercices
	}
	if err != nil {
		return nil, fmt.Errorf("query next exercise: %w", err)
	}

	// Charge l'exercice complet
	return store.FindExercise(exerciseID)
}

// StopSession : Alias pour EndSession
func (s *SessionService) StopSession(sessionID int64) error {
	return s.EndSession(sessionID)
}

// ============================================
// HELPERS
// ============================================

func sortByPriority(exercises []models.Exercise) []models.Exercise {
	now := time.Now()

	sort.Slice(exercises, func(i, j int) bool {
		a, b := exercises[i], exercises[j]

		// Priorité 1 : En retard (urgent)
		aOverdue := a.Done && a.NextReviewAt.Before(now)
		bOverdue := b.Done && b.NextReviewAt.Before(now)
		if aOverdue != bOverdue {
			return aOverdue
		}

		// Priorité 2 : À réviser aujourd'hui
		today := now.Truncate(24 * time.Hour)
		tomorrow := today.Add(24 * time.Hour)
		aToday := a.Done && a.NextReviewAt.After(today) && a.NextReviewAt.Before(tomorrow)
		bToday := b.Done && b.NextReviewAt.After(today) && b.NextReviewAt.Before(tomorrow)
		if aToday != bToday {
			return aToday
		}

		// Priorité 3 : Nouveaux
		if !a.Done && !b.Done {
			return a.ID < b.ID
		}

		// Priorité 4 : Par date
		return a.NextReviewAt.Before(b.NextReviewAt)
	})

	return exercises
}
