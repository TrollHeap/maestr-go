package domain

import (
	"time"

	"maestro/internal/models"
)

// Scheduler gère l'algorithme SM-2 (Spaced Repetition) AMÉLIORÉ pour ADHD
type Scheduler struct {
	initialEaseFactor float64
	minEaseFactor     float64
	maxEaseFactor     float64
}

// NewScheduler crée une nouvelle instance Scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		initialEaseFactor: 2.5,
		minEaseFactor:     1.3,
		maxEaseFactor:     3.0, // ← AMÉLIORÉ: Max augmentée à 3.0 (vs 2.5)
	}
}

// ReviewExercise applique l'algorithme SM-2 ADHD-friendly
// rating doit être entre 1 et 4:
//
//	1 = Complètement oublié (reset + review TODAY)
//	2 = Très difficile (EF - 0.1, revue dans 1 jour)
//	3 = Normal (EF constant, augmente interval normal)
//	4 = Facile (EF + 0.15 bonus ADHD, interval augmente plus)
func (s *Scheduler) ReviewExercise(ex *models.Exercise, rating int) {
	if rating < 1 || rating > 4 {
		return // Invalid rating
	}

	var newInterval int
	var newEF float64

	switch rating {
	case 4: // Facile ✅
		// AMÉLIORÉ: Bonus ADHD pour encourager
		if ex.IntervalDays == 0 {
			newInterval = 1
		} else {
			// Apply ease factor × 1.1 bonus (ADHD encouragement)
			newInterval = int(float64(ex.IntervalDays) * ex.EaseFactor * 1.1)
		}
		newEF = ex.EaseFactor + 0.15 // ← AMÉLIORÉ: +0.15 vs +0.1 (plus de bonus!)

	case 3: // Normal 📌
		if ex.IntervalDays == 0 {
			newInterval = 1
		} else {
			newInterval = int(float64(ex.IntervalDays) * ex.EaseFactor)
		}
		newEF = ex.EaseFactor

	case 2: // Difficile ⚠️
		// AMÉLIORÉ: Moins harsh pour ADHD
		// Au lieu de 0.5x, c'est 1 jour minimum
		newInterval = 1             // ← AMÉLIORÉ: Toujours au moins 1 jour (pas 0.5x harsh)
		newEF = ex.EaseFactor - 0.1 // ← AMÉLIORÉ: -0.1 vs -0.2 (moins harsh)

	case 1: // Oublié ❌
		// AMÉLIORÉ: Review TODAY + reset, mais pas trop harsh
		newInterval = 0             // ← CRUCIAL: 0 = REVIEW TODAY! (pas 1)
		newEF = ex.EaseFactor - 0.2 // ← Pénalité mais supportable
		ex.Completed = false        // ← CRUCIAL: Mark as incomplete! (pour revoir)
	}

	// Clamp EF between min and max
	if newEF < s.minEaseFactor {
		newEF = s.minEaseFactor
	}
	if newEF > s.maxEaseFactor {
		newEF = s.maxEaseFactor // ← AMÉLIORÉ: Use maxEaseFactor (3.0)
	}

	// Update exercise
	now := time.Now()
	ex.LastReviewed = &now
	ex.IntervalDays = newInterval
	ex.EaseFactor = newEF
	ex.Repetitions++
	ex.UpdatedAt = now

	// ✅ IMPORTANT: Si pas oublié (rating 2-4), marquer comme complété
	if rating != 1 {
		ex.Completed = true
	}
	// Si oublié (rating 1): déjà marqué incomplet au-dessus
}

// IsDueForReview vérifie si l'exercice doit être révisé
func (s *Scheduler) IsDueForReview(ex *models.Exercise) bool {
	// ✅ AMÉLIORÉ: Vérifier les 3 conditions:
	// 1. Pas jamais reviewé (LastReviewed == nil)
	// 2. Interval = 0 (oublié, review today)
	// 3. NextReview date passed

	if ex.LastReviewed == nil {
		// ✅ NEW: Nouveau exercice, peut être reviewé
		return true
	}

	// Si interval est 0, c'est urgent (oublié)
	if ex.IntervalDays == 0 {
		return true
	}

	// Sinon, check si dépassé la date
	nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
	return time.Now().After(nextReview)
}

// GetDaysUntilReview retourne le nombre de jours avant la prochaine révision
func (s *Scheduler) GetDaysUntilReview(ex *models.Exercise) int {
	if ex.LastReviewed == nil {
		return 0 // Nouveau, révision immédiate
	}

	// ✅ AMÉLIORÉ: Si interval est 0, c'est DUE TODAY
	if ex.IntervalDays == 0 {
		return 0 // Due now!
	}

	nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
	daysUntil := time.Until(nextReview)

	// Convert to days
	days := int(daysUntil.Hours() / 24)

	if days < 0 {
		return 0 // Overdue
	}
	return days
}

// GetNextReviewDate retourne la date exacte de la prochaine révision (NOUVEAU)
func (s *Scheduler) GetNextReviewDate(ex *models.Exercise) *time.Time {
	if ex.LastReviewed == nil {
		// Nouveau exercice, revoir immédiatement
		now := time.Now()
		return &now
	}

	nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
	return &nextReview
}

// GetReadableNextReview retourne un texte lisible pour l'interface (NOUVEAU)
func (s *Scheduler) GetReadableNextReview(ex *models.Exercise) string {
	if ex.LastReviewed == nil {
		return "Nouveau"
	}

	days := s.GetDaysUntilReview(ex)

	switch {
	case days == 0:
		return "Aujourd'hui"
	case days == 1:
		return "Demain"
	case days < 7:
		return "Cette semaine"
	case days < 30:
		return "Ce mois"
	default:
		return "Plus tard"
	}
}

// ============= HELPER FUNCTIONS =============

// CalculateMastery retourne le pourcentage de maîtrise (0-100) (NOUVEAU)
// Basé sur EF: 1.3 = 0%, 3.0 = 100%
func CalculateMastery(ef float64) int {
	if ef < 1.3 {
		return 0
	}
	if ef > 3.0 {
		return 100
	}

	// Linear scale from 1.3 → 3.0 = 0% → 100%
	percentage := ((ef - 1.3) / (3.0 - 1.3)) * 100
	return int(percentage)
}

// IsCompleted retourne si exercise est vraiment complété (NOUVEAU)
func IsCompleted(ex *models.Exercise) bool {
	return ex.Completed && ex.LastReviewed != nil
}

// IsReadyForReview retourne si exercise peut être reviewé (NOUVEAU)
func (s *Scheduler) IsReadyForReview(ex *models.Exercise) bool {
	// ✅ Ne revoir que si:
	// 1. Déjà reviewé ET
	// 2. Due for review

	if ex.LastReviewed == nil {
		// Nouveau = peut commencer
		return true
	}

	if !ex.Completed {
		// Si marked as incomplete (oublié), toujours reviewable
		return true
	}

	// Sinon, check due date
	return s.IsDueForReview(ex)
}
