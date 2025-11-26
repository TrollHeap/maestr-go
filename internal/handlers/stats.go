package handlers

import (
	"net/http"

	"maestro/v2-refacto/internal/models"
)

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// PAGE PRINCIPALE STATS
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func HandleStatsPage(w http.ResponseWriter, r *http.Request) {
	// Rend le layout complet (squelette avec zones HTMX)
	if err := Tmpl.ExecuteTemplate(w, "stats", nil); err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// FRAGMENT MÉTRIQUES (appelé par HTMX)
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func HandleStatsMetrics(w http.ResponseWriter, r *http.Request) {
	// Récupère les exercices (mock pour l'instant)
	exercises := getMockExercises()

	// Calcule les stats
	total := len(exercises)
	completed := 0
	wip := 0
	for _, ex := range exercises {
		if ex.Done {
			completed++
		} else if ex.ID%3 == 0 { // Mock WIP logic
			wip++
		}
	}
	todo := total - completed - wip
	completionRate := 0
	if total > 0 {
		completionRate = (completed * 100) / total
	}

	data := models.StatsMetrics{
		Total:          total,
		Completed:      completed,
		WIP:            wip,
		Todo:           todo,
		CompletionRate: completionRate,
	}

	// Rend UNIQUEMENT le fragment métriques
	if err := Tmpl.ExecuteTemplate(w, "stats-metrics", data); err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// FRAGMENT DOMAINES (appelé par HTMX)
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func HandleStatsDomains(w http.ResponseWriter, r *http.Request) {
	exercises := getMockExercises()

	// Calcule stats par domaine
	domainsMap := make(map[string]*models.DomainStat)
	for _, ex := range exercises {
		if domainsMap[ex.Domain] == nil {
			domainsMap[ex.Domain] = &models.DomainStat{
				Name: ex.Domain,
			}
		}
		domainsMap[ex.Domain].Total++
		if ex.Done {
			domainsMap[ex.Domain].Completed++
		}
	}

	// Calcule les pourcentages
	maxCount := 0
	for _, stat := range domainsMap {
		if stat.Total > maxCount {
			maxCount = stat.Total
		}
	}

	var domains []models.DomainStat
	for _, stat := range domainsMap {
		if maxCount > 0 {
			stat.Percentage = (stat.Total * 100) / maxCount
		}
		domains = append(domains, *stat)
	}

	data := struct {
		Domains []models.DomainStat
	}{
		Domains: domains,
	}

	// Rend UNIQUEMENT le fragment domaines
	if err := Tmpl.ExecuteTemplate(w, "stats-domains", data); err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// FRAGMENT DIFFICULTÉS (appelé par HTMX)
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func HandleStatsDifficulties(w http.ResponseWriter, r *http.Request) {
	exercises := getMockExercises()

	// Calcule stats par difficulté
	counts := make(map[int]int)
	for _, ex := range exercises {
		counts[ex.Difficulty]++
	}

	symbols := []string{"░░░", "▒▒▒", "▓▓▓", "███", "▉▉▉"}
	labels := []string{"FACILE", "MOYEN", "DIFFICILE", "EXPERT", "EXTRÊME"}

	var difficulties []models.DifficultyStat
	for level := 1; level <= 5; level++ {
		difficulties = append(difficulties, models.DifficultyStat{
			Level:  level,
			Symbol: symbols[level-1],
			Label:  labels[level-1],
			Count:  counts[level],
		})
	}

	data := struct {
		Difficulties []models.DifficultyStat
	}{
		Difficulties: difficulties,
	}

	// Rend UNIQUEMENT le fragment difficultés
	if err := Tmpl.ExecuteTemplate(w, "stats-difficulties", data); err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// FONCTIONS HELPER
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

// Mock data (à remplacer par DB)
func getMockExercises() []models.Exercise {
	return []models.Exercise{
		{ID: 1, Title: "Gestion Erreurs", Domain: "Go", Difficulty: 1, Done: true},
		{ID: 2, Title: "Validation Inputs", Domain: "Go", Difficulty: 2, Done: false},
		{ID: 3, Title: "Types Personnalisés", Domain: "Go", Difficulty: 3, Done: false},
		{ID: 4, Title: "SplitSeq Parsing", Domain: "Go", Difficulty: 4, Done: false},
		{ID: 5, Title: "Tri Rapide", Domain: "Algorithmes", Difficulty: 3, Done: true},
		{ID: 6, Title: "Arbre Binaire", Domain: "Algorithmes", Difficulty: 4, Done: false},
		{ID: 7, Title: "HTMX Filtres", Domain: "HTMX", Difficulty: 2, Done: true},
		{ID: 8, Title: "Templates Go", Domain: "HTMX", Difficulty: 2, Done: false},
		{ID: 9, Title: "Système Fichiers", Domain: "Linux", Difficulty: 3, Done: false},
		{ID: 10, Title: "Architecture MVC", Domain: "Architecture", Difficulty: 5, Done: false},
	}
}
