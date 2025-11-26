# 🚀 Maestro Backend - Setup Complet

## Structure des Fichiers Créés

```
maestro/
├── go.mod                          # Modules Go
│
├── internal/
│   ├── models/
│   │   └── exercise.go            # Structs: Exercise, ReviewInput, Stats, etc.
│   │
│   ├── storage/
│   │   ├── store.go               # Interface Store (abstraction)
│   │   └── jsonstore.go           # Implémentation JSONStore
│   │
│   ├── domain/
│   │   ├── scheduler.go           # SM-2 Algorithm (CORE LOGIC)
│   │   ├── recommender.go         # Recommandations intelligentes
│   │   └── scheduler_test.go      # Tests unitaires
│   │
│   └── api/
│       └── handlers.go            # HTTP Handlers pour les endpoints
│
├── cmd/maestro-server/
│   └── main.go                    # Entry point HTTP Server
│
├── Makefile                        # Build automation
├── go.mod                          # Dépendances Go
└── README.md                       # Cette file
```

---

## 📋 Installation & Démarrage

### 1️⃣ Cloner et Setup

```bash
# Créer le dossier
mkdir maestro && cd maestro

# Copier les fichiers:
# - go.mod
# - go.mod-starter (renommer en go.mod si nécessaire)
# - Tous les fichiers internal/**/*.go
# - cmd/maestro-server/main.go
# - Makefile

# Structurer correctement:
mkdir -p internal/{models,storage,domain,api} cmd/maestro-server

# Copier les fichiers aux bons endroits (voir liste ci-dessus)
```

### 2️⃣ Installer Dépendances

```bash
make deps
# ou
go mod download
go mod tidy
```

### 3️⃣ Build

```bash
make build
# ou
go build -o bin/maestro-server ./cmd/maestro-server
```

### 4️⃣ Run

```bash
make run
# ou
./bin/maestro-server -port 8080

# Avec data directory custom:
./bin/maestro-server -port 8080 -data-dir /custom/path
```

**Vous verrez:**
```
🎯 Maestro Backend listening on http://localhost:8080
📁 Data directory: /home/user/.maestro
📄 Exercises file: /home/user/.maestro/exercises.json

✨ Endpoints:
  GET  http://localhost:8080/api/health
  GET  http://localhost:8080/api/exercises
  GET  http://localhost:8080/api/recommended
  POST http://localhost:8080/api/rate
  GET  http://localhost:8080/api/stats

🌐 Web UI: http://localhost:8080
```

---

## 🧪 Tests

```bash
# Tous les tests
make test

# Avec couverture
make test-coverage

# Juste domain
make test-domain

# Juste storage
make test-storage
```

**Expected output:**
```
ok  	maestro/internal/domain	0.002s	coverage: 75.0% of statements
ok  	maestro/internal/storage	0.001s	coverage: 60.0% of statements
```

---

## 🔧 Commandes Makefile

```bash
make build           # Build le binaire
make run             # Build + run sur 8080
make run-port        # Demande le port interactif
make test            # Lancer tests
make test-coverage   # Tests + coverage
make fmt             # Formater code
make lint            # Linter code
make clean           # Nettoyer
make deps            # Télécharger deps
make dev             # Hot reload (nécessite air)
make help            # Voir toutes les commandes
```

---

## 📡 API Endpoints

### GET `/api/health`
Vérifier que le serveur est alive

```bash
curl http://localhost:8080/api/health

# Response:
{
  "status": "ok",
  "message": "Maestro Backend is running"
}
```

### GET `/api/exercises`
Retourner tous les exercices

```bash
curl http://localhost:8080/api/exercises

# Response:
[
  {
    "id": "go-001",
    "title": "Goroutines Basics",
    "description": "Learn how goroutines work",
    "domain": "golang",
    "difficulty": 1,
    "completed": false,
    "ease_factor": 2.5,
    "interval_days": 0,
    ...
  }
]
```

### GET `/api/recommended`
Retourner les 3 exercices recommandés

```bash
curl http://localhost:8080/api/recommended

# Response: Array of 3 exercises (due for review or new)
```

