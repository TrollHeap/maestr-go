package models

import "time"

// Exercise représente un exercice d'apprentissage
type Exercise struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Domain      string   `json:"domain"`
	Difficulty  int      `json:"difficulty"` // 1-5
	Steps       []string `json:"steps"`
	Content     string   `json:"content"`

	// Spaced Repetition (SM-2)
	Completed      bool       `json:"completed"`
	CompletedSteps []int      `json:"completed_steps"`
	LastReviewed   *time.Time `json:"last_reviewed"`
	EaseFactor     float64    `json:"ease_factor"` // 1.3 - 2.5
	IntervalDays   int        `json:"interval_days"`
	Repetitions    int        `json:"repetitions"`

	// ✅ AJOUT pour ADHD features
	SkippedCount int        `json:"skipped_count"`
	LastSkipped  *time.Time `json:"last_skipped"`
	Deleted      bool       `json:"deleted"`
	DeletedAt    *time.Time `json:"deleted_at"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReviewInput représente une révision d'exercice
type ReviewInput struct {
	Rating int `json:"rating"` // 1-4
}

// ReviewResponse représente la réponse après révision
type ReviewResponse struct {
	Exercise   Exercise   `json:"exercise"`
	NextReview *time.Time `json:"next_review"`
}

// Stats représente les statistiques globales
type Stats struct {
	Total      int                   `json:"total"`
	Completed  int                   `json:"completed"`
	InProgress int                   `json:"in_progress"`
	DueReview  int                   `json:"due_review"`
	ByDomain   map[string]DomainStat `json:"by_domain"`
}

// DomainStat représente les statistiques par domaine
type DomainStat struct {
	Completed int     `json:"completed"`
	Total     int     `json:"total"`
	Mastery   float64 `json:"mastery"`
}

// PlannedSession représente une session planifiée
type PlannedSession struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	TimeSlot    string    `json:"time_slot"` // morning, afternoon, evening
	ExerciseIDs []string  `json:"exercise_ids"`
	Duration    int       `json:"duration"` // minutes
	Status      string    `json:"status"`   // planned, completed, missed
	Notes       string    `json:"notes"`
}

// DailyPlan représente le plan d'une journée
type DailyPlan struct {
	Date         string           `json:"date"`
	Sessions     []PlannedSession `json:"sessions"`
	TotalMinutes int              `json:"total_minutes"`
	Completed    int              `json:"completed"`
	Total        int              `json:"total"`
}

// WeeklyPlan représente le plan d'une semaine
type WeeklyPlan struct {
	StartDate    string      `json:"start_date"`
	EndDate      string      `json:"end_date"`
	Days         []DailyPlan `json:"days"`
	TotalMinutes int         `json:"total_minutes"`
	Completed    int         `json:"completed"`
	Total        int         `json:"total"`
}

// PlannerStats représente les statistiques du planner
type PlannerStats struct {
	TodayPlanned    int     `json:"today_planned"`
	TodayCompleted  int     `json:"today_completed"`
	WeekPlanned     int     `json:"week_planned"`
	WeekCompleted   int     `json:"week_completed"`
	CompletionRate  float64 `json:"completion_rate"`
	AverageDuration int     `json:"average_duration"`
}
