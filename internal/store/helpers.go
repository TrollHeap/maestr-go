package store

import (
	"database/sql"
	"encoding/json"
	"time"

	"maestro/internal/models"
)

// ============================================
// DATE HELPERS (YYYYMMDD)
// ============================================

// toDateInt : Convertit time.Time en YYYYMMDD
func toDateInt(t time.Time) int {
	if t.IsZero() {
		return 0
	}
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}

// fromDateInt : Convertit YYYYMMDD en time.Time
func fromDateInt(dateInt int) time.Time {
	if dateInt == 0 {
		return time.Time{}
	}

	year := dateInt / 10000
	month := (dateInt % 10000) / 100
	day := dateInt % 100

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// todayInt : Date d'aujourd'hui en YYYYMMDD
func todayInt() int {
	return toDateInt(time.Now())
}

// addDays : Ajoute N jours à une date YYYYMMDD
func addDays(dateInt int, days int) int {
	t := fromDateInt(dateInt)
	t = t.AddDate(0, 0, days)
	return toDateInt(t)
}

// placeholders génère placeholders SQL
func placeholders(n int) string {
	if n == 0 {
		return ""
	}
	s := "?"
	for i := 1; i < n; i++ {
		s += ",?"
	}
	return s
}

// ============================================
// QUERY HELPERS
// ============================================

// queryExercisesLight : Requête light (liste)
func queryExercisesLight(query string, args ...interface{}) ([]models.Exercise, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var ex models.Exercise
		var stepsJSON, completedJSON string
		var nextReviewDate int

		rows.Scan(
			&ex.ID, &ex.Title, &ex.Domain, &ex.Difficulty,
			&ex.Done, &nextReviewDate, &completedJSON, &stepsJSON,
		)

		json.Unmarshal([]byte(stepsJSON), &ex.Steps)
		json.Unmarshal([]byte(completedJSON), &ex.CompletedSteps)
		ex.NextReviewAt = fromDateInt(nextReviewDate)

		exercises = append(exercises, ex)
	}

	return exercises, nil
}

// queryExercisesFull : Requête complète (détails)
func queryExercisesFull(query string, args ...interface{}) ([]models.Exercise, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var ex models.Exercise
		var stepsJSON, completedJSON, visualsJSON string
		var lastReviewedDate, lastSkippedDate, nextReviewDate, createdAt, updatedAt sql.NullInt64

		err := rows.Scan(
			&ex.ID, &ex.Title, &ex.Description, &ex.Domain, &ex.Difficulty,
			&ex.Content, &ex.Mnemonic, &visualsJSON,
			&stepsJSON, &completedJSON,
			&ex.Done, &lastReviewedDate, &nextReviewDate,
			&ex.EaseFactor, &ex.IntervalDays, &ex.Repetitions,
			&ex.SkippedCount, &lastSkippedDate,
			&ex.Deleted, &createdAt, &updatedAt,
		)
		if err != nil {
			continue
		}

		parseExerciseFields(&ex, stepsJSON, completedJSON, visualsJSON,
			lastReviewedDate, lastSkippedDate, nextReviewDate, createdAt, updatedAt)

		exercises = append(exercises, ex)
	}

	return exercises, nil
}

// parseExerciseFields : Parse JSON + dates
func parseExerciseFields(ex *models.Exercise,
	stepsJSON, completedJSON, visualsJSON string,
	lastReviewedDate, lastSkippedDate, nextReviewDate, createdAt, updatedAt sql.NullInt64,
) {
	json.Unmarshal([]byte(stepsJSON), &ex.Steps)
	json.Unmarshal([]byte(completedJSON), &ex.CompletedSteps)
	json.Unmarshal([]byte(visualsJSON), &ex.ConceptualVisuals)

	if lastReviewedDate.Valid && lastReviewedDate.Int64 > 0 {
		t := fromDateInt(int(lastReviewedDate.Int64))
		ex.LastReviewed = &t
	}
	if lastSkippedDate.Valid && lastSkippedDate.Int64 > 0 {
		t := fromDateInt(int(lastSkippedDate.Int64))
		ex.LastSkipped = &t
	}
	if nextReviewDate.Valid {
		ex.NextReviewAt = fromDateInt(int(nextReviewDate.Int64))
	}
	if createdAt.Valid {
		ex.CreatedAt = fromDateInt(int(createdAt.Int64))
	}
	if updatedAt.Valid {
		ex.UpdatedAt = fromDateInt(int(updatedAt.Int64))
	}
}

// GetFilteredByQuery : Exécute requête custom
func GetFilteredByQuery(query string, args ...interface{}) ([]models.Exercise, error) {
	return queryExercisesFull(query, args...)
}
