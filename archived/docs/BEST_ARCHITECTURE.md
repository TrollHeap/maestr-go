# 🏗️ Architecture Maestro - Analyse Complète de Maintenabilité

## 🎯 La Meilleure Architecture pour Maestro

Je vais vous montrer une architecture **vraiment maintenable** basée sur les bonnes pratiques Go.

---

## 1️⃣ Architecture Recommandée (BEST PRACTICE)

### Structure de Dossiers

```
maestro/
│
├── cmd/                          # Points d'entrée (entry points)
│   ├── maestro-server/
│   │   └── main.go              # HTTP Server + API
│   ├── maestro-cli/
│   │   └── main.go              # CLI Terminal (optionnel)
│   └── maestro-sync/
│       └── main.go              # Sync daemon (optionnel)
│
├── internal/                     # Code privé (non-importable)
│   ├── models/
│   │   ├── exercise.go          # struct Exercise
│   │   ├── stats.go             # struct Stats
│   │   └── types.go             # Types génériques
│   │
│   ├── domain/                  # LOGIQUE MÉTIER (indépendante)
│   │   ├── scheduler.go         # SM-2 Algorithm (CORE LOGIC)
│   │   ├── recommender.go       # Recommandations
│   │   ├── calculator.go        # Calculs progression
│   │   └── validator.go         # Validations
│   │
│   ├── storage/                 # Abstraction de persistance
│   │   ├── store.go             # Interface Store
│   │   ├── jsonstore.go         # Implémentation JSON
│   │   └── mockstore.go         # Mock pour tests
│   │
│   ├── api/                     # HTTP API (handlers)
│   │   ├── middleware.go        # CORS, logging, etc.
│   │   ├── exercises.go         # Handlers exercices
│   │   ├── recommendations.go   # Handlers recommandations
│   │   ├── stats.go             # Handlers stats
│   │   └── routes.go            # Route registration
│   │
│   ├── config/                  # Configuration
│   │   ├── config.go            # Structs config
│   │   └── loader.go            # Charge desde YAML/ENV
│   │
│   └── logger/                  # Logging structuré
│       └── logger.go            # Loggers centralisés
│
├── public/                       # Assets statiques
│   ├── index.html               # Frontend web
│   ├── css/
│   │   └── style.css            # Styles
│   └── js/
│       └── app.js               # Frontend logic
│
├── tests/                        # Tests intégration
│   ├── api_test.go
│   ├── domain_test.go
│   └── fixtures/
│       └── exercises.json
│
├── docs/                         # Documentation
│   ├── API.md                   # API spec
│   ├── ARCHITECTURE.md          # Architecture decision records
│   └── DEVELOPMENT.md           # Dev guide
│
├── go.mod
├── go.sum
├── Makefile                      # Build automation
├── docker-compose.yml           # Local dev setup
└── README.md
```

---

## 2️⃣ Principes Clés d'Architecture

### A. **Separation of Concerns (SoC)**

**MAUVAIS** 🔴 :
```go
// Tout mélangé
type Model struct {
    exercises []Exercise
    // HTTP stuff
    w http.ResponseWriter
    r *http.Request
    // Rendering
    buffer strings.Builder
    // Business logic
    schedulerData map[string]interface{}
}
```

**BON** ✅ :
```go
// Modèles de domaine purs (pas d'imports HTTP)
type Exercise struct {
    ID           string
    Title        string
    EaseFactor   float64
    IntervalDays int
}

// Logique métier (indépendante)
type Scheduler interface {
    CalculateNextReview(ex Exercise, rating int) Exercise
    GetNextDueExercises(exercises []Exercise, limit int) []Exercise
}

// API handlers (utilise le reste)
type ExerciseHandler struct {
    scheduler Scheduler
    store     Store
    logger    Logger
}

func (h *ExerciseHandler) RateExercise(w http.ResponseWriter, r *http.Request) {
    // Appelle la logique métier
}
```

### B. **Dependency Injection (DI)**

**MAUVAIS** 🔴 :
```go
// Global state - impossible à tester
var store Store
var scheduler Scheduler

func RateExercise(w http.ResponseWriter, r *http.Request) {
    // Utilise global store - tight coupling
}
```

**BON** ✅ :
```go
// Passer les dépendances en paramètre
type API struct {
    store      Store
    scheduler  Scheduler
    logger     Logger
    config     *Config
}

func NewAPI(store Store, scheduler Scheduler, logger Logger, cfg *Config) *API {
    return &API{
        store:     store,
        scheduler: scheduler,
        logger:    logger,
        config:    cfg,
    }
}

func (a *API) RateExercise(w http.ResponseWriter, r *http.Request) {
    // Utilise les dépendances injectées
}
```

