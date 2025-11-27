package srs

import (
	"math"
	"time"
)

// ReviewQuality représente la qualité de la révision
type ReviewQuality int

const (
	Again ReviewQuality = 0 // Oublié
	Hard  ReviewQuality = 1 // Difficile
	Good  ReviewQuality = 3 // Bien
	Easy  ReviewQuality = 5 // Facile
)

// ReviewResult contient le résultat d'un calcul SRS
type ReviewResult struct {
	IntervalDays int
	EaseFactor   float64
	Repetitions  int
	NextReview   time.Time
}

// CalculateNextReview applique l'algorithme SM-2
// Fonction PURE : pas d'effet de bord, testable facilement
func CalculateNextReview(
	quality ReviewQuality,
	currentInterval int,
	currentEase float64,
	currentReps int,
) ReviewResult {
	now := time.Now()
	result := ReviewResult{
		EaseFactor:  currentEase,
		Repetitions: currentReps,
	}

	switch quality {
	case Again: // Reset complet
		result.IntervalDays = 0
		result.Repetitions = 0
		result.EaseFactor = math.Max(1.3, currentEase-0.3)
		result.NextReview = now.Add(10 * time.Minute)

	case Hard:
		result.IntervalDays = 1
		result.Repetitions = currentReps + 1
		result.EaseFactor = math.Max(1.3, currentEase-0.2)
		result.NextReview = now.AddDate(0, 0, 1)

	case Good:
		if currentInterval == 0 {
			result.IntervalDays = 1
		} else {
			result.IntervalDays = currentInterval * 2
		}
		result.Repetitions = currentReps + 1
		result.EaseFactor = currentEase
		result.NextReview = now.AddDate(0, 0, result.IntervalDays)

	case Easy:
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
