package store

import (
	"time"

	"maestro/v2-refacto/internal/models"
)

// Variable globale pour la démo (En vrai : Database struct)
var exercises = []models.Exercise{
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
		CompletedSteps: []int{}, // Aucune étape faite pour l'instant
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
		CompletedSteps: []int{}, // Aucune étape faite pour l'instant
		Content:        "# Filtrage avec Go Generics\n\nDepuis Go 1.18...",
		Done:           false,
		CreatedAt:      time.Now(),
	},
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
