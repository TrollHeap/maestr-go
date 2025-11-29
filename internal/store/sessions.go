package store

import (
	"database/sql"
	"fmt"
	"time"

	"maestro/internal/models"
)

// StartSession : Crée nouvelle session
func StartSession(energy models.EnergyLevel, exercises []models.Exercise) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	config := models.SessionConfigs[energy]
	result, err := tx.Exec(`
        INSERT INTO sessions (started_at, energy_level, mode)
        VALUES (?, ?, ?)
    `, time.Now().Unix(), energyToString(energy), config.Mode)
	if err != nil {
		return 0, fmt.Errorf("insert session: %w", err)
	}

	sessionID, _ := result.LastInsertId()

	stmt, _ := tx.Prepare(`
        INSERT INTO session_exercises (session_id, exercise_id, position)
        VALUES (?, ?, ?)
    `)
	defer stmt.Close()

	for i, ex := range exercises {
		_, err := stmt.Exec(sessionID, ex.ID, i)
		if err != nil {
			return 0, fmt.Errorf("insert session exercise: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return sessionID, nil
}

// CompleteSessionExercise : Marque exercice complété
func CompleteSessionExercise(sessionID int64, exerciseID int, quality int) error {
	query := `UPDATE session_exercises SET
        completed = 1,
        quality = ?,
        reviewed_at = ?
    WHERE session_id = ? AND exercise_id = ?`

	_, err := db.Exec(query, quality, time.Now().Unix(), sessionID, exerciseID)
	return err
}

// EndSession : Termine session
func EndSession(sessionID int64) error {
	var startedAt int64
	db.QueryRow("SELECT started_at FROM sessions WHERE id = ?", sessionID).Scan(&startedAt)

	var completedCount int
	db.QueryRow(`
        SELECT COUNT(*) FROM session_exercises 
        WHERE session_id = ? AND completed = 1
    `, sessionID).Scan(&completedCount)

	durationMin := int(time.Since(time.Unix(startedAt, 0)).Minutes())

	query := `UPDATE sessions SET
        ended_at = ?,
        completed_count = ?,
        duration_min = ?
    WHERE id = ?`

	_, err := db.Exec(query, time.Now().Unix(), completedCount, durationMin, sessionID)

	if err == nil {
		updateAnalytics(completedCount, durationMin)
	}

	return err
}

// GetActiveSession : Session en cours
func GetActiveSession() (int64, error) {
	var sessionID int64
	err := db.QueryRow(`
        SELECT id FROM sessions 
        WHERE ended_at IS NULL 
        ORDER BY started_at DESC 
        LIMIT 1
    `).Scan(&sessionID)

	if err == sql.ErrNoRows {
		return 0, nil
	}

	return sessionID, err
}

func energyToString(e models.EnergyLevel) string {
	switch e {
	case models.EnergyLow:
		return "low"
	case models.EnergyMedium:
		return "medium"
	case models.EnergyHigh:
		return "high"
	default:
		return "medium"
	}
}
