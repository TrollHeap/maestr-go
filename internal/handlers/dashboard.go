package handlers

import (
	"log"
	"net/http"
	"time"

	"maestro/internal/service"
	"maestro/internal/views/pages"
)

var dashboardService *service.DashboardService

func init() {
	dashboardService = service.NewDashboardService()
	plannerService = service.NewPlannerService()
}

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 Dashboard: rendering with templ")

	// Récupère stats
	stats := dashboardService.GetDashboardStats()
	todayReviews := plannerService.GetReviewsForDate(time.Now())
	overdueReviews := plannerService.GetOverdueReviews()
	upcomingReviews := plannerService.GetUpcomingReviews(5)

	log.Printf("📊 Stats: today=%d, overdue=%d, upcoming=%d",
		len(todayReviews), len(overdueReviews), len(upcomingReviews))

	// ✅ Pass models.Exercise slices directement
	component := pages.Dashboard(
		stats,
		len(todayReviews),
		len(overdueReviews),
		len(upcomingReviews),
		overdueReviews,  // []models.Exercise
		upcomingReviews, // []models.Exercise
	)

	// Render component
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("❌ Error rendering dashboard: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Dashboard rendered successfully")
}
