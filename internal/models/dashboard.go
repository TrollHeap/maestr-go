package models

import "time"

// DashboardStats - Stats complètes
type DashboardStats struct {
	// Existant
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

	// ✅ NOUVEAUX pour analyse avancée
	AverageEaseFactor float64
	RetentionRate     int // % de reviews réussies (quality >= 2)
}

// ✅ FailurePattern - Pattern d'échecs répétés
type FailurePattern struct {
	ExerciseID int
	Title      string
	Domain     string
	FailCount  int // Nombre d'échecs (quality 0-1)
	EaseFactor float64
}

// ✅ RepetitionStat - Exercices les plus révisés
type RepetitionStat struct {
	ExerciseID  int
	Title       string
	Domain      string
	ReviewCount int
}

// ✅ DomainStrength - Force par domaine
type DomainStrength struct {
	Name            string
	TotalCount      int
	MasteredCount   int
	AvgEaseFactor   float64
	StrengthPercent int
}
