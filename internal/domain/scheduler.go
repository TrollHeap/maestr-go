package domain

import (
	"time"

	"maestro/internal/models"
)

// Scheduler gère l'algorithme de répétition espacée (SM-2)
type Scheduler struct {
	minEaseFactor float64
	maxInterval   int
}

// NewScheduler crée un nouveau Scheduler avec les paramètres par défaut
func NewScheduler() *Scheduler {
	return &Scheduler{
		minEaseFactor: 1.3,
		maxInterval:   365, // 1 an maximum
	}
}

// IsDueForReview vérifie si un exercice doit être révisé
func (s *Scheduler) IsDueForReview(ex *models.Exercise) bool {
	// Si jamais révisé, c'est une nouvelle carte
	if ex.LastReviewed == nil {
		return true
	}

	// Si déjà complété et pas de révision nécessaire
	if ex.Completed && ex.IntervalDays == 0 {
		return false
	}

	// Calculer la date de prochaine révision
	nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)

	// Due si la date de révision est passée
	return time.Now().After(nextReview) || time.Now().Equal(nextReview)
}

// ReviewExercise met à jour l'exercice après une révision (SM-2 algorithm)
// rating: 1 = Oublié, 2 = Difficile, 3 = Normal, 4 = Facile
func (s *Scheduler) ReviewExercise(ex *models.Exercise, rating int) {
	if rating < 1 || rating > 4 {
		return
	}

	var newInterval int
	var newEF float64

	// ✅ CORRECTION: Gérer le cas initial (IntervalDays = 0)
	isFirstReview := ex.IntervalDays == 0

	switch rating {
	case 4: // Facile
		if isFirstReview {
			newInterval = 1 // ✅ Premier intervalle = 1 jour
		} else {
			newInterval = int(float64(ex.IntervalDays) * ex.EaseFactor)
		}
		newEF = ex.EaseFactor + 0.1

	case 3: // Normal
		if isFirstReview {
			newInterval = 1
		} else {
			newInterval = int(float64(ex.IntervalDays) * ex.EaseFactor)
		}
		newEF = ex.EaseFactor

	case 2: // Difficile
		if isFirstReview {
			newInterval = 1
		} else {
			newInterval = int(float64(ex.IntervalDays) * 0.5)
			if newInterval < 1 {
				newInterval = 1 // ✅ Minimum 1 jour
			}
		}
		newEF = ex.EaseFactor - 0.2

	case 1: // Oublié
		newInterval = 1
		newEF = ex.EaseFactor - 0.5
	}

	// ✅ Clamp EaseFactor entre min et max
	if newEF < s.minEaseFactor {
		newEF = s.minEaseFactor
	}
	if newEF > 2.5 {
		newEF = 2.5
	}

	// ✅ Limiter l'intervalle maximum
	if newInterval > s.maxInterval {
		newInterval = s.maxInterval
	}

	// Mettre à jour l'exercice
	now := time.Now()
	ex.LastReviewed = &now
	ex.IntervalDays = newInterval
	ex.EaseFactor = newEF
	ex.Repetitions++
	ex.UpdatedAt = now
}

// GetNextReviewDate retourne la prochaine date de révision
func (s *Scheduler) GetNextReviewDate(ex *models.Exercise) *time.Time {
	if ex.LastReviewed == nil {
		return nil
	}

	nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
	return &nextReview
}

// ResetProgress réinitialise la progression d'un exercice
func (s *Scheduler) ResetProgress(ex *models.Exercise) {
	ex.LastReviewed = nil
	ex.IntervalDays = 0
	ex.EaseFactor = 2.5
	ex.Repetitions = 0
	ex.Completed = false
	ex.CompletedSteps = []int{}
}
