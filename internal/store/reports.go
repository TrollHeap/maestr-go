package store

import (
	"fmt"
	"log"

	"maestro/internal/models"
)

func GetTodayReport() (models.SessionReport, []models.Exercise, error) {
	today := todayInt()

	log.Printf("🔍 [GetTodayReport] today = %d (%s)", today, formatDateInt(today))

	report := models.SessionReport{}

	// 1. Compte exercices dus AUJOURD'HUI ou EN RETARD (ignore done)
	err := db.QueryRow(`
        SELECT COUNT(*) FROM exercises 
        WHERE deleted = 0 
        AND next_review_date > 0
        AND next_review_date <= ?
    `, today).Scan(&report.TodayDue)
	if err != nil {
		log.Printf("🔍 [ERREUR] TodayDue: %v", err)
	}
	log.Printf("🔍 [GetTodayReport] TodayDue (retard+aujourd'hui) = %d ✅", report.TodayDue)

	// 2. Nouveaux (jamais révisés)
	err = db.QueryRow(`
        SELECT COUNT(*) FROM exercises 
        WHERE deleted = 0 
        AND last_reviewed_date IS NULL
    `).Scan(&report.TodayNew)
	if err != nil {
		log.Printf("🔍 [ERREUR] TodayNew: %v", err)
	}
	log.Printf("🔍 [GetTodayReport] TodayNew = %d", report.TodayNew)

	report.TotalAvailable = report.TodayDue + report.TodayNew
	log.Printf("🔍 [SESSION] Total disponible = %d 🚀", report.TotalAvailable)

	// 3. Liste des exercices (tous ceux à réviser AUJOURD'HUI ou EN RETARD)
	query := `
        SELECT id, title, description, domain, difficulty,
               content, mnemonic, conceptual_visuals,
               steps, completed_steps, done, 
               last_reviewed_date, next_review_date,
               ease_factor, interval_days, repetitions,
               skipped_count, last_skipped_date,
               deleted, created_at, updated_at
        FROM exercises 
        WHERE deleted = 0 
        AND next_review_date > 0
        AND next_review_date <= ?
        ORDER BY next_review_date ASC
    `

	exercises, err := queryExercisesFull(query, today)
	if err != nil {
		log.Printf("🔍 [ERREUR] queryExercisesFull: %v", err)
	}
	log.Printf("🔍 [SESSION] %d exercices à pratiquer", len(exercises))

	return report, exercises, nil
}

func getUpcomingReviews(days int) []models.UpcomingReview {
	today := todayInt()
	future := addDays(today, days)

	query := `
        SELECT next_review_date, id, title
        FROM exercises
        WHERE deleted = 0 
        AND done = 1 
        AND next_review_date > ? 
        AND next_review_date <= ?
        ORDER BY next_review_date ASC
        LIMIT 10
    `

	rows, err := db.Query(query, today, future)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var reviews []models.UpcomingReview
	for rows.Next() {
		var dateInt int
		var id int
		var title string

		rows.Scan(&dateInt, &id, &title)

		reviews = append(reviews, models.UpcomingReview{
			Date:          fromDateInt(dateInt),
			ExerciseID:    id,
			ExerciseTitle: title,
		})
	}

	return reviews
}

func formatDateInt(dateInt int) string {
	if dateInt == 0 {
		return "N/A"
	}
	year := dateInt / 10000
	month := (dateInt % 10000) / 100
	day := dateInt % 100
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}
