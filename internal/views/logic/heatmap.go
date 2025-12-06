// internal/views/logic/heatmap.go
package logic

import (
	"time"
)

// HeatmapDay - Single day data
type HeatmapDay struct {
	Date      string // ✅ KEEP as string for map lookup
	Count     int
	DayOfWeek int
	IsToday   bool
}

// GenerateHeatmapDays - Generate heatmap data for N weeks
func GenerateHeatmapDays(reviewCounts map[string]int, weeks int) []HeatmapDay {
	now := time.Now()
	days := []HeatmapDay{}

	// Start from N weeks ago
	startDate := now.AddDate(0, 0, -(weeks * 7))

	// Align to Sunday
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	// Generate all days
	totalDays := weeks * 7
	for i := 0; i < totalDays; i++ {
		date := startDate.AddDate(0, 0, i)
		dateKey := date.Format("2006-01-02")

		day := HeatmapDay{
			Date:      dateKey, // ✅ Store as string
			Count:     reviewCounts[dateKey],
			DayOfWeek: int(date.Weekday()),
			IsToday:   isSameDay(date, now),
		}

		days = append(days, day)
	}

	return days
}

func isSameDay(d1, d2 time.Time) bool {
	return d1.Year() == d2.Year() &&
		d1.Month() == d2.Month() &&
		d1.Day() == d2.Day()
}
