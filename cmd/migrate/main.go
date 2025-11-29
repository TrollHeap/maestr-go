package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"maestro/internal/store"

	_ "modernc.org/sqlite"
)

// JSONExercise : Structure JSON source (FORMAT DATES YYYYMMDD)
type JSONExercise struct {
	ID                int         `json:"id"`
	Title             string      `json:"title"`
	Description       string      `json:"description"`
	Domain            string      `json:"domain"`
	Difficulty        int         `json:"difficulty"`
	Steps             []string    `json:"steps"`
	Content           string      `json:"content"`
	ConceptualVisuals []VisualAid `json:"conceptual_visuals"`
	Mnemonic          string      `json:"mnemonic"`
	Done              bool        `json:"done"`
	CompletedSteps    []int       `json:"completed_steps"`

	// Dates en format YYYYMMDD
	LastReviewedDate *int `json:"last_reviewed_date"`
	NextReviewDate   int  `json:"next_review_date"`
	LastSkippedDate  *int `json:"last_skipped_date"`
	CreatedAt        int  `json:"created_at"`
	UpdatedAt        int  `json:"updated_at"`

	EaseFactor   float64 `json:"ease_factor"`
	IntervalDays int     `json:"interval_days"`
	Repetitions  int     `json:"repetitions"`
	SkippedCount int     `json:"skipped_count"`
	Deleted      bool    `json:"deleted"`
}

type VisualAid struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Caption string `json:"caption"`
}

func main() {
	log.Println("🚀 Migration JSON → SQLite (Dates YYYYMMDD)")

	// 1. Init DB
	if err := store.InitDB("data/maestro.db"); err != nil {
		log.Fatal("Erreur init DB:", err)
	}
	defer store.CloseDB()

	// 2. Lit exercises.json
	data, err := os.ReadFile("data/exercises.json")
	if err != nil {
		log.Fatal("Erreur lecture exercises.json:", err)
	}

	var exercises []JSONExercise
	if err := json.Unmarshal(data, &exercises); err != nil {
		log.Fatal("Erreur parse JSON:", err)
	}

	log.Printf("📦 %d exercices trouvés dans JSON\n", len(exercises))

	// 3. Prépare statement INSERT
	db := store.GetDB()
	stmt, err := db.Prepare(`
		INSERT INTO exercises (
			id, title, description, domain, difficulty,
			content, mnemonic, conceptual_visuals,
			steps, completed_steps,
			done, last_reviewed_date, next_review_date,
			ease_factor, interval_days, repetitions,
			skipped_count, last_skipped_date,
			deleted, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Fatal("Erreur préparation statement:", err)
	}
	defer stmt.Close()

	// 4. Insert chaque exercice
	for _, ex := range exercises {
		// Serialize JSON fields
		stepsJSON, _ := json.Marshal(ex.Steps)
		completedJSON, _ := json.Marshal(ex.CompletedSteps)
		visualsJSON, _ := json.Marshal(ex.ConceptualVisuals)

		// Execute INSERT (les dates sont déjà en format YYYYMMDD depuis le JSON)
		_, err := stmt.Exec(
			ex.ID, ex.Title, ex.Description, ex.Domain, ex.Difficulty,
			ex.Content, ex.Mnemonic, visualsJSON,
			stepsJSON, completedJSON,
			boolToInt(ex.Done),  // ✅ Conversion bool → int
			ex.LastReviewedDate, // Déjà *int depuis JSON
			ex.NextReviewDate,   // Déjà int depuis JSON
			ex.EaseFactor, ex.IntervalDays, ex.Repetitions,
			ex.SkippedCount,
			ex.LastSkippedDate,    // Déjà *int depuis JSON
			boolToInt(ex.Deleted), // ✅ Conversion bool → int
			ex.CreatedAt,          // Déjà int depuis JSON
			ex.UpdatedAt,          // Déjà int depuis JSON
		)
		if err != nil {
			log.Printf("❌ Erreur insert exercice #%d (%s): %v", ex.ID, ex.Title, err)
			continue
		}

		log.Printf("✅ Migré: #%d - %s", ex.ID, ex.Title)
		log.Printf(
			"   └─ Next review: %d (%s)",
			ex.NextReviewDate,
			formatDateInt(ex.NextReviewDate),
		)
	}

	// 5. Vérification
	var count int
	db.QueryRow("SELECT COUNT(*) FROM exercises WHERE deleted = 0").Scan(&count)

	log.Printf("\n🎉 Migration terminée : %d/%d exercices dans la DB\n", count, len(exercises))

	// 6. Affiche stats
	printStats(db)
}

// ============================================
// HELPERS
// ============================================

// boolToInt : Convertit bool en 0/1 pour SQLite
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// toDateInt : Convertit time.Time en YYYYMMDD
func toDateInt(t time.Time) int {
	if t.IsZero() {
		return 0
	}
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}

// formatDateInt : YYYYMMDD → "2025-11-29"
func formatDateInt(dateInt int) string {
	if dateInt == 0 {
		return "N/A"
	}
	year := dateInt / 10000
	month := (dateInt % 10000) / 100
	day := dateInt % 100
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

// printStats : Affiche statistiques post-migration
func printStats(db *sql.DB) {
	log.Println("\n📊 STATISTIQUES POST-MIGRATION")
	log.Println("═══════════════════════════════")

	// Total par domaine
	rows, _ := db.Query(`
		SELECT domain, COUNT(*) as count, 
		       SUM(CASE WHEN done = 1 THEN 1 ELSE 0 END) as mastered
		FROM exercises 
		WHERE deleted = 0 
		GROUP BY domain
	`)
	defer rows.Close()

	log.Println("\n📚 Par domaine :")
	for rows.Next() {
		var domain string
		var total, mastered int
		rows.Scan(&domain, &total, &mastered)
		log.Printf("   %s: %d total (%d maîtrisés)", domain, total, mastered)
	}

	// Exercices dus aujourd'hui
	today := toDateInt(time.Now())
	var dueToday int
	db.QueryRow(`
		SELECT COUNT(*) FROM exercises 
		WHERE deleted = 0 AND done = 1 AND next_review_date <= ?
	`, today).Scan(&dueToday)

	log.Printf("\n⏰ Dus aujourd'hui : %d", dueToday)

	// Nouveaux exercices
	var newExercises int
	db.QueryRow(`
		SELECT COUNT(*) FROM exercises 
		WHERE deleted = 0 AND done = 0 AND last_reviewed_date IS NULL
	`).Scan(&newExercises)

	log.Printf("🆕 Nouveaux (jamais révisés) : %d", newExercises)

	// Prochaine date de révision
	var nextDate int
	err := db.QueryRow(`
		SELECT MIN(next_review_date) FROM exercises 
		WHERE deleted = 0 AND done = 1 AND next_review_date > ?
	`, today).Scan(&nextDate)

	if err == nil && nextDate > 0 {
		log.Printf("📅 Prochaine révision : %s", formatDateInt(nextDate))
	}

	log.Println("\n═══════════════════════════════\n")
}
