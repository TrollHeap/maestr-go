package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"maestro/internal/domain/exercise"
	"maestro/internal/domain/srs"
	"maestro/internal/store"
)

func HandleReview(w http.ResponseWriter, r *http.Request) {
	// 1. Parse params
	id, _ := strconv.Atoi(r.PathValue("id"))
	quality, _ := strconv.Atoi(r.URL.Query().Get("quality"))
	fromSession := r.URL.Query().Get("from") == "session"
	sessionIDStr := r.URL.Query().Get("session")

	log.Printf("🔍 [Review] id=%d, quality=%d, fromSession=%v, sessionID=%s",
		id, quality, fromSession, sessionIDStr)

	// 2. Validation ID
	if err := exercise.ValidateID(id); err != nil {
		log.Printf("❌ Validation ID failed: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 3. Validation quality
	if err := exercise.ValidateQuality(quality); err != nil {
		log.Printf("❌ Validation Quality failed: %v", err)
		http.Error(w, "Quality invalide", http.StatusBadRequest)
		return
	}

	// 4. Applique SRS
	ex, err := exerciseService.ReviewExercise(id, srs.ReviewQuality(quality))
	if err != nil {
		log.Printf("❌ ReviewExercise error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Review applied: ease=%.2f, nextReview=%s",
		ex.EaseFactor, ex.NextReviewAt.Format("2006-01-02"))

	// 5. Marque DONE si quality >= 1
	if quality >= 1 {
		ex.Done = true
		if err := store.SaveExercise(ex); err != nil {
			log.Printf("❌ SaveExercise error: %v", err)
			http.Error(w, "Erreur sauvegarde", http.StatusInternalServerError)
			return
		}
		log.Printf("✅ Exercise marked DONE")
	}

	// 6. MODE SESSION uniquement
	if fromSession && sessionIDStr != "" {
		sessionID, _ := strconv.ParseInt(sessionIDStr, 10, 64)

		log.Printf("🔄 Session mode: sessionID=%d", sessionID)

		// a) Enregistre dans session
		if err := sessionService.CompleteExercise(sessionID, id, quality); err != nil {
			log.Printf("❌ CompleteExercise error: %v", err)
		} else {
			log.Printf("✅ Exercise completed in session")
		}

		// b) Prochain exercice
		nextEx, err := sessionService.GetNextExercise(sessionID)
		if err != nil {
			log.Printf("❌ GetNextExercise error: %v", err)
		}

		if nextEx != nil {
			// → Redirection HTMX vers exercice suivant
			redirectURL := fmt.Sprintf("/exercise/%d?from=session&session=%d",
				nextEx.ID, sessionID)
			log.Printf("➡️ HX-Redirect to: %s", redirectURL)

			w.Header().Set("HX-Redirect", redirectURL)
			w.WriteHeader(http.StatusOK)
			return
		}

		// → Session terminée
		log.Println("✅ Session complete, no more exercises")

		if err := sessionService.EndSession(sessionID); err != nil {
			log.Printf("❌ EndSession error: %v", err)
		}

		w.Header().Set("HX-Redirect", fmt.Sprintf("/session/complete?id=%d", sessionID))
		w.WriteHeader(http.StatusOK)
		return
	}

	// 7. Si quelqu'un appelle /review sans from=session → 400
	log.Println("⚠️ Review appelée hors mode session")
	http.Error(w, "Review disponible uniquement en mode session", http.StatusBadRequest)
}
