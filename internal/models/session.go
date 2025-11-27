package models

import (
	"time"
)

type EnergyLevel int

const (
	EnergyLow    EnergyLevel = 1 // Faible (5-15min)
	EnergyMedium EnergyLevel = 2 // Moyenne (20-30min)
	EnergyHigh   EnergyLevel = 3 // Haute (45-60min)
)

type SessionMode string

const (
	MicroSession    SessionMode = "micro"    // 5-15min
	StandardSession SessionMode = "standard" // 20-30min
	DeepSession     SessionMode = "deep"     // 45-60min
)

type AdaptiveSession struct {
	Mode          SessionMode     `json:"mode"`
	EnergyLevel   EnergyLevel     `json:"energy_level"`
	EstimatedTime time.Duration   `json:"estimated_time"`
	Exercises     []Exercise      `json:"exercises"`
	BreakSchedule []time.Duration `json:"break_schedule"`
	StartedAt     time.Time       `json:"started_at"`
	CurrentIndex  int             `json:"current_index"`
}

// Config pour chaque mode
var SessionConfigs = map[EnergyLevel]struct {
	Mode          SessionMode
	Duration      time.Duration
	ExerciseCount int
	Breaks        []time.Duration
}{
	EnergyLow: {
		Mode:          MicroSession,
		Duration:      10 * time.Minute,
		ExerciseCount: 1,
		Breaks:        []time.Duration{},
	},
	EnergyMedium: {
		Mode:          StandardSession,
		Duration:      25 * time.Minute,
		ExerciseCount: 2,
		Breaks:        []time.Duration{12 * time.Minute},
	},
	EnergyHigh: {
		Mode:          DeepSession,
		Duration:      50 * time.Minute,
		ExerciseCount: 3,
		Breaks:        []time.Duration{17 * time.Minute, 34 * time.Minute},
	},
}

// ActiveSession représente une session en cours
type ActiveSession struct {
	ID           string          `json:"id"`
	Session      AdaptiveSession `json:"session"`
	CurrentIndex int             `json:"current_index"`
	StartedAt    time.Time       `json:"started_at"`
	CompletedIDs []int           `json:"completed_ids"`
}

// MarkCompleted marque un exercice comme complété
func (as *ActiveSession) MarkCompleted(exerciseID int) {
	as.CompletedIDs = append(as.CompletedIDs, exerciseID)
}

// NextExercise retourne le prochain exercice
func (as *ActiveSession) NextExercise() *Exercise {
	as.CurrentIndex++
	if as.CurrentIndex >= len(as.Session.Exercises) {
		return nil // Session terminée
	}
	return &as.Session.Exercises[as.CurrentIndex]
}

// IsLast vérifie si c'est le dernier exercice
func (as *ActiveSession) IsLast() bool {
	return as.CurrentIndex >= len(as.Session.Exercises)-1
}
