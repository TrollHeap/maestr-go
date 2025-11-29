package models

import "time"

type ReviewResult struct {
	EaseFactor   float64
	IntervalDays int
	Repetitions  int
	NextReview   time.Time
}
