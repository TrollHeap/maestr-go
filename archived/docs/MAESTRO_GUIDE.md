# 🎯 Maestro - Ultra-Learning Practice Tool pour ADHD

**L'outil parfait pour apprendre par la pratique en Go, avec support ADHD natif.**

Maestro est un **CLI TUI** (Terminal User Interface) conçu selon les principes scientifiques d'apprentissage ultra-rapide et l'ultra-learning, spécialisé pour les personnes avec ADHD.

## 🚀 Caractéristiques Principales

### 1. **Sessions Flash Anti-Procrastination**
- **15-30 minutes** de pratique focalisée max
- Lanceables en **1 seule commande**
- Compte à rebours visuel
- 3 exercices recommandés seulement

### 2. **Spaced Repetition (Algorithme SM-2)**
- Intervalles optimaux : 1, 3, 7, 14, 30 jours
- Tracking automatique de progression
- Révisions intelligentes basées sur difficulté
- Calcul d'intervalle basé sur performance

### 3. **Visual Progress Indicators**
```
Streak: ✓✓✓✓✓  (5 jours)
Mastery:  
  Golang     [████░░░░░] 40%
  Linux      [██░░░░░░░] 20%
  Architecture [███░░░░░░] 30%
```

### 4. **Système de Katas en Go**
- Exercices progressifs (Facile → Moyen → Difficile)
- Templates de code pratiques
- Checklists visuelles de progression
- Domaines : Go, Linux, Architecture Système

### 5. **Logiques Visuelles pour Comprendre**
ASCII art des concepts complexes :
```
Memory Hierarchy:
┌─────────────────┐
│   L1 Cache (32KB)│  ← Rapide, petit
│  ┌───────────┐  │
│  │ L2 Cache  │  │
│  │ (256KB)   │  │
│  │ ┌───────┐ │  │
│  │ │ RAM   │ │  │  Lent, énorme
│  │ │ (16GB)│ │  │
│  │ │ DISK  │ │  │
│  │ │ (1TB) │ │  │
└─────────────────┘
```

## 📋 Installation

### Prérequis
- **Go 1.21+**
- Terminal compatible (Linux/macOS/Windows avec Git Bash)

### 1. Cloner le repo
```bash
git clone https://github.com/yourusername/maestro.git
cd maestro
```

### 2. Installer les dépendances
```bash
go mod download
go mod tidy
```

### 3. Compiler
```bash
go build -o maestro .
```

### 4. Lancer
```bash
./maestro
```

## 🎮 Utilisation

### Dashboard Principal
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  🎯 MAESTRO - Ultra-Learning
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  Streak: ✓✓✓✓✓ (5 jours)
  Aujourd'hui: 2/3 exercices
  
  ⏱  Prochaine session: 8:45 AM
  
  Recommandé: Go - Goroutines Basics (Débutant)
  
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  [q] Quick Start  [b] Browse  [s] Stats  [q] Quit
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Keybindings
| Touche | Action |
|--------|--------|
| `q` | Quick Start 15-min session |
| `b` | Browse all exercises |
| `s` | View statistics & progress |
| `d` | Domain filter |
| `↑/↓` `j/k` | Navigate |
| `↵` `Enter` | Select/Confirm |
| `1/2/3/4` | Rate exercise (1=forgot, 4=easy) |
| `esc` `q` | Back/Quit |

## 🏗️ Architecture du Code

```
maestro/
├── main.go                 # Entry point
├── go.mod
├── go.sum
├── models/
│   ├── exercise.go        # Exercise data structure
│   ├── stats.go           # User statistics
│   └── spaced_rep.go      # SM-2 algorithm
├── storage/
│   ├── json_store.go      # JSON persistence
│   └── exercises.json     # Default exercises
├── ui/
│   ├── dashboard.go       # Main dashboard view
│   ├── browser.go         # Exercise browser
│   ├── practice.go        # Practice mode
│   ├── styles.go          # Lipgloss styling
│   └── visual_models.go   # ASCII art diagrams
├── logic/
│   ├── scheduler.go       # Spaced repetition scheduling
│   ├── session.go         # Session management
│   └── progress.go        # Progress calculations
└── README.md
```

## 📊 Fichier `exercises.json`

```json
{
  "exercises": [
    {
      "id": "go-001",
      "title": "Goroutines Basics",
      "description": "Learn how goroutines work...",
      "domain": "golang",
      "difficulty": 1,
      "steps": ["Create goroutine", "Use WaitGroup", "Understand scheduling"],
      "content": "package main\n\nimport (\n\t\"sync\"\n)\n\nfunc main() {\n\tvar wg sync.WaitGroup\n\t// TODO: Add goroutines\n}\n",
      "completed": false,
      "last_reviewed": "2025-11-16",
      "ease_factor": 2.5,
      "interval_days": 0,
      "repetitions": 0
    }
  ],
  "user_stats": {
    "current_streak": 5,
    "total_completed": 4,
    "total_reviews": 9,
    "last_session": "2025-11-17"
  }
}
```

