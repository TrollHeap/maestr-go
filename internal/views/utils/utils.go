package utils

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"maestro/internal/models"
)

func GetDomainBadgeClasses(domain string) string {
	base := "text-xs border bg-slate-900/80 border-slate-600 text-slate-200"

	// Tu peux spécialiser par domaine plus tard si besoin
	return base
}

// Helper function
func CalculateProgressPercent(completed, total int) int {
	if total == 0 {
		return 0
	}
	return (completed * 100) / total
}

func GetExerciseStatusBadgeClass(done bool) string {
	if done {
		return "border border-emerald-500/60 text-emerald-300 bg-emerald-500/10"
	}
	return "border border-amber-500/60 text-amber-300 bg-amber-500/10"
}

func GetDomainBadgeClass(domain string) string {
	// Tu peux spécialiser par domaine si besoin
	return "border border-slate-600 text-slate-200 bg-slate-900/80"
}

func GetDifficultyBadgeClass(difficulty int) string {
	base := "border text-xs"

	switch difficulty {
	case 1, 2:
		return base + " border-emerald-500/60 text-emerald-300 bg-emerald-500/10"
	case 3:
		return base + " border-amber-500/70 text-amber-300 bg-amber-500/10"
	default:
		return base + " border-rose-500/70 text-rose-300 bg-rose-500/10"
	}
}

func GetProgressTextColor(completed, total int) string {
	percent := CalculateProgressPercent(completed, total)

	switch {
	case percent == 100:
		// Vert émeraude (100% comme ton badge "Maîtrisé")
		return "text-emerald-300"
	case percent >= 67:
		// Turquoise (proche de la fin, même teinte que ta barre)
		return "text-sky-300"
	case percent > 0:
		// Jaune chaud (en cours, cohérent avec les steps non finis)
		return "text-amber-300"
	default:
		return "text-emerald-300"
	}
}

func BuildReviewURL(exerciseID, quality int, fromSession bool, sessionID string) string {
	url := fmt.Sprintf("/exercise/%d/review?quality=%d", exerciseID, quality)
	if fromSession && sessionID != "" {
		url += fmt.Sprintf("&from=session&session=%s", sessionID)
	}
	return url
}

func GetReviewButtonClassesGradient(quality int) string {
	switch quality {
	case 0:
		return "border-rose-500/60 bg-gradient-to-br from-rose-950/70 to-slate-950/70 hover:from-rose-900/80 hover:to-slate-900/80 text-rose-100"
	case 1:
		return "border-orange-500/60 bg-gradient-to-br from-orange-950/70 to-slate-950/70 hover:from-orange-900/80 hover:to-slate-900/80 text-orange-100"
	case 2:
		return "border-amber-500/60 bg-gradient-to-br from-amber-950/70 to-slate-950/70 hover:from-amber-900/80 hover:to-slate-900/80 text-amber-100"
	case 3:
		return "border-emerald-500/60 bg-gradient-to-br from-emerald-950/70 to-slate-950/70 hover:from-emerald-900/80 hover:to-slate-900/80 text-emerald-100"
	default:
		return "border-slate-700 bg-slate-900/70 hover:bg-slate-800/80 text-slate-100"
	}
}

func GetIntervalBadgeClasses(quality int) string {
	switch quality {
	case 0:
		return "bg-rose-500/20 text-rose-200 border border-rose-500/40"
	case 1:
		return "bg-orange-500/20 text-orange-200 border border-orange-500/40"
	case 2:
		return "bg-amber-500/20 text-amber-200 border border-amber-500/40"
	case 3:
		return "bg-emerald-500/20 text-emerald-200 border border-emerald-500/40"
	default:
		return "bg-slate-700/40 text-slate-300 border border-slate-600"
	}
}

func GetStatusBadgeClassesGlow(done bool) string {
	if done {
		// Maîtrisé : vert avec glow fort
		return "border border-emerald-500/70 text-emerald-200 bg-emerald-500/15 shadow-[0_0_12px_rgba(16,185,129,0.5)] hover:shadow-[0_0_20px_rgba(16,185,129,0.7)] hover:border-emerald-400"
	}
	// En cours : jaune avec glow fort
	return "border border-amber-500/70 text-amber-200 bg-amber-500/15 shadow-[0_0_12px_rgba(251,191,36,0.5)] hover:shadow-[0_0_20px_rgba(251,191,36,0.7)] hover:border-amber-400"
}

// Helper function
func IsStepCompleted(completedSteps []int, index int) bool {
	return slices.Contains(completedSteps, index)
}

func GetStepItemClasses(completedSteps []int, index int) string {
	if IsStepCompleted(completedSteps, index) {
		return "border-emerald-600/40 bg-emerald-950/20 hover:border-emerald-500/60"
	}
	return "border-slate-700 bg-slate-900/40 hover:border-slate-600"
}

func GetCheckboxClasses(completedSteps []int, index int) string {
	if IsStepCompleted(completedSteps, index) {
		return "border-emerald-500 bg-emerald-500/20"
	}
	return "border-slate-600 bg-slate-900/60 peer-hover:border-emerald-400/60"
}

func GetStepNumberClasses(completedSteps []int, index int) string {
	if IsStepCompleted(completedSteps, index) {
		return "text-emerald-400"
	}
	return "text-slate-500"
}

func GetStepTextClasses(completedSteps []int, index int) string {
	if IsStepCompleted(completedSteps, index) {
		return "text-emerald-200 line-through decoration-emerald-400/50"
	}
	return "text-slate-200"
}

func GetStatusBadgeClasses(done bool) string {
	if done {
		return "border border-emerald-500/60 text-emerald-300 bg-emerald-500/10"
	}
	return "border border-amber-500/60 text-amber-300 bg-amber-500/10"
}

