package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"maestro/internal/models"
	"os"
)

// JSONStore implémente Store avec JSON persistence
type JSONStore struct {
	filepath string
}

// NewJSONStore crée une nouvelle instance JSONStore
func NewJSONStore(filepath string) *JSONStore {
	return &JSONStore{filepath: filepath}
}

// Load charge tous les exercices depuis le fichier JSON
func (s *JSONStore) Load(ctx context.Context) ([]models.Exercise, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data, err := os.ReadFile(s.filepath)
	if err != nil {
		// Si le fichier n'existe pas, retourner une liste vide
		if os.IsNotExist(err) {
			return []models.Exercise{}, nil
		}
		return nil, fmt.Errorf("read file: %w", err)
	}

	var exercises []models.Exercise
	if err := json.Unmarshal(data, &exercises); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	return exercises, nil
}

// Save persiste les exercices dans le fichier JSON
func (s *JSONStore) Save(ctx context.Context, exercises []models.Exercise) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	data, err := json.MarshalIndent(exercises, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	if err := os.WriteFile(s.filepath, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// GetByID récupère un exercice par ID
func (s *JSONStore) GetByID(ctx context.Context, id string) (*models.Exercise, error) {
	exercises, err := s.Load(ctx)
	if err != nil {
		return nil, err
	}

	for i, ex := range exercises {
		if ex.ID == id {
			return &exercises[i], nil
		}
	}

	return nil, fmt.Errorf("exercise not found: %s", id)
}

// Update met à jour un exercice
func (s *JSONStore) Update(ctx context.Context, ex *models.Exercise) error {
	exercises, err := s.Load(ctx)
	if err != nil {
		return err
	}

	for i := range exercises {
		if exercises[i].ID == ex.ID {
			exercises[i] = *ex
			return s.Save(ctx, exercises)
		}
	}

	return fmt.Errorf("exercise not found: %s", ex.ID)
}
