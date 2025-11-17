package domain

import (
	"testing"
	"time"

	"maestro/internal/models"
)

func TestSM2Algorithm_RatingFour(t *testing.T) {
	scheduler := NewScheduler()

	ex := &models.Exercise{
		ID:           "test-1",
		Title:        "Test",
		EaseFactor:   2.5,
		IntervalDays: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	scheduler.ReviewExercise(ex, 4) // Easy

	if ex.EaseFactor != 2.6 {
		t.Fatalf("Expected EF 2.6, got %f", ex.EaseFactor)
	}

	if ex.IntervalDays != 1 {
		t.Fatalf("Expected interval 1, got %d", ex.IntervalDays)
	}

	if !ex.Completed {
		t.Fatal("Exercise should be marked as completed")
	}
}

func TestSM2Algorithm_RatingOne(t *testing.T) {
	scheduler := NewScheduler()

	ex := &models.Exercise{
		ID:           "test-2",
		Title:        "Test",
		EaseFactor:   2.5,
		IntervalDays: 5,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	scheduler.ReviewExercise(ex, 1) // Forgot

	if ex.EaseFactor != 2.0 {
		t.Fatalf("Expected EF 2.0, got %f", ex.EaseFactor)
	}

	if ex.IntervalDays != 1 {
		t.Fatalf("Expected interval reset to 1, got %d", ex.IntervalDays)
	}
}

func TestIsDueForReview_NewExercise(t *testing.T) {
	scheduler := NewScheduler()

	ex := &models.Exercise{
		ID:           "test-3",
		LastReviewed: nil,
	}

	// ✅ FIX: ex est déjà un pointer, pas &ex
	if scheduler.IsDueForReview(ex) {
		t.Fatal("New exercise should not be due yet")
	}
}

func TestIsDueForReview_PastDate(t *testing.T) {
	scheduler := NewScheduler()

	now := time.Now()
	past := now.AddDate(0, 0, -2) // 2 days ago

	ex := &models.Exercise{
		ID:           "test-4",
		LastReviewed: &past,
		IntervalDays: 1, // Due 1 day after last review = yesterday
	}

	// ✅ FIX: ex est déjà un pointer, pas &ex
	if !scheduler.IsDueForReview(ex) {
		t.Fatal("Exercise should be due for review")
	}
}

func TestGetDaysUntilReview(t *testing.T) {
	scheduler := NewScheduler()

	now := time.Now()
	// ✅ FIX: Supprimé la variable non-utilisée 'future'

	ex := &models.Exercise{
		ID:           "test-5",
		LastReviewed: &now,
		IntervalDays: 5,
	}

	// Roughly check (may be off by 1 due to time difference)
	// ✅ FIX: ex est déjà un pointer, pas &ex
	days := scheduler.GetDaysUntilReview(ex)
	if days < 4 || days > 5 {
		t.Fatalf("Expected 4-5 days, got %d", days)
	}
}
