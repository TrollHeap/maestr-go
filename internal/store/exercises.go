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
	today := todayInt()

	// Filtres dynamiques (YYYYMMDD)
	switch filter.View {
	case "urgent":
		query += " AND done = 1 AND next_review_date < ?"
		args = append(args, today)
	case "today":
		tomorrow := addDays(today, 1)
		query += " AND done = 1 AND next_review_date >= ? AND next_review_date < ?"
		args = append(args, today, tomorrow)
	case "upcoming":
		tomorrow := addDays(today, 1)
		in3days := addDays(today, 3)
		query += " AND done = 1 AND next_review_date >= ? AND next_review_date < ?"
		args = append(args, tomorrow, in3days)
	case "active":
		query += " AND done = 0 AND completed_steps != '[]'"
	case "new":
		query += " AND done = 0 AND (completed_steps = '[]' OR completed_steps IS NULL)"
	}

	if filter.Domain != "" {
		query += " AND domain = ?"
		args = append(args, filter.Domain)
	}

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

// CountByView : Compte exercices par vue
func CountByView(view string) int {
	exercises, _ := GetFiltered(models.ExerciseFilter{View: view})
	return len(exercises)
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
