# 📦 Starter Kit Complet - Fichiers à Récupérer

## ✅ Fichiers Prêts (Vous avez tout)

### Configuration
- **go.mod** ← Modules Go

### Models
- **internal-models-exercise.go** → Mettre dans `internal/models/exercise.go`

### Storage Layer
- **internal-storage-store.go** → Mettre dans `internal/storage/store.go`
- **internal-storage-jsonstore.go** → Mettre dans `internal/storage/jsonstore.go`

### Domain Layer (LOGIQUE MÉTIER)
- **internal-domain-scheduler.go** → Mettre dans `internal/domain/scheduler.go`
- **internal-domain-recommender.go** → Mettre dans `internal/domain/recommender.go`
- **internal-domain-scheduler_test.go** → Mettre dans `internal/domain/scheduler_test.go`

### API Layer
- **internal-api-handlers.go** → Mettre dans `internal/api/handlers.go`

### Entry Point
- **cmd-maestro-server-main.go** → Mettre dans `cmd/maestro-server/main.go`

### Build & Config
- **Makefile** ← Automation
- **SETUP_GUIDE.md** ← Ce guide

---

## 🎯 Step-by-Step Setup (5 minutes)

```bash
# 1. Créer la structure
mkdir -p maestro/internal/{models,storage,domain,api} maestro/cmd/maestro-server
cd maestro

# 2. Copier les fichiers
# → go.mod → go.mod
# → internal-models-exercise.go → internal/models/exercise.go
# → internal-storage-store.go → internal/storage/store.go
# → internal-storage-jsonstore.go → internal/storage/jsonstore.go
# → internal-domain-scheduler.go → internal/domain/scheduler.go
# → internal-domain-recommender.go → internal/domain/recommender.go
# → internal-domain-scheduler_test.go → internal/domain/scheduler_test.go
# → internal-api-handlers.go → internal/api/handlers.go
# → cmd-maestro-server-main.go → cmd/maestro-server/main.go
# → Makefile → Makefile

# 3. Télécharger les dépendances
go mod download
go mod tidy

# 4. Build
make build

# 5. Run
make run
```

**Résultat:**
```
🎯 Maestro Backend listening on http://localhost:8080
📁 Data directory: /home/user/.maestro
📄 Exercises file: /home/user/.maestro/exercises.json

✨ Endpoints ready:
  GET  http://localhost:8080/api/health
  GET  http://localhost:8080/api/exercises
  GET  http://localhost:8080/api/recommended
  POST http://localhost:8080/api/rate
  GET  http://localhost:8080/api/stats

🌐 Web UI: http://localhost:8080
```

---

## 🧪 Test Immédiatement

```bash
# Health check
curl http://localhost:8080/api/health

# Voir les exercises (vides pour l'instant)
curl http://localhost:8080/api/exercises

# Voir les statistiques
curl http://localhost:8080/api/stats
```

---

## 📁 Arborescence Finale

```
maestro/
├── go.mod
├── Makefile
├── SETUP_GUIDE.md
│
├── internal/
│   ├── models/
│   │   └── exercise.go
│   ├── storage/
│   │   ├── store.go
│   │   └── jsonstore.go
│   ├── domain/
│   │   ├── scheduler.go
│   │   ├── recommender.go
│   │   └── scheduler_test.go
│   └── api/
│       └── handlers.go
│
├── cmd/
│   └── maestro-server/
│       └── main.go
│
└── bin/
    └── maestro-server  ← Binary après `make build`
```

---

## 🚀 Prochaines Étapes

### Étape 1: Ajouter les Exercices de Départ
```bash
# Créer ~/.maestro/exercises.json avec les 10 exercices starter
# (Voir MAESTRO_QUICKSTART.md pour le format)
```

### Étape 2: Frontend Web
```bash
# Créer public/index.html qui consomme l'API
# Voir les fichiers générés précédemment
```

### Étape 3: Configuration (Optionnel)
```bash
# Créer ~/.maestro/config.yml
# Pour customizer sans recompiler
```

---

## 🔥 Architecture Recap

```
User Browser (Frontend)
        ↓
    HTTP API (Go Backend)
        ↓
    Domain Logic (SM-2, Recommender)
        ↓
    Storage (JSON File)
```

**100% Clean** ✨
- ✅ Go Backend logique métier
- ✅ Frontend juste consommateur
- ✅ Testable
- ✅ Maintenable
- ✅ Extensible

---

## 💡 Quick Commands

```bash
make build           # Compiler
make run             # Compiler + lancer
make test            # Tester
make test-coverage   # Couverture tests
make fmt             # Formater
make lint            # Linter
make clean           # Nettoyer
make help            # Aide
```

---

**Vous êtes prêt ! 🎉**

Lancez `make run` et commencez à utiliser Maestro !
