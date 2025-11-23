package main

import "time"

// --- 1. MODÈLES & DONNÉES (Simulation DB) ---
type Exercise struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Domain         string    `json:"domain"`
	Difficulty     int       `json:"difficulty"`
	Steps          []string  `json:"steps"`
	CompletedSteps []int     `json:"completed_steps"`
	Content        string    `json:"content"`
	Done           bool      `json:"done"`
	CreatedAt      time.Time `json:"created_at"`
}
