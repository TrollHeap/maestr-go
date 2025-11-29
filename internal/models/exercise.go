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
	View       string // "all", "urgent", "today", "upcoming", "active", "new"
	Domain     string
	Difficulty int
}

// ========================================
// 2. MÉTHODES DE LA STRUCT (Juste en dessous)
// ========================================

// AllStepsCompleted vérifie si toutes les étapes sont validées
func (e *Exercise) AllStepsCompleted() bool {
	if len(e.Steps) == 0 {
		return false
	}
	return len(e.CompletedSteps) == len(e.Steps)
}

// IsDueForReview vérifie si l'exercice doit être révisé aujourd'hui (Spaced Repetition)
func (e *Exercise) IsDueForReview() bool {
	if e.LastReviewed == nil {
		return true // Jamais révisé = due
	}
	nextReview := e.LastReviewed.AddDate(0, 0, e.IntervalDays)
	return time.Now().After(nextReview)
}

// IsAtRisk détecte si l'exercice est ignoré depuis trop longtemps (ADHD flag)
func (e *Exercise) IsAtRisk() bool {
	if e.LastSkipped == nil {
		return false
	}
	return time.Since(*e.LastSkipped) > 7*24*time.Hour // 7 jours sans toucher
}

// Dans models/exercise.go
// Dans models/exercise.go
func (e *Exercise) MasteryLevel() string {
	if e.Repetitions == 0 {
		return "[░░░░░] 0%" // Jamais révisé
	}
	if e.Repetitions >= 6 && e.EaseFactor >= 2.3 {
		return "[█████] 100%" // Maîtrise totale
	}
	if e.Repetitions >= 4 && e.EaseFactor >= 2.1 {
		return "[████░] 80%"
	}
	if e.Repetitions >= 2 && e.EaseFactor >= 1.9 {
		return "[███░░] 60%"
	}
	if e.Repetitions >= 1 {
		return "[██░░░] 40%"
	}
	return "[█░░░░] 20%" // Premier essai
}
