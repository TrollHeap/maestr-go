// internal/models/dashboard.go
package models

import "time"

// DashboardStats - Stats complètes
type DashboardStats struct {
	// Core metrics
	TotalExercises    int
	CompletedCount    int
	InProgressCount   int
	TodoCount         int
	OverdueCount      int
	TotalMastered     int
	StreakDays        int
	WeeklyReviews     int
	AvgSessionTime    time.Duration
	CompletionRate    int
	AverageInterval   int
	NextReviewDate    time.Time
	TotalSessionTime  time.Duration
	SessionCount      int
	AverageDifficulty float64
	TopDomain         string
	DomainBreakdown   map[string]int

	// Advanced analytics
	AverageEaseFactor float64
	RetentionRate     int // % de reviews réussies (quality >= 2)

	// ✅ NEW: Learning velocity
	TotalReviews     int     // Total reviews (30d)
	PeakDailyReviews int     // Max reviews in a single day
	VelocityTrend    float64 // +/- trend percentage

	// ✅ NEW: Focus recommendations
	WeakestDomain          string
	WeakestDomainRetention int
	WeakestDomainCount     int
	ShortIntervalCount     int // Exercises with interval < 7d
	LowEaseCount           int // Exercises with ease < 2.3
}

// FailurePattern - Pattern d'échecs répétés
type FailurePattern struct {
	ExerciseID int
	Title      string
	Domain     string
	FailCount  int // Nombre d'échecs (quality 0-1)
	EaseFactor float64
}

// RepetitionStat - Exercices les plus révisés
type RepetitionStat struct {
	ExerciseID  int
	Title       string
	Domain      string
	ReviewCount int
}

// DomainStrength - Force par domaine
type DomainStrength struct {
	Name            string
	TotalCount      int
	MasteredCount   int
	AvgEaseFactor   float64
	StrengthPercent int
}
