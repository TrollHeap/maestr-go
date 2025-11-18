package domain

import (
	"time"

	"maestro/internal/models"
)

// ============= ADVANCED STATISTICS ANALYZER =============

type DifficultyLevel string

const (
	VeryEasy DifficultyLevel = "very_easy" // EF > 2.8, Interval > 30
	Easy     DifficultyLevel = "easy"      // EF 2.5-2.8, Interval 10-30
	Medium   DifficultyLevel = "medium"    // EF 2.0-2.5, Interval 3-10
	Hard     DifficultyLevel = "hard"      // EF 1.5-2.0, Interval 1-3
	VeryHard DifficultyLevel = "very_hard" // EF < 1.5, Interval 0-1
)

type ExerciseAnalysis struct {
	ExerciseID       string          `json:"exercise_id"`
	Title            string          `json:"title"`
	Domain           string          `json:"domain"`
	LevelDetected    DifficultyLevel `json:"level_detected"`
	EF               float64         `json:"ef"`
	Interval         int             `json:"interval"`
	Repetitions      int             `json:"repetitions"`
	LastReviewedAt   *time.Time      `json:"last_reviewed_at"`
	IsOverdue        bool            `json:"is_overdue"`
	DaysOverdue      int             `json:"days_overdue"`
	Mastery          int             `json:"mastery"` // 0-100%
	NeedsImprovement bool            `json:"needs_improvement"`
	ConfidenceLevel  string          `json:"confidence_level"` // Low/Medium/High
	Recommendation   string          `json:"recommendation"`
}

type StatisticsAnalyzer struct {
	scheduler *Scheduler
}

func NewStatisticsAnalyzer(scheduler *Scheduler) *StatisticsAnalyzer {
	return &StatisticsAnalyzer{
		scheduler: scheduler,
	}
}

// ============= MAIN ANALYSIS FUNCTIONS =============

// AnalyzeExercise détaille l'état d'un exercice
func (sa *StatisticsAnalyzer) AnalyzeExercise(ex *models.Exercise) *ExerciseAnalysis {
	now := time.Now()

	analysis := &ExerciseAnalysis{
		ExerciseID:     ex.ID,
		Title:          ex.Title,
		Domain:         ex.Domain,
		EF:             ex.EaseFactor,
		Interval:       ex.IntervalDays,
		Repetitions:    ex.Repetitions,
		LastReviewedAt: ex.LastReviewed,
	}

	// 1. Détecter le niveau de difficulté (CLÉE!)
	analysis.LevelDetected = sa.detectDifficultyLevel(ex)

	// 2. Calculer la maîtrise (%)
	analysis.Mastery = sa.calculateMastery(ex.EaseFactor)

	// 3. Vérifier si overdue
	if ex.LastReviewed != nil && ex.IntervalDays > 0 {
		nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
		analysis.IsOverdue = now.After(nextReview)

		if analysis.IsOverdue {
			analysis.DaysOverdue = int(now.Sub(nextReview).Hours() / 24)
		}
	}

	// 4. Détecter si besoin amélioration
	analysis.NeedsImprovement = sa.needsImprovement(ex)

	// 5. Évaluer le niveau de confiance
	analysis.ConfidenceLevel = sa.getConfidenceLevel(ex)

	// 6. Générer une recommandation
	analysis.Recommendation = sa.generateRecommendation(ex, analysis)

	return analysis
}

// GetStruggling retourne les exercices avec lesquels vous avez du mal
func (sa *StatisticsAnalyzer) GetStruggling(exercises []models.Exercise) []ExerciseAnalysis {
	var struggling []ExerciseAnalysis

	for _, ex := range exercises {
		if ex.Deleted {
			continue
		}

		// Considérer comme "struggling" si:
		// 1. EF très bas (< 1.7) = Vraiment difficile
		// 2. Très de répétitions (> 10) mais toujours difficile = Besoin de pratique
		// 3. Overdue = Négligé

		if ex.EaseFactor < 1.7 ||
			(ex.Repetitions > 10 && ex.EaseFactor < 2.0) ||
			(ex.LastReviewed != nil && sa.scheduler.IsDueForReview(&ex)) {

			analysis := sa.AnalyzeExercise(&ex)
			if analysis.NeedsImprovement {
				struggling = append(struggling, *analysis)
			}
		}
	}

	return struggling
}

