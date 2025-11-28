package handlers

import (
	"net/http"
	"strconv"
	"time"

	"maestro/internal/service"
)

var plannerService *service.PlannerService

func init() {
	plannerService = service.NewPlannerService()
}

// HandlePlannerPage affiche la page principale du planner
func HandlePlannerPage(w http.ResponseWriter, r *http.Request) {
	today := time.Now()

	data := map[string]any{
		"CurrentDate": today,
		"Reviews":     plannerService.GetReviewsForDate(today),
		"Upcoming":    plannerService.GetUpcomingReviews(10),
		"Overdue":     plannerService.GetOverdueReviews(),
	}

	if err := Tmpl.ExecuteTemplate(w, "planner-page", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandlePlannerDay affiche les révisions d'un jour (Fragment HTMX)
func HandlePlannerDay(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = time.Now()
	}

	reviews := plannerService.GetReviewsForDate(date)

	data := map[string]any{
		"Date":    date,
		"Reviews": reviews,
		"Count":   len(reviews),
	}

	Tmpl.ExecuteTemplate(w, "planner-day-view", data)
}

// HandlePlannerWeek affiche la semaine (Fragment HTMX)
func HandlePlannerWeek(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")

	var startDate time.Time
	if weekStr != "" {
		parsed, err := time.Parse("2006-01-02", weekStr)
		if err == nil {
			startDate = parsed
		}
	}

	if startDate.IsZero() {
		// Commence au lundi de cette semaine
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Dimanche
		}
		startDate = now.AddDate(0, 0, -(weekday - 1))
	}

	schedule := plannerService.GetWeekSchedule(startDate)

	data := map[string]any{
		"StartDate": startDate,
		"Schedule":  schedule,
	}

	Tmpl.ExecuteTemplate(w, "planner-week-view", data)
}

// HandlePlannerMonth affiche le mois (Fragment HTMX)
func HandlePlannerMonth(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	now := time.Now()
	year := now.Year()
	month := now.Month()

	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	if monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil {
			month = time.Month(m)
		}
	}

	counts := plannerService.GetMonthSchedule(year, month)

	// ✅ Calcule le nombre de jours dans le mois
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	daysInMonth := lastDay.Day()

	// Crée la slice des jours
	days := make([]int, daysInMonth)
	for i := 0; i < daysInMonth; i++ {
		days[i] = i + 1
	}

	data := map[string]any{
		"Year":   year,
		"Month":  month,
		"Counts": counts,
		"Days":   days, // ← AJOUTE
	}

	Tmpl.ExecuteTemplate(w, "planner-month-view", data)
}
