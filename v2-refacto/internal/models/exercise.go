package models

import "time"

// Exercise représente un exercice d'apprentissage avec Spaced Repetition
type Exercise struct {
	// Identité
	ID          int      `json:"id"` // On garde int pour la simplicité routing
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Domain      string   `json:"domain"`
	Difficulty  int      `json:"difficulty"` // 1-5
	Steps       []string `json:"steps"`
	Content     string   `json:"content"`

	// Progression Utilisateur
	Done           bool  `json:"done"`            // Marqué manuellement (ton système actuel)
	CompletedSteps []int `json:"completed_steps"` // Indices des étapes validées

	// 🔥 Spaced Repetition (SM-2 Algorithm)
	LastReviewed *time.Time `json:"last_reviewed"` // Dernière révision
	NextReviewAt time.Time  `json:"next_review_at"`
	EaseFactor   float64    `json:"ease_factor"`   // 1.3 - 2.5 (facilité mémorisation)
	IntervalDays int        `json:"interval_days"` // Prochaine révision dans X jours
	Repetitions  int        `json:"repetitions"`   // Nombre de révisions réussies

	// 🔥 ADHD Features (Anti-Blocage)
	SkippedCount int        `json:"skipped_count"` // Combien de fois ignoré
	LastSkipped  *time.Time `json:"last_skipped"`  // Dernière fois ignoré (flag rouge si > 7 jours)

	// Soft Delete (Archivage)
	Deleted   bool       `json:"deleted"`
	DeletedAt *time.Time `json:"deleted_at"`

	// Timestamps (Audit)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExerciseFilter struct {
	Domain     string // "Algorithmes", "Go", "" (tous)
	Status     string // "done", "todo", "" (tous)
	Difficulty int    // 0 (tous), 1-5
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
