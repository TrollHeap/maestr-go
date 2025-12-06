package logic

import "time"

func GetEmptyDaysBefore(firstDay time.Time) int {
	weekday := int(firstDay.Weekday())

	// Dimanche = 0 → 7 (ISO week)
	if weekday == 0 {
		weekday = 7
	}

	// Lundi = 1, donc on veut weekday-1 espaces vides
	return weekday - 1
}

// getWeekDayCardClasses – déjà utilisé dans le template
func GetWeekDayCardClasses(date time.Time, count int) string {
	base := "border-slate-800"

	if isToday(date) {
		base += " border-emerald-500/80 shadow-[0_0_0_1px_rgba(16,185,129,0.4)]"
	} else if count > 0 {
		base += " border-emerald-600/40"
	}

	return base
}

// getMonthDayClasses – déjà utilisé dans le template, on le laisse ici
func GetMonthDayClasses(year int, month time.Month, day int, count int) string {
	base := "border-slate-800 bg-slate-900/60"

	if isDayToday(year, month, day) {
		base += " border-emerald-500/80 bg-emerald-950/40 shadow-[0_0_0_1px_rgba(16,185,129,0.4)]"
	} else if count > 0 {
		base += " border-emerald-600/40 bg-emerald-950/20"
	}

	return base
}

// isToday - utilisé par WeekDayCard
func isToday(date time.Time) bool {
	now := time.Now()
	return date.Year() == now.Year() && date.YearDay() == now.YearDay()
}

// isDayToday - Vérifie si c'est aujourd'hui (pour MonthDayCell)
func isDayToday(year int, month time.Month, day int) bool {
	now := time.Now()
	return year == now.Year() && month == now.Month() && day == now.Day()
}

// truncate - utilisé pour les titres des exos
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
