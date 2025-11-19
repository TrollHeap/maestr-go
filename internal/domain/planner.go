package domain

import (
	"fmt"
	"time"

	"maestro/internal/models"
)

// Planner gère la planification des sessions d'exercices
type Planner struct {
	sessions []models.PlannedSession
}

// NewPlanner crée un nouveau Planner
func NewPlanner() *Planner {
	return &Planner{
		sessions: []models.PlannedSession{},
	}
}

// AddSession ajoute une session planifiée
func (p *Planner) AddSession(session models.PlannedSession) {
	p.sessions = append(p.sessions, session)
}

// UpdateSession met à jour une session existante
func (p *Planner) UpdateSession(session models.PlannedSession) error {
	for i := range p.sessions {
		if p.sessions[i].ID == session.ID {
			p.sessions[i] = session
			return nil
		}
	}
	return fmt.Errorf("session not found: %s", session.ID)
}

// DeleteSession supprime une session
func (p *Planner) DeleteSession(id string) error {
	for i := range p.sessions {
		if p.sessions[i].ID == id {
			p.sessions = append(p.sessions[:i], p.sessions[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("session not found: %s", id)
}

// ✅ HELPER: Filtrer par date exacte
func (p *Planner) filterSessionsByDate(date time.Time) []models.PlannedSession {
	dateKey := date.Format("2006-01-02")
	filtered := []models.PlannedSession{}

	for _, session := range p.sessions {
		if session.Date.Format("2006-01-02") == dateKey {
			filtered = append(filtered, session)
		}
	}

	return filtered
}

// ✅ HELPER: Filtrer par plage de dates
func (p *Planner) filterSessionsByDateRange(start, end time.Time) []models.PlannedSession {
	filtered := []models.PlannedSession{}

	startDate := start.Truncate(24 * time.Hour)
	endDate := end.Truncate(24 * time.Hour)

	for _, session := range p.sessions {
		sessionDate := session.Date.Truncate(24 * time.Hour)

		if (sessionDate.Equal(startDate) || sessionDate.After(startDate)) &&
			sessionDate.Before(endDate.AddDate(0, 0, 1)) {
			filtered = append(filtered, session)
		}
	}

	return filtered
}

// ✅ HELPER: Compter sessions complétées
func countCompleted(sessions []models.PlannedSession) int {
	count := 0
	for _, s := range sessions {
		if s.Status == "completed" {
			count++
		}
	}
	return count
}

// ✅ HELPER: Calculer durée totale
func totalDuration(sessions []models.PlannedSession) int {
	total := 0
	for _, s := range sessions {
		total += s.Duration
	}
	return total
}

// ✅ HELPER: Obtenir lundi de la semaine
func getMonday(date time.Time) time.Time {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7 // Dimanche = 7
	}
	return date.AddDate(0, 0, -(weekday - 1))
}

// GetToday retourne le plan du jour
func (p *Planner) GetToday() models.DailyPlan {
	today := time.Now().Truncate(24 * time.Hour)
	sessions := p.filterSessionsByDate(today)

	return models.DailyPlan{
		Date:         today.Format("2006-01-02"),
		Sessions:     sessions,
		TotalMinutes: totalDuration(sessions),
		Completed:    countCompleted(sessions),
		Total:        len(sessions),
	}
}

// GetWeek retourne le plan de la semaine
func (p *Planner) GetWeek() models.WeeklyPlan {
	today := time.Now().Truncate(24 * time.Hour)
	monday := getMonday(today)
	sunday := monday.AddDate(0, 0, 6)

	weekSessions := p.filterSessionsByDateRange(monday, sunday)

	// Générer les plans journaliers
	days := []models.DailyPlan{}
	for i := 0; i < 7; i++ {
		date := monday.AddDate(0, 0, i)
		daySessions := p.filterSessionsByDate(date)

		days = append(days, models.DailyPlan{
			Date:         date.Format("2006-01-02"),
			Sessions:     daySessions,
			TotalMinutes: totalDuration(daySessions),
			Completed:    countCompleted(daySessions),
			Total:        len(daySessions),
		})
	}

	return models.WeeklyPlan{
		StartDate:    monday.Format("2006-01-02"),
		EndDate:      sunday.Format("2006-01-02"),
		Days:         days,
		TotalMinutes: totalDuration(weekSessions),
		Completed:    countCompleted(weekSessions),
		Total:        len(weekSessions),
	}
}

// GetStats retourne les statistiques du planner (✅ refactoré)
func (p *Planner) GetStats() models.PlannerStats {
	today := time.Now().Truncate(24 * time.Hour)
	monday := getMonday(today)
	sunday := monday.AddDate(0, 0, 6)

	// ✅ Utiliser helpers
	todaySessions := p.filterSessionsByDate(today)
	weekSessions := p.filterSessionsByDateRange(monday, sunday)

	todayCompleted := countCompleted(todaySessions)
	weekCompleted := countCompleted(weekSessions)

	// Calculer taux de complétion
	var completionRate float64
	if len(weekSessions) > 0 {
		completionRate = float64(weekCompleted) / float64(len(weekSessions)) * 100
	}

	// Calculer durée moyenne
	var avgDuration int
	if weekCompleted > 0 {
		totalTime := 0
		for _, s := range weekSessions {
			if s.Status == "completed" {
				totalTime += s.Duration
			}
		}
		avgDuration = totalTime / weekCompleted
	}

	return models.PlannerStats{
		TodayPlanned:    len(todaySessions),
		TodayCompleted:  todayCompleted,
		WeekPlanned:     len(weekSessions),
		WeekCompleted:   weekCompleted,
		CompletionRate:  completionRate,
		AverageDuration: avgDuration,
	}
}
