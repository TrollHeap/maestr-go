package handlers

import (
	"log"
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
// HandlePlannerWeek affiche la semaine (Fragment HTMX)
func HandlePlannerWeek(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")

	var startDate time.Time
	if weekStr != "" {
		parsed, err := time.Parse("2006-01-02", weekStr)
		if err == nil {
			startDate = parsed
		} else {
			log.Printf("⚠️  Erreur parse date semaine: %v (reçu: %s)", err, weekStr)
		}
	}

	// Si pas de date ou erreur de parse, commence au lundi actuel
	if startDate.IsZero() {
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Dimanche = 7
		}
		startDate = now.AddDate(0, 0, -(weekday - 1))
	}

	log.Printf("🔍 [Planner] Semaine demandée: %s → StartDate: %s",
		weekStr, startDate.Format("2006-01-02"))

	schedule := plannerService.GetWeekSchedule(startDate)

	// Calcule les dates de navigation
	prevWeek := startDate.AddDate(0, 0, -7).Format("2006-01-02")
	nextWeek := startDate.AddDate(0, 0, 7).Format("2006-01-02")

	log.Printf("📅 Navigation: Prev=%s, Current=%s, Next=%s",
		prevWeek, startDate.Format("2006-01-02"), nextWeek)

	data := map[string]any{
		"StartDate": startDate,
		"Schedule":  schedule,
		"PrevWeek":  prevWeek,
		"NextWeek":  nextWeek,
	}

	if err := Tmpl.ExecuteTemplate(w, "planner-week-view", data); err != nil {
		log.Printf("❌ Erreur template planner-week-view: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	// Calcule le nombre de jours dans le mois
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	daysInMonth := lastDay.Day()

	// Crée la slice des jours
	days := make([]int, daysInMonth)
	for i := range daysInMonth {
		days[i] = i + 1
	}

	// Navigation : mois précédent/suivant
	prevMonth := month - 1
	prevYear := year
	if prevMonth < 1 {
		prevMonth = 12
		prevYear--
	}

	nextMonth := month + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}

	// Nom du mois en français
	monthNames := map[time.Month]string{
		time.January:   "Janvier",
		time.February:  "Février",
		time.March:     "Mars",
		time.April:     "Avril",
		time.May:       "Mai",
		time.June:      "Juin",
		time.July:      "Juillet",
		time.August:    "Août",
		time.September: "Septembre",
		time.October:   "Octobre",
		time.November:  "Novembre",
		time.December:  "Décembre",
	}

	data := map[string]any{
		"Year":      year,
		"Month":     int(month),
		"MonthName": monthNames[month],
		"Counts":    counts,
		"Days":      days,
		"PrevYear":  prevYear,
		"PrevMonth": int(prevMonth),
		"NextYear":  nextYear,
		"NextMonth": int(nextMonth),
		"WeekDays":  []string{"Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"},
	}

	Tmpl.ExecuteTemplate(w, "planner-month-view", data)
}
