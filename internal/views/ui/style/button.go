package style

func GetReviewButtonClassesGradient(quality int) string {
	switch quality {
	case 0:
		return "border-rose-500/60 bg-gradient-to-br from-rose-950/70 to-slate-950/70 hover:from-rose-900/80 hover:to-slate-900/80 text-rose-100"
	case 1:
		return "border-orange-500/60 bg-gradient-to-br from-orange-950/70 to-slate-950/70 hover:from-orange-900/80 hover:to-slate-900/80 text-orange-100"
	case 2:
		return "border-amber-500/60 bg-gradient-to-br from-amber-950/70 to-slate-950/70 hover:from-amber-900/80 hover:to-slate-900/80 text-amber-100"
	case 3:
		return "border-emerald-500/60 bg-gradient-to-br from-emerald-950/70 to-slate-950/70 hover:from-emerald-900/80 hover:to-slate-900/80 text-emerald-100"
	default:
		return "border-slate-700 bg-slate-900/70 hover:bg-slate-800/80 text-slate-100"
	}
}
