package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"maestro/internal/service"
)

// ============================================
// SERVICE GLOBAL
// ============================================

func init() {
	plannerService = service.NewPlannerService()
}

// ============================================
// 1️⃣ PAGE PRINCIPALE PLANNER
// ============================================

func HandlePlannerPage(w http.ResponseWriter, r *http.Request) {
	today := time.Now()

	log.Printf("🔍 PlannerPage: date=%s", today.Format("2006-01-02"))

	// 1. Récupère données
	reviews := plannerService.GetReviewsForDate(today)
	upcoming := plannerService.GetUpcomingReviews(10)
	overdue := plannerService.GetOverdueReviews()

	// 2. Structure données
	data := map[string]any{
		"CurrentDate": today,
		"Reviews":     reviews,
		"Upcoming":    upcoming,
		"Overdue":     overdue,
	}

	log.Printf("✅ Planner data: reviews=%d, upcoming=%d, overdue=%d",
		len(reviews), len(upcoming), len(overdue))

	// 3. ✅ Render avec helper
	RenderTemplateOrError(w, "planner.html", data)
}

// ============================================
// 2️⃣ FRAGMENT : Vue jour (HTMX)
// ============================================

func HandlePlannerDay(w http.ResponseWriter, r *http.Request) {
	// 1. Parse date query param
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("⚠️ Invalid date '%s', using today", dateStr)
		date = time.Now()
	}

	log.Printf("🔍 PlannerDay: date=%s", date.Format("2006-01-02"))

	// 2. Récupère reviews pour ce jour
	reviews := plannerService.GetReviewsForDate(date)

	// 3. Structure données
	data := map[string]any{
		"Date":    date,
		"Reviews": reviews,
		"Count":   len(reviews),
	}

	log.Printf("✅ Day view: %d reviews", len(reviews))

	// 4. ✅ Render fragment
	RenderTemplateOrError(w, "planner-day-view", data)
}

// ============================================
// 3️⃣ FRAGMENT : Vue semaine (HTMX)
// ============================================

func HandlePlannerWeek(w http.ResponseWriter, r *http.Request) {
	// 1. Parse week query param
	weekStr := r.URL.Query().Get("week")
	var startDate time.Time

	if weekStr != "" {
		parsed, err := time.Parse("2006-01-02", weekStr)
		if err != nil {
			log.Printf("⚠️ Invalid week date '%s': %v", weekStr, err)
		} else {
			startDate = parsed
		}
	}

	// 2. Fallback : début de semaine courante (lundi)
	if startDate.IsZero() {
		now := time.Now()
		weekday := int(now.Weekday())

		// Dimanche = 0 → 7 (ISO week)
		if weekday == 0 {
			weekday = 7
		}

		// Calcule lundi de la semaine
		startDate = now.AddDate(0, 0, -(weekday - 1))
	}

	log.Printf("🔍 PlannerWeek: startDate=%s", startDate.Format("2006-01-02"))

	// 3. Récupère schedule semaine
	schedule := plannerService.GetWeekSchedule(startDate)

	// 4. Calcule prev/next week
	prevWeek := startDate.AddDate(0, 0, -7)
	nextWeek := startDate.AddDate(0, 0, 7)

	// 5. Structure données
	data := map[string]any{
		"StartDate": startDate,
		"Schedule":  schedule,
		"PrevWeek":  prevWeek.Format("2006-01-02"),
		"NextWeek":  nextWeek.Format("2006-01-02"),
	}

	log.Printf("✅ Week view: %d days scheduled", len(schedule))

	// 6. ✅ Render fragment
	RenderTemplateOrError(w, "planner-week-view", data)
}

// ============================================
// 4️⃣ FRAGMENT : Vue mois (HTMX)
// ============================================

func HandlePlannerMonth(w http.ResponseWriter, r *http.Request) {
	// 1. Parse year/month query params
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	now := time.Now()
	year := now.Year()
	month := now.Month()

	if y, err := strconv.Atoi(yearStr); err == nil && y > 0 {
		year = y
	}
	if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
		month = time.Month(m)
	}

	log.Printf("🔍 PlannerMonth: year=%d, month=%d", year, month)

	// 2. Récupère counts mois
	counts := plannerService.GetMonthSchedule(year, month)

	// 3. Calcule métadonnées mois
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)
	daysInMonth := lastDay.Day()

	// 4. Génère liste jours [1, 2, ..., 31]
	days := make([]int, daysInMonth)
	for i := range daysInMonth {
		days[i] = i + 1
	}

	// 5. Noms mois en français
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

	// 6. Calcule prev/next month
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

	// 7. Structure données
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
		"FirstDay":  firstDay,
	}

	log.Printf("✅ Month view: %s %d (%d days)", monthNames[month], year, daysInMonth)

	// 8. ✅ Render fragment
	RenderTemplateOrError(w, "planner-month-view", data)
}
