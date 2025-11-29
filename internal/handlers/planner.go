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

func HandlePlannerDay(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = time.Now()
	}

	data := map[string]any{
		"Date":    date,
		"Reviews": plannerService.GetReviewsForDate(date),
		"Count":   len(plannerService.GetReviewsForDate(date)),
	}

	Tmpl.ExecuteTemplate(w, "planner-day-view", data)
}

func HandlePlannerWeek(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")
	var startDate time.Time

	if weekStr != "" {
		parsed, err := time.Parse("2006-01-02", weekStr)
		if err != nil {
			log.Printf("⚠️  Erreur parse date semaine: %v (reçu: %s)", err, weekStr)
		} else {
			startDate = parsed
		}
	}

	if startDate.IsZero() {
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startDate = now.AddDate(0, 0, -(weekday - 1))
	}

	log.Printf("🔍 [Planner] Semaine début: %s", startDate.Format("2006-01-02"))

	data := map[string]any{
		"StartDate": startDate,
		"Schedule":  plannerService.GetWeekSchedule(startDate),
		"PrevWeek":  startDate.AddDate(0, 0, -7).Format("2006-01-02"),
		"NextWeek":  startDate.AddDate(0, 0, 7).Format("2006-01-02"),
	}

	if err := Tmpl.ExecuteTemplate(w, "planner-week-view", data); err != nil {
		log.Printf("❌ Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandlePlannerMonth(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	now := time.Now()
	year := now.Year()
	month := now.Month()

	if y, err := strconv.Atoi(yearStr); err == nil {
		year = y
	}
	if m, err := strconv.Atoi(monthStr); err == nil {
		month = time.Month(m)
	}

	counts := plannerService.GetMonthSchedule(year, month)

	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	daysInMonth := lastDay.Day()

	days := make([]int, daysInMonth)
	for i := 0; i < daysInMonth; i++ {
		days[i] = i + 1
	}

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
		"PrevYear":  year,
		"PrevMonth": int(month),
		"NextYear":  year,
		"NextMonth": int(month),
		"WeekDays":  []string{"Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"},
	}

	if err := Tmpl.ExecuteTemplate(w, "planner-month-view", data); err != nil {
		log.Printf("❌ Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
