package domain

import (
	"fmt"
	"time"

	"maestro/internal/models"
)

// StreakManager gère les streaks (jours consécutifs)
type StreakManager struct {
	lastSessionDate *time.Time
	currentStreak   int
}

// NewStreakManager crée une nouvelle instance
func NewStreakManager() *StreakManager {
	return &StreakManager{
		currentStreak: 0,
	}
}

// UpdateStreak met à jour le streak après une session
func (sm *StreakManager) UpdateStreak(today time.Time) int {
	if sm.lastSessionDate == nil {
		// Première session
		sm.currentStreak = 1
		sm.lastSessionDate = &today
		return 1
	}

	lastDate := *sm.lastSessionDate
	daysDiff := int(today.Sub(lastDate).Hours() / 24)

	switch daysDiff {
	case 0:
		// Même jour, pas de changement
		return sm.currentStreak
	case 1:
		// Jour suivant, streak continue
		sm.currentStreak++
		sm.lastSessionDate = &today
	default:
		// Streak cassé, reset
		sm.currentStreak = 1
		sm.lastSessionDate = &today
	}

	return sm.currentStreak
}

// GetCurrentStreak retourne le streak actuel
func (sm *StreakManager) GetCurrentStreak() int {
	return sm.currentStreak
}

// GetStreakDisplay retourne une représentation visuelle du streak
func (sm *StreakManager) GetStreakDisplay() string {
	display := ""
	for i := 0; i < sm.currentStreak && i < 30; i++ {
		display += "✓"
	}
	return display
}

// CalculateNextReviewDates calcule les dates de révision pour une liste d'exercices
func CalculateNextReviewDates(exercises []models.Exercise) map[string]string {
	dates := make(map[string]string)

	for _, ex := range exercises {
		if ex.LastReviewed == nil {
			dates[ex.ID] = "Nouveau"
			continue
		}

		lastReview := *ex.LastReviewed
		nextReview := lastReview.AddDate(0, 0, ex.IntervalDays)

		now := time.Now()
		diffDays := int(nextReview.Sub(now).Hours() / 24)

		if diffDays < 0 {
			dates[ex.ID] = "À réviser maintenant"
		} else if diffDays == 0 {
			dates[ex.ID] = "Aujourd'hui"
		} else if diffDays == 1 {
			dates[ex.ID] = "Demain"
		} else if diffDays <= 7 {
			// ✅ FIX: Convertir int en string avec fmt.Sprintf
			dates[ex.ID] = fmt.Sprintf("%d jours", diffDays)
		} else {
			dates[ex.ID] = nextReview.Format("02 Jan 2006")
		}
	}

	return dates
}

// GetStreakStats retourne les stats du streak
type StreakStats struct {
	CurrentStreak int        `json:"current_streak"`
	StreakDisplay string     `json:"streak_display"`
	LastSessionAt *time.Time `json:"last_session_at"`
}

func (sm *StreakManager) GetStats() StreakStats {
	return StreakStats{
		CurrentStreak: sm.currentStreak,
		StreakDisplay: sm.GetStreakDisplay(),
		LastSessionAt: sm.lastSessionDate,
	}
}
