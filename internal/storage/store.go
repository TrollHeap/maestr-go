package storage

import (
	"context"
	"maestro/internal/models"
)

// Store est l'interface pour la persistance
type Store interface {
	// Load charge tous les exercices
	Load(ctx context.Context) ([]models.Exercise, error)

	// Save persiste les exercices
	Save(ctx context.Context, exercises []models.Exercise) error

	// GetByID récupère un exercice par ID
	GetByID(ctx context.Context, id string) (*models.Exercise, error)

	// Update met à jour un exercice
	Update(ctx context.Context, ex *models.Exercise) error
}
