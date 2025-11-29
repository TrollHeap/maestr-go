package store

import (
	"database/sql"
	"fmt"
	"time"

	"maestro/internal/domain/session" // ✅ NOUVEAU
	"maestro/internal/models"
)

// ============================================
// SESSION CRUD
// ============================================

// StartSession : Crée nouvelle session en DB
func StartSession(energy models.EnergyLevel, exercises []models.Exercise) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Récupère config depuis domain
	config := session.GetConfig(energy)

	// Insert session
	result, err := tx.Exec(`
        INSERT INTO sessions (started_at, energy_level, mode)
        VALUES (?, ?, ?)
    `, time.Now().Unix(), energyToString(energy), config.Mode)
	if err != nil {
		return 0, fmt.Errorf("insert session: %w", err)
	}

	sessionID, _ := result.LastInsertId()

	// Insert exercices de la session
	stmt, err := tx.Prepare(`
        INSERT INTO session_exercises (session_id, exercise_id, position)
        VALUES (?, ?, ?)
    `)
	if err != nil {
		return 0, fmt.Errorf("prepare session exercises: %w", err)
	}
	defer stmt.Close()

	for i, ex := range exercises {
		_, err := stmt.Exec(sessionID, ex.ID, i)
		if err != nil {
			return 0, fmt.Errorf("insert session exercise %d: %w", ex.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit session: %w", err)
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

	result, err := db.Exec(query, quality, time.Now().Unix(), sessionID, exerciseID)
	if err != nil {
		return fmt.Errorf("update session exercise: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("exercise %d not found in session %d", exerciseID, sessionID)
	}

	return nil
}

// EndSession : Termine session
func EndSession(sessionID int64) error {
	// Récupère heure de début
	var startedAt int64
	err := db.QueryRow("SELECT started_at FROM sessions WHERE id = ?", sessionID).Scan(&startedAt)
	if err != nil {
		return fmt.Errorf("query session start time: %w", err)
	}

	// Compte exercices complétés
	var completedCount int
	err = db.QueryRow(`
        SELECT COUNT(*) FROM session_exercises 
        WHERE session_id = ? AND completed = 1
    `, sessionID).Scan(&completedCount)
	if err != nil {
		return fmt.Errorf("count completed exercises: %w", err)
	}

	// Calcule durée
	durationMin := int(time.Since(time.Unix(startedAt, 0)).Minutes())

	// Update session
	query := `UPDATE sessions SET
        ended_at = ?,
        completed_count = ?,
        duration_min = ?
    WHERE id = ?`

	_, err = db.Exec(query, time.Now().Unix(), completedCount, durationMin, sessionID)
	if err != nil {
		return fmt.Errorf("update session end: %w", err)
	}

	// Update analytics (non-bloquant)
	if err := updateAnalytics(completedCount, durationMin); err != nil {
		fmt.Printf("⚠️ Update analytics failed: %v\n", err)
	}

	return nil
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

	if err != nil {
		return 0, fmt.Errorf("query active session: %w", err)
	}

	return sessionID, nil
}

// GetNextSessionExercise : Prochain exercice non complété dans la session
func GetNextSessionExercise(sessionID int64) (int, error) {
	query := `SELECT exercise_id
              FROM session_exercises
              WHERE session_id = ? AND completed = 0
              ORDER BY position ASC
              LIMIT 1`

	var exerciseID int
	err := db.QueryRow(query, sessionID).Scan(&exerciseID)

	if err == sql.ErrNoRows {
		return 0, nil // Plus d'exercices
	}

	if err != nil {
		return 0, fmt.Errorf("query next session exercise: %w", err)
	}

	return exerciseID, nil
}

// GetSessionResult : Récupère résultat d'une session terminée
func GetSessionResult(sessionID int64) (*models.SessionResult, error) {
	query := `SELECT 
        completed_count, 
        duration_min, 
        ended_at
    FROM sessions WHERE id = ?`

	var completedCount, durationMin int
	var endedAt int64

	err := db.QueryRow(query, sessionID).Scan(
		&completedCount, &durationMin, &endedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("query session result: %w", err)
	}

	// Récupère IDs des exercices complétés
	exerciseQuery := `SELECT exercise_id, quality 
                      FROM session_exercises 
                      WHERE session_id = ? AND completed = 1
                      ORDER BY position`

	rows, err := db.Query(exerciseQuery, sessionID)
	if err != nil {
		return nil, fmt.Errorf("query session exercises: %w", err)
	}
	defer rows.Close()

	var exerciseIDs []int
	qualities := make(map[int]int)

	for rows.Next() {
		var id, quality int
		if err := rows.Scan(&id, &quality); err != nil {
			return nil, fmt.Errorf("scan session exercise: %w", err)
		}
		exerciseIDs = append(exerciseIDs, id)
		qualities[id] = quality
	}

	return &models.SessionResult{
		SessionID:      sessionID,
		CompletedCount: completedCount,
		Duration:       time.Duration(durationMin) * time.Minute,
		CompletedAt:    time.Unix(endedAt, 0),
		Exercises:      exerciseIDs,
		Qualities:      qualities,
	}, nil
}

// ============================================
// HELPERS
// ============================================

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
