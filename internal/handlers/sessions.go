package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"maestro/internal/domain/session"
	"maestro/internal/models"
	"maestro/internal/service"
	"maestro/internal/store"
	"maestro/internal/views/pages"
)

// ============================================
// SERVICE GLOBAL
// ============================================

var sessionService *service.SessionService

func init() {
	sessionService = service.NewSessionService()
}

// ============================================
// 1️⃣ SESSION BUILDER (Choix énergie)
// ============================================

func HandleSessionBuilder(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 SessionBuilder: show energy selection")

	// ✅ Utilise domain configs au lieu de models
	configs := []session.Config{
		session.GetConfig(models.EnergyLow),
		session.GetConfig(models.EnergyMedium),
		session.GetConfig(models.EnergyHigh),
	}

	// ✅ CHANGEMENT : Render avec templ
	component := pages.SessionBuilder(configs)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// 2️⃣ SESSION START (Démarrage)
// ============================================

func HandleStartSession(w http.ResponseWriter, r *http.Request) {
	// 1. PARSE energy (LOGIQUE IDENTIQUE)
	energyStr := r.URL.Query().Get("energy")
	energy, err := strconv.Atoi(energyStr)
	if err != nil || energy < 1 || energy > 3 {
		energy = 2 // Default medium
	}

	energyLevel := models.EnergyLevel(energy)
	log.Printf("🔍 START SESSION: energy=%d", energy)

	// 2. RÉCUPÈRE EXERCICES DISPONIBLES (LOGIQUE IDENTIQUE)
	report, exercises, err := store.GetTodayReport()
	if err != nil {
		log.Printf("❌ GetTodayReport failed: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("🔍 [SESSION] Disponibles: %d dus + %d nouveaux = %d total",
		report.TodayDue, report.TodayNew, len(exercises))

	// 3. AUCUN EXERCICE ? Affiche rapport (LOGIQUE IDENTIQUE)
	if len(exercises) == 0 {
		component := pages.NoExercisesToday(report)
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("❌ Render error: %v", err)
			http.Error(w, "Erreur affichage", http.StatusInternalServerError)
		}
		return
	}

	// 4. ✅ APPLIQUE LIMITE ÉNERGIE (LOGIQUE IDENTIQUE)
	exerciseIDs := make([]int, len(exercises))
	for i, ex := range exercises {
		exerciseIDs[i] = ex.ID
	}

	limitedIDs := session.LimitExercises(exerciseIDs, energyLevel)

	log.Printf("🔍 [SESSION] Limité à %d exercices (max=%d pour energy=%d)",
		len(limitedIDs),
		session.GetMaxExercises(energyLevel),
		energy,
	)

	// 5. CRÉE SESSION (LOGIQUE IDENTIQUE)
	sessionID, sessionData, err := sessionService.StartSession(energyLevel, limitedIDs)
	if err != nil {
		log.Printf("❌ StartSession failed: %v", err)
		http.Error(w, "Erreur création session", http.StatusInternalServerError)
		return
	}

	// 6. REDIRIGE vers premier exercice (LOGIQUE IDENTIQUE)
	if len(sessionData.Exercises) == 0 {
		log.Printf("⚠️ Session created but no exercises")
		component := pages.NoExercisesToday(report)
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("❌ Render error: %v", err)
			http.Error(w, "Erreur affichage", http.StatusInternalServerError)
		}
		return
	}

	firstExerciseID := sessionData.Exercises[0]
	redirectURL := fmt.Sprintf("/exercise/%d?from=session&session=%d",
		firstExerciseID, sessionID)

	log.Printf("🚀 Session %d started → exo #%d", sessionID, firstExerciseID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// ============================================
// 3️⃣ SESSION COMPLETE (Résultats)
// ============================================

func HandleSessionComplete(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.URL.Query().Get("id")

	// Si pas d'ID fourni, essaie de récupérer la dernière session active (LOGIQUE IDENTIQUE)
	if sessionIDStr == "" {
		sessionID, err := sessionService.GetActiveSession()
		if err != nil || sessionID == 0 {
			log.Printf("⚠️ No active session found")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		sessionIDStr = fmt.Sprintf("%d", sessionID)
	}

	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("❌ Invalid session ID: %s", sessionIDStr)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Récupère résultat depuis service (LOGIQUE IDENTIQUE)
	result, err := sessionService.GetSessionResult(sessionID)
	if err != nil {
		log.Printf("❌ GetSessionResult failed: %v", err)

		// Fallback : affiche page vide
		component := pages.SessionCompleteEmpty()
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("❌ Render error: %v", err)
			http.Error(w, "Erreur affichage", http.StatusInternalServerError)
		}
		return
	}

	// ✅ CHANGEMENT : Render avec templ
	component := pages.SessionComplete(
		sessionID,
		result.CompletedCount,
		result.Duration.Round(time.Second),
		result.CompletedAt,
		result.Exercises,
	)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// 4️⃣ SESSION STOP (Arrêt manuel)
// ============================================

func HandleStopSession(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.PathValue("id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("❌ Invalid session ID: %s", sessionIDStr)
		http.Error(w, "ID de session invalide", http.StatusBadRequest)
		return
	}

	// Termine la session (LOGIQUE IDENTIQUE)
	if err := sessionService.StopSession(sessionID); err != nil {
		log.Printf("❌ StopSession failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Session %d stopped manually", sessionID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
