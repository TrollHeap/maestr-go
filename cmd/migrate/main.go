package main

import (
	"encoding/json"
	"fmt"
	"os"

	"maestro/internal/models"
)

func main() {
	// 1. Lit l'ancien exercises.json
	data, _ := os.ReadFile("data/exercises.json")

	var allExercises []models.Exercise
	json.Unmarshal(data, &allExercises)

	// 2. Crée meta.json
	var metaList models.ExerciseList
	for _, ex := range allExercises {
		meta := models.ExerciseMeta{
			ID:         ex.ID,
			Title:      ex.Title,
			Domain:     ex.Domain,
			Difficulty: ex.Difficulty,
			FilePath:   fmt.Sprintf("data/exercises/%s_%d.json", ex.Domain, ex.ID),
		}
		metaList = append(metaList, meta)
	}

	// Sauve meta.json
	metaJSON, _ := json.MarshalIndent(metaList, "", "  ")
	os.WriteFile("data/exercises/meta.json", metaJSON, 0o644)

	// 3. Crée fichiers individuels
	for _, ex := range allExercises {
		path := fmt.Sprintf("data/exercises/%s_%d.json", ex.Domain, ex.ID)
		exJSON, _ := json.MarshalIndent(ex, "", "  ")
		os.WriteFile(path, exJSON, 0o644)
		fmt.Printf("✅ Créé : %s\n", path)
	}

	// 4. Extrait progress → user_data/progress.json
	var progressData []struct {
		ExerciseID   int     `json:"exercise_id"`
		Done         bool    `json:"done"`
		IntervalDays int     `json:"interval_days"`
		EaseFactor   float64 `json:"ease_factor"`
	}

	for _, ex := range allExercises {
		progressData = append(progressData, struct {
			ExerciseID   int     `json:"exercise_id"`
			Done         bool    `json:"done"`
			IntervalDays int     `json:"interval_days"`
			EaseFactor   float64 `json:"ease_factor"`
		}{
			ExerciseID:   ex.ID,
			Done:         ex.Done,
			IntervalDays: ex.IntervalDays,
			EaseFactor:   ex.EaseFactor,
		})
	}

	progressJSON, _ := json.MarshalIndent(progressData, "", "  ")
	os.WriteFile("data/user_data/progress.json", progressJSON, 0o644)

	fmt.Println("✅ Migration terminée !")
}
