package store

import (
	"time"
)

// GetAnalytics : Métriques globales
func GetAnalytics() (map[string]any, error) {
	query := `SELECT 
        avg_session_length_min,
        current_streak,
        longest_streak,
        total_sessions,
        total_exercises_done
    FROM analytics WHERE id = 1`

	var avgLength float64
	var currentStreak, longestStreak, totalSessions, totalExercises int

	err := db.QueryRow(query).Scan(
		&avgLength, &currentStreak, &longestStreak,
		&totalSessions, &totalExercises,
	)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"avg_session_length": avgLength,
		"current_streak":     currentStreak,
		"longest_streak":     longestStreak,
		"total_sessions":     totalSessions,
		"total_exercises":    totalExercises,
	}, nil
}

func updateAnalytics(completedCount, durationMin int) error {
	query := `UPDATE analytics SET
        total_sessions = total_sessions + 1,
        total_exercises_done = total_exercises_done + ?,
        last_session_date = ?,
        updated_at = ?
    WHERE id = 1`

	_, err := db.Exec(query, completedCount, time.Now().Unix(), time.Now().Unix())
	return err
}
