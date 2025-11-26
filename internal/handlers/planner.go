package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"maestro/internal/models"
)

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// PAGE PRINCIPALE (Layout complet)
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func HandlePlannerPage(w http.ResponseWriter, r *http.Request) {
	data := models.DayData{
		Date:      time.Now(),
		TimeSlots: generateDefaultTimeSlots(),
	}

	// Rend le layout complet avec vue JOUR par défaut
	if err := Tmpl.ExecuteTemplate(w, "planner", data); err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// FRAGMENT JOUR (appelé par HTMX)
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func HandlePlannerDay(w http.ResponseWriter, r *http.Request) {
	// Parse et valide la date
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = time.Now() // Fallback sur aujourd'hui
	}

	data := models.DayData{
		Date:      date,
		TimeSlots: generateTimeSlots(date),
	}

	// Rend UNIQUEMENT le fragment day-schedule
	if err := Tmpl.ExecuteTemplate(w, "day-schedule", data); err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// FRAGMENT SEMAINE (appelé par HTMX)
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func HandlePlannerWeek(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")

	// Récupère la semaine actuelle par défaut
	now := time.Now()
	_, currentWeek := now.ISOWeek()
	weekNum := currentWeek

	// Parse le numéro de semaine si fourni
	if weekStr != "" {
		if parsed, err := strconv.Atoi(weekStr); err == nil {
			// Validation : entre 1 et 53
			if parsed >= 1 && parsed <= 53 {
				weekNum = parsed
			}
		}
	}

	// Calcule les dates de la semaine
	startDate, endDate := getWeekDates(weekNum, now.Year())

	data := models.WeekData{
		WeekNumber: weekNum,
		StartDate:  startDate,
		EndDate:    endDate,
		Days:       generateWeekDays(startDate),
	}

	// Rend UNIQUEMENT la grille semaine
	if err := Tmpl.ExecuteTemplate(w, "week-grid", data); err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// FRAGMENT MOIS (appelé par HTMX)
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

func HandlePlannerMonth(w http.ResponseWriter, r *http.Request) {
	monthStr := r.URL.Query().Get("month")

	// Parse le mois (format "2025-11")
	targetDate := time.Now()
	if monthStr != "" {
		if parsed, err := time.Parse("2006-01", monthStr); err == nil {
			targetDate = parsed
		}
	}

	monthName := strings.ToUpper(targetDate.Format("January"))
	monthNameFr := translateMonth(monthName)

	data := models.MonthData{
		Month:    monthNameFr,
		Year:     targetDate.Year(),
		MonthNum: int(targetDate.Month()),
		Days:     generateMonthDays(targetDate),
	}

	// Rend UNIQUEMENT le calendrier mensuel
	if err := Tmpl.ExecuteTemplate(w, "month-calendar", data); err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// FONCTIONS HELPER
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

// Génère les créneaux horaires par défaut
func generateDefaultTimeSlots() []models.TimeSlot {
	return []models.TimeSlot{
		{
			StartTime: "08:00",
			EndTime:   "10:00",
			Tasks: []models.Task{
				{Title: "Révision Goroutines & Channels", Priority: "high", Completed: false},
			},
		},
		{
			StartTime: "10:00",
			EndTime:   "12:00",
			Tasks: []models.Task{
				{Title: "Lecture: Error Handling", Priority: "medium", Completed: false},
			},
		},
		{
			StartTime: "14:00",
			EndTime:   "16:00",
			Tasks: []models.Task{
				{Title: "HTMX: Filtres Dynamiques", Priority: "low", Completed: false},
			},
		},
		{
			StartTime: "16:00",
			EndTime:   "18:00",
			Tasks: []models.Task{
				{Title: "Review: Code Security", Priority: "high", Completed: false},
			},
		},
	}
}

// Génère les créneaux pour une date spécifique
func generateTimeSlots(date time.Time) []models.TimeSlot {
	// TODO: Récupérer depuis la base de données
	// Pour l'instant, retourne des créneaux par défaut
	return generateDefaultTimeSlots()
}

// Calcule les dates de début et fin d'une semaine ISO
func getWeekDates(weekNum, year int) (time.Time, time.Time) {
	// Trouver le premier lundi de l'année
	jan1 := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Calculer le premier lundi
	daysUntilMonday := (8 - int(jan1.Weekday())) % 7
	firstMonday := jan1.AddDate(0, 0, daysUntilMonday)

	// Si le 1er janvier est après jeudi, la semaine 1 commence le lundi suivant
	if jan1.Weekday() > time.Thursday {
		firstMonday = firstMonday.AddDate(0, 0, 7)
	}

	// Calculer le début de la semaine demandée
	startDate := firstMonday.AddDate(0, 0, (weekNum-1)*7)
	endDate := startDate.AddDate(0, 0, 6)

	return startDate, endDate
}

// Génère les 7 jours de la semaine
func generateWeekDays(startDate time.Time) []models.WeekDay {
	days := make([]models.WeekDay, 7)
	today := time.Now()

	dayNames := []string{"LUN", "MAR", "MER", "JEU", "VEN", "SAM", "DIM"}

	for i := 0; i < 7; i++ {
		currentDay := startDate.AddDate(0, 0, i)
		isToday := currentDay.Year() == today.Year() &&
			currentDay.Month() == today.Month() &&
			currentDay.Day() == today.Day()

		days[i] = models.WeekDay{
			DayName:   dayNames[i],
			DayNumber: currentDay.Day(),
			IsToday:   isToday,
			Tasks:     generateTasksForDay(currentDay),
		}
	}

	return days
}

// Génère les tâches pour un jour (mock)
func generateTasksForDay(date time.Time) []models.Task {
	// TODO: Récupérer depuis la base de données
	// Pour l'instant, retourne quelques tâches mock
	return []models.Task{
		{Title: "08:00 Goroutines", Priority: "high", Completed: false},
		{Title: "14:00 HTMX", Priority: "medium", Completed: false},
	}
}

// Génère tous les jours du mois (avec jours adjacents)
func generateMonthDays(current time.Time) []models.MonthDay {
	days := make([]models.MonthDay, 0, 42) // 6 semaines max
	today := time.Now()

	firstDay := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, current.Location())
	lastDay := firstDay.AddDate(0, 1, -1)

	// Jours du mois précédent pour remplir la première semaine
	weekday := int(firstDay.Weekday())
	if weekday == 0 {
		weekday = 7 // Dimanche = 7 en format ISO
	}

	for i := weekday - 1; i > 0; i-- {
		prevDay := firstDay.AddDate(0, 0, -i)
		days = append(days, models.MonthDay{
			Number:       prevDay.Day(),
			IsOtherMonth: true,
			IsToday:      false,
			TaskCount:    0,
		})
	}

	// Jours du mois courant
	for d := 1; d <= lastDay.Day(); d++ {
		isToday := d == today.Day() &&
			current.Month() == today.Month() &&
			current.Year() == today.Year()

		days = append(days, models.MonthDay{
			Number:       d,
			IsOtherMonth: false,
			IsToday:      isToday,
			TaskCount:    calculateTaskCount(current.Year(), current.Month(), d),
		})
	}

	// Jours du mois suivant pour compléter la grille (6 semaines = 42 cases)
	nextMonthDay := 1
	for len(days) < 42 {
		days = append(days, models.MonthDay{
			Number:       nextMonthDay,
			IsOtherMonth: true,
			IsToday:      false,
			TaskCount:    0,
		})
		nextMonthDay++
	}

	return days
}

// Calcule le nombre de tâches pour un jour donné
func calculateTaskCount(year int, month time.Month, day int) int {
	// TODO: Récupérer le vrai compte depuis la base de données
	// Pour l'instant, retourne un nombre aléatoire mock
	if day%5 == 0 {
		return 3
	} else if day%3 == 0 {
		return 2
	} else if day%2 == 0 {
		return 1
	}
	return 0
}

// Traduit les noms de mois anglais vers français
func translateMonth(englishMonth string) string {
	months := map[string]string{
		"JANUARY":   "JANVIER",
		"FEBRUARY":  "FÉVRIER",
		"MARCH":     "MARS",
		"APRIL":     "AVRIL",
		"MAY":       "MAI",
		"JUNE":      "JUIN",
		"JULY":      "JUILLET",
		"AUGUST":    "AOÛT",
		"SEPTEMBER": "SEPTEMBRE",
		"OCTOBER":   "OCTOBRE",
		"NOVEMBER":  "NOVEMBRE",
		"DECEMBER":  "DÉCEMBRE",
	}

	if french, ok := months[englishMonth]; ok {
		return french
	}
	return englishMonth
}
