package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"maestro/internal/models"
)

// JSONStore implémente Store en utilisant des fichiers JSON
type JSONStore struct {
	mu       sync.RWMutex
	dataPath string
}

// NewJSONStore crée un nouveau JSONStore
func NewJSONStore(dataPath string) (*JSONStore, error) {
	// Créer le dossier si nécessaire
	if err := os.MkdirAll(dataPath, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &JSONStore{
		dataPath: dataPath,
	}, nil
}

// Load charge tous les exercices depuis le fichier JSON
func (s *JSONStore) Load(ctx context.Context) ([]models.Exercise, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := filepath.Join(s.dataPath, "exercises.json")

	// Si le fichier n'existe pas, retourner liste vide
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []models.Exercise{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var exercises []models.Exercise
	if err := json.Unmarshal(data, &exercises); err != nil {
		return nil, fmt.Errorf("failed to unmarshal exercises: %w", err)
	}

	return exercises, nil
}

// Save sauvegarde tous les exercices dans le fichier JSON
func (s *JSONStore) Save(ctx context.Context, exercises []models.Exercise) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(s.dataPath, "exercises.json")

	data, err := json.MarshalIndent(exercises, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal exercises: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetByID récupère un exercice par son ID
func (s *JSONStore) GetByID(ctx context.Context, id string) (*models.Exercise, error) {
	exercises, err := s.Load(ctx)
	if err != nil {
		return nil, err
	}

	for i := range exercises {
		if exercises[i].ID == id {
			return &exercises[i], nil
		}
	}

	return nil, fmt.Errorf("exercise not found: %s", id)
}

// Update met à jour un exercice existant
func (s *JSONStore) Update(ctx context.Context, ex *models.Exercise) error {
	exercises, err := s.Load(ctx)
	if err != nil {
		return err
	}

	found := false
	for i := range exercises {
		if exercises[i].ID == ex.ID {
			exercises[i] = *ex
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("exercise not found: %s", ex.ID)
	}

	return s.Save(ctx, exercises)
}

// Delete supprime un exercice par son ID
func (s *JSONStore) Delete(ctx context.Context, id string) error {
	exercises, err := s.Load(ctx)
	if err != nil {
		return err
	}

	filtered := make([]models.Exercise, 0, len(exercises))
	found := false

	for i := range exercises {
		if exercises[i].ID != id {
			filtered = append(filtered, exercises[i])
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("exercise not found: %s", id)
	}

	return s.Save(ctx, filtered)
}
