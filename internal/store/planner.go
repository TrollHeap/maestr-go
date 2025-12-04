package store

import "maestro/internal/models"

// store/planner.go (nouveau fichier ou dans reports.go)

// GetPlannerExercises : filtre par date (pour Planner uniquement)
func GetPlannerExercises(view string) ([]models.Exercise, error) {
	query := `SELECT id, title, domain, difficulty, done, 
                     next_review_date, completed_steps, steps
              FROM exercises WHERE deleted = 0`

	args := []interface{}{}
	today := todayInt()

	switch view {
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
	case "new":
		query += " AND last_reviewed_date IS NULL"
	default:
		// "all" → tous les exercices à réviser (done=1 ou jamais révisés)
		query += " AND (done = 1 OR last_reviewed_date IS NULL)"
	}

	query += " ORDER BY next_review_date ASC, id ASC"
	return queryExercisesLight(query, args...)
}
