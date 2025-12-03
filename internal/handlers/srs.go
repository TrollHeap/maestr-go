package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"maestro/internal/domain/exercise"
	"maestro/internal/domain/srs"
	"maestro/internal/store"
	"maestro/internal/views/components"
	"maestro/internal/views/pages"
)

func HandleReview(w http.ResponseWriter, r *http.Request) {
	// 1. Parse params (LOGIQUE IDENTIQUE)
	id, _ := strconv.Atoi(r.PathValue("id"))
	quality, _ := strconv.Atoi(r.URL.Query().Get("quality"))
	fromSession := r.URL.Query().Get("from") == "session"
	sessionIDStr := r.URL.Query().Get("session")

	log.Printf("🔍 [Review] id=%d, quality=%d, fromSession=%v, sessionID=%s",
		id, quality, fromSession, sessionIDStr)

	// 2. Validation ID (LOGIQUE IDENTIQUE)
	if err := exercise.ValidateID(id); err != nil {
		log.Printf("❌ Validation ID failed: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 3. Validation quality (LOGIQUE IDENTIQUE)
	if err := exercise.ValidateQuality(quality); err != nil {
		log.Printf("❌ Validation Quality failed: %v", err)
		http.Error(w, "Quality invalide", http.StatusBadRequest)
		return
	}

	// 4. Applique algorithme SRS (LOGIQUE IDENTIQUE)
	ex, err := exerciseService.ReviewExercise(id, srs.ReviewQuality(quality))
	if err != nil {
		log.Printf("❌ ReviewExercise error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Review applied: ease=%.2f, nextReview=%s",
		ex.EaseFactor, ex.NextReviewAt.Format("2006-01-02"))

	// 5. Marque DONE si quality >= 1 (LOGIQUE IDENTIQUE)
	if quality >= 1 {
		ex.Done = true
		if err := store.SaveExercise(ex); err != nil {
			log.Printf("❌ SaveExercise error: %v", err)
			http.Error(w, "Erreur sauvegarde", http.StatusInternalServerError)
			return
		}
		log.Printf("✅ Exercise marked DONE")
	}

	// 6. MODE SESSION : Flow exercice suivant (LOGIQUE IDENTIQUE)
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
		} else {
			// → Session terminée
			log.Println("✅ Session complete, no more exercises")

			if err := sessionService.EndSession(sessionID); err != nil {
				log.Printf("❌ EndSession error: %v", err)
			}

			w.Header().Set("HX-Redirect", fmt.Sprintf("/session/complete?id=%d", sessionID))
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// 7. MODE LIBRE : soit fragment HTMX, soit full page
	log.Println("🔄 Free mode, reload detail")

	// Requête HTMX ? (clic sur bouton Review avec hx-post)
	if r.Header.Get("HX-Request") == "true" {
		// On renvoie uniquement le panneau Review
		component := components.ReviewPanel(*ex, fromSession, sessionIDStr)
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("❌ Render error: %v", err)
			http.Error(w, "Erreur affichage", http.StatusInternalServerError)
		}
		return
	}

	// Sinon, navigation classique → full page
	component := pages.ExerciseDetail(*ex, fromSession, sessionIDStr)
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// func HandleNextExercise(w http.ResponseWriter, r *http.Request) {
// 	log.Println("🔍 HandleNextExercise: Free mode")
//
// 	// 1. Récupère rapport + exercices disponibles (LOGIQUE IDENTIQUE)
// 	report, exercises, err := store.GetTodayReport()
// 	if err != nil {
// 		log.Printf("❌ GetTodayReport error: %v", err)
// 		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
// 		return
// 	}
//
// 	// 2. Aucun exercice disponible (LOGIQUE IDENTIQUE)
// 	if len(exercises) == 0 {
// 		log.Println("ℹ️ No exercises available, showing report")
//
// 		// ✅ CHANGEMENT : Render avec templ
// 		component := pages.NoExercisesToday(report)
//
// 		if err := component.Render(r.Context(), w); err != nil {
// 			log.Printf("❌ Render error: %v", err)
// 			http.Error(w, "Erreur affichage", http.StatusInternalServerError)
// 		}
// 		return
// 	}
//
// 	// 3. Exercice(s) disponible(s) (LOGIQUE IDENTIQUE)
// 	nextExercise := exercises[0]
// 	redirectURL := fmt.Sprintf("/exercise/%d", nextExercise.ID)
//
// 	log.Printf("➡️ Redirect to exercise #%d: %s", nextExercise.ID, nextExercise.Title)
// 	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
// }
