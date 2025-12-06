// internal/service/dashboard.go
package service

import (
	"time"

	"maestro/internal/models"
	"maestro/internal/store"
	"maestro/internal/views/logic"
)

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// GetDashboardStats - Stats principales
func (s *DashboardService) GetDashboardStats() models.DashboardStats {
	allExercises := store.GetAll()
	now := time.Now()

	stats := models.DashboardStats{
		TotalExercises:  len(allExercises),
		DomainBreakdown: make(map[string]int),
	}

	var totalDifficulty int
	var nextReviews []time.Time
	var totalInterval int
	var intervalCount int
	var totalEaseFactor float64
	var easeCount int
	var successfulReviews int
	var totalReviews int
	weekAgo := now.AddDate(0, 0, -7)
	var weeklyReviewCount int

	for _, ex := range allExercises {
		// Compte les états
		if ex.Done {
			stats.CompletedCount++
			stats.TotalMastered++
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

		// Weekly reviews
		if ex.LastReviewed != nil && !ex.LastReviewed.IsZero() && ex.LastReviewed.After(weekAgo) {
			weeklyReviewCount++
		}

		totalDifficulty += ex.Difficulty

		// Intervalle SRS
		if ex.IntervalDays > 0 {
			totalInterval += ex.IntervalDays
			intervalCount++
		}

		// ✅ EaseFactor moyen
		if ex.EaseFactor > 0 {
			totalEaseFactor += ex.EaseFactor
			easeCount++
		}

		// ✅ Retention rate (reviews réussies)
		if ex.Repetitions > 0 {
			totalReviews += ex.Repetitions
			// Approximation: si EF >= 2.5, c'est "réussi"
			if ex.EaseFactor >= 2.5 {
				successfulReviews += ex.Repetitions
			}
		}

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

	if easeCount > 0 {
		stats.AverageEaseFactor = totalEaseFactor / float64(easeCount)
	}

	if totalReviews > 0 {
		stats.RetentionRate = (successfulReviews * 100) / totalReviews
	}

	stats.WeeklyReviews = weeklyReviewCount
	stats.StreakDays = calculateStreak(allExercises)
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

// ✅ GetHeatmapData - Données pour le heatmap GitHub-style
func (s *DashboardService) GetHeatmapData(weeks int) []logic.HeatmapDay {
	allExercises := store.GetAll()

	// Compte les reviews par date
	reviewCounts := make(map[string]int)

	for _, ex := range allExercises {
		if ex.LastReviewed != nil && !ex.LastReviewed.IsZero() {
			dateKey := ex.LastReviewed.Format("2006-01-02")
			reviewCounts[dateKey]++
		}
	}

	// Génère les jours via logic helper
	return logic.GenerateHeatmapDays(reviewCounts, weeks)
}

// ✅ GetWeakExercises - Exercices avec EaseFactor faible
func (s *DashboardService) GetWeakExercises(limit int) []models.Exercise {
	allExercises := store.GetAll()
	var weak []models.Exercise

	// Seuil: EaseFactor < 2.3 ET pas encore maîtrisé
	for _, ex := range allExercises {
		if !ex.Done && ex.EaseFactor < 2.3 && ex.EaseFactor > 0 {
			weak = append(weak, ex)
		}
	}

	// Trie par EaseFactor croissant (plus faibles en premier)
	for i := 0; i < len(weak)-1; i++ {
		for j := i + 1; j < len(weak); j++ {
			if weak[i].EaseFactor > weak[j].EaseFactor {
				weak[i], weak[j] = weak[j], weak[i]
			}
		}
	}

	// Limite à N résultats
	if len(weak) > limit {
		weak = weak[:limit]
	}

	return weak
}

// ✅ GetFailurePatterns - Exercices avec échecs répétés
func (s *DashboardService) GetFailurePatterns(limit int) []models.FailurePattern {
	allExercises := store.GetAll()
	var patterns []models.FailurePattern

	for _, ex := range allExercises {
		// Critère: EaseFactor bas + répétitions élevées = échecs répétés
		if !ex.Done && ex.EaseFactor < 2.3 && ex.Repetitions >= 3 {
			failCount := ex.Repetitions // Approximation
			patterns = append(patterns, models.FailurePattern{
				ExerciseID: ex.ID,
				Title:      ex.Title,
				Domain:     ex.Domain,
				FailCount:  failCount,
				EaseFactor: ex.EaseFactor,
			})
		}
	}

	// Trie par FailCount décroissant
	for i := 0; i < len(patterns)-1; i++ {
		for j := i + 1; j < len(patterns); j++ {
			if patterns[i].FailCount < patterns[j].FailCount {
				patterns[i], patterns[j] = patterns[j], patterns[i]
			}
		}
	}

	if len(patterns) > limit {
		patterns = patterns[:limit]
	}

	return patterns
}

// ✅ GetRepetitionStats - Exercices les plus révisés
func (s *DashboardService) GetRepetitionStats(limit int) []models.RepetitionStat {
	allExercises := store.GetAll()
	var stats []models.RepetitionStat

	for _, ex := range allExercises {
		if ex.Repetitions > 0 {
			stats = append(stats, models.RepetitionStat{
				ExerciseID:  ex.ID,
				Title:       ex.Title,
				Domain:      ex.Domain,
				ReviewCount: ex.Repetitions,
			})
		}
	}

	// Trie par ReviewCount décroissant
	for i := 0; i < len(stats)-1; i++ {
		for j := i + 1; j < len(stats); j++ {
			if stats[i].ReviewCount < stats[j].ReviewCount {
				stats[i], stats[j] = stats[j], stats[i]
			}
		}
	}

	if len(stats) > limit {
		stats = stats[:limit]
	}

	return stats
}

// ✅ GetDomainStrengths - Analyse force par domaine
func (s *DashboardService) GetDomainStrengths() []models.DomainStrength {
	allExercises := store.GetAll()
	domainMap := make(map[string]*models.DomainStrength)

	for _, ex := range allExercises {
		if domainMap[ex.Domain] == nil {
			domainMap[ex.Domain] = &models.DomainStrength{
				Name: ex.Domain,
			}
		}

		ds := domainMap[ex.Domain]
		ds.TotalCount++

		if ex.Done {
			ds.MasteredCount++
		}

		if ex.EaseFactor > 0 {
			ds.AvgEaseFactor += ex.EaseFactor
		}
	}

	var strengths []models.DomainStrength
	for _, ds := range domainMap {
		if ds.TotalCount > 0 {
			ds.AvgEaseFactor /= float64(ds.TotalCount)
			ds.StrengthPercent = (ds.MasteredCount * 100) / ds.TotalCount
		}
		strengths = append(strengths, *ds)
	}

	// Trie par StrengthPercent décroissant
	for i := 0; i < len(strengths)-1; i++ {
		for j := i + 1; j < len(strengths); j++ {
			if strengths[i].StrengthPercent < strengths[j].StrengthPercent {
				strengths[i], strengths[j] = strengths[j], strengths[i]
			}
		}
	}

	return strengths
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func calculateStreak(exercises []models.Exercise) int {
	if len(exercises) == 0 {
		return 0
	}

	reviewDates := make(map[string]bool)

	for _, ex := range exercises {
		if ex.LastReviewed != nil && !ex.LastReviewed.IsZero() {
			dateKey := ex.LastReviewed.Format("2006-01-02")
			reviewDates[dateKey] = true
		}
	}

	streak := 0
	currentDate := time.Now()

	for {
		dateKey := currentDate.Format("2006-01-02")
		if !reviewDates[dateKey] {
			break
		}
		streak++
		currentDate = currentDate.AddDate(0, 0, -1)

		if streak > 365 {
			break
		}
	}

	return streak
}

func getSessionStats() (int, time.Duration) {
	// Version simple sans table sessions
	return 0, 0
}
