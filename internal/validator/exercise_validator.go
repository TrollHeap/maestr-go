package validator

import (
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

// ValidateQuality vérifie qu'une quality SRS est valide
func ValidateQuality(quality int) error {
	validQualities := map[int]bool{
		0: true, // Again
		1: true, // Hard
		3: true, // Good
		5: true, // Easy
	}

	if !validQualities[quality] {
		return fmt.Errorf("quality invalide: %d (doit être 0, 1, 3, ou 5)", quality)
	}
	return nil
}
