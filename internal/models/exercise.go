package models

import "time"

// Exercise représente un exercice d'apprentissage
type Exercise struct {
	// Core fields
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Domain      string   `json:"domain"`
	Difficulty  int      `json:"difficulty"`
	Steps       []string `json:"steps"`
	Content     string   `json:"content"`

	// Spaced Repetition (SM-2)
	Completed    bool       `json:"completed"`
	LastReviewed *time.Time `json:"last_reviewed"`
	EaseFactor   float64    `json:"ease_factor"`
	IntervalDays int        `json:"interval_days"`
	Repetitions  int        `json:"repetitions"`

	// ============= NOUVELLES FONCTIONNALITÉS v3.0 =============

	// Tracking steps (persistence)
	CompletedSteps []int `json:"completed_steps"`

	// Skip tracking (no penalty)
	SkippedCount int        `json:"skipped_count"`
	LastSkipped  *time.Time `json:"last_skipped"`

	// Soft delete (undo capability)
	Deleted   bool       `json:"deleted"`
	DeletedAt *time.Time `json:"deleted_at"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReviewInput représente l'input pour noter un exercice
type ReviewInput struct {
	ExerciseID string `json:"exercise_id"`
	Rating     int    `json:"rating"` // 1-4
}

// ReviewResponse représente la réponse après une note
type ReviewResponse struct {
	Exercise     *Exercise `json:"exercise"`
	NextReviewIn int       `json:"next_review_in"`
	Message      string    `json:"message"`
}
