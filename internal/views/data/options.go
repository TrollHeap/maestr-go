package data

import "sort"

// ────────────────────────────────────────────────────
// STATUS OPTIONS - Pour pills
// ────────────────────────────────────────────────────

// StatusOption représente une option de statut
type StatusOption struct {
	Value string
	Label string
}

// GetStatusOptions - Retourne les options de statut dans l'ordre
func GetStatusOptions() []StatusOption {
	return []StatusOption{
		{"", "Tous"},
		{"inprogress", "En cours"},
		{"mastered", "Maîtrisés"},
	}
}

// ────────────────────────────────────────────────────
// DIFFICULTY FILTER OPTIONS - Pour pills
// ────────────────────────────────────────────────────

// DifficultyFilterOption représente une option de filtrage par difficulté
type DifficultyFilterOption struct {
	Value int
	Label string
}

// GetDifficultyFilterOptions - Retourne les options de difficultés dans l'ordre
func GetDifficultyFilterOptions() []DifficultyFilterOption {
	return []DifficultyFilterOption{
		{0, "Toutes diff."},
		{2, "D1-D2"},
		{3, "D3"},
		{4, "D4"},
	}
}

// ────────────────────────────────────────────────────
// DOMAINS - Ordre alphabétique
// ────────────────────────────────────────────────────

// GetDomains - Liste des domaines triée alphabétiquement
func GetDomains() []string {
	domains := []string{
		"Go",
		"Algorithms",
		"Database",
		"Architecture",
		"Security",
		"DevOps",
		"Frontend",
		"Networking",
	}

	// ✅ Tri alphabétique
	sort.Strings(domains)

	return domains
}

// GetDomainOptionsOrdered - Slice ordonnée pour FilterDropdown
// "Tous domaines" en premier, puis alphabétique
func GetDomainOptionsOrdered() []struct {
	Value string
	Label string
} {
	// ✅ "Tous domaines" en premier (reset option)
	ordered := []struct {
		Value string
		Label string
	}{
		{"", "Tous domaines"},
	}

	// ✅ Ajoute domaines triés
	for _, domain := range GetDomains() {
		ordered = append(ordered, struct {
			Value string
			Label string
		}{domain, domain})
	}

	return ordered
}

// ────────────────────────────────────────────────────
// DIFFICULTIES - Ordre croissant
// ────────────────────────────────────────────────────

// GetDifficulties - Niveaux de difficulté (ordre fixe 1→4)
func GetDifficulties() []struct {
	Value int
	Label string
} {
	return []struct {
		Value int
		Label string
	}{
		{1, "1 - Facile"},
		{2, "2 - Moyen"},
		{3, "3 - Difficile"},
		{4, "4 - Expert"},
	}
}

// ────────────────────────────────────────────────────
// SORT OPTIONS - Ordre logique
// ────────────────────────────────────────────────────

// GetSortOptionsOrdered - Options de tri avec ordre logique
func GetSortOptionsOrdered() []struct {
	Value string
	Label string
} {
	return []struct {
		Value string
		Label string
	}{
		{"recent", "Plus récents"},
		{"difficulty", "Difficulté"},
		{"domain", "Domaine"},
	}
}
