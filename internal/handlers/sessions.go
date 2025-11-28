package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"maestro/internal/models"
	"maestro/internal/service"
)

func init() {
	exerciseService = service.NewExerciseService()
	sessionService = service.NewSessionService() // ← AJOUTE CETTE LIGNE
	log.Println("✅ SessionService initialisé")   // ← AJOUTE
}

// Page sélection énergie
func HandleSessionBuilder(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Configs": models.SessionConfigs,
	}

	if err := Tmpl.ExecuteTemplate(w, "session-builder", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Démarre une session
func HandleStartSession(w http.ResponseWriter, r *http.Request) {
	log.Println("🔥 HandleStartSession appelé")

	energyStr := r.URL.Query().Get("energy")
	log.Printf("Energy reçu: %s", energyStr)

	energy, err := strconv.Atoi(energyStr)
	if err != nil || energy < 1 || energy > 3 {
		log.Printf("❌ Énergie invalide: %v", err)
		http.Error(w, "Niveau d'énergie invalide", http.StatusBadRequest)
		return
	}

	log.Println("🚀 Appel StartSession...")
	sessionID, session, err := sessionService.StartSession(models.EnergyLevel(energy))
	if err != nil {
		log.Printf("❌ Erreur StartSession: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Session créée: %s", sessionID)

	if len(session.Exercises) == 0 {
		log.Println("❌ Aucun exercice disponible")
		// ✅ Affiche la belle page au lieu d'une erreur brute
		Tmpl.ExecuteTemplate(w, "no-exercises", nil)
		return
	}

	firstExercise := session.Exercises[0]
	redirectURL := fmt.Sprintf("/exercise/%d?from=session&sid=%s", firstExercise.ID, sessionID)
	log.Printf("➡️ Redirection vers: %s", redirectURL)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// Affiche la session en cours
func HandleCurrentSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")

	activeSession := sessionService.GetActiveSession()
	if activeSession == nil || activeSession.ID != sessionID {
		http.NotFound(w, r)
		return
	}

	data := map[string]any{
		"Session":  activeSession,
		"Exercise": &activeSession.Session.Exercises[activeSession.CurrentIndex],
	}

	if err := Tmpl.ExecuteTemplate(w, "session-current", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Complète un exercice de session
func HandleCompleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")

	nextEx, err := sessionService.CompleteExercise(sessionID, 0) // TODO: récupérer exerciseID
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if nextEx != nil {
		http.Redirect(
			w,
			r,
			fmt.Sprintf("/exercise/%d?from=session&sid=%s", nextEx.ID, sessionID),
			http.StatusSeeOther,
		)
	} else {
		// Session terminée
		sessionService.StopSession(sessionID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Arrête une session
func HandleStopSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")

	if err := sessionService.StopSession(sessionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
