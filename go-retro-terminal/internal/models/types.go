package models

type StatsData struct {
	CPU     float64
	Memory  float64
	Uptime  string
	Status  string
	Version string
}

type TerminalOutput struct {
	Command string
	Output  string
	Time    string
}