func GetDifficultyBadgeClasses(difficulty int) string {
	base := "border text-xs"

	switch difficulty {
	case 1, 2:
		return base + " border-emerald-500/60 text-emerald-300 bg-emerald-500/10"
	case 3:
		return base + " border-amber-500/70 text-amber-300 bg-amber-500/10"
	default:
		return base + " border-rose-500/70 text-rose-300 bg-rose-500/10"
	}
}

func GetStatCardBorderClass(cssClass string) string {
	switch cssClass {
	case "urgent":
		return "border-rose-500/40"
	case "today":
		return "border-sky-500/40"
	case "upcoming":
		return "border-amber-500/40"
	case "active":
		return "border-emerald-500/40"
	case "new":
		return "border-purple-500/40"
	default:
		return "border-slate-700"
	}
}

func GetStatCardTextClass(cssClass string) string {
	switch cssClass {
	case "urgent":
		return "text-rose-400"
	case "today":
		return "text-sky-400"
	case "upcoming":
		return "text-amber-400"
	case "active":
		return "text-emerald-400"
	case "new":
		return "text-purple-400"
	default:
		return "text-slate-300"
	}
}

func GetEnergyCardClass(level int) string {
	base := "group relative rounded-2xl border bg-slate-900/70 backdrop-blur-xl " +
		"p-5 shadow-lg shadow-slate-900/40 cursor-pointer " +
		"transition-all duration-200 hover:-translate-y-1 hover:shadow-2xl"

	switch level {
	case 1:
		// Faible : bleu/teal doux
		return base + " border-emerald-500/30 hover:border-emerald-400 " +
			"hover:bg-gradient-to-br hover:from-slate-900 hover:to-emerald-900/40"
	case 2:
		// Moyen : ambre/orange
		return base + " border-amber-500/30 hover:border-amber-400 " +
			"hover:bg-gradient-to-br hover:from-slate-900 hover:to-amber-900/40"
	case 3:
		// Élevé : rouge
		return base + " border-rose-500/30 hover:border-rose-400 " +
			"hover:bg-gradient-to-br hover:from-slate-900 hover:to-rose-900/40"
	default:
		return base + " border-slate-700 hover:border-slate-500"
	}
}

// Helper
func FormatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	if minutes > 0 {
		return fmt.Sprintf("%d min %d s", minutes, seconds)
	}
	return fmt.Sprintf("%d s", seconds)
}

func GetFormTitle(ex *models.Exercise) string {
	if ex == nil {
		return "Nouvel exercice"
	}
	return "Éditer exercice"
}

func GetFormAction(ex *models.Exercise) string {
	if ex == nil {
		return "/exercises/create"
	}
	return "/exercise/" + strconv.Itoa(ex.ID) + "/update"
}

func GetStringValue(ex *models.Exercise, field string) string {
	if ex == nil {
		return ""
	}

	switch field {
	case "title":
		return ex.Title
	case "description":
		return ex.Description
	case "content":
		return ex.Content
	case "mnemonic":
		return ex.Mnemonic
	default:
		return ""
	}
}

func GetStepsValue(ex *models.Exercise) string {
	if ex == nil || len(ex.Steps) == 0 {
		return ""
	}
	return strings.Join(ex.Steps, "\n")
}

func GetVisualsValue(ex *models.Exercise) string {
	if ex == nil || len(ex.ConceptualVisuals) == 0 {
		return ""
	}

	var result []string
	for _, visual := range ex.ConceptualVisuals {
		visualText := visual.Content
		if visual.Caption != "" {
			visualText += "\nCaption: " + visual.Caption
		}
		result = append(result, visualText)
	}

	return strings.Join(result, "\n---\n")
}

func IsDomainSelected(ex *models.Exercise, domain string) bool {
	return ex != nil && ex.Domain == domain
}

func IsDifficultySelected(ex *models.Exercise, difficulty int) bool {
	return ex != nil && ex.Difficulty == difficulty
}

func GetCancelURL(ex *models.Exercise) string {
	if ex == nil {
		return "/exercises"
	}
	return "/exercise/" + strconv.Itoa(ex.ID)
}

// GetFilterPillClass retourne les classes Tailwind pour un filtre "pill".
func GetFilterPillClass(active bool) string {
	base := "inline-flex items-center gap-2 rounded-full px-4 py-1.5 " +
		"text-[0.65rem] font-mono tracking-widest uppercase transition-colors"

	if active {
		return base + " bg-slate-800 text-sky-300 border border-sky-500/60"
	}

	return base + " text-slate-300 border border-slate-700 hover:bg-slate-800"
}

// GetFilterDropdownSelectClass : style du <select> pour les filtres.
func GetFilterDropdownSelectClass() string {
	return "rounded-full border border-slate-700 bg-slate-900/80 px-3 py-1.5 " +
		"text-[0.7rem] font-mono tracking-widest uppercase text-slate-200 " +
		"hover:bg-slate-800 focus:outline-none focus:border-sky-500 " +
		"focus:ring-1 focus:ring-sky-500/40"
}

func BuildFilterURL(basePath, param, value, q, status, domain, difficulty, sort string) string {
	v := url.Values{}
	if q != "" {
		v.Set("q", q)
	}
	if status != "" {
		v.Set("status", status)
	}
	if difficulty != "" {
		v.Set("difficulty", difficulty)
	}
	if sort != "" {
		v.Set("sort", sort)
	}
	// domaine géré via param/value
	if param != "domain" && domain != "" {
		v.Set("domain", domain)
	}
	if value != "" {
		v.Set(param, value)
	}
	u := basePath
	if encoded := v.Encode(); encoded != "" {
		u += "?" + encoded
	}
	return u
}
