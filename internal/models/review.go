package models

// ReviewInput représente l'input pour noter un exercice
type ReviewInput struct {
	ExerciseID string `json:"exercise_id"`
	Rating     int    `json:"rating"` // 1-4
}

// ReviewResponse représente la réponse après une note
type ReviewResponse struct {
	Exercise     *Exercise `json:"exercise"`
	NextReviewIn int       `json:"next_review_in"`
	Message      string    `json:"message"`
}
