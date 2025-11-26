package validator

import (
	"errors"
	"fmt"
)

// Valide qu'un ID est positif et raisonnable
func ValidateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("ID invalide: doit être positif")
	}
	if id > 1000000 { // Limite arbitraire pour éviter les abus
		return fmt.Errorf("ID suspect: trop grand")
	}
	return nil
}

// Valide que l'étape est dans les bornes de l'exercice
func ValidateStep(stepIndex int, maxSteps int) error {
	if stepIndex < 0 || stepIndex >= maxSteps {
		return errors.New("étape hors limites")
	}
	return nil
}
