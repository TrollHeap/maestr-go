package models

import "time"

type Exercise struct {
	// Identification
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	Difficulty  int    `json:"difficulty"`

	// Contenu
	Steps          []string `json:"steps"`
	CompletedSteps []int    `json:"completed_steps"`
	Content        string   `json:"content"`

	// État
	Completed bool `json:"completed"`
	Deleted   bool `json:"deleted"`

	// Spaced Repetition
	LastReviewed *time.Time `json:"last_reviewed"` // ← POINTEUR
	NextReview   time.Time  `json:"next_review"`
	Repetitions  int        `json:"repetitions"`
	EaseFactor   float64    `json:"ease_factor"`
	Interval     int        `json:"interval"`
	IntervalDays int        `json:"interval_days"` // ← AJOUTÉ

	// ADHD Features
	SkippedCount int        `json:"skipped_count"` // ← AJOUTÉ
	LastSkipped  *time.Time `json:"last_skipped"`  // ← AJOUTÉ

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // ← AJOUTÉ
}

type ExerciseList struct {
	Exercises   []Exercise        `json:"exercises"`
	Total       int               `json:"total"`
	Page        int               `json:"page"`
	PageSize    int               `json:"page_size"`
	ReviewDates map[string]string `json:"review_dates"`
}
