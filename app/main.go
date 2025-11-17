package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Exercise représente un exercice d'apprentissage
type Exercise struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Domain       string     `json:"domain"`     // golang, linux, architecture
	Difficulty   int        `json:"difficulty"` // 1-3
	Steps        []string   `json:"steps"`
	Content      string     `json:"content"`
	Completed    bool       `json:"completed"`
	LastReviewed *time.Time `json:"last_reviewed"`
	EaseFactor   float64    `json:"ease_factor"`
	IntervalDays int        `json:"interval_days"`
	Repetitions  int        `json:"repetitions"`
}

// UserStats représente les statistiques utilisateur
type UserStats struct {
	CurrentStreak  int        `json:"current_streak"`
	TotalCompleted int        `json:"total_completed"`
	TotalReviews   int        `json:"total_reviews"`
	LastSession    *time.Time `json:"last_session"`
}

// Store gère la persistance des données
type Store struct {
	Exercises Exercise
	Stats     UserStats
	filepath  string
}

// Model représente l'état de l'application
type Model struct {
	store        *Store
	currentView  string // "dashboard", "browser", "practice", "stats"
	selectedIdx  int
	exercises    []Exercise
	sessionTimer int
	filterDomain string
	message      string
}

// ✅ MÉTHODE MANQUANTE - C'EST LA CLÉ !
// Init initialise le modèle (requis par tea.Model)
func (m Model) Init() tea.Cmd {
	return nil
}

// View affiche l'interface utilisateur
func (m Model) View() string {
	switch m.currentView {
	case "dashboard":
		return m.renderDashboard()
	case "browser":
		return m.renderBrowser()
	case "practice":
		return m.renderPractice()
	case "stats":
		return m.renderStats()
	default:
		return m.renderDashboard()
	}
}

// renderDashboard affiche le dashboard principal
func (m Model) renderDashboard() string {
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("5")).
		Bold(true).
		Render("🎯 MAESTRO - Ultra-Learning")

	streak := lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Render(fmt.Sprintf("Streak: %s %d jours", repeatString("✓", m.store.Stats.CurrentStreak), m.store.Stats.CurrentStreak))

	// Calculer prochains exercices à réviser
	nextDue := m.getNextDueExercises()
	var nextExercise string
	if len(nextDue) > 0 {
		nextExercise = fmt.Sprintf("Recommandé: %s", nextDue[0].Title)
	} else {
		nextExercise = "Aucun exercice recommandé pour maintenant"
	}

	progress := lipgloss.NewStyle().
		Foreground(lipgloss.Color("3")).
		Render(fmt.Sprintf("Aujourd'hui: %d exercices", m.store.Stats.TotalCompleted))

	footer := lipgloss.NewStyle().
		Faint(true).
		Render("[q] Quick Start  [b] Browse  [s] Stats  [q] Quit")

	content := fmt.Sprintf(
		"%s\n\n%s\n%s\n\n%s\n\n%s",
		header,
		streak,
		progress,
		nextExercise,
		footer,
	)
	return content
}

// renderBrowser affiche la liste des exercices
func (m Model) renderBrowser() string {
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("5")).
		Bold(true).
		Render("📚 Browse Exercises")

	var items []string
	for i, ex := range m.exercises {
		selected := ""
		if i == m.selectedIdx {
			selected = " > "
		}

		status := "○"
		if ex.Completed {
			status = "✓"
		} else if m.isDueForReview(ex) {
			status = "⏱"
		}

		diffStr := fmt.Sprintf("[D%d]", ex.Difficulty)
		item := fmt.Sprintf("%s %s %s %s", selected, status, diffStr, ex.Title)
		items = append(items, item)
	}

	list := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1).
		Render(fmt.Sprintf("%s\n%s", header, fmt.Sprintf("%v", items)))

	return list
}

