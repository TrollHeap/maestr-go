package exercise

import (
	"errors"
	"fmt"
)

// ValidateID vérifie qu'un ID est valide
func ValidateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("ID invalide: %d", id)
	}
	return nil
}

// ValidateStep vérifie qu'un step est valide pour un exercice donné
func ValidateStep(step, maxSteps int) error {
	if step < 0 {
		return fmt.Errorf("step invalide: %d (doit être >= 0)", step)
	}
	if step >= maxSteps {
		return fmt.Errorf("step invalide: %d (doit être < %d)", step, maxSteps)
	}
	return nil
}

// ValidateQuality : Vérifie qualité SRS (0-3)
func ValidateQuality(quality int) error {
	if quality < 0 || quality > 3 {
		return fmt.Errorf("quality must be 0-3 (Again/Hard/Good/Easy), got %d", quality)
	}
	return nil
}

func ValidateExerciseInput(title string, difficulty int, domain string) error {
	if len(title) == 0 || len(title) > 200 {
		return errors.New("title: length must be 1-200")
	}

	if difficulty < 1 || difficulty > 5 {
		return errors.New("difficulty: must be 1-5")
	}

	if len(domain) == 0 {
		return errors.New("domain: required")
	}

	return nil
}
