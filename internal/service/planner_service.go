package service

import (
	"log"
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

// normalizeDate supprime l'heure et force UTC
func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// GetReviewsForDate récupère les exercices à réviser pour une date
func (s *PlannerService) GetReviewsForDate(date time.Time) []models.Exercise {
	allExercises := store.GetAll()
	var reviews []models.Exercise

	// Normalise la date cible (ignore l'heure, force UTC)
	targetDate := normalizeDate(date)

	log.Printf("🔍 [PlannerService] GetReviewsForDate pour %s (normalisé: %s)",
		date.Format("2006-01-02"), targetDate.Format("2006-01-02"))
	log.Printf("📚 [PlannerService] %d exercices au total dans le store", len(allExercises))

	for _, ex := range allExercises {
		// Skip si pas de date de révision
		if ex.NextReviewAt.IsZero() {
			continue
		}

		// Normalise la date de révision (ignore l'heure, force UTC)
		reviewDate := normalizeDate(ex.NextReviewAt)

		// Debug pour chaque exercice
		log.Printf("   ├─ Exo #%d '%s': NextReviewAt=%s (normalisé=%s) | Match=%v",
			ex.ID, ex.Title,
			ex.NextReviewAt.Format("2006-01-02 15:04"),
			reviewDate.Format("2006-01-02"),
			reviewDate.Equal(targetDate))

		// Compare les dates normalisées
		if reviewDate.Equal(targetDate) {
			reviews = append(reviews, ex)
			log.Printf("   └─ ✅ Match ! Ajouté à la liste")
		}
	}

	log.Printf("✅ [PlannerService] %d révision(s) trouvée(s) pour %s",
		len(reviews), targetDate.Format("2006-01-02"))

	// Trie par difficulté (urgent d'abord)
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].Difficulty > reviews[j].Difficulty
	})

	return reviews
}

// GetWeekSchedule récupère le planning de la semaine
func (s *PlannerService) GetWeekSchedule(startDate time.Time) []DaySchedule {
	log.Printf("📅 [PlannerService] GetWeekSchedule pour semaine du %s",
		startDate.Format("2006-01-02"))

	schedule := make([]DaySchedule, 7)

	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i)
		exercises := s.GetReviewsForDate(date)

		schedule[i] = DaySchedule{
			Date:      date,
			Exercises: exercises,
			Count:     len(exercises),
		}

		log.Printf("   Jour %d (%s): %d révision(s)",
			i+1, date.Format("Mon 02 Jan"), len(exercises))
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
