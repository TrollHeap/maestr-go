package domain

import (
	"sort"

	"maestro/internal/models"
)

// Recommender recommande les prochains exercices à faire
type Recommender struct {
	scheduler *Scheduler
}

// NewRecommender crée un nouveau Recommender
func NewRecommender(scheduler *Scheduler) *Recommender {
	return &Recommender{
		scheduler: scheduler,
	}
}

// exerciseWithPriority associe un exercice avec sa priorité
type exerciseWithPriority struct {
	exercise models.Exercise
	priority int // 1 = due review, 2 = new, 3 = completed
}

// ✅ OPTIMISÉ: Une seule passe O(n log n)
func (r *Recommender) GetNextExercises(exercises []models.Exercise, limit int) []models.Exercise {
	// Une seule passe pour catégoriser
	categorized := make([]exerciseWithPriority, 0, len(exercises))

	for _, ex := range exercises {
		var priority int

		if r.scheduler.IsDueForReview(&ex) {
			priority = 1 // Highest priority - révisions dues
		} else if !ex.Completed {
			priority = 2 // Nouveaux exercices
		} else {
			priority = 3 // Complétés - skip
		}

		categorized = append(categorized, exerciseWithPriority{
			exercise: ex,
			priority: priority,
		})
	}

	// Trier par priorité (puis par difficulté comme tiebreaker)
	sort.Slice(categorized, func(i, j int) bool {
		if categorized[i].priority != categorized[j].priority {
			return categorized[i].priority < categorized[j].priority
		}
		// Même priorité: trier par difficulté (plus facile d'abord)
		return categorized[i].exercise.Difficulty < categorized[j].exercise.Difficulty
	})

	// Extraire les N premiers (limit)
	result := make([]models.Exercise, 0, limit)
	for i := 0; i < len(categorized) && len(result) < limit; i++ {
		if categorized[i].priority < 3 { // Skip complétés
			result = append(result, categorized[i].exercise)
		}
	}

	return result
}

// CalculateStats calcule les statistiques globales (✅ float64 mastery)
func (r *Recommender) CalculateStats(exercises []models.Exercise) models.Stats {
	stats := models.Stats{
		Total:      len(exercises),
		Completed:  0,
		InProgress: 0,
		DueReview:  0,
		ByDomain:   make(map[string]models.DomainStat), // ✅ INITIALISÉ
	}

	for _, ex := range exercises {
		// Skip deleted
		if ex.Deleted {
			continue
		}

		// Statistiques de domaine
		domainStat := stats.ByDomain[ex.Domain]
		domainStat.Total++

		if ex.Completed {
			stats.Completed++
			domainStat.Completed++

			// Calculer mastery (0-100)
			mastery := ((ex.EaseFactor - 1.3) / (2.5 - 1.3)) * 100
			if mastery < 0 {
				mastery = 0
			}
			if mastery > 100 {
				mastery = 100
			}

			// Moyenne pondérée
			if domainStat.Mastery == 0 {
				domainStat.Mastery = mastery
			} else {
				domainStat.Mastery = (domainStat.Mastery + mastery) / 2
			}
		} else {
			stats.InProgress++
		}

		if r.scheduler.IsDueForReview(&ex) {
			stats.DueReview++
		}

		// ✅ IMPORTANT: Sauvegarder les modifications dans la map
		stats.ByDomain[ex.Domain] = domainStat
	}

	return stats
}
