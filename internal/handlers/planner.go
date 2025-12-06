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

func HandlePlannerPage(w http.ResponseWriter, r *http.Request) {
	today := time.Now()

	log.Printf("🔍 PlannerPage: date=%s", today.Format("2006-01-02"))

	// 1. Récupère données
	reviews := plannerService.GetReviewsForDate(today)
	upcoming := plannerService.GetUpcomingReviews(10)
	overdue := plannerService.GetOverdueReviews()

	log.Printf("✅ Planner data: reviews=%d, upcoming=%d, overdue=%d",
		len(reviews), len(upcoming), len(overdue))

	// 2. Render page complète
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
	// Parse date
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("⚠️ Invalid date '%s', using today", dateStr)
		date = time.Now()
	}

	log.Printf("🔍 PlannerDay: date=%s", date.Format("2006-01-02"))

	// Récupère reviews
	reviews := plannerService.GetReviewsForDate(date)

	log.Printf("✅ Day view: %d reviews", len(reviews))

	// Render fragment (OLD component - keep as is)
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
	// Parse week query param
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

	// Fallback : début de semaine courante (lundi)
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

	// ✅ FIX: Use NEW Enhanced component (single param)
	component := components.PlannerWeekView(startDate)

	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Render error: %v", err)
		http.Error(w, "Erreur affichage", http.StatusInternalServerError)
	}
}

// ============================================
// 4️⃣ FRAGMENT : Vue mois (HTMX)
// ============================================

func HandlePlannerMonth(w http.ResponseWriter, r *http.Request) {
	// Parse year/month query params
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

	// Build date for this month
	currentDate := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)

	log.Printf("✅ Month view: %s", currentDate.Format("January 2006"))

	// ✅ FIX: Use NEW Enhanced component (single param)
	component := components.PlannerMonthView(currentDate)

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