### POST `/api/rate`
Noter un exercice (applique SM-2)

```bash
curl -X POST http://localhost:8080/api/rate \
  -H "Content-Type: application/json" \
  -d '{
    "exercise_id": "go-001",
    "rating": 4
  }'

# Response:
{
  "exercise": { ... },
  "next_review_in_days": 1,
  "message": "🔥 Excellent! Parfaitement maîtrisé!"
}
```

### GET `/api/stats`
Retourner les statistiques

```bash
curl http://localhost:8080/api/stats

# Response:
{
  "total_completed": 4,
  "total_reviews": 9,
  "domain_stats": {
    "golang": {
      "completed": 2,
      "total": 4,
      "mastery": 75
    },
    "linux": {
      "completed": 1,
      "total": 3,
      "mastery": 50
    }
  }
}
```

---

## 📂 Fichier de Données

Les exercices sont persistés dans:
```
~/.maestro/exercises.json
```

Format:
```json
[
  {
    "id": "go-001",
    "title": "Goroutines Basics",
    "description": "...",
    "domain": "golang",
    "difficulty": 1,
    "steps": ["Step 1", "Step 2"],
    "content": "package main\n...",
    "completed": false,
    "last_reviewed": null,
    "ease_factor": 2.5,
    "interval_days": 0,
    "repetitions": 0,
    "created_at": "2025-11-17T20:00:00Z",
    "updated_at": "2025-11-17T20:00:00Z"
  }
]
```

---

## 🧠 Architecture Expliquée

### Layers

1. **Models** (`internal/models/`)
   - Pure data structures
   - NO business logic
   - NO HTTP stuff

2. **Domain** (`internal/domain/`)
   - SM-2 Algorithm
   - Recommender Logic
   - Testable WITHOUT database
   - Testable WITHOUT HTTP

3. **Storage** (`internal/storage/`)
   - Interface Store (abstraction)
   - JSONStore implementation
   - Easy to swap for Database later

4. **API** (`internal/api/`)
   - HTTP handlers only
   - Uses domain logic
   - Uses storage

5. **Main** (`cmd/maestro-server/`)
   - Entry point
   - Wires everything together

### Data Flow

```
HTTP Request
    ↓
API Handler (exerciseHandler.RateExercise)
    ↓
Domain Logic (scheduler.ReviewExercise)
    ↓
Storage (store.Update)
    ↓
JSON File (exercises.json)
    ↓
HTTP Response
```

---

## 🚦 Tester Rapidement

```bash
# Build
make build

# Run
make run &

# Dans un autre terminal:

# 1. Health check
curl http://localhost:8080/api/health

# 2. Voir les exercices
curl http://localhost:8080/api/exercises | jq

# 3. Voir les recommandés
curl http://localhost:8080/api/recommended | jq

# 4. Noter un exercice
curl -X POST http://localhost:8080/api/rate \
  -H "Content-Type: application/json" \
  -d '{"exercise_id":"go-001","rating":4}' | jq

# 5. Voir les stats
curl http://localhost:8080/api/stats | jq

# Arrêter le serveur
pkill maestro-server
```

---

## 🎯 Prochaines Étapes

1. **Ajouter exercices de départ** → Créer `exercises.json` avec contenu initial
2. **Frontend Web** → Créer `public/index.html` qui consomme l'API
3. **Configuration** → Ajouter `config.yml` optionnel
4. **CLI Tool** → Créer `cmd/maestro-cli/main.go` pour Terminal UI

---

## 💡 Tips

### Hot Reload Development

```bash
go install github.com/cosmtrek/air@latest
make dev
```

### Format + Lint avant Commit

```bash
make fmt
make lint
git add .
git commit -m "message"
```

### Debug Requests

```bash
# Avec verbose output
curl -v http://localhost:8080/api/exercises

# Pretty print JSON
curl http://localhost:8080/api/exercises | jq .

# Suivre les redirect
curl -L http://localhost:8080/api/exercises
```

---

**Vous avez maintenant une architecture Go professionnelle prête pour la production !** 🎉
