package style

import "strings"

// GetHeaderClasses - Terminal header badge styling
func GetHeaderClasses(theme string) string {
	base := []string{
		"inline-flex",
		"items-center",
		"gap-3",
		"px-4",
		"py-2",
		"rounded-full",
		"text-xs",
		"font-mono",
		"tracking-widest",
	}

	// Theme-specific colors
	themeClasses := map[string][]string{
		"sky": {
			"bg-sky-500/10",
			"border",
			"border-sky-400/40",
			"text-sky-300",
		},
		"emerald": {
			"bg-emerald-500/10",
			"border",
			"border-emerald-400/40",
			"text-emerald-300",
		},
		"purple": {
			"bg-purple-500/10",
			"border",
			"border-purple-400/40",
			"text-purple-300",
		},
		"amber": {
			"bg-amber-500/10",
			"border",
			"border-amber-400/40",
			"text-amber-300",
		},
	}

	classes := append(base, themeClasses[theme]...)
	return strings.Join(classes, " ")
}

func GetPulseDotClasses(theme string) string {
	base := []string{
		"inline-block",
		"h-2",
		"w-2",
		"rounded-full",
		"animate-pulse",
	}

	// Theme-specific dot colors with glow
	dotClasses := map[string][]string{
		"sky": {
			"bg-sky-400",
			"shadow-[0_0_6px_rgba(56,189,248,0.8)]",
		},
		"emerald": {
			"bg-emerald-400",
			"shadow-[0_0_6px_rgba(52,211,153,0.8)]",
		},
		"purple": {
			"bg-purple-400",
			"shadow-[0_0_6px_rgba(192,132,252,0.8)]",
		},
		"amber": {
			"bg-amber-400",
			"shadow-[0_0_6px_rgba(251,191,36,0.8)]",
		},
	}

	classes := append(base, dotClasses[theme]...)
	return strings.Join(classes, " ")
}
