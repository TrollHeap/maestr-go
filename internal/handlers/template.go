package handlers

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

var Tmpl *template.Template

func InitTemplates() error {
	var err error

	Tmpl = template.New("").Funcs(template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("02 Jan 2006")
		},
		"upper": strings.ToUpper,
		"add": func(a, b int) int {
			return a + b
		},
		"repeat": strings.Repeat,
		"lower":  strings.ToLower,
		// ✨ NOUVEAU : Génère range [0, n-1]
		"until": func(n int) []int {
			result := make([]int, n)
			for i := range n {
				result[i] = i
			}
			return result
		},
	})

	Tmpl, err = Tmpl.ParseGlob("templates/pages/*.html")
	if err != nil {
		return fmt.Errorf("pages: %w", err)
	}

	Tmpl, err = Tmpl.ParseGlob("templates/partials/*.html")
	if err != nil {
		return fmt.Errorf("partials: %w", err)
	}

	Tmpl, err = Tmpl.ParseGlob("templates/components/*/*.html")
	if err != nil {
		return fmt.Errorf("components: %w", err)
	}

	return nil
}
