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

	// Création du template avec fonctions custom (optionnel)
	Tmpl = template.New("").Funcs(template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"upper": strings.ToUpper,
		"add": func(a, b int) int {
			return a + b
		},
		"repeat": strings.Repeat,
		"lower":  strings.ToLower,
	})

	Tmpl, err = Tmpl.ParseGlob("templates/pages/*.html")
	if err != nil {
		return fmt.Errorf("pages: %w", err)
	}

	Tmpl, err = Tmpl.ParseGlob("templates/partials/*.html")
	if err != nil {
		return fmt.Errorf("partials: %w", err)
	}

	Tmpl, err = Tmpl.ParseGlob("templates/components/*.html")
	if err != nil {
		return fmt.Errorf("components: %w", err)
	}

	return nil
}
