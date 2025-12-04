package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"maestro/internal/models"
)

func GetFiltered(filter models.ExerciseFilter) ([]models.Exercise, error) {
	query := `SELECT id, title, domain, difficulty, done, 
                     next_review_date, completed_steps, steps
              FROM exercises WHERE deleted = 0`

	args := []interface{}{}

	// ✅ FILTRES CONTENU (pas de temps)

	// 1. Statut (done)
	if filter.Status != "" {
		switch filter.Status {
		case "in_progress":
			query += " AND done = 0"
		case "mastered":
			query += " AND done = 1"
		}
	}

	// 2. Domaine
	if filter.Domain != "" {
		query += " AND domain = ?"
		args = append(args, filter.Domain)
	}

	// 3. Difficulté
	if filter.Difficulty > 0 {
		query += " AND difficulty = ?"
		args = append(args, filter.Difficulty)
	}

	query += " ORDER BY id ASC"
	return queryExercisesLight(query, args...)
}

// FindExercise : Lecture complète d'un exercice
func FindExercise(id int) (*models.Exercise, error) {
	query := `SELECT 
        id, title, description, domain, difficulty,
        content, mnemonic, conceptual_visuals,
        steps, completed_steps,
        done, last_reviewed_date, next_review_date,
        ease_factor, interval_days, repetitions,
        skipped_count, last_skipped_date,
        deleted, created_at, updated_at
    FROM exercises 
    WHERE id = ? AND deleted = 0`

	var ex models.Exercise
	var stepsJSON, completedJSON, visualsJSON string
	var lastReviewed, lastSkipped, nextReview, createdAt, updatedAt sql.NullInt64

	err := db.QueryRow(query, id).Scan(
		&ex.ID, &ex.Title, &ex.Description, &ex.Domain, &ex.Difficulty,
		&ex.Content, &ex.Mnemonic, &visualsJSON,
		&stepsJSON, &completedJSON,
		&ex.Done, &lastReviewed, &nextReview,
		&ex.EaseFactor, &ex.IntervalDays, &ex.Repetitions,
		&ex.SkippedCount, &lastSkipped,
		&ex.Deleted, &createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find exercise: %w", err)
	}

	// Parse JSON + timestamps
	parseExerciseFields(&ex, stepsJSON, completedJSON, visualsJSON,
		lastReviewed, lastSkipped, nextReview, createdAt, updatedAt)

	return &ex, nil
}

// SaveExercise : UPDATE atomique
func SaveExercise(ex *models.Exercise) error {
	stepsJSON, _ := json.Marshal(ex.Steps)
	completedJSON, _ := json.Marshal(ex.CompletedSteps)
	visualsJSON, _ := json.Marshal(ex.ConceptualVisuals)

	var lastReviewedDate sql.NullInt64
	if ex.LastReviewed != nil && !ex.LastReviewed.IsZero() {
		lastReviewedDate = sql.NullInt64{Int64: int64(toDateInt(*ex.LastReviewed)), Valid: true}
	}

	var lastSkippedDate sql.NullInt64
	if ex.LastSkipped != nil && !ex.LastSkipped.IsZero() {
		lastSkippedDate = sql.NullInt64{Int64: int64(toDateInt(*ex.LastSkipped)), Valid: true}
	}

	nextReviewDate := toDateInt(ex.NextReviewAt)
	updatedAt := todayInt()

	query := `UPDATE exercises SET
        title = ?, description = ?, content = ?,
        mnemonic = ?, conceptual_visuals = ?,
        steps = ?, completed_steps = ?,
        done = ?, last_reviewed_date = ?, next_review_date = ?,
        ease_factor = ?, interval_days = ?, repetitions = ?,
        skipped_count = ?, last_skipped_date = ?,
        updated_at = ?
    WHERE id = ?`

	_, err := db.Exec(query,
		ex.Title, ex.Description, ex.Content,
		ex.Mnemonic, visualsJSON,
		stepsJSON, completedJSON,
		ex.Done, lastReviewedDate, nextReviewDate,
		ex.EaseFactor, ex.IntervalDays, ex.Repetitions,
		ex.SkippedCount, lastSkippedDate,
		updatedAt,
		ex.ID,
	)

	return err
}

// GetAll : Tous les exercices
func GetAll() []models.Exercise {
	exercises, _ := GetFiltered(models.ExerciseFilter{})
	return exercises
}

// CreateExercise : INSERT nouveau + RETURNING id
func CreateExercise(ex *models.Exercise) error {
	// 1. Serialize JSON
	stepsJSON, err := json.Marshal(ex.Steps)
	if err != nil {
		return fmt.Errorf("marshal steps: %w", err)
	}

	visualsJSON := "[]" // Vide par défaut
	if len(ex.ConceptualVisuals) > 0 {
		data, _ := json.Marshal(ex.ConceptualVisuals)
		visualsJSON = string(data)
	}

	// 2. Dates
	now := todayInt() // YYYYMMDD

	// 3. INSERT avec RETURNING id (SQLite 3.35+)
	query := `
        INSERT INTO exercises (
            title, description, domain, difficulty,
            content, mnemonic, conceptual_visuals,
            steps, completed_steps,
            done, next_review_date,
            ease_factor, interval_days, repetitions,
            deleted, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        RETURNING id
    `

	err = db.QueryRow(query,
		ex.Title, ex.Description, ex.Domain, ex.Difficulty,
		ex.Content, ex.Mnemonic, visualsJSON,
		stepsJSON, "[]", // completed_steps vide
		0,         // done = false
		now,       // next_review_date = aujourd'hui
		2.5, 0, 0, // ease_factor, interval, repetitions (défauts SRS)
		0, now, now, // deleted, created_at, updated_at
	).Scan(&ex.ID)
	if err != nil {
		// Check contrainte UNIQUE (si tu l'ajoutes)
		if err.Error() == "UNIQUE constraint failed: exercises.title" {
			return fmt.Errorf("un exercice avec ce titre existe déjà")
		}
		return fmt.Errorf("insert exercise: %w", err)
	}

	return nil
}

