package session

import (
	"sort"
	"time"

	"maestro/internal/models"
)

// ============================================
// RÈGLES DE PRIORITÉ (Domain Logic)
// ============================================

// SortByPriority : Trie exercices par priorité (en retard → aujourd'hui → nouveaux)
func SortByPriority(exercises []models.Exercise) []models.Exercise {
	// now := time.Now()
	// today := now.Truncate(24 * time.Hour)
	// tomorrow := today.Add(24 * time.Hour)

	sort.Slice(exercises, func(i, j int) bool {
		a, b := exercises[i], exercises[j]

		// Priorité 1 : En retard (urgent)
		aOverdue := IsOverdue(a)
		bOverdue := IsOverdue(b)
		if aOverdue != bOverdue {
			return aOverdue
		}

		// Priorité 2 : À réviser aujourd'hui
		aToday := IsDueToday(a)
		bToday := IsDueToday(b)
		if aToday != bToday {
			return aToday
		}

		// Priorité 3 : Nouveaux (jamais révisés)
		aNew := IsNew(a)
		bNew := IsNew(b)
		if aNew && bNew {
			return a.ID < b.ID // Ordre de création
		}
		if aNew != bNew {
			return aNew
		}

		// Priorité 4 : Par date de révision
		return a.NextReviewAt.Before(b.NextReviewAt)
	})

	return exercises
}

// IsOverdue : Exercice en retard
func IsOverdue(ex models.Exercise) bool {
	return ex.Done && ex.NextReviewAt.Before(time.Now())
}

// IsDueToday : Exercice à réviser aujourd'hui
func IsDueToday(ex models.Exercise) bool {
	if !ex.Done {
		return false
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	return ex.NextReviewAt.After(today) && ex.NextReviewAt.Before(tomorrow)
}

// IsNew : Exercice jamais révisé
func IsNew(ex models.Exercise) bool {
	return !ex.Done && ex.LastReviewed == nil
}

// GetPriorityLabel : Label de priorité pour affichage
func GetPriorityLabel(ex models.Exercise) string {
	if IsOverdue(ex) {
		return "🔴 En retard"
	}
	if IsDueToday(ex) {
		return "🟡 Aujourd'hui"
	}
	if IsNew(ex) {
		return "🆕 Nouveau"
	}
	return "🟢 À venir"
}
