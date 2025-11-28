package service

import (
	"sort"
	"time"

	"maestro/internal/models"
	"maestro/internal/store"
)

// PlannerService gère la logique du calendrier
type PlannerService struct{}

// NewPlannerService crée une nouvelle instance
func NewPlannerService() *PlannerService {
	return &PlannerService{}
}

// DaySchedule représente les exercices d'un jour
type DaySchedule struct {
	Date      time.Time
	Exercises []models.Exercise
	Count     int
}

// GetReviewsForDate récupère les exercices à réviser pour une date
func (s *PlannerService) GetReviewsForDate(date time.Time) []models.Exercise {
	allExercises := store.GetAll()
	var reviews []models.Exercise

	// Normalise la date (ignore l'heure)
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	for _, ex := range allExercises {
		if ex.NextReviewAt.IsZero() {
			continue
		}

		// Normalise la date de révision
		reviewDate := time.Date(
			ex.NextReviewAt.Year(),
			ex.NextReviewAt.Month(),
			ex.NextReviewAt.Day(),
			0, 0, 0, 0,
			ex.NextReviewAt.Location(),
		)

		// Compare les dates
		if reviewDate.Equal(targetDate) {
			reviews = append(reviews, ex)
		}
	}

	// Trie par difficulté (urgent d'abord)
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].Difficulty > reviews[j].Difficulty
	})

	return reviews
}

// GetWeekSchedule récupère le planning de la semaine
func (s *PlannerService) GetWeekSchedule(startDate time.Time) []DaySchedule {
	schedule := make([]DaySchedule, 7)

	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i)
		exercises := s.GetReviewsForDate(date)

		schedule[i] = DaySchedule{
			Date:      date,
			Exercises: exercises,
			Count:     len(exercises),
		}
	}

	return schedule
}

// GetMonthSchedule récupère le planning du mois
func (s *PlannerService) GetMonthSchedule(year int, month time.Month) map[int]int {
	counts := make(map[int]int)
	allExercises := store.GetAll()

	for _, ex := range allExercises {
		if ex.NextReviewAt.IsZero() {
			continue
		}

		// Si la révision est dans ce mois
		if ex.NextReviewAt.Year() == year && ex.NextReviewAt.Month() == month {
			day := ex.NextReviewAt.Day()
			counts[day]++
		}
	}

	return counts
}

// GetUpcomingReviews récupère les N prochaines révisions
func (s *PlannerService) GetUpcomingReviews(limit int) []models.Exercise {
	allExercises := store.GetAll()
	var upcoming []models.Exercise

	now := time.Now()

	for _, ex := range allExercises {
		if !ex.NextReviewAt.IsZero() && ex.NextReviewAt.After(now) {
			upcoming = append(upcoming, ex)
		}
	}

	// Trie par date de révision
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].NextReviewAt.Before(upcoming[j].NextReviewAt)
	})

	// Limite le nombre
	if len(upcoming) > limit {
		upcoming = upcoming[:limit]
	}

	return upcoming
}

// GetOverdueReviews récupère les révisions en retard
func (s *PlannerService) GetOverdueReviews() []models.Exercise {
	allExercises := store.GetAll()
	var overdue []models.Exercise

	now := time.Now()

	for _, ex := range allExercises {
		if !ex.NextReviewAt.IsZero() && ex.NextReviewAt.Before(now) && !ex.Done {
			overdue = append(overdue, ex)
		}
	}

	// Trie par ancienneté (plus vieux d'abord)
	sort.Slice(overdue, func(i, j int) bool {
		return overdue[i].NextReviewAt.Before(overdue[j].NextReviewAt)
	})

	return overdue
}
