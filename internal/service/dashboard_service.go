package service

import (
	"time"

	"maestro/internal/store"
)

// DashboardService gère les stats du dashboard
type DashboardService struct{}

// NewDashboardService crée une nouvelle instance
func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// DashboardStats regroupe toutes les statistiques
type DashboardStats struct {
	TotalExercises    int
	CompletedCount    int
	InProgressCount   int
	TodoCount         int
	OverdueCount      int
	CompletionRate    int // %
	AverageInterval   int // jours
	NextReviewDate    time.Time
	TotalSessionTime  time.Duration
	SessionCount      int
	AverageDifficulty float64
	TopDomain         string
	DomainBreakdown   map[string]int
}

// GetDashboardStats calcule toutes les stats
func (s *DashboardService) GetDashboardStats() DashboardStats {
	allExercises := store.GetAll()
	now := time.Now()

	stats := DashboardStats{
		TotalExercises:  len(allExercises),
		DomainBreakdown: make(map[string]int),
	}

	var totalDifficulty int
	var nextReviews []time.Time
	var totalInterval int
	var intervalCount int

	for _, ex := range allExercises {
		// Compte les états
		if ex.Done {
			stats.CompletedCount++
		} else if len(ex.CompletedSteps) > 0 {
			stats.InProgressCount++
		} else {
			stats.TodoCount++
		}

		// Overdue
		if !ex.NextReviewAt.IsZero() && ex.NextReviewAt.Before(now) && !ex.Done {
			stats.OverdueCount++
		}

		// Prochaines révisions
		if !ex.NextReviewAt.IsZero() && ex.NextReviewAt.After(now) {
			nextReviews = append(nextReviews, ex.NextReviewAt)
		}

		// Stats de difficulté
		totalDifficulty += ex.Difficulty

		// Stats d'intervalle SRS
		if ex.IntervalDays > 0 {
			totalInterval += ex.IntervalDays
			intervalCount++
		}

		// Breakdown par domaine
		stats.DomainBreakdown[ex.Domain]++
	}

	// Calcule les moyennes
	if stats.TotalExercises > 0 {
		stats.CompletionRate = (stats.CompletedCount * 100) / stats.TotalExercises
		stats.AverageDifficulty = float64(totalDifficulty) / float64(stats.TotalExercises)
	}

	if intervalCount > 0 {
		stats.AverageInterval = totalInterval / intervalCount
	}

	// Prochaine révision
	if len(nextReviews) > 0 {
		stats.NextReviewDate = nextReviews[0]
	}

	// Top domaine
	maxCount := 0
	for domain, count := range stats.DomainBreakdown {
		if count > maxCount {
			maxCount = count
			stats.TopDomain = domain
		}
	}

	return stats
}
