package storage

import (
	"context"

	"maestro/internal/models"
)

// Store est l'interface pour la persistance
type Store interface {
	// Méthodes avec context
	Load(ctx context.Context) ([]models.Exercise, error)
	Save(ctx context.Context, exercises []models.Exercise) error
	GetByID(ctx context.Context, id string) (*models.Exercise, error)
	Update(ctx context.Context, ex *models.Exercise) error

	// Méthodes sans context (pour les handlers API)
	GetExercise(id string) (*models.Exercise, error)
	UpdateExercise(exercise *models.Exercise) error
}
