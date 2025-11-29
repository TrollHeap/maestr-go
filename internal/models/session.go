package models

import (
	"fmt"
	"time"
)

// ============================================
// SESSION CONFIGURATION
// ============================================

type EnergyLevel int

const (
	EnergyLow    EnergyLevel = 1
	EnergyMedium EnergyLevel = 2
	EnergyHigh   EnergyLevel = 3
)

type SessionConfig struct {
	Mode          string
	Duration      time.Duration
	ExerciseCount int
	Breaks        []time.Duration
}

var SessionConfigs = map[EnergyLevel]SessionConfig{
	EnergyLow: {
		Mode:          "micro",
		Duration:      15 * time.Minute,
		ExerciseCount: 3,
		Breaks:        []time.Duration{5 * time.Minute},
	},
	EnergyMedium: {
		Mode:          "standard",
		Duration:      30 * time.Minute,
		ExerciseCount: 5,
		Breaks:        []time.Duration{5 * time.Minute, 10 * time.Minute},
	},
	EnergyHigh: {
		Mode:          "deep",
		Duration:      60 * time.Minute,
		ExerciseCount: 10,
		Breaks:        []time.Duration{5 * time.Minute, 10 * time.Minute, 15 * time.Minute},
	},
}

// ============================================
// SESSION STRUCTURES
// ============================================

type AdaptiveSession struct {
	Mode          string
	EnergyLevel   EnergyLevel
	EstimatedTime time.Duration
	Exercises     []Exercise
	BreakSchedule []time.Duration
	StartedAt     time.Time
	CurrentIndex  int
}

type SessionResult struct {
	CompletedCount int
	Duration       time.Duration
	CompletedAt    time.Time
	Exercises      []int // IDs des exercices complétés
}

// ============================================
// SESSION REPORT
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

// NoExercisesTodayError : Erreur custom avec rapport
type NoExercisesTodayError struct {
	Report SessionReport
}

func (e *NoExercisesTodayError) Error() string {
	if e.Report.NextReviewDate.IsZero() {
		return "Aucun exercice disponible aujourd'hui. Aucune révision programmée."
	}
	return fmt.Sprintf("Aucun exercice disponible aujourd'hui. Prochaine révision : %s",
		e.Report.NextReviewDate.Format("2006-01-02"))
}
