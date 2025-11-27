package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"maestro/internal/models"
)

var exercises = []models.Exercise{}

const dataFile = "data/exercises.json"

// Load charge les données depuis le fichier JSON
func Load() error {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Fichier absent = normal au premier lancement
		}
		return fmt.Errorf("lecture data: %w", err)
	}
	return json.Unmarshal(data, &exercises)
}

func Save() error {
	os.MkdirAll("data", 0o755)
	data, _ := json.MarshalIndent(exercises, "", "  ")
	return os.WriteFile(dataFile, data, 0o644)
}

func GetFiltered(filter models.ExerciseFilter) []models.Exercise {
	var result []models.Exercise
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	upcoming := today.AddDate(0, 0, 3) // +3 jours

	for _, ex := range exercises {
		included := true

		// ✅ FILTRE PAR VUE SRS
		if filter.View != "" && filter.View != "all" {
			switch filter.View {
			case "urgent":
				// Révisions en retard (overdue)
				if !ex.Done || ex.NextReviewAt.After(now) {
					included = false
				}

			case "today":
				// Révisions dues aujourd'hui
				if !ex.Done || ex.NextReviewAt.After(today) ||
					ex.NextReviewAt.Before(now.AddDate(0, 0, -1)) {
					included = false
				}

			case "upcoming":
				// Révisions dans les 3 prochains jours
				if !ex.Done || ex.NextReviewAt.Before(today) || ex.NextReviewAt.After(upcoming) {
					included = false
				}

			case "active":
				// En cours (WIP) - pas done, au moins 1 étape complétée
				if ex.Done || len(ex.CompletedSteps) == 0 {
					included = false
				}

			case "new":
				// Nouveaux (TODO) - pas done, aucune étape complétée
				if ex.Done || len(ex.CompletedSteps) > 0 {
					included = false
				}
			}
		}

		// Filtre par domaine
		if filter.Domain != "" && ex.Domain != filter.Domain {
			included = false
		}

		// ✅ NOUVEAU : Filtre par difficulté
		if filter.Difficulty > 0 && ex.Difficulty != filter.Difficulty {
			included = false
		}

		if included {
			result = append(result, ex)
		}
	}

	// ✅ TRI INTELLIGENT PAR PRIORITÉ
	result = sortByPriority(result, filter.View)

	return result
}

// Tri par priorité SRS
func sortByPriority(exercises []models.Exercise, view string) []models.Exercise {
	sort.Slice(exercises, func(i, j int) bool {
		a, b := exercises[i], exercises[j]

		if view == "urgent" || view == "today" || view == "upcoming" {
			// Tri par date de révision (plus ancien en premier)
			return a.NextReviewAt.Before(b.NextReviewAt)
		}

		if view == "active" {
			// Tri par progression (plus avancé en premier)
			return len(a.CompletedSteps) > len(b.CompletedSteps)
		}

		// Par défaut : tri par ID
		return a.ID < b.ID
	})

	return exercises
}

// findExercise retourne un pointeur pour pouvoir modifier l'original
func FindExercise(id int) *models.Exercise {
	for i := range exercises {
		if exercises[i].ID == id {
			return &exercises[i]
		}
	}
	return nil
}

func GetAll() []models.Exercise {
	return exercises
}

// NOTE: Just en mode Dev pour faciliter les tests
// InitDefaultExercises peuple le store avec des exercices par défaut
// Appelée au premier lancement si le store est vide
func InitDefaultExercises() error {
	if len(exercises) > 0 {
		return nil // Déjà initialisé
	}

	exercises = []models.Exercise{
		{
			ID:          1,
			Title:       "Tri rapide (Quicksort)",
			Description: "Implémenter l'algorithme de tri rapide en Go",
			Domain:      "Algorithmes",
			Difficulty:  4,
			Steps: []string{
				"Comprendre le principe du pivot",
				"Implémenter la partition",
				"Récursion gauche et droite",
				"Tester avec des cas limites",
			},
			CompletedSteps: []int{},
			Content:        "# Quicksort...",
			Done:           false,
			CreatedAt:      time.Now(),
		},
		{
			ID:          2,
			Title:       "Filtrage de slice",
			Description: "Créer une fonction générique de filtrage",
			Domain:      "Go",
			Difficulty:  2,
			Steps: []string{
				"Définir la signature avec generics",
				"Implémenter la boucle de filtrage",
				"Écrire les tests unitaires",
			},
			CompletedSteps: []int{},
			Content:        "# Filtrage avec Go Generics...",
			Done:           false,
			CreatedAt:      time.Now(),
		},
	}

	return Save() // Persiste immédiatement
}

func CountByView(view string) int {
	filter := models.ExerciseFilter{View: view}
	return len(GetFiltered(filter))
}

func BuildAdaptiveSession(energy models.EnergyLevel) models.AdaptiveSession {
	config := models.SessionConfigs[energy]
	var selectedExercises []models.Exercise

	// 1. DEBUG: On prend simplement tous les exercices non supprimés
	allExercises := GetFiltered(models.ExerciseFilter{View: "all"})
	fmt.Println("  All exercises loaded:", len(allExercises))
	for i := 0; i < config.ExerciseCount && i < len(allExercises); i++ {
		if !allExercises[i].Deleted {
			selectedExercises = append(selectedExercises, allExercises[i])
		}
	}

	fmt.Println("  Exercises selected for session:", len(selectedExercises))
	return models.AdaptiveSession{
		Mode:          config.Mode,
		EnergyLevel:   energy,
		EstimatedTime: config.Duration,
		Exercises:     selectedExercises,
		BreakSchedule: config.Breaks,
	}
}