## 🧠 Principes ADHD Intégrés

### 1. Réduction de la Surcharge Cognitive
- ✅ Pas de menu de 50 options
- ✅ Choix limités (3 recommandations max)
- ✅ Interface épurée et claire

### 2. Gratification Immédiate
- ✅ Streaks visuels (✓✓✓✓✓)
- ✅ Compteurs de progression
- ✅ Feedback immédiat après chaque exercice
- ✅ Messages encourageants

### 3. Chunking (Décomposition)
- ✅ Exercices en 15-30 min max
- ✅ Tâches divisées en steps
- ✅ Progrès visible par étape

### 4. Momentum Building
- ✅ Sessions courtes construisent les streaks
- ✅ Streaks génèrent de la motivation
- ✅ Recommandations basées sur capacité actuelle

## 🛠️ Dépendances Go

```go
// go.mod
require (
    github.com/charmbracelet/bubbletea v0.24.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/charmbracelet/huh v0.3.0
)
```

## 📚 Exercices Inclus

### Golang (4 exercices)
- Goroutines & Concurrency
- Channels & Communication
- Interfaces & Polymorphism
- Error Handling Patterns

### Linux (3 exercices)
- Tmux Window Management
- Shell Scripting Fundamentals
- File Permissions & Ownership

### Architecture Système (3 exercices)
- Memory Hierarchy & Caches
- Process vs Threads Model
- Virtual Memory & Paging

## 🎯 Cas d'Usage

### 1. **Apprendre Go rapidement**
```bash
./maestro q
# → Lancer une session 15-min, 3 exercices Go recommandés
```

### 2. **Reviser régulièrement**
```bash
./maestro b
# → Voir exercices dues pour révision (marquées ⏱)
```

### 3. **Tracker progression**
```bash
./maestro s
# → Voir stats complètes et graphique de mastery
```

## 🔄 Algorithme Spaced Repetition (SM-2)

L'application utilise l'algorithme SM-2 optimisé :

```
Intervals (days): 1, 3, 7, 14, 30
EaseFactor = initial 2.5
On rating (1-4):
  - Rating 4 (Facile): interval *= EF, EF += 0.1
  - Rating 3 (Normal): interval *= EF
  - Rating 2 (Difficile): interval *= 0.5, EF -= 0.2
  - Rating 1 (Oublié): reset à 1 jour
```

## 💾 Persistence des Données

Toutes les données sont stockées en **JSON local** :
- `~/.maestro/exercises.json` - Exercices et progress
- `~/.maestro/stats.json` - Statistiques utilisateur
- `~/.maestro/sessions.json` - Historique sessions

### Format complet d'un exercice persisté
```json
{
  "id": "go-001",
  "title": "Goroutines",
  "completed": true,
  "last_reviewed": "2025-11-17T10:30:00Z",
  "ease_factor": 2.8,
  "interval_days": 7,
  "repetitions": 3,
  "next_review": "2025-11-24",
  "review_history": [
    {"date": "2025-11-15", "rating": 3},
    {"date": "2025-11-17", "rating": 4}
  ]
}
```

## 🚀 Prochaines Étapes - Roadmap

- [ ] Exercices interactifs avec exécution Go en temps réel
- [ ] Synchronisation cloud pour multi-device
- [ ] Dashboard web pour visualisation
- [ ] Export statistiques (CSV/PDF)
- [ ] Système de badges et récompenses
- [ ] Intégration Anki pour flashcards
- [ ] Mobile app (Flutter)

## 📖 Ressources Pédagogiques Intégrées

Chaque exercice inclut :
- Description claire du concept
- Template de code commenté
- Visualisation ASCII du concept
- Checklist de compréhension
- Liens vers ressources externes

## 🎨 Customization

### Thème personnalisé
```go
// ui/styles.go
var customTheme = Theme{
    Primary:   lipgloss.Color("#5D4E60"),
    Success:   lipgloss.Color("#90EE90"),
    Warning:   lipgloss.Color("#FFD700"),
    Error:     lipgloss.Color("#FF6B6B"),
}
```

### Ajouter des exercices
```bash
# Editer exercises.json
{
  "id": "custom-001",
  "title": "Mon exercice",
  "domain": "golang",
  "difficulty": 2,
  ...
}
```

## 📞 Support & Contributions

Ce projet est open-source. Contributions bienvenues !

```bash
git clone <your-fork>
git checkout -b feature/mon-feature
# Faire changements
git push origin feature/mon-feature
# Créer Pull Request
```

## 📄 Licence

MIT - Voir LICENSE pour détails

---

## ⚡ TL;DR - Get Started in 2 Minutes

```bash
# Clone
git clone https://github.com/yourusername/maestro.git && cd maestro

# Build
go build -o maestro .

# Launch
./maestro

# Press 'q' for Quick Start session
```

**Prêt à maîtriser Go, Linux et l'architecture système ? Let's go ! 🚀**
