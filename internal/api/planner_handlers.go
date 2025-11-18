package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"maestro/internal/domain"
	"maestro/internal/storage"
)

// PlannerHandler gère les endpoints du planner
type PlannerHandler struct {
	planner *domain.Planner
	store   *storage.JSONStore
}

// NewPlannerHandler crée un nouveau handler planner
func NewPlannerHandler(planner *domain.Planner, store *storage.JSONStore) *PlannerHandler {
	return &PlannerHandler{
		planner: planner,
		store:   store,
	}
}

// CreateSession crée une nouvelle session planifiée
// POST /api/planner/session
func (h *PlannerHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Date        string   `json:"date"`      // YYYY-MM-DD
		TimeSlot    string   `json:"time_slot"` // morning, afternoon, evening
		ExerciseIDs []string `json:"exercise_ids"`
		Duration    int      `json:"duration"` // minutes
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format (use YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Validate time slot
	if req.TimeSlot != "morning" && req.TimeSlot != "afternoon" && req.TimeSlot != "evening" {
		http.Error(
			w,
			"Invalid time_slot (must be: morning, afternoon, or evening)",
			http.StatusBadRequest,
		)
		return
	}

	// Create session
	session := h.planner.CreateSession(date, req.TimeSlot, req.ExerciseIDs, req.Duration)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// UpdateSession met à jour une session
// PUT /api/planner/session/{id}
func (h *PlannerHandler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract session ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/planner/session/")
	sessionID := path

	var req struct {
		Status string `json:"status"` // planned, completed, skipped
		Notes  string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update session
	if err := h.planner.UpdateSession(sessionID, req.Status, req.Notes); err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Session updated"})
}

// DeleteSession supprime une session
// DELETE /api/planner/session/{id}
func (h *PlannerHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract session ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/planner/session/")
	sessionID := path

	// Delete session
	if err := h.planner.DeleteSession(sessionID); err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Session deleted"})
}

// GetToday retourne le plan du jour
// GET /api/planner/today
func (h *PlannerHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	today := time.Now()
	plan := h.planner.GetDailyPlan(today)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

// GetWeek retourne le plan de la semaine
// GET /api/planner/week?start=YYYY-MM-DD
func (h *PlannerHandler) GetWeek(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")

	var startDate time.Time
	var err error

	if startStr != "" {
		startDate, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			http.Error(w, "Invalid start date format (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now()
	}

	plan := h.planner.GetWeeklyPlan(startDate)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

// GetMonth retourne le plan du mois (4 semaines)
// GET /api/planner/month?start=YYYY-MM-DD
func (h *PlannerHandler) GetMonth(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")

	var startDate time.Time
	var err error

	if startStr != "" {
		startDate, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			http.Error(w, "Invalid start date format (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now()
	}

	weeks := h.planner.GetMonthlyPlan(startDate)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"weeks": weeks,
	})
}

// GetStats retourne les statistiques du planner
// GET /api/planner/stats
func (h *PlannerHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := h.planner.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// RegisterRoutes enregistre les routes du planner
func (h *PlannerHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/api/planner/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateSession(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/api/planner/session/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			h.UpdateSession(w, r)
		} else if r.Method == http.MethodDelete {
			h.DeleteSession(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/api/planner/today", h.GetToday)
	router.HandleFunc("/api/planner/week", h.GetWeek)
	router.HandleFunc("/api/planner/month", h.GetMonth)
	router.HandleFunc("/api/planner/stats", h.GetStats)
}
