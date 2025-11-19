package domain

import (
	"time"

	"maestro/internal/models"
)

// Planner gère la planification des sessions d'apprentissage
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

// GetToday retourne le plan du jour
func (p *Planner) GetToday() models.DailyPlan {
	today := time.Now().Format("2006-01-02")

	todaySessions := []models.PlannedSession{}
	totalMinutes := 0
	completed := 0

	for _, session := range p.sessions {
		sessionDate := session.Date.Format("2006-01-02")
		if sessionDate == today {
			todaySessions = append(todaySessions, session)
			totalMinutes += session.Duration
			if session.Status == "completed" {
				completed++
			}
		}
	}

	return models.DailyPlan{
		Date:         today,
		Sessions:     todaySessions,
		TotalMinutes: totalMinutes,
		Completed:    completed,
		Total:        len(todaySessions),
	}
}

// GetWeek retourne le plan de la semaine
func (p *Planner) GetWeek() models.WeeklyPlan {
	now := time.Now()

	// Trouver le lundi de cette semaine
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Dimanche = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))
	sunday := monday.AddDate(0, 0, 6)

	startDate := monday.Format("2006-01-02")
	endDate := sunday.Format("2006-01-02")

	// Créer les jours de la semaine
	days := []models.DailyPlan{}
	totalMinutes := 0
	totalCompleted := 0
	totalSessions := 0

	for i := 0; i < 7; i++ {
		currentDay := monday.AddDate(0, 0, i)
		dayStr := currentDay.Format("2006-01-02")

		daySessions := []models.PlannedSession{}
		dayMinutes := 0
		dayCompleted := 0

		for _, session := range p.sessions {
			sessionDate := session.Date.Format("2006-01-02")
			if sessionDate == dayStr {
				daySessions = append(daySessions, session)
				dayMinutes += session.Duration
				if session.Status == "completed" {
					dayCompleted++
				}
			}
		}

		days = append(days, models.DailyPlan{
			Date:         dayStr,
			Sessions:     daySessions,
			TotalMinutes: dayMinutes,
			Completed:    dayCompleted,
			Total:        len(daySessions),
		})

		totalMinutes += dayMinutes
		totalCompleted += dayCompleted
		totalSessions += len(daySessions)
	}

	return models.WeeklyPlan{
		StartDate:    startDate,
		EndDate:      endDate,
		Days:         days,
		TotalMinutes: totalMinutes,
		Completed:    totalCompleted,
		Total:        totalSessions,
	}
}

// UpdateSession met à jour une session existante
func (p *Planner) UpdateSession(session models.PlannedSession) error {
	for i, s := range p.sessions {
		if s.ID == session.ID {
			p.sessions[i] = session
			return nil
		}
	}
	return ErrSessionNotFound
}

// DeleteSession supprime une session
func (p *Planner) DeleteSession(sessionID string) error {
	for i, s := range p.sessions {
		if s.ID == sessionID {
			p.sessions = append(p.sessions[:i], p.sessions[i+1:]...)
			return nil
		}
	}
	return ErrSessionNotFound
}

// GetStats retourne les statistiques du planner
func (p *Planner) GetStats() models.PlannerStats {
	now := time.Now()
	today := now.Format("2006-01-02")

	// Stats du jour
	todayPlanned := 0
	todayCompleted := 0

	// Stats de la semaine
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))
	weekStart := monday.Format("2006-01-02")

	weekPlanned := 0
	weekCompleted := 0
	totalDuration := 0
	sessionCount := 0

	for _, session := range p.sessions {
		sessionDate := session.Date.Format("2006-01-02")

		// Aujourd'hui
		if sessionDate == today {
			todayPlanned++
			if session.Status == "completed" {
				todayCompleted++
			}
		}

		// Cette semaine
		if sessionDate >= weekStart {
			weekPlanned++
			if session.Status == "completed" {
				weekCompleted++
				totalDuration += session.Duration
				sessionCount++
			}
		}
	}

	// Calculer le taux de complétion
	completionRate := 0.0
	if weekPlanned > 0 {
		completionRate = float64(weekCompleted) / float64(weekPlanned) * 100
	}

	// Calculer la durée moyenne
	averageDuration := 0
	if sessionCount > 0 {
		averageDuration = totalDuration / sessionCount
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

// ErrSessionNotFound est retourné quand une session n'existe pas
var ErrSessionNotFound = error(nil)

func init() {
	ErrSessionNotFound = &PlannerError{Message: "session not found"}
}

// PlannerError représente une erreur du planner
type PlannerError struct {
	Message string
}

func (e *PlannerError) Error() string {
	return e.Message
}
