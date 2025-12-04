package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"maestro/internal/domain/exercise"
	"maestro/internal/models"
	"maestro/internal/service"
	"maestro/internal/views/components"
	"maestro/internal/views/pages"
)

// ============================================
// SERVICES GLOBAUX
// ============================================

var exerciseService *service.ExerciseService

func init() {
	exerciseService = service.NewExerciseService()
	sessionService = service.NewSessionService()
}

func HandleExercisesPage(w http.ResponseWriter, r *http.Request) {
	// 1. Parse filtres contenu
	q := r.URL.Query().Get("q")
	status := r.URL.Query().Get("status")                          // "in_progress", "mastered", ""
	domain := r.URL.Query().Get("domain")                          // "Go", "Algorithms", ""
	difficulty, _ := strconv.Atoi(r.URL.Query().Get("difficulty")) // 1-4, 0
	sort := r.URL.Query().
		Get("sort")
		// "", "title", "difficulty", "domain"

	filter := models.ExerciseFilter{
		Status:     status,
		Domain:     domain,
		Difficulty: difficulty,
		Query:      q,
		Sort:       sort,
	}

	// 2. Récupère exercices filtrés
	exercises, err := exerciseService.GetFilteredExercises(filter)
	if err != nil {
		log.Printf("❌ GetFilteredExercises error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// 3. Compte total (sans filtre) pour l’info "X / Y"
	allExercises, err := exerciseService.GetAllExercises()
	if err != nil {
		log.Printf("❌ GetAllExercises error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	total := len(allExercises)

	// 4. Render page complète
	component := pages.ExerciseListPage(exercises, filter, total)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

func HandleListExercice(w http.ResponseWriter, r *http.Request) {
	// Même parsing que pour la page, mais ne renvoie que le fragment liste
	q := r.URL.Query().Get("q")
	status := r.URL.Query().Get("status")
	domain := r.URL.Query().Get("domain")
	difficulty, _ := strconv.Atoi(r.URL.Query().Get("difficulty"))
	sort := r.URL.Query().Get("sort")

	filter := models.ExerciseFilter{
		Status:     status,
		Domain:     domain,
		Difficulty: difficulty,
		Query:      q,
		Sort:       sort,
	}

	filteredList, err := exerciseService.GetFilteredExercises(filter)
	if err != nil {
		log.Printf("❌ GetFilteredExercises error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	component := components.ExerciseList(filteredList)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// 6️⃣ CRÉATION : Formulaire nouveau
// ============================================

func HandleExerciseCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("📝 ExerciseCreate: render form")

	// Render formulaire vide (mode création)
	component := pages.ExerciseForm(nil) // nil = nouveau

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// 7️⃣ CRÉATION : Soumission formulaire
// ============================================

func HandleExerciseSubmit(w http.ResponseWriter, r *http.Request) {
	log.Println("📝 ExerciseSubmit: processing form")

	// 1. Parse form
	if err := r.ParseForm(); err != nil {
		log.Printf("❌ Parse form error: %v", err)
		http.Error(w, "Erreur formulaire", http.StatusBadRequest)
		return
	}

	// 2. Parse difficulty
	difficulty, err := strconv.Atoi(r.FormValue("difficulty"))
	if err != nil {
		log.Printf("❌ Invalid difficulty: %v", err)
		http.Error(w, "Difficulté invalide", http.StatusBadRequest)
		return
	}

	// 3. Parse steps (textarea → slice)
	stepsText := r.FormValue("steps")
	steps := parseSteps(stepsText)

	// ✅ 4. Parse conceptual_visuals (NOUVEAU)
	visualsText := r.FormValue("conceptual_visuals")
	visuals := parseConceptualVisuals(visualsText)

	// 5. Build exercise model
	ex := &models.Exercise{
		Title:             r.FormValue("title"),
		Description:       r.FormValue("description"),
		Domain:            r.FormValue("domain"),
		Difficulty:        difficulty,
		Content:           r.FormValue("content"),
		Mnemonic:          r.FormValue("mnemonic"),
		Steps:             steps,
		ConceptualVisuals: visuals, // ✅ AJOUTÉ
	}

	// 6. Create via service
	if err := exerciseService.CreateExercise(ex); err != nil {
		log.Printf("❌ CreateExercise error: %v", err)

		component := components.FormError(err.Error())
		if renderErr := component.Render(r.Context(), w); renderErr != nil {
			http.Error(w, "Erreur création", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("✅ Exercise #%d created: %s", ex.ID, ex.Title)

	// 7. Redirect vers détail (HTMX)
	w.Header().Set("HX-Redirect", fmt.Sprintf("/exercise/%d", ex.ID))
	w.WriteHeader(http.StatusOK)
}

// DeleteExercise : soft delete via POST
func HandleExerciseDelete(w http.ResponseWriter, r *http.Request) {
	// 1. Parse ID (depuis path /exercise/{id}/delete par ex.)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("❌ Invalid ID: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	log.Printf("🗑️ ExerciseDelete: id=%d", id)

	// 2. Validation domaine
	if err := exercise.ValidateID(id); err != nil {
		log.Printf("❌ Validation error: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 3. Delete via service
	if err := exerciseService.DeleteExercise(id); err != nil {
		log.Printf("❌ DeleteExercise error: %v", err)
		component := components.FormError(err.Error())
		if renderErr := component.Render(r.Context(), w); renderErr != nil {
			http.Error(w, "Erreur suppression", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("✅ Exercise #%d deleted", id)

	// 4. Redirect vers liste (HTMX friendly)
	w.Header().Set("HX-Redirect", "/exercises")
	w.WriteHeader(http.StatusOK)
}

// ============================================
// 8️⃣ ÉDITION : Formulaire éditer
// ============================================

func HandleExerciseEdit(w http.ResponseWriter, r *http.Request) {
	// 1. Parse ID
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("❌ Invalid ID: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	log.Printf("✏️ ExerciseEdit: id=%d", id)

	// 2. Validation ID
	if err := exercise.ValidateID(id); err != nil {
		log.Printf("❌ Validation error: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 3. Récupère exercice existant
	ex, err := exerciseService.GetExerciseWithMarkdown(id)
	if err != nil {
		log.Printf("❌ Exercise #%d not found: %v", id, err)
		http.NotFound(w, r)
		return
	}

	// 4. Render formulaire pré-rempli (mode édition)
	component := pages.ExerciseForm(ex) // ex != nil = édition

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// 9️⃣ ÉDITION : Mise à jour
// ============================================

func HandleExerciseUpdate(w http.ResponseWriter, r *http.Request) {
	// 1-3. Parse ID, form, difficulty (INCHANGÉ)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("❌ Invalid ID: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	log.Printf("✏️ ExerciseUpdate: id=%d", id)

	if err := r.ParseForm(); err != nil {
		log.Printf("❌ Parse form error: %v", err)
		http.Error(w, "Erreur formulaire", http.StatusBadRequest)
		return
	}

	difficulty, err := strconv.Atoi(r.FormValue("difficulty"))
	if err != nil {
		log.Printf("❌ Invalid difficulty: %v", err)
		http.Error(w, "Difficulté invalide", http.StatusBadRequest)
		return
	}

	// 4. Parse steps
	stepsText := r.FormValue("steps")
	steps := parseSteps(stepsText)

	// ✅ 5. Parse conceptual_visuals (NOUVEAU)
	visualsText := r.FormValue("conceptual_visuals")
	visuals := parseConceptualVisuals(visualsText)

	// 6. Build exercise model (avec ID)
	ex := &models.Exercise{
		ID:                id,
		Title:             r.FormValue("title"),
		Description:       r.FormValue("description"),
		Domain:            r.FormValue("domain"),
		Difficulty:        difficulty,
		Content:           r.FormValue("content"),
		Mnemonic:          r.FormValue("mnemonic"),
		Steps:             steps,
		ConceptualVisuals: visuals, // ✅ AJOUTÉ
	}

	// 7. Update via service
	if err := exerciseService.UpdateExercise(ex); err != nil {
		log.Printf("❌ UpdateExercise error: %v", err)

		component := components.FormError(err.Error())
		if renderErr := component.Render(r.Context(), w); renderErr != nil {
			http.Error(w, "Erreur mise à jour", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("✅ Exercise #%d updated: %s", ex.ID, ex.Title)

	w.Header().Set("HX-Redirect", fmt.Sprintf("/exercise/%d", ex.ID))
	w.WriteHeader(http.StatusOK)
}

// ============================================
// 🔧 HELPER : Parse steps textarea
// ============================================

// parseSteps : Convertit textarea multi-lignes en slice
func parseSteps(text string) []string {
	if text == "" {
		return []string{}
	}

	lines := strings.Split(text, "\n")
	steps := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			steps = append(steps, line)
		}
	}

	return steps
}

// ============================================
// 1️⃣ PAGE PRINCIPALE EXERCICES
// ============================================

// ============================================
// 3️⃣ PAGE DÉTAIL EXERCICE
// ============================================

func HandleDetailExercice(w http.ResponseWriter, r *http.Request) {
	// 1. Parse params (LOGIQUE IDENTIQUE)
	id, _ := strconv.Atoi(r.PathValue("id"))
	fromSession := r.URL.Query().Get("from") == "session"
	sessionIDStr := r.URL.Query().Get("session")

	log.Printf("🔍 DetailExercice: id=%d, fromSession=%v, sessionID=%s",
		id, fromSession, sessionIDStr)

	// 2. Validation ID (LOGIQUE IDENTIQUE)
	if err := exercise.ValidateID(id); err != nil {
		log.Printf("❌ Invalid ID: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 3. Récupère exercice (LOGIQUE IDENTIQUE)
	ex, err := exerciseService.GetExerciseWithMarkdown(id)
	if err != nil {
		log.Printf("❌ Exercice #%d non trouvé: %v", id, err)
		http.NotFound(w, r)
		return
	}

	// 4. ✅ CHANGEMENT : Render avec templ (params typés)
	component := pages.ExerciseDetail(*ex, fromSession, sessionIDStr)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// 4️⃣ ACTION : Toggle Done
// ============================================

func HandleToggleDone(w http.ResponseWriter, r *http.Request) {
	// 1. Parse params (LOGIQUE IDENTIQUE)
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	fromSession := r.URL.Query().Get("from") == "session"
	sessionIDStr := r.URL.Query().Get("session")

	log.Printf("🔄 ToggleDone: id=%d, fromSession=%v, sessionID=%s",
		id, fromSession, sessionIDStr)

	// 2. Validation (LOGIQUE IDENTIQUE)
	if err := exercise.ValidateID(id); err != nil {
		log.Printf("❌ Validation ID failed: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 3. Toggle via service (LOGIQUE IDENTIQUE)
	ex, err := exerciseService.ToggleExerciseDone(id)
	if err != nil {
		log.Printf("❌ ToggleExerciseDone error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Exercice #%d toggled: Done=%v", id, ex.Done)

	// 4. MODE SESSION : Gestion flow (LOGIQUE IDENTIQUE)
	if fromSession && ex.Done && sessionIDStr != "" {
		sessionID, _ := strconv.ParseInt(sessionIDStr, 10, 64)

		// a) Marque exercice complété
		if err := sessionService.CompleteExercise(sessionID, id, 3); err != nil {
			log.Printf("❌ CompleteExercise error: %v", err)
		}

		// b) Récupère prochain exercice
		nextEx, err := sessionService.GetNextExercise(sessionID)
		if err != nil {
			log.Printf("❌ GetNextExercise error: %v", err)
		}

		if nextEx != nil {
			// → Redirige vers exercice suivant
			redirectURL := fmt.Sprintf("/exercise/%d?from=session&session=%d",
				nextEx.ID, sessionID)
			log.Printf("➡️ Redirect to: %s", redirectURL)

			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		} else {
			// → Session terminée
			log.Println("✅ Session complete, no more exercises")

			if err := sessionService.EndSession(sessionID); err != nil {
				log.Printf("❌ EndSession error: %v", err)
			}

			http.Redirect(w, r, fmt.Sprintf("/session/complete?id=%d", sessionID),
				http.StatusSeeOther)
			return
		}
	}

	// 5. ✅ CHANGEMENT : Render fragment templ
	component := components.StatusIndicator(*ex)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// 5️⃣ ACTION : Toggle Step
// ============================================

func HandleToggleStep(w http.ResponseWriter, r *http.Request) {
	// 1. Parse params (LOGIQUE IDENTIQUE)
	id, _ := strconv.Atoi(r.PathValue("id"))
	step, _ := strconv.Atoi(r.URL.Query().Get("step"))

	log.Printf("🔄 ToggleStep: id=%d, step=%d", id, step)

	// 2. Validation ID (LOGIQUE IDENTIQUE)
	if err := exercise.ValidateID(id); err != nil {
		log.Printf("❌ Invalid ID: %v", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 3. Récupère exercice (LOGIQUE IDENTIQUE)
	ex, err := exerciseService.GetExerciseWithMarkdown(id)
	if err != nil {
		log.Printf("❌ Exercise #%d not found: %v", id, err)
		http.NotFound(w, r)
		return
	}

	// 4. Validation step (LOGIQUE IDENTIQUE)
	if err := exercise.ValidateStep(step, len(ex.Steps)); err != nil {
		log.Printf("❌ Invalid step: %v", err)
		http.Error(w, "Step invalide", http.StatusBadRequest)
		return
	}

	// 5. Toggle step (LOGIQUE IDENTIQUE)
	ex, err = exerciseService.ToggleExerciseStep(id, step)
	if err != nil {
		log.Printf("❌ ToggleExerciseStep error: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Step #%d toggled", step)

	// 6. ✅ CHANGEMENT : Render fragment templ
	component := components.StepsFragmentWrapper(*ex)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// parseConceptualVisuals : Convertit textarea en []models.VisualAid
func parseConceptualVisuals(text string) []models.VisualAid {
	if text == "" {
		return []models.VisualAid{}
	}

	blocks := strings.Split(text, "---")
	visuals := make([]models.VisualAid, 0, len(blocks))

	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		var content, caption string

		lowerBlock := strings.ToLower(block)
		if idx := strings.Index(lowerBlock, "\ncaption:"); idx != -1 {
			content = strings.TrimSpace(block[:idx])
			caption = strings.TrimSpace(block[idx+9:]) // +9 = len("\nCaption:")
		} else {
			content = block
		}

		if content != "" {
			visuals = append(visuals, models.VisualAid{
				Type:    "ascii",
				Content: content,
				Caption: caption,
			})
		}
	}

	return visuals
}
