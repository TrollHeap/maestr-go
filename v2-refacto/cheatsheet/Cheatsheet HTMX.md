# Cheatsheet HTMX — Rappels et Syntaxe clés


***

## 1. Inclusion de la bibliothèque HTMX

Dans ton fichier HTML principal, ajoute :

```html
<script src="https://unpkg.com/htmx.org@1.9.10"></script>
```


***

## 2. Attributs principaux d’HTMX

| Attribut | Description | Exemple |
| :-- | :-- | :-- |
| `hx-get` | Envoie une requête HTTP GET vers une URL | `<button hx-get="/list" hx-target="#out"></button>` |
| `hx-post` | Envoie une requête HTTP POST | `<form hx-post="/save"></form>` |
| `hx-target` | Détermine l’élément HTML à mettre à jour avec la réponse | `hx-target="#container"` |
| `hx-swap` | Détermine la façon dont la réponse est insérée (`innerHTML`, `outerHTML`, `beforeend`, etc) | `hx-swap="innerHTML"` |
| `hx-trigger` | Déclenchement de la requête (clic, change, délai, etc) | `hx-trigger="click"` |
| `hx-include` | Inclut des éléments dans la requête (ex: champs formulaire) | `hx-include="#input1"` |


***

## 3. Exemple simple : charger du contenu à la demande

HTML :

```html
<button hx-get="/list" hx-target="#container" hx-swap="innerHTML">
  Charger la liste
</button>

<div id="container"></div>
```

- Au clic sur le bouton, HTMX fait un `GET /list`
- La réponse (HTML partiel) remplace le contenu de `#container`

***

## 4. Exemple avec formulaire et mise à jour partielle

```html
<form hx-post="/submit" hx-target="#result" hx-swap="innerHTML">
  <input type="text" name="name" />
  <button type="submit">Envoyer</button>
</form>

<div id="result"></div>
```

- Le formulaire POST `/submit`
- Le serveur renvoie un fragment HTML pour remplacer `#result`

***

## 5. Étapes pour usage en Go

- Côté serveur, crée un handler HTTP qui renvoie **du HTML partiel** (ex: rendu d’un template nommé).
- Dans le template, place les attributs HTMX.
- Lors des interactions utilisateurs, HTMX se charge d’envoyer/reçevoir, et mettre à jour le DOM automatiquement.

***

## 6. Résumé visuel

| Action utilisateur | Action HTMX déclenchée | Requête HTTP | Réponse attendue | Mise à jour DOM |
| :-- | :-- | :-- | :-- | :-- |
| Clic sur bouton | hx-get | GET /url | HTML partiel | Remplace contenu ciblé |
| Envoi formulaire | hx-post | POST /url | HTML partiel | Remplace contenu ciblé |


***

## 7. Conseils pratiques

- Utilise HTMX pour **simplifier le JavaScript** et décharger côté serveur.
- Rends le backend capable de produire des **fragments HTML réutilisables** via Go templates.
