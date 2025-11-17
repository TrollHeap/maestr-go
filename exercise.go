package models

import "time"

// Exercise représente un exercice d'apprentissage
type Exercise struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Domain      string   `json:"domain"`     // golang, linux, architecture
	Difficulty  int      `json:"difficulty"` // 1-3
	Steps       []string `json:"steps"`
	Content     string   `json:"content"`

	// Spaced Repetition
	Completed    bool       `json:"completed"`
	LastReviewed *time.Time `json:"last_reviewed"`
	EaseFactor   float64    `json:"ease_factor"`
	IntervalDays int        `json:"interval_days"`
	Repetitions  int        `json:"repetitions"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReviewInput est la requête pour noter un exercice
type ReviewInput struct {
	ExerciseID string `json:"exercise_id"`
	Rating     int    `json:"rating"` // 1-4
}

// ReviewResponse est la réponse après notation
type ReviewResponse struct {
	Exercise     *Exercise `json:"exercise"`
	NextReviewIn int       `json:"next_review_in_days"`
	Message      string    `json:"message"`
}

// Stats représente les statistiques utilisateur
type Stats struct {
	TotalCompleted int                   `json:"total_completed"`
	TotalReviews   int                   `json:"total_reviews"`
	DomainStats    map[string]DomainStat `json:"domain_stats"`
}

// DomainStat représente les stats par domaine
type DomainStat struct {
	Completed int `json:"completed"`
	Total     int `json:"total"`
	Mastery   int `json:"mastery"` // 0-100
}
