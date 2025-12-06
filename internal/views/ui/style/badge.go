package style

import "strings"

// GetBadgeClasses - Compute badge Tailwind classes
func GetBadgeClasses(variant string, size string) string {
	base := []string{
		"inline-flex",
		"items-center",
		"gap-1.5",
		"rounded-lg",
		"font-mono",
		"uppercase",
		"tracking-wider",
		"transition-all",
	}

	// Size variants
	sizeClasses := map[string][]string{
		"xs": {"px-1.5", "py-0.5", "text-[0.55rem]"},
		"sm": {"px-2", "py-0.5", "text-[0.6rem]"},
		"md": {"px-2.5", "py-1", "text-[0.65rem]"},
	}

	// Variant colors
	variantClasses := map[string][]string{
		"status": {
			"bg-emerald-900/50",
			"border",
			"border-emerald-500/60",
			"text-emerald-200",
		},
		"domain": {
			"bg-sky-900/50",
			"border",
			"border-sky-500/60",
			"text-sky-200",
		},
		"difficulty": {
			"bg-purple-900/50",
			"border",
			"border-purple-500/60",
			"text-purple-200",
		},
		"system": {
			"bg-sky-500/20",
			"border",
			"border-sky-400/60",
			"text-sky-200",
			"shadow-[0_0_10px_rgba(56,189,248,0.2)]",
		},
		"count": {
			"bg-emerald-500/20",
			"border",
			"border-emerald-400/60",
			"text-emerald-200",
			"font-semibold",
		},
	}

	classes := append(base, sizeClasses[size]...)
	classes = append(classes, variantClasses[variant]...)

	return strings.Join(classes, " ")
}

func GetIntervalBadgeClasses(quality int) string {
	switch quality {
	case 0:
		return "bg-rose-500/20 text-rose-200 border border-rose-500/40"
	case 1:
		return "bg-orange-500/20 text-orange-200 border border-orange-500/40"
	case 2:
		return "bg-amber-500/20 text-amber-200 border border-amber-500/40"
	case 3:
		return "bg-emerald-500/20 text-emerald-200 border border-emerald-500/40"
	default:
		return "bg-slate-700/40 text-slate-300 border border-slate-600"
	}
}

func GetStatusBadgeClassesGlow(done bool) string {
	if done {
		// Maîtrisé : vert avec glow fort
		return "border border-emerald-500/70 text-emerald-200 bg-emerald-500/15 shadow-[0_0_12px_rgba(16,185,129,0.5)] hover:shadow-[0_0_20px_rgba(16,185,129,0.7)] hover:border-emerald-400"
	}
	// En cours : jaune avec glow fort
	return "border border-amber-500/70 text-amber-200 bg-amber-500/15 shadow-[0_0_12px_rgba(251,191,36,0.5)] hover:shadow-[0_0_20px_rgba(251,191,36,0.7)] hover:border-amber-400"
}