### C. **Interfaces au Lieu de Concrétion**

**MAUVAIS** 🔴 :
```go
// Dépendance sur l'implémentation concrète
type ExerciseHandler struct {
    store *JSONStore  // Couplé à JSON
}
```

**BON** ✅ :
```go
// Dépendance sur l'interface
type ExerciseHandler struct {
    store Store  // Abstraction
}

// Interface définie
type Store interface {
    Load(ctx context.Context) ([]Exercise, error)
    Save(ctx context.Context, exercises []Exercise) error
    GetByID(ctx context.Context, id string) (*Exercise, error)
    Update(ctx context.Context, ex *Exercise) error
}

// N'importe quelle implémentation peut être utilisée
type JSONStore struct { ... }
type DatabaseStore struct { ... }
type MemoryStore struct { ... }
```

### D. **Error Handling Approprié**

**MAUVAIS** 🔴 :
```go
func LoadExercises() []Exercise {
    data := ioutil.ReadFile("file.json")  // Ignore l'erreur!
    var exercises []Exercise
    json.Unmarshal(data, &exercises)
    return exercises
}
```

**BON** ✅ :
```go
func (s *JSONStore) Load(ctx context.Context) ([]Exercise, error) {
    data, err := os.ReadFile(s.filepath)
    if err != nil {
        return nil, fmt.Errorf("load exercises: %w", err)
    }
    
    var exercises []Exercise
    if err := json.Unmarshal(data, &exercises); err != nil {
        return nil, fmt.Errorf("parse exercises: %w", err)
    }
    
    return exercises, nil
}
```

### E. **Context Awareness**

**MAUVAIS** 🔴 :
```go
func (a *API) GetExercises(w http.ResponseWriter, r *http.Request) {
    // Pas de timeout, pas d'annulation possible
    exercises := a.store.Load()  // Bloque indéfiniment
}
```

**BON** ✅ :
```go
func (a *API) GetExercises(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // Récupère le context de la request
    
    exercises, err := a.store.Load(ctx)
    if err != nil {
        http.Error(w, "load failed", http.StatusInternalServerError)
        return
    }
}

// Le client peut annuler:
// curl --max-time 5 http://localhost:8080/api/exercises
```

---

## 3️⃣ Code d'Architecture Complète

### models/exercise.go
```go
package models

import "time"

// Exercise représente un exercice d'apprentissage
type Exercise struct {
    ID           string     `json:"id"`
    Title        string     `json:"title"`
    Description  string     `json:"description"`
    Domain       string     `json:"domain"`           // golang, linux, architecture
    Difficulty   int        `json:"difficulty"`       // 1-3
    Steps        []string   `json:"steps"`
    Content      string     `json:"content"`
    
    // Spaced Repetition
    Completed    bool       `json:"completed"`
    LastReviewed *time.Time `json:"last_reviewed"`
    EaseFactor   float64    `json:"ease_factor"`
    IntervalDays int        `json:"interval_days"`
    Repetitions  int        `json:"repetitions"`
    
    // Timestamps
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

// ReviewInput est la requête pour noter un exercice
type ReviewInput struct {
    ExerciseID string `json:"exercise_id"`
    Rating     int    `json:"rating"`  // 1-4
}

// ReviewResponse est la réponse après notation
type ReviewResponse struct {
    Exercise      *Exercise `json:"exercise"`
    NextReviewIn  int       `json:"next_review_in_days"`
    Message       string    `json:"message"`
}
```

### internal/storage/store.go (INTERFACE)
```go
package storage

import (
    "context"
    "maestro/internal/models"
)

// Store est l'interface pour la persistance
type Store interface {
    // Load charge tous les exercices
    Load(ctx context.Context) ([]models.Exercise, error)
    
    // Save persiste les exercices
    Save(ctx context.Context, exercises []models.Exercise) error
    
    // GetByID récupère un exercice par ID
    GetByID(ctx context.Context, id string) (*models.Exercise, error)
    
    // Update met à jour un exercice
    Update(ctx context.Context, ex *models.Exercise) error
}
```

