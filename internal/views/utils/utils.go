package utils

import (
	"slices"
	"strconv"
	"strings"

	"maestro/internal/models"
)

// Helper function

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

// Helper

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