// renderPractice affiche le mode pratique
func (m Model) renderPractice() string {
	if m.selectedIdx >= len(m.exercises) {
		return "Aucun exercice sélectionné"
	}

	ex := m.exercises[m.selectedIdx]

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("5")).
		Bold(true).
		Render(fmt.Sprintf("📝 %s", ex.Title))

	description := lipgloss.NewStyle().
		Render(fmt.Sprintf("Description: %s\n", ex.Description))

	// Afficher les steps
	stepsStr := "Steps:\n"
	for i, step := range ex.Steps {
		stepsStr += fmt.Sprintf("  %d. %s\n", i+1, step)
	}

	progress := lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Render(fmt.Sprintf("Progress: [%s%s] %d/%d\n",
			repeatString("█", 2), repeatString("░", 8), 2, len(ex.Steps)))

	content := fmt.Sprintf("%s\n%s\n%s\n%s", header, description, stepsStr, progress)
	return content
}

// renderStats affiche les statistiques
func (m Model) renderStats() string {
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("5")).
		Bold(true).
		Render("📊 Statistics")

	stats := fmt.Sprintf(
		"Streak: %d jours\nTotal Complétés: %d\nTotal Révisions: %d\n",
		m.store.Stats.CurrentStreak,
		m.store.Stats.TotalCompleted,
		m.store.Stats.TotalReviews,
	)

	return fmt.Sprintf("%s\n\n%s", header, stats)
}

// Update gère les inputs utilisateur
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.currentView == "dashboard" {
				return m, tea.Quit
			}
			m.currentView = "dashboard"
		case "b":
			m.currentView = "browser"
			m.selectedIdx = 0
		case "s":
			m.currentView = "stats"
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down", "j":
			if m.selectedIdx < len(m.exercises)-1 {
				m.selectedIdx++
			}
		case "enter":
			m.currentView = "practice"
		}
	}
	return m, nil
}

// Fonctions utilitaires

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func (m Model) isDueForReview(ex Exercise) bool {
	if ex.LastReviewed == nil {
		return false
	}
	nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
	return time.Now().After(nextReview)
}

func (m Model) getNextDueExercises() []Exercise {
	var due []Exercise
	for _, ex := range m.exercises {
		if m.isDueForReview(ex) {
			due = append(due, ex)
		}
	}
	// Limiter à 3 recommandations
	if len(due) > 3 {
		due = due[:3]
	}
	return due
}

// Gestion du store

func NewStore(dataDir string) (*Store, error) {
	filepath := filepath.Join(dataDir, "exercises.json")

	store := &Store{
		filepath: filepath,
	}

	// Charger les données existantes
	data, err := os.ReadFile(filepath)
	if err == nil {
		var container struct {
			Exercises Exercise  `json:"exercises"`
			Stats     UserStats `json:"user_stats"`
		}
		json.Unmarshal(data, &container)
		store.Exercises = container.Exercises
		store.Stats = container.Stats
	}

	return store, nil
}

func (s *Store) Save() error {
	container := map[string]interface{}{
		"exercises":  s.Exercises,
		"user_stats": s.Stats,
	}

	data, err := json.MarshalIndent(container, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filepath, data, 0o644)
}

// Main
func main() {
	// Créer répertoire de données
	dataDir := filepath.Join(os.Getenv("HOME"), ".maestro")
	os.MkdirAll(dataDir, 0o755)

	// Charger le store
	store, _ := NewStore(dataDir)

	// Créer exercices de test
	exercises := []Exercise{
		{
			ID:          "go-001",
			Title:       "Goroutines Basics",
			Description: "Learn how goroutines work",
			Domain:      "golang",
			Difficulty:  1,
			Steps:       []string{"Create goroutine", "Use WaitGroup", "Understand scheduling"},
			Completed:   false,
		},
		{
			ID:          "go-002",
			Title:       "Channels and Communication",
			Description: "Master channel patterns",
			Domain:      "golang",
			Difficulty:  2,
			Steps:       []string{"Create channels", "Send/receive", "Producer-consumer"},
			Completed:   true,
		},
		{
			ID:          "linux-001",
			Title:       "Tmux Window Management",
			Description: "Master terminal multiplexing",
			Domain:      "linux",
			Difficulty:  1,
			Steps:       []string{"Create sessions", "Split windows", "Navigate"},
			Completed:   false,
		},
	}

	// Trier exercices
	sort.Slice(exercises, func(i, j int) bool {
		return exercises[i].ID < exercises[j].ID
	})

	model := Model{
		store:       store,
		currentView: "dashboard",
		exercises:   exercises,
		selectedIdx: 0,
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Erreur: %v\n", err)
		os.Exit(1)
	}
}
