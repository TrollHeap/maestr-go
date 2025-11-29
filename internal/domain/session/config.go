// internal/domain/session/config.go
package session

import (
	"time"

	"maestro/internal/models"
)

// ============================================
// SESSION CONFIGURATION (Règles Métier)
// ============================================

// Config : Configuration d'une session par niveau d'énergie
type Config struct {
	Level         models.EnergyLevel
	Mode          string
	Duration      time.Duration
	MaxExercises  int
	BreakSchedule []time.Duration
	Description   string
}

// Configs : Map des configurations par niveau d'énergie
var Configs = map[models.EnergyLevel]Config{
	models.EnergyLow: {
		Level:         models.EnergyLow,
		Mode:          "micro",
		Duration:      15 * time.Minute,
		MaxExercises:  2, // ✅ CORRIGÉ (était 3)
		BreakSchedule: []time.Duration{5 * time.Minute},
		Description:   "Session courte (1-2 exos, 15min)",
	},
	models.EnergyMedium: {
		Level:         models.EnergyMedium,
		Mode:          "standard",
		Duration:      30 * time.Minute,
		MaxExercises:  4, // ✅ CORRIGÉ (était 5)
		BreakSchedule: []time.Duration{5 * time.Minute, 10 * time.Minute},
		Description:   "Session moyenne (2-4 exos, 30min)",
	},
	models.EnergyHigh: {
		Level:         models.EnergyHigh,
		Mode:          "deep",
		Duration:      60 * time.Minute,
		MaxExercises:  8, // ✅ CORRIGÉ (était 10)
		BreakSchedule: []time.Duration{5 * time.Minute, 10 * time.Minute, 15 * time.Minute},
		Description:   "Session longue (4-8 exos, 60min)",
	},
}

// ============================================
// RÈGLES MÉTIER (Domain Logic)
// ============================================

// GetConfig : Récupère config pour un niveau d'énergie
func GetConfig(energy models.EnergyLevel) Config {
	config, exists := Configs[energy]
	if !exists {
		return Configs[models.EnergyMedium] // Défaut
	}
	return config
}

// GetMaxExercises : Retourne max exercices pour un niveau d'énergie
func GetMaxExercises(energy models.EnergyLevel) int {
	return GetConfig(energy).MaxExercises
}

// LimitExercises : Limite une liste d'exercices selon l'énergie
func LimitExercises(exerciseIDs []int, energy models.EnergyLevel) []int {
	maxExercises := GetMaxExercises(energy)

	if len(exerciseIDs) <= maxExercises {
		return exerciseIDs
	}

	return exerciseIDs[:maxExercises]
}

// EstimateSessionTime : Estime durée selon nombre d'exercices
func EstimateSessionTime(exerciseCount int, energy models.EnergyLevel) time.Duration {
	config := GetConfig(energy)

	if exerciseCount == 0 {
		return 0
	}

	// Temps par exercice (moyenne)
	timePerExercise := config.Duration / time.Duration(config.MaxExercises)

	return timePerExercise * time.Duration(exerciseCount)
}

// ShouldTakeBreak : Règle "prendre une pause après X exercices ?"
func ShouldTakeBreak(completedCount int, energy models.EnergyLevel) bool {
	switch energy {
	case models.EnergyLow:
		return false // Pas de pause en mode micro
	case models.EnergyMedium:
		return completedCount%2 == 0 && completedCount > 0 // Pause tous les 2 exos
	case models.EnergyHigh:
		return completedCount%3 == 0 && completedCount > 0 // Pause tous les 3 exos
	default:
		return false
	}
}

// GetBreakDuration : Durée de pause selon progression
func GetBreakDuration(completedCount int, energy models.EnergyLevel) time.Duration {
	config := GetConfig(energy)

	if len(config.BreakSchedule) == 0 {
		return 0
	}

	// Cycle dans les durées de pause
	index := (completedCount - 1) % len(config.BreakSchedule)
	return config.BreakSchedule[index]
}
