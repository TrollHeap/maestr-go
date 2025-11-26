package domain

import (
	"math"
	"time"

	"maestro/internal/models"
)

// ============= SESSION MANAGER (ADHD-Friendly) =============

type SessionStatus int

const (
	SessionActive SessionStatus = iota
	SessionWarning
	SessionEnded
)

type SessionManager struct {
	SessionDuration  time.Duration
	WarningThreshold time.Duration
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		SessionDuration:  15 * time.Minute,
		WarningThreshold: 5 * time.Minute,
	}
}

func (s *SessionManager) GetStatus(elapsed time.Duration) SessionStatus {
	remaining := s.SessionDuration - elapsed

	if remaining <= 0 {
		return SessionEnded
	}

	if remaining <= s.WarningThreshold {
		return SessionWarning
	}

	return SessionActive
}

func (s *SessionManager) GetTimeRemaining(elapsed time.Duration) time.Duration {
	remaining := s.SessionDuration - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ============= ADHD RECOMMENDER =============

type ADHDRecommender struct {
	maxChoices int
}

func NewADHDRecommender() *ADHDRecommender {
	return &ADHDRecommender{
		maxChoices: 1, // Force 1 choice only
	}
}

// GetNextExercise returns THE exercise to do (or nil)
func (r *ADHDRecommender) GetNextExercise(exercises []models.Exercise) *models.Exercise {
	now := time.Now()

	// Priority 1: Overdue (urgent!)
	for _, ex := range exercises {
		if r.isOverdue(&ex, now) {
			return r.copyExercise(&ex)
		}
	}

	// Priority 2: Due today
	for _, ex := range exercises {
		if r.isDueToday(&ex, now) {
			return r.copyExercise(&ex)
		}
	}

	// Priority 3: Incomplete + easiest
	for _, ex := range exercises {
		if !ex.Completed {
			return r.copyExercise(&ex)
		}
	}

	// Nothing to do!
	return nil
}

// GetNextExercises limits to max 3
func (r *ADHDRecommender) GetNextExercises(
	exercises []models.Exercise,
	limit int,
) []models.Exercise {
	if limit > 3 {
		limit = 3
	}

	next := r.GetNextExercise(exercises)
	if next == nil {
		return []models.Exercise{}
	}

	return []models.Exercise{*next}
}

func (r *ADHDRecommender) isOverdue(ex *models.Exercise, now time.Time) bool {
	if !ex.Completed || ex.LastReviewed == nil {
		return false
	}

	nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
	return now.After(nextReview)
}

func (r *ADHDRecommender) isDueToday(ex *models.Exercise, now time.Time) bool {
	if !ex.Completed || ex.LastReviewed == nil {
		return false
	}

	nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)

	// Same day (00:00 to 23:59)
	return now.Year() == nextReview.Year() &&
		now.YearDay() == nextReview.YearDay()
}

func (r *ADHDRecommender) copyExercise(ex *models.Exercise) *models.Exercise {
	copy := *ex
	return &copy
}

// ============= REWARD ENGINE =============

type RewardEvent int

const (
	EventFirstStep RewardEvent = iota
	EventStepCompleted
	EventExerciseCompleted
	EventStreakDay
	EventStreakWeek
	EventMastery50
	EventMastery70
)

type RewardEngine struct {
	messages map[RewardEvent]string
}

func NewRewardEngine() *RewardEngine {
	return &RewardEngine{
		messages: map[RewardEvent]string{
			EventFirstStep:         "💪 Let's go! First step!",
			EventStepCompleted:     "✓ Step done! Momentum!",
			EventExerciseCompleted: "🔥 Exercice complété!",
			EventStreakDay:         "✓ 1 day streak!",
			EventStreakWeek:        "🎯 7 days! Unstoppable!",
			EventMastery50:         "📈 50% mastery!",
			EventMastery70:         "🏆 Domaine maîtrisé!",
		},
	}
}

func (e *RewardEngine) GetMessage(event RewardEvent) string {
	if msg, ok := e.messages[event]; ok {
		return msg
	}
	return "Keep going!"
}

// ============= ADHD SM-2 SCHEDULER =============

type ADHDScheduler struct {
	initialEF float64
	minEF     float64
	easyBias  float64
}

func NewADHDScheduler() *ADHDScheduler {
	return &ADHDScheduler{
		initialEF: 2.5,
		minEF:     1.3,
		easyBias:  1.1, // +10% bonus for consistency
	}
}

// ReviewExercise applique SM-2 ADHD-friendly
func (s *ADHDScheduler) ReviewExercise(ex *models.Exercise, rating int) {
	if rating < 1 || rating > 4 {
		return
	}

	var newInterval int
	var newEF float64

	switch rating {
	case 4: // Facile - huge bonus!
		newInterval = int(float64(ex.IntervalDays) * ex.EaseFactor * s.easyBias)
		newEF = ex.EaseFactor + 0.15

	case 3: // Normal
		newInterval = int(float64(ex.IntervalDays) * ex.EaseFactor)
		newEF = ex.EaseFactor

	case 2: // Difficile - gentle penalty
		newInterval = 1 // Reset to 1 day (not 0.5x)
		newEF = ex.EaseFactor - 0.1

	case 1: // Oublié - review TODAY
		newInterval = 0 // Review immediately
		newEF = ex.EaseFactor - 0.2
	}

	if newInterval < 1 && rating > 1 {
		newInterval = 1
	}

	ex.EaseFactor = math.Max(s.minEF, newEF)
	ex.IntervalDays = newInterval
	now := time.Now()
	ex.LastReviewed = &now
	ex.Repetitions++
	ex.Completed = true
}

// ============= EXERCISE ACTIONS =============

// UncompleteExercise resets an exercise to initial state
func UncompleteExercise(ex *models.Exercise) {
	ex.Completed = false
	ex.LastReviewed = nil
	ex.EaseFactor = 2.5
	ex.IntervalDays = 0
	ex.Repetitions = 0
	ex.CompletedSteps = []int{}
}

// SkipExercise marks for later review (no penalty)
func SkipExercise(ex *models.Exercise) {
	ex.SkippedCount++
	now := time.Now()
	ex.LastSkipped = &now
	// NO change to rating/interval
}

// ToggleStep toggles a step completion
func ToggleStep(ex *models.Exercise, stepIndex int) {
	idx := indexOfInt(ex.CompletedSteps, stepIndex)
	if idx >= 0 {
		// Remove
		ex.CompletedSteps = append(
			ex.CompletedSteps[:idx],
			ex.CompletedSteps[idx+1:]...,
		)
	} else {
		// Add
		ex.CompletedSteps = append(ex.CompletedSteps, stepIndex)
	}
	ex.UpdatedAt = time.Now()
}

// ============= HELPER FUNCTIONS =============

func indexOfInt(slice []int, val int) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}
