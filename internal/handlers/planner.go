package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"maestro/internal/service"
	"maestro/internal/views/components"
	"maestro/internal/views/pages"
)

// ============================================
// SERVICE GLOBAL
// ============================================

var plannerService *service.PlannerService

func init() {
	plannerService = service.NewPlannerService()
}

// ============================================
// 1️⃣ PAGE PRINCIPALE PLANNER
// ============================================
//

// func HandlePlannerPage(w http.ResponseWriter, r *http.Request) {
// 	view := r.URL.Query().Get("view") // "urgent", "today", "upcoming", "new", "all"
//
// 	if view == "" {
// 		view = "all"
// 	}
//
// 	// Récupère via service Planner
// 	exercises, err := exerciseService.GetPlannerExercises(view)
// 	if err != nil {
// 		log.Printf("❌ GetPlannerExercises error: %v", err)
// 		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
// 		return
// 	}
//
// 	// Stats temps
// 	stats := map[string]int{
// 		"urgent":   store.CountPlannerView("urgent"),
// 		"today":    store.CountPlannerView("today"),
// 		"upcoming": store.CountPlannerView("upcoming"),
// 		"new":      store.CountPlannerView("new"),
// 	}
//
// 	component := pages.PlannerPage(exercises, stats, view)
//
// 	if err := component.Render(r.Context(), w); err != nil {
// 		log.Printf("❌ Render error: %v", err)
// 		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
// 	}
// }

func HandlePlannerPage(w http.ResponseWriter, r *http.Request) {
	today := time.Now()

	log.Printf("🔍 PlannerPage: date=%s", today.Format("2006-01-02"))

	// 1. Récupère données (LOGIQUE IDENTIQUE)
	reviews := plannerService.GetReviewsForDate(today)
	upcoming := plannerService.GetUpcomingReviews(10)
	overdue := plannerService.GetOverdueReviews()

	log.Printf("✅ Planner data: reviews=%d, upcoming=%d, overdue=%d",
		len(reviews), len(upcoming), len(overdue))

	// 2. ✅ CHANGEMENT : Render avec templ
	component := pages.PlannerPage(today, reviews, upcoming, overdue)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// 2️⃣ FRAGMENT : Vue jour (HTMX)
// ============================================

func HandlePlannerDay(w http.ResponseWriter, r *http.Request) {
	// 1. Parse date query param (LOGIQUE IDENTIQUE)
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("⚠️ Invalid date '%s', using today", dateStr)
		date = time.Now()
	}

	log.Printf("🔍 PlannerDay: date=%s", date.Format("2006-01-02"))

	// 2. Récupère reviews pour ce jour (LOGIQUE IDENTIQUE)
	reviews := plannerService.GetReviewsForDate(date)

	log.Printf("✅ Day view: %d reviews", len(reviews))

	// 3. ✅ CHANGEMENT : Render fragment templ
	component := components.PlannerDayView(date, reviews)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// 3️⃣ FRAGMENT : Vue semaine (HTMX)
// ============================================

func HandlePlannerWeek(w http.ResponseWriter, r *http.Request) {
	// 1. Parse week query param (LOGIQUE IDENTIQUE)
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

	// 2. Fallback : début de semaine courante (lundi) (LOGIQUE IDENTIQUE)
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

	// 3. Récupère schedule semaine (LOGIQUE IDENTIQUE)
	schedule := plannerService.GetWeekSchedule(startDate)

	// 4. Calcule prev/next week (LOGIQUE IDENTIQUE)
	prevWeek := startDate.AddDate(0, 0, -7)
	nextWeek := startDate.AddDate(0, 0, 7)

	log.Printf("✅ Week view: %d days scheduled", len(schedule))

	// 5. ✅ CHANGEMENT : Render fragment templ
	component := components.PlannerWeekView(startDate, schedule, prevWeek, nextWeek)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// 4️⃣ FRAGMENT : Vue mois (HTMX)
// ============================================

func HandlePlannerMonth(w http.ResponseWriter, r *http.Request) {
	// 1. Parse year/month query params (LOGIQUE IDENTIQUE)
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

	// 2. Récupère counts mois (LOGIQUE IDENTIQUE)
	counts := plannerService.GetMonthSchedule(year, month)

	// 3. Calcule métadonnées mois (LOGIQUE IDENTIQUE)
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)
	daysInMonth := lastDay.Day()

	// 4. Calcule prev/next month (LOGIQUE IDENTIQUE)
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

	monthName := getMonthName(month)

	log.Printf("✅ Month view: %s %d (%d days)", monthName, year, daysInMonth)

	// 5. ✅ CHANGEMENT : Render fragment templ
	component := components.PlannerMonthView(
		year, month, monthName,
		counts, firstDay, daysInMonth,
		prevYear, prevMonth, nextYear, nextMonth,
	)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// HELPERS
// ============================================

func getMonthName(month time.Month) string {
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
	return monthNames[month]
}
