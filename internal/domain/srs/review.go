package srs

import (
	"math"
	"time"

	"maestro/internal/models"
)

// ReviewQuality : Qualité de révision (ALIGNED avec DB CHECK 0-3)
type ReviewQuality int

const (
	Again ReviewQuality = 0 // Oublié (reset complet)
	Hard  ReviewQuality = 1 // Difficile (1 jour)
	Good  ReviewQuality = 2 // Bien (intervalle × 2)
	Easy  ReviewQuality = 3 // Facile (intervalle × 3)
)

// CalculateNextReview : Algorithme SM-2 adapté
func CalculateNextReview(
	quality ReviewQuality,
	currentInterval int,
	currentEase float64,
	currentReps int,
) models.ReviewResult {
	now := time.Now()
	result := models.ReviewResult{
		EaseFactor:  currentEase,
		Repetitions: currentReps,
	}

	switch quality {
	case Again: // 0 - Oublié
		result.IntervalDays = 0
		result.Repetitions = 0
		result.EaseFactor = math.Max(1.3, currentEase-0.3)
		result.NextReview = now.Add(10 * time.Minute)

	case Hard: // 1 - Difficile
		result.IntervalDays = 1
		result.Repetitions = currentReps + 1
		result.EaseFactor = math.Max(1.3, currentEase-0.2)
		result.NextReview = now.AddDate(0, 0, 1)

	case Good: // 2 - Bien (était 3)
		if currentInterval == 0 {
			result.IntervalDays = 1
		} else {
			result.IntervalDays = currentInterval * 2
		}
		result.Repetitions = currentReps + 1
		result.EaseFactor = currentEase
		result.NextReview = now.AddDate(0, 0, result.IntervalDays)

	case Easy: // 3 - Facile (était 5)
		if currentInterval == 0 {
			result.IntervalDays = 4
		} else {
			result.IntervalDays = currentInterval * 3
		}
		result.Repetitions = currentReps + 1
		result.EaseFactor = math.Min(2.5, currentEase+0.1)
		result.NextReview = now.AddDate(0, 0, result.IntervalDays)
	}

	return result
}

// IsDueForReview vérifie si une révision est due
func IsDueForReview(nextReview time.Time) bool {
	return time.Now().After(nextReview)
}