### internal/storage/jsonstore.go (IMPLÉMENTATION)
```go
package storage

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "maestro/internal/models"
)

type JSONStore struct {
    filepath string
}

func NewJSONStore(filepath string) *JSONStore {
    return &JSONStore{filepath: filepath}
}

func (s *JSONStore) Load(ctx context.Context) ([]models.Exercise, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    data, err := os.ReadFile(s.filepath)
    if err != nil {
        return nil, fmt.Errorf("read file: %w", err)
    }
    
    var exercises []models.Exercise
    if err := json.Unmarshal(data, &exercises); err != nil {
        return nil, fmt.Errorf("parse JSON: %w", err)
    }
    
    return exercises, nil
}

func (s *JSONStore) Save(ctx context.Context, exercises []models.Exercise) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    data, err := json.MarshalIndent(exercises, "", "  ")
    if err != nil {
        return fmt.Errorf("marshal JSON: %w", err)
    }
    
    if err := os.WriteFile(s.filepath, data, 0644); err != nil {
        return fmt.Errorf("write file: %w", err)
    }
    
    return nil
}

func (s *JSONStore) GetByID(ctx context.Context, id string) (*models.Exercise, error) {
    exercises, err := s.Load(ctx)
    if err != nil {
        return nil, err
    }
    
    for _, ex := range exercises {
        if ex.ID == id {
            return &ex, nil
        }
    }
    
    return nil, fmt.Errorf("exercise not found: %s", id)
}

func (s *JSONStore) Update(ctx context.Context, ex *models.Exercise) error {
    exercises, err := s.Load(ctx)
    if err != nil {
        return err
    }
    
    for i := range exercises {
        if exercises[i].ID == ex.ID {
            exercises[i] = *ex
            return s.Save(ctx, exercises)
        }
    }
    
    return fmt.Errorf("exercise not found: %s", ex.ID)
}
```

### internal/domain/scheduler.go (LOGIQUE MÉTIER - PURE)
```go
package domain

import (
    "time"
    "maestro/internal/models"
)

// Scheduler gère l'algorithme SM-2
type Scheduler struct {
    initialEaseFactor float64
    minEaseFactor     float64
}

func NewScheduler() *Scheduler {
    return &Scheduler{
        initialEaseFactor: 2.5,
        minEaseFactor:     1.3,
    }
}

// ReviewExercise applique l'algorithme SM-2
func (s *Scheduler) ReviewExercise(ex *models.Exercise, rating int) {
    if rating < 1 || rating > 4 {
        return  // Invalid rating
    }
    
    var newInterval int
    var newEF float64
    
    switch rating {
    case 4:  // Facile
        newInterval = int(float64(ex.IntervalDays) * ex.EaseFactor)
        newEF = ex.EaseFactor + 0.1
    case 3:  // Normal
        newInterval = int(float64(ex.IntervalDays) * ex.EaseFactor)
        newEF = ex.EaseFactor
    case 2:  // Difficile
        newInterval = int(float64(ex.IntervalDays) * 0.5)
        newEF = ex.EaseFactor - 0.2
    case 1:  // Oublié
        newInterval = 1
        newEF = ex.EaseFactor - 0.5
    }
    
    // Clamp EF
    if newEF < s.minEaseFactor {
        newEF = s.minEaseFactor
    }
    
    // Update exercise
    now := time.Now()
    ex.LastReviewed = &now
    ex.IntervalDays = newInterval
    ex.EaseFactor = newEF
    ex.Repetitions++
    ex.UpdatedAt = now
}

// IsDueForReview vérifie si l'exercice est à réviser
func (s *Scheduler) IsDueForReview(ex *models.Exercise) bool {
    if ex.LastReviewed == nil {
        return false
    }
    nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
    return time.Now().After(nextReview)
}
```

### internal/domain/recommender.go
```go
package domain

import (
    "maestro/internal/models"
)

// Recommender suggère les exercices à faire
type Recommender struct {
    scheduler *Scheduler
}

func NewRecommender(scheduler *Scheduler) *Recommender {
    return &Recommender{scheduler: scheduler}
}

// GetNextExercises retourne les exercices à faire ensuite
func (r *Recommender) GetNextExercises(exercises []models.Exercise, limit int) []models.Exercise {
    var recommended []models.Exercise
    
    // Priorité 1: Exercices dus pour révision
    for _, ex := range exercises {
        if r.scheduler.IsDueForReview(&ex) {
            recommended = append(recommended, ex)
        }
    }
    
    // Si pas assez, ajouter des nouveaux
    if len(recommended) < limit {
        for _, ex := range exercises {
            if !ex.Completed && r.scheduler.IsDueForReview(&ex) == false {
                recommended = append(recommended, ex)
            }
            if len(recommended) >= limit {
                break
            }
        }
    }
    
    return recommended[:limit]
}
```

