package domain

import (
	"crypto/rand"
	"encoding/hex"
	"sort"
	"time"

	"maestro/internal/models"
)

// Planner gère la planification des sessions
type Planner struct {
	sessions []models.PlannedSession
}

// NewPlanner crée un nouveau planner
func NewPlanner() *Planner {
	return &Planner{
		sessions: []models.PlannedSession{},
	}
}

// LoadSessions charge les sessions depuis le storage
func (p *Planner) LoadSessions(sessions []models.PlannedSession) {
	p.sessions = sessions
}

// GetSessions retourne toutes les sessions
func (p *Planner) GetSessions() []models.PlannedSession {
	return p.sessions
}

// CreateSession crée une nouvelle session planifiée
func (p *Planner) CreateSession(
	date time.Time,
	timeSlot string,
	exerciseIDs []string,
	duration int,
) models.PlannedSession {
	now := time.Now()

	session := models.PlannedSession{
		ID:          generateID(),
		Date:        date,
		TimeSlot:    timeSlot,
		ExerciseIDs: exerciseIDs,
		Duration:    duration,
		Status:      "planned",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	p.sessions = append(p.sessions, session)
	return session
}

// UpdateSession met à jour une session
func (p *Planner) UpdateSession(sessionID string, status string, notes string) error {
	for i := range p.sessions {
		if p.sessions[i].ID == sessionID {
			p.sessions[i].Status = status
			p.sessions[i].Notes = notes
			p.sessions[i].UpdatedAt = time.Now()

			if status == "completed" {
				now := time.Now()
				p.sessions[i].CompletedAt = &now
			}

			return nil
		}
	}
	return nil
}

// DeleteSession supprime une session
func (p *Planner) DeleteSession(sessionID string) error {
	for i, session := range p.sessions {
		if session.ID == sessionID {
			p.sessions = append(p.sessions[:i], p.sessions[i+1:]...)
			return nil
		}
	}
	return nil
}

// GetDailyPlan retourne le plan d'une journée
func (p *Planner) GetDailyPlan(date time.Time) models.DailyPlan {
	dateKey := date.Format("2006-01-02")

	// Filter sessions for this day
	daySessions := []models.PlannedSession{}
	totalMinutes := 0
	completed := 0

	for _, session := range p.sessions {
		if session.Date.Format("2006-01-02") == dateKey {
			daySessions = append(daySessions, session)
			totalMinutes += session.Duration
			if session.Status == "completed" {
				completed++
			}
		}
	}

	// Sort by time slot
	sort.Slice(daySessions, func(i, j int) bool {
		order := map[string]int{"morning": 1, "afternoon": 2, "evening": 3}
		return order[daySessions[i].TimeSlot] < order[daySessions[j].TimeSlot]
	})

	return models.DailyPlan{
		Date:         dateKey,
		Sessions:     daySessions,
		TotalMinutes: totalMinutes,
		Completed:    completed,
		Total:        len(daySessions),
	}
}

// GetWeeklyPlan retourne le plan d'une semaine
func (p *Planner) GetWeeklyPlan(startDate time.Time) models.WeeklyPlan {
	// Ensure start is Monday
	weekday := int(startDate.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday
	}
	monday := startDate.AddDate(0, 0, -(weekday - 1))
	sunday := monday.AddDate(0, 0, 6)

	days := []models.DailyPlan{}
	totalMinutes := 0
	totalCompleted := 0
	totalSessions := 0

	// Generate daily plans for each day of week
	for d := monday; !d.After(sunday); d = d.AddDate(0, 0, 1) {
		dailyPlan := p.GetDailyPlan(d)
		days = append(days, dailyPlan)
		totalMinutes += dailyPlan.TotalMinutes
		totalCompleted += dailyPlan.Completed
		totalSessions += dailyPlan.Total
	}

	return models.WeeklyPlan{
		StartDate:    monday.Format("2006-01-02"),
		EndDate:      sunday.Format("2006-01-02"),
		Days:         days,
		TotalMinutes: totalMinutes,
		Completed:    totalCompleted,
		Total:        totalSessions,
	}
}

// GetMonthlyPlan retourne le plan d'un mois (4 semaines)
func (p *Planner) GetMonthlyPlan(startDate time.Time) []models.WeeklyPlan {
	weeks := []models.WeeklyPlan{}

	for i := 0; i < 4; i++ {
		weekStart := startDate.AddDate(0, 0, i*7)
		week := p.GetWeeklyPlan(weekStart)
		weeks = append(weeks, week)
	}

	return weeks
}

// GetStats retourne les statistiques du planner
func (p *Planner) GetStats() models.PlannerStats {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Get week range
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := today.AddDate(0, 0, -(weekday - 1))
	sunday := monday.AddDate(0, 0, 6)

	todayPlanned := 0
	todayCompleted := 0
	weekPlanned := 0
	weekCompleted := 0
	totalDuration := 0
	totalSessions := 0

	for _, session := range p.sessions {
		sessionDate := time.Date(
			session.Date.Year(),
			session.Date.Month(),
			session.Date.Day(),
			0,
			0,
			0,
			0,
			time.UTC,
		)

		// Today stats
		if sessionDate.Equal(today) {
			todayPlanned++
			if session.Status == "completed" {
				todayCompleted++
			}
		}

		// Week stats
		if (sessionDate.Equal(monday) || sessionDate.After(monday)) &&
			sessionDate.Before(sunday.AddDate(0, 0, 1)) {
			weekPlanned++
			if session.Status == "completed" {
				weekCompleted++
				totalDuration += session.Duration
				totalSessions++
			}
		}
	}

	completionRate := 0.0
	if weekPlanned > 0 {
		completionRate = float64(weekCompleted) / float64(weekPlanned) * 100
	}

	averageDuration := 0
	if totalSessions > 0 {
		averageDuration = totalDuration / totalSessions
	}

	return models.PlannerStats{
		TodayPlanned:    todayPlanned,
		TodayCompleted:  todayCompleted,
		WeekPlanned:     weekPlanned,
		WeekCompleted:   weekCompleted,
		CompletionRate:  completionRate,
		AverageDuration: averageDuration,
	}
}

// GetTodaySessions retourne les sessions du jour
func (p *Planner) GetTodaySessions() []models.PlannedSession {
	today := time.Now()
	return p.GetSessionsByDate(today)
}

// GetSessionsByDate retourne les sessions d'une date
func (p *Planner) GetSessionsByDate(date time.Time) []models.PlannedSession {
	dateKey := date.Format("2006-01-02")
	sessions := []models.PlannedSession{}

	for _, session := range p.sessions {
		if session.Date.Format("2006-01-02") == dateKey {
			sessions = append(sessions, session)
		}
	}

	// Sort by time slot
	sort.Slice(sessions, func(i, j int) bool {
		order := map[string]int{"morning": 1, "afternoon": 2, "evening": 3}
		return order[sessions[i].TimeSlot] < order[sessions[j].TimeSlot]
	})

	return sessions
}

// generateID génère un ID unique
func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
