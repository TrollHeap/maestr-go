package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"
)

// ============================================
// VARIABLES GLOBALES
// ============================================

var (
	Tmpl *template.Template
	once sync.Once // ✅ Thread-safe init (une seule fois)
)

// ============================================
// INIT TEMPLATES (Thread-safe)
// ============================================

func InitTemplates() error {
	var err error

	once.Do(func() {
		// ✅ 1. Crée template avec FuncMap
		Tmpl = template.New("").Funcs(createFuncMap())

		// ✅ 2. Parse tous les patterns (accumulation)
		patterns := []string{
			"templates/layouts/*.html",
			"templates/pages/*.html",
			"templates/partials/*.html",
			"templates/components/*/*.html",
		}

		for _, pattern := range patterns {
			Tmpl, err = Tmpl.ParseGlob(pattern)
			if err != nil {
				err = fmt.Errorf("❌ pattern %s: %w", pattern, err)
				return
			}
		}

		// ✅ 3. Log templates chargés (dev/debug)
		if Tmpl != nil {
			templateNames := make([]string, 0)
			for _, t := range Tmpl.Templates() {
				templateNames = append(templateNames, t.Name())
			}
			log.Printf("✅ %d templates chargés: %v", len(templateNames), templateNames)
		}
	})

	return err
}

// ============================================
// FUNCMAP (helpers templates)
// ============================================

func createFuncMap() template.FuncMap {
	return template.FuncMap{
		// === DATE & TIME ===
		"formatDate": func(t time.Time) string {
			return t.Format("02 Jan 2006")
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("02/01/2006 15:04")
		},
		"daysUntil": func(t time.Time) int {
			diff := time.Until(t)
			days := int(diff.Hours() / 24)
			if days < 0 {
				return 0
			}
			return days
		},
		"isToday": func(t time.Time) bool {
			now := time.Now()
			return t.Year() == now.Year() && t.YearDay() == now.YearDay()
		},
		"isPast": func(t time.Time) bool {
			return t.Before(time.Now())
		},

		// === STRING MANIPULATION ===
		"upper":  strings.ToUpper,
		"lower":  strings.ToLower,
		"repeat": strings.Repeat,
		"trim":   strings.TrimSpace,

		// === MATH ===
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a % b
		},

		// === UTILITY ===
		"string": func(v any) string {
			return fmt.Sprintf("%v", v)
		},
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"until": func(n int) []int {
			result := make([]int, n)
			for i := range n {
				result[i] = i
			}
			return result
		},
		"seq": func(start, end int) []int {
			if start > end {
				return []int{}
			}
			result := make([]int, end-start+1)
			for i := range result {
				result[i] = start + i
			}
			return result
		},

		// === CONDITIONAL ===
		"default": func(defaultVal, val any) any {
			if val == nil || val == "" || val == 0 {
				return defaultVal
			}
			return val
		},
		"ternary": func(condition bool, trueVal, falseVal any) any {
			if condition {
				return trueVal
			}
			return falseVal
		},

		// === COLLECTIONS ===
		"len": func(v any) int {
			switch val := v.(type) {
			case string:
				return len(val)
			case []any:
				return len(val)
			case map[string]any:
				return len(val)
			default:
				return 0
			}
		},
		"contains": func(slice []string, val string) bool {
			return slices.Contains(slice, val)
		},
	}
}

// ============================================
// RENDER HELPERS
// ============================================

// RenderTemplate : Helper pour render avec error handling
func RenderTemplate(w http.ResponseWriter, name string, data any) error {
	// ✅ Cherche nom exact
	tmpl := Tmpl.Lookup(name)

	// ✅ Fallback : essaie avec .html
	if tmpl == nil && !strings.HasSuffix(name, ".html") {
		tmpl = Tmpl.Lookup(name + ".html")
	}

	if tmpl == nil {
		return fmt.Errorf("template '%s' non trouvé", name)
	}

	return Tmpl.ExecuteTemplate(w, name, data)
}

// RenderTemplateOrError : Render + log error + HTTP 500
func RenderTemplateOrError(w http.ResponseWriter, name string, data any) {
	if err := RenderTemplate(w, name, data); err != nil {
		log.Printf("❌ Erreur render template '%s': %v", name, err)
		http.Error(w, "Erreur affichage page", http.StatusInternalServerError)
	}
}

// ============================================
// DEBUG HELPERS (dev uniquement)
// ============================================

// ListTemplates : Retourne noms de tous les templates chargés
func ListTemplates() []string {
	if Tmpl == nil {
		return []string{}
	}

	names := make([]string, 0)
	for _, t := range Tmpl.Templates() {
		names = append(names, t.Name())
	}
	return names
}

// HasTemplate : Vérifie si template existe
func HasTemplate(name string) bool {
	if Tmpl == nil {
		return false
	}
	return Tmpl.Lookup(name) != nil
}