// GetMastered retourne les exercices maîtrisés
func (sa *StatisticsAnalyzer) GetMastered(exercises []models.Exercise) []ExerciseAnalysis {
	var mastered []ExerciseAnalysis

	for _, ex := range exercises {
		if ex.Deleted {
			continue
		}

		// Maîtrisé si:
		// 1. EF > 2.8 ET
		// 2. Repetitions > 3 ET
		// 3. Interval > 30

		if ex.EaseFactor > 2.8 && ex.Repetitions > 3 && ex.IntervalDays > 30 {
			mastered = append(mastered, *sa.AnalyzeExercise(&ex))
		}
	}

	return mastered
}

// GetNeedsPractice retourne les exercices qui ont besoin de pratique
func (sa *StatisticsAnalyzer) GetNeedsPractice(exercises []models.Exercise) []ExerciseAnalysis {
	var needsPractice []ExerciseAnalysis

	for _, ex := range exercises {
		if ex.Deleted {
			continue
		}

		// Besoin de pratique si:
		// 1. Peu de répétitions (< 3)
		// 2. OU EF entre 1.5-2.0 (moyen-difficile)
		// 3. OU jamais revisité

		if ex.Repetitions < 3 ||
			(ex.EaseFactor > 1.5 && ex.EaseFactor < 2.0) ||
			ex.LastReviewed == nil {

			analysis := sa.AnalyzeExercise(&ex)
			if !analysis.NeedsImprovement { // Pas struggling
				needsPractice = append(needsPractice, *analysis)
			}
		}
	}

	return needsPractice
}

// ============= HELPER FUNCTIONS =============

// detectDifficultyLevel détecte le niveau d'un exercice
func (sa *StatisticsAnalyzer) detectDifficultyLevel(ex *models.Exercise) DifficultyLevel {
	// Basé sur EF (Ease Factor) ET Interval

	if ex.EaseFactor > 2.8 && ex.IntervalDays > 30 {
		return VeryEasy
	}
	if ex.EaseFactor >= 2.5 && ex.IntervalDays >= 10 {
		return Easy
	}
	if ex.EaseFactor >= 2.0 && ex.IntervalDays >= 3 {
		return Medium
	}
	if ex.EaseFactor >= 1.5 && ex.IntervalDays >= 1 {
		return Hard
	}

	return VeryHard // EF < 1.5 ou Interval 0-1
}

// calculateMastery convertit EF en pourcentage de maîtrise (0-100%)
func (sa *StatisticsAnalyzer) calculateMastery(ef float64) int {
	// Range: 1.3 → 3.0 = 0% → 100%

	if ef < 1.3 {
		return 0
	}
	if ef > 3.0 {
		return 100
	}

	// Linear scale
	percentage := ((ef - 1.3) / (3.0 - 1.3)) * 100
	return int(percentage)
}

// needsImprovement détecte si exercice a vraiment besoin de travail
func (sa *StatisticsAnalyzer) needsImprovement(ex *models.Exercise) bool {
	// Vrai si:
	// 1. EF < 1.7 (struggling!)
	// 2. Rating 1 récent (oublié)
	// 3. Overdue par > 5 jours
	// 4. Repetitions high (> 15) mais EF low (< 1.8) = récalcitrant

	if ex.EaseFactor < 1.7 {
		return true
	}

	// Vérifier si overdue longtemps
	if ex.LastReviewed != nil && ex.IntervalDays > 0 {
		nextReview := ex.LastReviewed.AddDate(0, 0, ex.IntervalDays)
		daysOverdue := int(time.Until(nextReview).Hours() / -24)
		if daysOverdue > 5 {
			return true
		}
	}

	// Vérifier le pattern "exercice récalcitrant"
	if ex.Repetitions > 15 && ex.EaseFactor < 1.8 {
		return true
	}

	return false
}

