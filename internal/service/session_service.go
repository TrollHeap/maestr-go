package service

import (
	"fmt"
	"log"
	"time"

	"maestro/internal/models"
	"maestro/internal/store"

	"github.com/google/uuid"
)

// SessionService gère la logique métier des sessions
type SessionService struct{}

// NewSessionService crée une nouvelle instance
func NewSessionService() *SessionService {
	return &SessionService{}
}

// BuildAdaptiveSession construit une session selon le niveau d'énergie
func (s *SessionService) BuildAdaptiveSession(
	energy models.EnergyLevel,
) (*models.AdaptiveSession, error) {
	log.Printf("📍 [SessionService] BuildAdaptiveSession energy=%d", energy) // ← AJOUTE

	config, exists := models.SessionConfigs[energy]
	if !exists {
		log.Printf("❌ [SessionService] Config introuvable pour energy=%d", energy) // ← AJOUTE
		return nil, fmt.Errorf("niveau d'énergie invalide: %d", energy)
	}
	log.Printf("✅ [SessionService] Config trouvée: %+v", config) // ← AJOUTE

	session := &models.AdaptiveSession{
		Mode:          config.Mode,
		EnergyLevel:   energy,
		EstimatedTime: config.Duration,
		BreakSchedule: config.Breaks,
		StartedAt:     time.Now(),
		CurrentIndex:  0,
	}

	// Sélectionne les exercices selon le niveau
	log.Println("🔍 [SessionService] Appel pickDueExercises...") // ← AJOUTE
	exercises := s.pickDueExercises(config.ExerciseCount)
	log.Printf("✅ [SessionService] %d exercices sélectionnés", len(exercises)) // ← AJOUTE

	session.Exercises = exercises

	return session, nil
}

// StartSession démarre une nouvelle session
// StartSession démarre une nouvelle session
func (s *SessionService) StartSession(
	energy models.EnergyLevel,
) (string, *models.AdaptiveSession, error) {
	log.Println("📍 [SessionService] StartSession début") // ← AJOUTE

	session, err := s.BuildAdaptiveSession(energy)
	if err != nil {
		log.Printf("❌ [SessionService] BuildAdaptiveSession erreur: %v", err) // ← AJOUTE
		return "", nil, err
	}
	log.Printf("✅ [SessionService] Session construite: %+v", session) // ← AJOUTE

	// Génère un ID unique
	sessionID := uuid.New().String()
	log.Printf("🆔 [SessionService] ID généré: %s", sessionID) // ← AJOUTE

	// Crée la session active
	activeSession := &models.ActiveSession{
		ID:           sessionID,
		Session:      *session,
		CurrentIndex: 0,
		StartedAt:    time.Now(),
		CompletedIDs: []int{},
	}
	log.Printf("🔧 [SessionService] ActiveSession créée") // ← AJOUTE

	// Sauvegarde dans le store
	log.Println("💾 [SessionService] Appel store.CreateSession...") // ← AJOUTE
	if err := store.CreateSession(sessionID, activeSession); err != nil {
		log.Printf("❌ [SessionService] store.CreateSession erreur: %v", err) // ← AJOUTE
		return "", nil, fmt.Errorf("erreur création session: %w", err)
	}
	log.Println("✅ [SessionService] Session sauvegardée dans store") // ← AJOUTE

	log.Printf("🎉 [SessionService] StartSession terminé avec succès, ID=%s", sessionID) // ← AJOUTE
	return sessionID, session, nil
}

// GetActiveSession récupère la session active
func (s *SessionService) GetActiveSession() *models.ActiveSession {
	return store.GetActiveSession()
}

// CompleteExercise marque un exercice comme complété
func (s *SessionService) CompleteExercise(
	sessionID string,
	exerciseID int,
) (*models.Exercise, error) {
	session, err := store.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	session.MarkCompleted(exerciseID)

	if err := store.UpdateSession(sessionID, session); err != nil {
		return nil, fmt.Errorf("erreur mise à jour session: %w", err)
	}

	return session.NextExercise(), nil
}

// StopSession arrête une session
func (s *SessionService) StopSession(sessionID string) error {
	return store.DeleteSession(sessionID)
}

// ClearAllSessions supprime toutes les sessions
func (s *SessionService) ClearAllSessions() error {
	return store.ClearActiveSession()
}

// pickDueExercises sélectionne les N exercices les plus urgents
func (s *SessionService) pickDueExercises(count int) []models.Exercise {
	allExercises := store.GetAll()
	now := time.Now()

	var due []models.Exercise
	for _, ex := range allExercises {
		if !ex.Done && ex.NextReviewAt.Before(now) {
			due = append(due, ex)
			if len(due) >= count {
				break
			}
		}
	}

	// Si pas assez d'exercices dus, prendre des nouveaux
	if len(due) < count {
		for _, ex := range allExercises {
			if !ex.Done && len(due) < count {
				due = append(due, ex)
			}
		}
	}

	return due
}
