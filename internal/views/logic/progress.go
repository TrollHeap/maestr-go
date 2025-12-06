package logic

func GetProgressTextColor(completed, total int) string {
	percent := CalculateProgressPercent(completed, total)

	switch {
	case percent == 100:
		// Vert émeraude (100% comme ton badge "Maîtrisé")
		return "text-emerald-300"
	case percent >= 67:
		// Turquoise (proche de la fin, même teinte que ta barre)
		return "text-sky-300"
	case percent > 0:
		// Jaune chaud (en cours, cohérent avec les steps non finis)
		return "text-amber-300"
	default:
		return "text-emerald-300"
	}
}

func CalculateProgressPercent(completed, total int) int {
	if total == 0 {
		return 0
	}
	return (completed * 100) / total
}
