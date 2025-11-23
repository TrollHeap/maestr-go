package models

import "time"

type Exercise struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Domain         string    `json:"domain"`
	Content        string    `json:"content"`
	Difficulty     int       `json:"difficulty"`
	Steps          []string  `json:"steps"`
	CompletedSteps []int     `json:"completed_steps"`
	Done           bool      `json:"done"`
	CreatedAt      time.Time `json:"created_at"`
}

type ExerciseFilter struct {
	Domain     string // "Algorithmes", "Go", "" (tous)
	Status     string // "done", "todo", "" (tous)
	Difficulty int    // 0 (tous), 1-5
}
