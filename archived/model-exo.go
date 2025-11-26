package models

import "time"

// ✨ NOUVEAU : Signature visuelle pour encodage hippocampique
type VisualSignature struct {
	IconEmoji    string `json:"icon"`     // "🌳" pour algo arbre
	ColorHex     string `json:"color"`    // "#FF5733" unique
	Mnemonic     string `json:"mnemonic"` // "QuickSort = Chef orchestre"
	ASCIIDiagram string `json:"ascii"`    // Diagramme minimaliste
}

// Exercise représente un exercice d'apprentissage avec Spaced Repetition
type Exercise struct {
	// Identité
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Domain      string   `json:"domain"`
	Difficulty  int      `json:"difficulty"` // 1-5
	Steps       []string `json:"steps"`
	Content     string   `json:"content"`

	// ✨ NOUVEAU : Signature visuelle
	Visual VisualSignature `json:"visual"`

	// Progression Utilisateur
	Done           bool  `json:"done"`
	CompletedSteps []int `json:"completed_steps"`

	// Spaced Repetition (SM-2 Algorithm)
	LastReviewed *time.Time `json:"last_reviewed"`
	NextReviewAt time.Time  `json:"next_review_at"`
	EaseFactor   float64    `json:"ease_factor"`
	IntervalDays int        `json:"interval_days"`
	Repetitions  int        `json:"repetitions"`

	// ADHD Features
	SkippedCount int        `json:"skipped_count"`
	LastSkipped  *time.Time `json:"last_skipped"`

	// Soft Delete
	Deleted   bool       `json:"deleted"`
	DeletedAt *time.Time `json:"deleted_at"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExerciseFilter struct {
	View       string // "all", "urgent", "today", "upcoming", "active", "new"
	Domain     string // "Go", "Algorithmes", etc.
	Difficulty int    // 1-5
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
