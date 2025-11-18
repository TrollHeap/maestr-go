package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"maestro/internal/models"
)

// JSONStore implémente Store avec JSON persistence
type JSONStore struct {
	filepath  string
	exercises []models.Exercise
	mu        sync.RWMutex
}

// NewJSONStore crée une nouvelle instance JSONStore
func NewJSONStore(filepath string) *JSONStore {
	return &JSONStore{
		filepath:  filepath,
		exercises: []models.Exercise{},
	}
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

	s.mu.Lock()
	s.exercises = exercises
	s.mu.Unlock()

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

	if err := os.WriteFile(s.filepath, data, 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	s.mu.Lock()
	s.exercises = exercises
	s.mu.Unlock()

	return nil
}

// saveToFile sauvegarde l'état actuel dans le fichier
func (s *JSONStore) saveToFile() error {
	data, err := json.MarshalIndent(s.exercises, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	if err := os.WriteFile(s.filepath, data, 0o644); err != nil {
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

// GetExercise récupère un exercice par ID (version sans context)
func (s *JSONStore) GetExercise(id string) (*models.Exercise, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.exercises {
		if s.exercises[i].ID == id {
			return &s.exercises[i], nil
		}
	}
	return nil, fmt.Errorf("exercise not found: %s", id)
}

// UpdateExercise met à jour un exercice existant (version sans context)
func (s *JSONStore) UpdateExercise(exercise *models.Exercise) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.exercises {
		if s.exercises[i].ID == exercise.ID {
			exercise.UpdatedAt = time.Now()
			s.exercises[i] = *exercise
			return s.saveToFile()
		}
	}
	return fmt.Errorf("exercise not found: %s", exercise.ID)
}
