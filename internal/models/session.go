// internal/models/session.go
package models

import "time"

// ============================================
// ENERGY LEVEL (Enum)
// ============================================

type EnergyLevel int

const (
	EnergyLow    EnergyLevel = 1
	EnergyMedium EnergyLevel = 2
	EnergyHigh   EnergyLevel = 3
)

// ============================================
// SESSION STRUCTURES (Data)
// ============================================

// AdaptiveSession : Session en cours
type AdaptiveSession struct {
	ID            int64
	Mode          string
	EnergyLevel   EnergyLevel
	EstimatedTime time.Duration
	Exercises     []int // IDs des exercices
	BreakSchedule []time.Duration
	StartedAt     time.Time
	CurrentIndex  int
	Completed     []int // IDs des exercices complétés
}

// SessionResult : Résultat final d'une session
type SessionResult struct {
	SessionID      int64
	CompletedCount int
	Duration       time.Duration
	CompletedAt    time.Time
	Exercises      []int       // IDs des exercices complétés
	Qualities      map[int]int // exerciseID → quality
}

// ============================================
// SESSION REPORT (Dashboard)
// ============================================

// SessionReport : Rapport de disponibilité des exercices
type SessionReport struct {
	TodayDue        int              `json:"today_due"`
	TodayNew        int              `json:"today_new"`
	TotalAvailable  int              `json:"total_available"`
	NextReviewDate  time.Time        `json:"next_review_date"`
	UpcomingReviews []UpcomingReview `json:"upcoming_reviews"`
}

// UpcomingReview : Exercice à réviser dans le futur
type UpcomingReview struct {
	Date          time.Time `json:"date"`
	ExerciseID    int       `json:"exercise_id"`
	ExerciseTitle string    `json:"exercise_title"`
}
