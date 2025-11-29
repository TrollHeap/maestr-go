package models

import "time"

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// MODÈLES PLANNER
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

// Task représente une tâche
type Task struct {
	Title     string
	Priority  string // "high", "medium", "low"
	Completed bool
}

// TimeSlot représente un créneau horaire avec ses tâches
type TimeSlot struct {
	StartTime string
	EndTime   string
	Tasks     []Task
}

// DayData contient les données pour la vue JOUR
type DayData struct {
	Date      time.Time
	TimeSlots []TimeSlot
}

// DaySchedule représente les exercices et données d'un jour
type DaySchedule struct {
	Date      time.Time
	Exercises []Exercise
	Count     int
}

// WeekDay représente un jour dans la vue SEMAINE
type WeekDay struct {
	DayName   string
	DayNumber int
	IsToday   bool
	Tasks     []Task
}

// WeekData contient les données pour la vue SEMAINE
type WeekData struct {
	WeekNumber int
	StartDate  time.Time
	EndDate    time.Time // ← Vérifie que ce champ existe
	Days       []WeekDay
}

// MonthDay représente un jour dans le calendrier mensuel
type MonthDay struct {
	Number       int
	IsToday      bool
	IsOtherMonth bool
	TaskCount    int
}

// MonthData contient les données pour la vue MOIS
type MonthData struct {
	Month    string
	Year     int
	MonthNum int
	Days     []MonthDay // ← AJOUTE CE CHAMP ICI
}
