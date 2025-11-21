# Cheatsheet Pédagogique Go + Templates (sans HTMX)


***

## 1. Chargement des templates

```go
tmpl := template.Must(template.ParseGlob("templates/*.html"))
```

- Charge tous les fichiers `*.html` définis dans `templates/`
- `Must` panique s’il y a une erreur

***

## 2. Exécution d’un template

```go
tmpl.Execute(w, data)                  // Rendu principal
tmpl.ExecuteTemplate(w, "tmplName", data) // Template nommé spécifique
```

- `w` est `http.ResponseWriter` (réponse HTTP)
- `data` est la variable contextuelle (map, struct, slice...)

***

## 3. Définir un template nommé (composant réutilisable)

```gohtml
{{define "list"}}
<ul>
{{range .}}
  <li>{{.}}</li>
{{end}}
</ul>
{{end}}
```

- Définit un bloc HTML réutilisable nommé `"list"`
- La variable contextuelle `.` est la donnée transmise (ici une slice)

***

## 4. Inclure un template dans un autre

Dans un template parent :

```gohtml
{{template "list" .Names}}
```

- Insère le template `"list"` avec `Names` comme contexte
- Favorise la modularité et évite la duplication

***

## 5. Utilisation des variables dans templates

- `.` : variable courante (globale ou dans la portée)
- `{{.Field}}` : accès à un champ (struct) ou clé (map)
- `{{range .}} ... {{end}}` : boucle sur slice/map avec `.` comme élément courant
- `{{if .Condition}} ... {{else}} ... {{end}}` : conditionnelle

***

## 6. Exemple minimal complet:

### `main.go`

```go
package main

import (
    "html/template"
    "net/http"
)

func main() {
    tmpl := template.Must(template.ParseGlob("templates/*.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := map[string]any{
            "Title": "Ma page",
            "Names": []string{"Alice", "Bob", "Charlie"},
        }
        tmpl.ExecuteTemplate(w, "base", data)
    })

    http.ListenAndServe(":8080", nil)
}
```


### `templates/base.html`

```gohtml
{{define "base"}}
<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
<h1>{{.Title}}</h1>
{{template "list" .Names}}
</body>
</html>
{{end}}
```


### `templates/list.html`

```gohtml
{{define "list"}}
<ul>
{{range .}}
<li>{{.}}</li>
{{end}}
</ul>
{{end}}
```


***

## Résumé de la logique

- Le **serveur Go** charge des **templates HTML avec des noms**.
- Les **données Go** (clés, slices, structs) sont transmises aux templates.
- Le template utilise `{{...}}` pour afficher ou boucler sur ces données.
- `{{define}}` permet de créer des morceaux réutilisables.
- Les templates sont combinés via `{{template}}`, favorisant la structure claire et modulaire.
