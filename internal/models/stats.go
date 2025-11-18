package models

import "time"

// Stats représente les statistiques globales
type Stats struct {
	TotalCompleted int                   `json:"total_completed"`
	TotalExercises int                   `json:"total_exercises"`
	AverageMastery float64               `json:"average_mastery"`
	CurrentStreak  int                   `json:"current_streak"`
	StreakDisplay  string                `json:"streak_display"`
	LastSessionAt  *time.Time            `json:"last_session_at"`
	TotalReviews   int                   `json:"total_reviews"`
	DomainStats    map[string]DomainStat `json:"domain_stats"`
}

// DomainStat représente les stats pour un domaine
type DomainStat struct {
	Domain      string  `json:"domain"`
	Completed   int     `json:"completed"`
	Total       int     `json:"total"`
	Percentage  int     `json:"percentage"`
	Mastery     float64 `json:"mastery"`
	EaseFactor  float64 `json:"ease_factor"`
	Repetitions int     `json:"repetitions"`
}