// UpdateExercise : UPDATE contenu SANS toucher aux données SRS
func UpdateExercise(ex *models.Exercise) error {
	// 1. Serialize JSON
	stepsJSON, _ := json.Marshal(ex.Steps)
	visualsJSON, _ := json.Marshal(ex.ConceptualVisuals)

	// 2. UPDATE (trigger met à jour updated_at automatiquement)
	query := `
        UPDATE exercises SET
            title = ?,
            description = ?,
            domain = ?,
            difficulty = ?,
            content = ?,
            mnemonic = ?,
            conceptual_visuals = ?,
            steps = ?
        WHERE id = ? AND deleted = 0
    `

	result, err := db.Exec(query,
		ex.Title, ex.Description, ex.Domain, ex.Difficulty,
		ex.Content, ex.Mnemonic, visualsJSON,
		stepsJSON,
		ex.ID,
	)
	if err != nil {
		return fmt.Errorf("update exercise: %w", err)
	}

	// 3. Check si exercice trouvé
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("exercice %d introuvable ou supprimé", ex.ID)
	}

	return nil
}

// GetNextDueExercise : Prochain exercice à réviser
func GetNextDueExercise(fromSession bool, sessionExercises []int) (*models.Exercise, error) {
	today := todayInt() // ✅ YYYYMMDD
	var query string
	var args []interface{}

	if fromSession && len(sessionExercises) > 0 {
		query = `SELECT id FROM exercises 
                 WHERE id IN (` + placeholders(len(sessionExercises)) + `)
                 AND deleted = 0
                 AND (done = 0 OR next_review_date <= ?)
                 ORDER BY next_review_date ASC, id ASC
                 LIMIT 1`
		for _, id := range sessionExercises {
			args = append(args, id)
		}
		args = append(args, today)
	} else {
		query = `SELECT id FROM exercises 
                 WHERE deleted = 0
                 AND (done = 0 OR next_review_date <= ?)
                 ORDER BY next_review_date ASC, id ASC
                 LIMIT 1`
		args = append(args, today)
	}

	var exerciseID int
	err := db.QueryRow(query, args...).Scan(&exerciseID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return FindExercise(exerciseID)
}

// LogProgress : Enregistre révision dans l'historique
func LogProgress(exerciseID int, quality int, ex *models.Exercise) error {
	query := `INSERT INTO progress_log (
        exercise_id, reviewed_at, quality,
        ease_factor, interval_days, repetitions
    ) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query,
		exerciseID, time.Now().Unix(), quality,
		ex.EaseFactor, ex.IntervalDays, ex.Repetitions,
	)

	return err
}

// GetProgressHistory : Historique révisions
func GetProgressHistory(exerciseID int, limit int) ([]map[string]interface{}, error) {
	query := `SELECT reviewed_at, quality, ease_factor, interval_days
              FROM progress_log
              WHERE exercise_id = ?
              ORDER BY reviewed_at DESC
              LIMIT ?`

	rows, err := db.Query(query, exerciseID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var reviewedAt int64
		var quality, intervalDays int
		var easeFactor float64

		rows.Scan(&reviewedAt, &quality, &easeFactor, &intervalDays)

		history = append(history, map[string]interface{}{
			"reviewed_at":   time.Unix(reviewedAt, 0),
			"quality":       quality,
			"ease_factor":   easeFactor,
			"interval_days": intervalDays,
		})
	}

	return history, nil
}

// DeleteExercise : soft delete (marque deleted = 1, deleted_at = today)
func DeleteExercise(id int) error {
	today := todayInt()

	query := `
        UPDATE exercises
        SET deleted = 1,
            updated_at = ?
        WHERE id = ? AND deleted = 0
    `

	result, err := db.Exec(query, today, id)
	if err != nil {
		return fmt.Errorf("delete exercise: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("exercise %d not found or already deleted", id)
	}

	return nil
}

// RestoreExercise : restaure un exercice soft-deleted (optionnel)
func RestoreExercise(id int) error {
	today := todayInt()

	query := `
        UPDATE exercises
        SET deleted = 0,
            updated_at = ?
        WHERE id = ? AND deleted = 1
    `

	result, err := db.Exec(query, today, id)
	if err != nil {
		return fmt.Errorf("restore exercise: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("exercise %d not found or not deleted", id)
	}

	return nil
}

// HardDeleteExercise : suppression définitive (si un jour tu veux purger)
func HardDeleteExercise(id int) error {
	query := `DELETE FROM exercises WHERE id = ?`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("hard delete exercise: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("exercise %d not found", id)
	}

	return nil
}
