package models

import "time"

// VisualAid représente un diagramme ou visuel pédagogique
type VisualAid struct {
	Type    string `json:"type"`    // "ascii", "svg", "mermaid"
	Content string `json:"content"` // Le diagramme lui-même
	Caption string `json:"caption"` // Description courte
}
type ExerciseView struct {
	Exercise    *Exercise
	FromSession bool
}

type Exercise struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Domain      string   `json:"domain"`
	Difficulty  int      `json:"difficulty"`
	Steps       []string `json:"steps"`
	Content     string   `json:"content"`

	// Visuels pédagogiques (nouveau format structuré)
	ConceptualVisuals []VisualAid `json:"conceptual_visuals"`
	Mnemonic          string      `json:"mnemonic"`

	// Ancien format visual (compatibilité)
	Visual map[string]string `json:"visual,omitempty"`

	// SRS tracking
	Done           bool       `json:"done"`
	CompletedSteps []int      `json:"completed_steps"`
	LastReviewed   *time.Time `json:"last_reviewed_date,omitempty"`
	NextReviewAt   time.Time  `json:"next_review_date"`
	EaseFactor     float64    `json:"ease_factor"`
	IntervalDays   int        `json:"interval_days"`
	Repetitions    int        `json:"repetitions"`
	SkippedCount   int        `json:"skipped_count"`
	LastSkipped    *time.Time `json:"last_skipped_date,omitempty"`

	// Soft delete
	Deleted   bool       `json:"deleted"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExerciseFilter struct {
	Status     string // "in_progress", "mastered", "" (tous)
	Domain     string // "Go", "Algorithms", "" (tous)
	Difficulty int    // 1-4, 0 = tous

	Query string // 🔍 texte de recherche
	Sort  string // "title", "difficulty", "domain", "" (default)
}
