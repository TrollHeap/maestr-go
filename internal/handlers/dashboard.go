package handlers

import (
	"log"
	"net/http"
	"time"

	"maestro/internal/service"
)

var dashboardService *service.DashboardService

func init() {
	dashboardService = service.NewDashboardService()
}

// HandleDashboard affiche le dashboard
func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	// Récupère les stats
	stats := dashboardService.GetDashboardStats()

	// Récupère les révisions d'aujourd'hui
	todayReviews := plannerService.GetReviewsForDate(time.Now())

	// Récupère les révisions en retard
	overdueReviews := plannerService.GetOverdueReviews()

	// Récupère les prochaines révisions
	upcomingReviews := plannerService.GetUpcomingReviews(5)

	data := map[string]any{
		"Stats":         stats,
		"TodayCount":    len(todayReviews),
		"OverdueCount":  len(overdueReviews),
		"UpcomingCount": len(upcomingReviews),
		"Overdue":       overdueReviews,
		"Upcoming":      upcomingReviews,
		"Now":           time.Now(),
	}

	if err := Tmpl.ExecuteTemplate(w, "dashboard", data); err != nil {
		log.Printf("❌ Erreur template dashboard: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
