package domain

import (
	"time"

	"maestro/internal/models"
)

// Scheduler gère l'algorithme SM-2 (Spaced Repetition)
type Scheduler struct {
	initialEaseFactor float64
	minEaseFactor     float64
}

// NewScheduler crée une nouvelle instance Scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		initialEaseFactor: 2.5,
		minEaseFactor:     1.3,
	}
}

// ReviewExercise applique l'algorithme SM-2
// rating doit être entre 1 et 4:
//   1 = Complètement oublié (reset)
//   2 = Très difficile (EF - 0.2)
//   3 = Normal (EF constant)
//   4 = Facile (EF + 0.1)
func (s *Scheduler) ReviewExercise(ex *models.Exercise, rating int) {
	if rating < 1 || rating > 4 {
		return  // Invalid rating
	}

	var newInterval int
	var newEF float64

	switch rating {
	case 4:  // Facile
		if ex.IntervalDays == 0 {
			newInterval = 1
		} else {
			newInterval = int(float64(ex.IntervalDays) * ex.EaseFactor)
		}
		newEF = ex.EaseFactor + 0.1

	case 3:  // Normal
		if ex.IntervalDays == 0 {
			newInterval = 1
		} else {
			newInterval = int(float64(ex.IntervalDays) * ex.EaseFactor)
		}
		newEF = ex.EaseFactor

	case 2:  // Difficile
		newInterval = int(float64(ex.IntervalDays) * 0.5)
		if newInterval < 1 {
			newInterval = 1
		}
		newEF = ex.EaseFactor - 0.2

	case 1:  // Oublié (reset)
		newInterval = 1
		newEF = ex.EaseFactor - 0.5
	}

	// Clamp EF between min and max
	if newEF < s.minEaseFactor {
		newEF = s.minEaseFactor
	}
	if newEF > 2.5 {
		newEF = 2.5
	}

	// Update exercise - ✅ IMPORTANT: Mettre à jour ex directement!
	now := time.Now()
	ex.LastReviewed = &now
	ex.IntervalDays = newInterval
	ex.EaseFactor = newEF  // ← ✅ CETTE LIGNE EST CRUCIALE
	ex.Repetitions++
	ex.Completed = true
	ex.UpdatedAt = now
}

// IsDueForReview vérifie si l'exercice doit être révisé
func (s *Scheduler) IsDueForReview(ex *models.Exercise) bool {
	if ex.LastReviewed == nil {
		return false
	}
	nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
	return time.Now().After(nextReview)
}

// GetDaysUntilReview retourne le nombre de jours avant la prochaine révision
func (s *Scheduler) GetDaysUntilReview(ex *models.Exercise) int {
	if ex.LastReviewed == nil {
		return 0  // Nouveau, révision immédiate
	}

	nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
	days := int(time.Until(nextReview).Hours() / 24)

	if days < 0 {
		return 0  // Due
	}
	return days
}
