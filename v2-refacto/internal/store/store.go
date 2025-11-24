package store

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"maestro/v2-refacto/internal/models"
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
	var results []models.Exercise

	// Optimisation : make avec capacité 0 mais on append
	// Ou mieux : estimer la taille si possible, mais ici restons simples

	for _, ex := range exercises {
		// Filtre Domaine
		if filter.Domain != "" && ex.Domain != filter.Domain {
			continue
		}

		// Filtre Statut
		if filter.Status == "done" && !ex.Done {
			continue
		}
		if filter.Status == "todo" && ex.Done {
			continue
		}

		// Filtre Difficulté (si spécifié > 0)
		if filter.Difficulty > 0 && ex.Difficulty != filter.Difficulty {
			continue
		}

		results = append(results, ex)
	}

	return results
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
