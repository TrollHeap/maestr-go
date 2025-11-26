package storage

import (
	"context"

	"maestro/internal/models"
)

// Store définit l'interface pour la persistance des exercices
type Store interface {
	// Load charge tous les exercices depuis le stockage
	Load(ctx context.Context) ([]models.Exercise, error)

	// Save sauvegarde tous les exercices dans le stockage
	Save(ctx context.Context, exercises []models.Exercise) error

	// GetByID récupère un exercice par son ID
	GetByID(ctx context.Context, id string) (*models.Exercise, error)

	// Update met à jour un exercice existant
	Update(ctx context.Context, ex *models.Exercise) error

	// Delete supprime un exercice par son ID
	Delete(ctx context.Context, id string) error
}
