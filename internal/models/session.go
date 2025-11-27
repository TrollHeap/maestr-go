// internal/models/session.go
package models

import "time"

type EnergyLevel int

const (
	EnergyLow    EnergyLevel = 1 // 🔋 Faible (5-15min)
	EnergyMedium EnergyLevel = 2 // 🔋🔋 Moyenne (20-30min)
	EnergyHigh   EnergyLevel = 3 // 🔋🔋🔋 Haute (45-60min)
)

type SessionMode string

const (
	MicroSession    SessionMode = "micro"    // 5-15min
	StandardSession SessionMode = "standard" // 20-30min
	DeepSession     SessionMode = "deep"     // 45-60min
)

type AdaptiveSession struct {
	Mode          SessionMode
	EnergyLevel   EnergyLevel
	EstimatedTime time.Duration
	Exercises     []Exercise
	BreakSchedule []time.Duration
	StartedAt     time.Time
	CurrentIndex  int // Exercice actuel (0-based)
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
		Breaks:        []time.Duration{12 * time.Minute}, // Pause après ex1
	},
	EnergyHigh: {
		Mode:          DeepSession,
		Duration:      50 * time.Minute,
		ExerciseCount: 3,
		Breaks:        []time.Duration{17 * time.Minute, 34 * time.Minute},
	},
}
