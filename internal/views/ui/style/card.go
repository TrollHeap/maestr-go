package style

func GetEnergyCardClass(level int) string {
	base := "group relative rounded-2xl border bg-slate-900/70 backdrop-blur-xl " +
		"p-5 shadow-lg shadow-slate-900/40 cursor-pointer " +
		"transition-all duration-200 hover:-translate-y-1 hover:shadow-2xl"

	switch level {
	case 1:
		// Faible : bleu/teal doux
		return base + " border-emerald-500/30 hover:border-emerald-400 " +
			"hover:bg-gradient-to-br hover:from-slate-900 hover:to-emerald-900/40"
	case 2:
		// Moyen : ambre/orange
		return base + " border-amber-500/30 hover:border-amber-400 " +
			"hover:bg-gradient-to-br hover:from-slate-900 hover:to-amber-900/40"
	case 3:
		// Élevé : rouge
		return base + " border-rose-500/30 hover:border-rose-400 " +
			"hover:bg-gradient-to-br hover:from-slate-900 hover:to-rose-900/40"
	default:
		return base + " border-slate-700 hover:border-slate-500"
	}
}
