package logic

import (
	"fmt"
	"time"
)

// HeatmapDay représente un jour dans le calendrier
type HeatmapDay struct {
	Date    time.Time
	Count   int
	Weekday string
}

// GetWeekdayLabel retourne le label court du jour (Mon, Tue, etc.)
func GetWeekdayLabel(weekday int) string {
	labels := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	if weekday < 0 || weekday >= len(labels) {
		return ""
	}
	return labels[weekday]
}

// GetHeatmapCellClasses retourne les classes Tailwind pour une cellule selon le count
func GetHeatmapCellClasses(count int) string {
	base := "w-full h-3 rounded transition-all hover:scale-110 cursor-pointer"

	switch {
	case count == 0:
		return base + " border border-slate-800 bg-slate-900/40"
	case count <= 2:
		return base + " bg-emerald-900/40 border border-emerald-800/60"
	case count <= 5:
		return base + " bg-emerald-700/60 border border-emerald-600/60"
	case count <= 10:
		return base + " bg-emerald-500/80 border border-emerald-400/60 shadow-[0_0_4px_rgba(52,211,153,0.4)]"
	default:
		return base + " bg-emerald-400 border border-emerald-300 shadow-[0_0_8px_rgba(52,211,153,0.6)]"
	}
}

// GetDayForWeek récupère le jour correspondant dans la slice
func GetDayForWeek(days []HeatmapDay, weekday, week int) HeatmapDay {
	index := weekday + (week * 7)
	if index < 0 || index >= len(days) {
		return HeatmapDay{
			Date:    time.Time{},
			Count:   0,
			Weekday: GetWeekdayLabel(weekday),
		}
	}
	return days[index]
}

// GenerateHeatmapDays génère les N derniers jours pour le heatmap
func GenerateHeatmapDays(reviewCounts map[string]int, weeks int) []HeatmapDay {
	days := make([]HeatmapDay, weeks*7)
	now := time.Now()

	for i := 0; i < weeks*7; i++ {
		date := now.AddDate(0, 0, -(weeks*7 - 1 - i))
		dateKey := date.Format("2006-01-02")

		days[i] = HeatmapDay{
			Date:    date,
			Count:   reviewCounts[dateKey],
			Weekday: date.Weekday().String()[:3],
		}
	}

	return days
}

// GetHeatmapTooltip génère le texte du tooltip pour une cellule
func GetHeatmapTooltip(day HeatmapDay) string {
	if day.Date.IsZero() {
		return "No data"
	}

	reviews := "reviews"
	if day.Count == 1 {
		reviews = "review"
	}

	return fmt.Sprintf("%s: %d %s",
		day.Date.Format("Jan 02, 2006"),
		day.Count,
		reviews,
	)
}