// getConfidenceLevel évalue comment vous maîtrisez bien
func (sa *StatisticsAnalyzer) getConfidenceLevel(ex *models.Exercise) string {
	// Basé sur:
	// 1. Nombre de révisions (plus = plus de confiance)
	// 2. EF (plus haut = plus de confiance)
	// 3. Pas overdue

	if ex.Repetitions < 2 || ex.EaseFactor < 1.8 {
		return "Low"
	}

	if ex.Repetitions >= 2 && ex.Repetitions < 5 && ex.EaseFactor >= 1.8 && ex.EaseFactor < 2.3 {
		return "Medium"
	}

	if ex.Repetitions >= 5 && ex.EaseFactor >= 2.3 {
		return "High"
	}

	return "Medium"
}

// generateRecommendation génère un conseil basé sur l'analyse
func (sa *StatisticsAnalyzer) generateRecommendation(
	ex *models.Exercise,
	analysis *ExerciseAnalysis,
) string {
	switch analysis.LevelDetected {
	case VeryEasy:
		return "✓ Bien maîtrisé! Revoir occasionnellement."

	case Easy:
		return "✓ Bon progrès! Continue à pratiquer occasionnellement."

	case Medium:
		return "→ Pratique régulière recommandée pour maintenir."

	case Hard:
		if analysis.Repetitions < 3 {
			return "⚠️ Difficile! Pratiquer davantage. Continuez!"
		}
		return "⚠️ Très difficile. Besoin de pratique intensive."

	case VeryHard:
		if analysis.IsOverdue {
			return "🔴 URGENT! Overdue de " + string(
				rune(analysis.DaysOverdue),
			) + " jours. Revoir MAINTENANT!"
		}
		if ex.Repetitions > 10 {
			return "🔴 PROBLÉMATIQUE! Malgré " + string(
				rune(ex.Repetitions),
			) + " révisions, c'est très difficile. Besoin de stratégie différente!"
		}
		return "🔴 TRÈS DIFFICILE! Besoin de pratique intensive et rythme + fréquent."
	}

	return "À déterminer"
}

// ============= STATISTICS SUMMARY =============

type DomainAnalysis struct {
	Domain         string  `json:"domain"`
	Total          int     `json:"total"`
	Completed      int     `json:"completed"`
	Mastery        float64 `json:"mastery"` // Average EF
	AvgRepetitions float64 `json:"avg_repetitions"`
	Struggling     int     `json:"struggling"`     // Count
	Mastered       int     `json:"mastered"`       // Count
	NeedsPractice  int     `json:"needs_practice"` // Count
	Recommendation string  `json:"recommendation"`
}

// AnalyzeDomain analyse un domaine entier
func (sa *StatisticsAnalyzer) AnalyzeDomain(
	domain string,
	exercises []models.Exercise,
) *DomainAnalysis {
	domainEx := []models.Exercise{}

	// Filter exercises by domain
	for _, ex := range exercises {
		if ex.Domain == domain && !ex.Deleted {
			domainEx = append(domainEx, ex)
		}
	}

	if len(domainEx) == 0 {
		return &DomainAnalysis{Domain: domain}
	}

	analysis := &DomainAnalysis{
		Domain: domain,
		Total:  len(domainEx),
	}

	// Calculate aggregates
	var totalEF float64
	var totalReps float64

	for _, ex := range domainEx {
		if ex.Completed {
			analysis.Completed++
		}

		totalEF += ex.EaseFactor
		totalReps += float64(ex.Repetitions)

		// Count struggling/mastered/needs practice
		exAnalysis := sa.AnalyzeExercise(&ex)

		if exAnalysis.NeedsImprovement {
			analysis.Struggling++
		} else if exAnalysis.LevelDetected == VeryEasy || exAnalysis.LevelDetected == Easy {
			analysis.Mastered++
		} else {
			analysis.NeedsPractice++
		}
	}

	analysis.Mastery = totalEF / float64(len(domainEx))
	analysis.AvgRepetitions = totalReps / float64(len(domainEx))

	// Generate recommendation
	if analysis.Struggling > 0 {
		analysis.Recommendation = "⚠️ " + string(
			rune(analysis.Struggling),
		) + " exercices difficiles - Pratiquez plus!"
	} else if analysis.Mastered > (len(domainEx) / 2) {
		analysis.Recommendation = "✓ Domaine bien maîtrisé!"
	} else if analysis.NeedsPractice > 0 {
		analysis.Recommendation = "→ Continuez à pratiquer régulièrement."
	}

	return analysis
}

