package models

import "time"

// PlannedSession représente une session de travail planifiée
type PlannedSession struct {
	ID          string     `json:"id"`
	Date        time.Time  `json:"date"`         // Date de la session
	TimeSlot    string     `json:"time_slot"`    // morning, afternoon, evening
	ExerciseIDs []string   `json:"exercise_ids"` // Exercices planifiés
	Duration    int        `json:"duration"`     // En minutes
	Status      string     `json:"status"`       // planned, completed, skipped
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Notes       string     `json:"notes,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DailyPlan représente le plan d'une journée
type DailyPlan struct {
	Date         string           `json:"date"` // YYYY-MM-DD
	Sessions     []PlannedSession `json:"sessions"`
	TotalMinutes int              `json:"total_minutes"`
	Completed    int              `json:"completed"`
	Total        int              `json:"total"`
}

// WeeklyPlan représente le plan d'une semaine
type WeeklyPlan struct {
	StartDate    string      `json:"start_date"` // YYYY-MM-DD (Monday)
	EndDate      string      `json:"end_date"`   // YYYY-MM-DD (Sunday)
	Days         []DailyPlan `json:"days"`
	TotalMinutes int         `json:"total_minutes"`
	Completed    int         `json:"completed"`
	Total        int         `json:"total"`
}

// PlannerStats statistiques du planner
type PlannerStats struct {
	TodayPlanned    int     `json:"today_planned"`
	TodayCompleted  int     `json:"today_completed"`
	WeekPlanned     int     `json:"week_planned"`
	WeekCompleted   int     `json:"week_completed"`
	CompletionRate  float64 `json:"completion_rate"`  // %
	AverageDuration int     `json:"average_duration"` // minutes
}
