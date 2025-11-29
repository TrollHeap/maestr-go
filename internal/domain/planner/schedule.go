package planner

import (
	"sort"
	"time"

	"maestro/internal/models"
)

// normalizeDate normalise une date à minuit UTC
func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// IsSameDay compare deux dates ignorants l'heure
func IsSameDay(a, b time.Time) bool {
	return normalizeDate(a).Equal(normalizeDate(b))
}

// GetReviewsForDate filtre les exercices dus à une date donnée
func GetReviewsForDate(exercises []models.Exercise, date time.Time) []models.Exercise {
	var reviews []models.Exercise
	targetDate := normalizeDate(date)

	for _, ex := range exercises {
		if ex.NextReviewAt.IsZero() {
			continue
		}
		if IsSameDay(ex.NextReviewAt, targetDate) {
			reviews = append(reviews, ex)
		}
	}

	// Trie par difficulté décroissante
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].Difficulty > reviews[j].Difficulty
	})

	return reviews
}

// GetOverdueReviews filtre les exercices en retard (non faits et date passée)
func GetOverdueReviews(exercises []models.Exercise) []models.Exercise {
	var overdue []models.Exercise
	now := time.Now()

	for _, ex := range exercises {
		if ex.NextReviewAt.IsZero() {
			continue
		}
		if ex.NextReviewAt.Before(now) && !ex.Done {
			overdue = append(overdue, ex)
		}
	}

	// Trie par date croissante (les plus anciens en premier)
	sort.Slice(overdue, func(i, j int) bool {
		return overdue[i].NextReviewAt.Before(overdue[j].NextReviewAt)
	})

	return overdue
}

// GetUpcomingReviews récupère les prochaines N révisions après aujourd'hui
func GetUpcomingReviews(exercises []models.Exercise, limit int) []models.Exercise {
	var upcoming []models.Exercise
	now := time.Now()

	for _, ex := range exercises {
		if ex.NextReviewAt.IsZero() {
			continue
		}
		if ex.NextReviewAt.After(now) {
			upcoming = append(upcoming, ex)
		}
	}

	// Trie par date croissante
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].NextReviewAt.Before(upcoming[j].NextReviewAt)
	})

	if len(upcoming) > limit {
		upcoming = upcoming[:limit]
	}

	return upcoming
}