### internal/api/exercises.go (HANDLERS)
```go
package api

import (
    "encoding/json"
    "net/http"
    "maestro/internal/domain"
    "maestro/internal/models"
    "maestro/internal/storage"
)

type ExerciseHandler struct {
    store      storage.Store
    scheduler  *domain.Scheduler
    recommender *domain.Recommender
}

func NewExerciseHandler(store storage.Store, scheduler *domain.Scheduler, recommender *domain.Recommender) *ExerciseHandler {
    return &ExerciseHandler{
        store:       store,
        scheduler:   scheduler,
        recommender: recommender,
    }
}

// GetExercises retourne tous les exercices
func (h *ExerciseHandler) GetExercises(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    exercises, err := h.store.Load(ctx)
    if err != nil {
        http.Error(w, "Failed to load exercises", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(exercises)
}

// RateExercise met à jour la note d'un exercice
func (h *ExerciseHandler) RateExercise(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    var input models.ReviewInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    exercise, err := h.store.GetByID(ctx, input.ExerciseID)
    if err != nil {
        http.Error(w, "Exercise not found", http.StatusNotFound)
        return
    }
    
    // Appliquer l'algorithme SM-2
    h.scheduler.ReviewExercise(exercise, input.Rating)
    
    // Persister
    if err := h.store.Update(ctx, exercise); err != nil {
        http.Error(w, "Failed to update", http.StatusInternalServerError)
        return
    }
    
    // Répondre
    response := models.ReviewResponse{
        Exercise:     exercise,
        NextReviewIn: exercise.IntervalDays,
        Message:      "✅ Exercice enregistré!",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

### cmd/maestro-server/main.go (ENTRY POINT)
```go
package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    
    "maestro/internal/api"
    "maestro/internal/domain"
    "maestro/internal/storage"
)

func main() {
    port := flag.String("port", "8080", "Port to listen on")
    dataDir := flag.String("data-dir", "", "Data directory")
    flag.Parse()
    
    // Setup data directory
    if *dataDir == "" {
        home, err := os.UserHomeDir()
        if err != nil {
            log.Fatal(err)
        }
        *dataDir = filepath.Join(home, ".maestro")
    }
    os.MkdirAll(*dataDir, 0755)
    
    // Initialize store
    store := storage.NewJSONStore(filepath.Join(*dataDir, "exercises.json"))
    
    // Initialize domain logic
    scheduler := domain.NewScheduler()
    recommender := domain.NewRecommender(scheduler)
    
    // Initialize handlers
    exerciseHandler := api.NewExerciseHandler(store, scheduler, recommender)
    
    // Setup routes
    http.HandleFunc("/api/exercises", exerciseHandler.GetExercises)
    http.HandleFunc("/api/rate", exerciseHandler.RateExercise)
    
    // Serve frontend
    fs := http.FileServer(http.Dir("public"))
    http.Handle("/", fs)
    
    fmt.Printf("🎯 Maestro listening on http://localhost:%s\n", *port)
    if err := http.ListenAndServe(":"+*port, nil); err != nil {
        log.Fatal(err)
    }
}
```

---

## 4️⃣ Testing Strategy

### internal/domain/scheduler_test.go
```go
package domain

import (
    "testing"
    "time"
    "maestro/internal/models"
)

func TestSM2Algorithm(t *testing.T) {
    scheduler := NewScheduler()
    
    ex := &models.Exercise{
        ID:         "test-1",
        EaseFactor: 2.5,
        IntervalDays: 0,
    }
    
    // Premier rating: facile
    scheduler.ReviewExercise(ex, 4)
    
    if ex.EaseFactor != 2.6 {
        t.Fatalf("Expected EF 2.6, got %f", ex.EaseFactor)
    }
}

func TestIsDueForReview(t *testing.T) {
    scheduler := NewScheduler()
    
    now := time.Now()
    past := now.AddDate(0, 0, -1)
    
    ex := &models.Exercise{
        LastReviewed: &past,
        IntervalDays: 0,  // Due since yesterday
    }
    
    if !scheduler.IsDueForReview(ex) {
        t.Fatal("Expected exercise to be due")
    }
}
```

---

## 5️⃣ Résumé: Pourquoi Cette Architecture?

| Aspect | Bénéfice |
|--------|----------|
| **Séparation Go/Web** | Frontend peut changer sans toucher logique |
| **Interfaces au lieu de concrétions** | Testable, pluggable |
| **Dependency Injection** | Pas de globals, facile à tester |
| **Pas d'imports circulaires** | Code clean, maintenable |
| **Context awareness** | Cancellation, timeouts, tracing |
| **Error handling** | Erreurs explicites, traçables |
| **Domain-driven** | Logique métier indépendante |

---

## 6️⃣ Makefile pour Simplifier

```makefile
.PHONY: build run test clean

build:
	go build -o bin/maestro-server ./cmd/maestro-server

run: build
	./bin/maestro-server -port 8080 -data-dir ~/.maestro

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -cover ./...

clean:
	rm -rf bin/

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

dev:
	air  # Hot reload
```

---

**C'est ça la vraie architecture pour la maintenabilité !** 🎯

Veux-tu que je te crée les fichiers Go complets pour démarrer avec cette architecture ?
