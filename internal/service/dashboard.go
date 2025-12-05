package service

import (
	"time"

	"maestro/internal/models"
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
	TotalExercises  int
	CompletedCount  int
	InProgressCount int
	TodoCount       int
	OverdueCount    int

	// ✅ NOUVEAUX CHAMPS
	TotalMastered  int           // Exercices maîtrisés (Done=true)
	StreakDays     int           // Jours consécutifs de révision
	WeeklyReviews  int           // Révisions cette semaine
	AvgSessionTime time.Duration // Durée moyenne session

	// Stats existantes
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

	// ✅ Variables pour nouveaux calculs
	weekAgo := now.AddDate(0, 0, -7)
	var weeklyReviewCount int

	for _, ex := range allExercises {
		// Compte les états
		if ex.Done {
			stats.CompletedCount++
			stats.TotalMastered++ // ✅ Maîtrisés = Done
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

		// ✅ FIX: Révisions cette semaine (check nil + type)
		if ex.LastReviewed != nil && !ex.LastReviewed.IsZero() && ex.LastReviewed.After(weekAgo) {
			weeklyReviewCount++
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

	// ✅ Weekly reviews
	stats.WeeklyReviews = weeklyReviewCount

	// ✅ Streak calculation
	stats.StreakDays = calculateStreak(allExercises)

	// ✅ Session stats (si tu as une table sessions)
	stats.SessionCount, stats.TotalSessionTime = getSessionStats()
	if stats.SessionCount > 0 {
		stats.AvgSessionTime = stats.TotalSessionTime / time.Duration(stats.SessionCount)
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

// ============================================
// HELPER FUNCTIONS
// ============================================

// calculateStreak calcule la séquence de jours consécutifs
func calculateStreak(exercises []models.Exercise) int {
	if len(exercises) == 0 {
		return 0
	}

	// Map des dates où au moins 1 exercice a été revu
	reviewDates := make(map[string]bool)

	for _, ex := range exercises {
		// ✅ FIX: Check nil avant IsZero
		if ex.LastReviewed != nil && !ex.LastReviewed.IsZero() {
			dateKey := ex.LastReviewed.Format("2006-01-02")
			reviewDates[dateKey] = true
		}
	}

	// Compte jours consécutifs depuis aujourd'hui
	streak := 0
	currentDate := time.Now()

	for {
		dateKey := currentDate.Format("2006-01-02")
		if !reviewDates[dateKey] {
			break
		}
		streak++
		currentDate = currentDate.AddDate(0, 0, -1)

		// Limite à 365 jours pour éviter boucle infinie
		if streak > 365 {
			break
		}
	}

	return streak
}

// getSessionStats récupère stats sessions (si tu as une table sessions)
func getSessionStats() (int, time.Duration) {
	// ✅ Version simple sans table sessions
	return 0, 0
}
