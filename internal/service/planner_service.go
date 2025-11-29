package service

import (
	"log"
	"time"

	"maestro/internal/domain/planner"
	"maestro/internal/models"
	"maestro/internal/store"
)

type PlannerService struct{}

func NewPlannerService() *PlannerService {
	return &PlannerService{}
}

func (s *PlannerService) GetReviewsForDate(date time.Time) []models.Exercise {
	allExercises := store.GetAll()
	reviews := planner.GetReviewsForDate(allExercises, date)

	log.Printf("🔍 [PlannerService] %d révision(s) pour %s", len(reviews), date.Format("2006-01-02"))
	return reviews
}

func (s *PlannerService) GetOverdueReviews() []models.Exercise {
	allExercises := store.GetAll()
	return planner.GetOverdueReviews(allExercises)
}

func (s *PlannerService) GetUpcomingReviews(limit int) []models.Exercise {
	allExercises := store.GetAll()
	return planner.GetUpcomingReviews(allExercises, limit)
}

func (s *PlannerService) GetWeekSchedule(startDate time.Time) []models.DaySchedule {
	schedule := make([]models.DaySchedule, 7)
	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i)
		schedule[i] = models.DaySchedule{
			Date:      date,
			Exercises: s.GetReviewsForDate(date),
			Count:     0, // Optionnel, ou len des exercices
		}
		schedule[i].Count = len(schedule[i].Exercises)
		log.Printf(
			"   Jour %d (%s): %d révision(s)",
			i+1,
			date.Format("Mon 02 Jan"),
			schedule[i].Count,
		)
	}
	return schedule
}

func (s *PlannerService) GetMonthSchedule(year int, month time.Month) map[int]int {
	counts := make(map[int]int)
	allExercises := store.GetAll()
	for _, ex := range allExercises {
		if ex.NextReviewAt.IsZero() {
			continue
		}
		if ex.NextReviewAt.Year() == year && ex.NextReviewAt.Month() == month {
			day := ex.NextReviewAt.Day()
			counts[day]++
		}
	}
	return counts
}
