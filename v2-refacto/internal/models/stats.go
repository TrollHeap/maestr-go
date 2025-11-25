package models

// StatsMetrics contient les métriques globales
type StatsMetrics struct {
	Total          int
	Completed      int
	WIP            int
	Todo           int
	CompletionRate int
}

// DomainStat représente les stats d'un domaine
type DomainStat struct {
	Name       string
	Total      int
	Completed  int
	Percentage int
}

// DifficultyStat représente les stats d'un niveau de difficulté
type DifficultyStat struct {
	Level  int
	Symbol string
	Label  string
	Count  int
}
