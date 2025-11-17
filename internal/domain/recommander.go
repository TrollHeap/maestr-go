package domain

import (
	"maestro/internal/models"
)

// Recommender suggère les exercices à faire
type Recommender struct {
	scheduler *Scheduler
}

// NewRecommender crée une nouvelle instance Recommender
func NewRecommender(scheduler *Scheduler) *Recommender {
	return &Recommender{scheduler: scheduler}
}

// GetNextExercises retourne les exercices à faire ensuite (priorisés)
// Priorité 1: Exercices dus pour révision
// Priorité 2: Nouveaux exercices
func (r *Recommender) GetNextExercises(exercises []models.Exercise, limit int) []models.Exercise {
	var recommended []models.Exercise

	// Priorité 1: Exercices dus pour révision
	for _, ex := range exercises {
		if r.scheduler.IsDueForReview(&ex) {
			recommended = append(recommended, ex)
		}
	}

	// Si pas assez, ajouter des nouveaux exercices
	if len(recommended) < limit {
		for _, ex := range exercises {
			if !ex.Completed && !r.scheduler.IsDueForReview(&ex) {
				recommended = append(recommended, ex)
			}
			if len(recommended) >= limit {
				break
			}
		}
	}

	// Limiter au nombre demandé
	if len(recommended) > limit {
		recommended = recommended[:limit]
	}

	return recommended
}

// CalculateStats calcule les statistiques globales
func (r *Recommender) CalculateStats(exercises []models.Exercise) models.Stats {
	stats := models.Stats{
		TotalCompleted: 0,
		TotalReviews:   0,
		DomainStats:    make(map[string]models.DomainStat),
	}

	// Compter par domaine
	for _, ex := range exercises {
		if ex.Completed {
			stats.TotalCompleted++
			stats.TotalReviews += ex.Repetitions
		}

		// Créer entry si n'existe pas
		if _, ok := stats.DomainStats[ex.Domain]; !ok {
			stats.DomainStats[ex.Domain] = models.DomainStat{
				Completed: 0,
				Total:     0,
				Mastery:   0,
			}
		}

		// Mettre à jour domaine
		domainStat := stats.DomainStats[ex.Domain]
		domainStat.Total++
		if ex.Completed {
			domainStat.Completed++
		}

		// Calculer mastery (0-100) basé sur EF
		// EF = 2.5 (max) → 100%
		// EF = 1.3 (min) → 50%
		if ex.Completed {
			mastery := int((ex.EaseFactor - 1.3) / (2.5 - 1.3) * 100)
			if mastery > domainStat.Mastery {
				domainStat.Mastery = mastery
			}
		}

		stats.DomainStats[ex.Domain] = domainStat
	}

	return stats
}
