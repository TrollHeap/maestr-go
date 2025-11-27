package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"maestro/internal/models"
)

var (
	sessions      = make(map[string]*models.ActiveSession)
	sessionsMutex sync.RWMutex
	sessionFile   = "data/sessions.json"
)

// LoadSessions charge les sessions depuis le fichier
func LoadSessions() error {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	data, err := os.ReadFile(sessionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("lecture sessions: %w", err)
	}

	return json.Unmarshal(data, &sessions)
}

// SaveSessions sauvegarde les sessions dans le fichier
// ⚠️ DOIT être appelé avec le mutex déjà acquis par le caller
func SaveSessions() error {
	os.MkdirAll("data", 0o755)
	data, _ := json.MarshalIndent(sessions, "", "  ")
	return os.WriteFile(sessionFile, data, 0o644)
}

// CreateSession crée une nouvelle session
func CreateSession(sessionID string, session *models.ActiveSession) error {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	sessions[sessionID] = session
	return SaveSessions()
}

// GetSession récupère une session par ID
func GetSession(sessionID string) (*models.ActiveSession, error) {
	sessionsMutex.RLock()
	defer sessionsMutex.RUnlock()

	session, exists := sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session %s introuvable", sessionID)
	}
	return session, nil
}

// GetActiveSession récupère la session active (dernière créée)
func GetActiveSession() *models.ActiveSession {
	sessionsMutex.RLock()
	defer sessionsMutex.RUnlock()

	// Retourne la première session (MVP mono-utilisateur)
	for _, session := range sessions {
		return session
	}
	return nil
}

// UpdateSession met à jour une session
func UpdateSession(sessionID string, session *models.ActiveSession) error {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	if _, exists := sessions[sessionID]; !exists {
		return fmt.Errorf("session %s introuvable", sessionID)
	}

	sessions[sessionID] = session
	return SaveSessions()
}

// DeleteSession supprime une session
func DeleteSession(sessionID string) error {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	delete(sessions, sessionID)
	return SaveSessions()
}

// ClearActiveSession supprime toutes les sessions actives
func ClearActiveSession() error {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	sessions = make(map[string]*models.ActiveSession)
	return SaveSessions()
}