// AnalyzeAllDomains analyse tous les domaines
func (sa *StatisticsAnalyzer) AnalyzeAllDomains(exercises []models.Exercise) []DomainAnalysis {
	// Collect unique domains
	domainMap := make(map[string]bool)
	for _, ex := range exercises {
		if !ex.Deleted {
			domainMap[ex.Domain] = true
		}
	}

	var results []DomainAnalysis
	for domain := range domainMap {
		results = append(results, *sa.AnalyzeDomain(domain, exercises))
	}

	return results
}

// ============= LEARNING INSIGHTS =============

type LearningInsights struct {
	StrongestDomain string  `json:"strongest_domain"` // Highest mastery
	WeakestDomain   string  `json:"weakest_domain"`   // Lowest mastery
	MostPracticed   string  `json:"most_practiced"`   // Highest avg reps
	Overdue         int     `json:"overdue"`          // Count
	SuccessRate     float64 `json:"success_rate"`     // % of exercises EF > 2.5
	AvgMastery      float64 `json:"avg_mastery"`      // Overall avg EF
	RecommendFocus  string  `json:"recommend_focus"`  // What to work on
}

// GenerateInsights génère des insights sur votre apprentissage
func (sa *StatisticsAnalyzer) GenerateInsights(exercises []models.Exercise) *LearningInsights {
	if len(exercises) == 0 {
		return &LearningInsights{}
	}

	domainAnalyses := sa.AnalyzeAllDomains(exercises)

	insights := &LearningInsights{}

	var bestMastery float64 = -1
	var worstMastery float64 = 999
	var mostReps float64 = 0
	var totalEF float64
	var successCount int
	var overdueCount int

	for _, domain := range domainAnalyses {
		if domain.Mastery > bestMastery {
			bestMastery = domain.Mastery
			insights.StrongestDomain = domain.Domain
		}
		if domain.Mastery < worstMastery && domain.Mastery > 0 {
			worstMastery = domain.Mastery
			insights.WeakestDomain = domain.Domain
		}
		if domain.AvgRepetitions > mostReps {
			mostReps = domain.AvgRepetitions
			insights.MostPracticed = domain.Domain
		}

		totalEF += domain.Mastery
	}

	// Count overdue & success
	for _, ex := range exercises {
		if ex.Deleted {
			continue
		}

		totalEF += ex.EaseFactor

		if ex.EaseFactor > 2.5 {
			successCount++
		}

		if sa.scheduler.IsDueForReview(&ex) {
			overdueCount++
		}
	}

	insights.AvgMastery = totalEF / float64(len(exercises))
	insights.SuccessRate = (float64(successCount) / float64(len(exercises))) * 100
	insights.Overdue = overdueCount

	// Recommend focus
	if insights.Overdue > 0 {
		insights.RecommendFocus = "🔴 " + string(
			rune(insights.Overdue),
		) + " exercices overdue - Rattrapez!"
	} else if insights.WeakestDomain != "" {
		insights.RecommendFocus = "→ Focus sur " + insights.WeakestDomain + " (mastery: " + string(rune(int(worstMastery))) + "%)"
	} else {
		insights.RecommendFocus = "✓ Bon travail! Continuez à pratiquer régulièrement."
	}

	return insights
}
