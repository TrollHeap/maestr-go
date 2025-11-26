# Guide Complet Go 1.25+ : Performance, Patterns et Diagnostic Personnel

## Table des Matières

### Partie I : Fondamentaux de Performance
1. [Gestion Mémoire](#1-gestion-mémoire)
2. [Slices et Maps](#2-slices-et-maps)
3. [Strings et I/O](#3-strings-et-io)
4. [Retours de Fonction](#4-retours-de-fonction)
5. [Profiling et Mesure](#5-profiling-et-mesure)

### Partie II : Concepts Avancés Go 1.24+
6. [strings.SplitSeq - Concept de Base](#6-stringssplitseq---concept-de-base)
7. [Quand Utiliser SplitSeq](#7-quand-utiliser-splitseq)
8. [Gains de Performance](#8-gains-de-performance)
9. [Limitations et Trade-offs](#9-limitations-et-trade-offs)
10. [Exercices Pratiques SplitSeq](#10-exercices-pratiques-splitseq)

### Partie III : Patterns Essentiels
11. [Gestion d'Erreurs Moderne](#11-gestion-derreurs-moderne)
12. [Pointeurs et Performance](#12-pointeurs-et-performance)
13. [Préallocation Optimale](#13-préallocation-optimale)
14. [Maps et Clés Composites](#14-maps-et-clés-composites)

### Partie IV : Diagnostic et Amélioration
15. [Diagnostic des Concepts Manquants](#15-diagnostic-des-concepts-manquants)
16. [Plan d'Apprentissage Structuré](#16-plan-dapprentissage-structuré)
17. [Auto-Diagnostic Personnel](#17-auto-diagnostic-personnel)

### Annexes
- [Index des Concepts](#index-des-concepts)
- [Références et Ressources](#références-et-ressources)

---

# Partie I : Fondamentaux de Performance

## 1. Gestion Mémoire

### Stack vs Heap : Où sont allouées les variables

**Analogie visuelle :** Imagine un bureau (stack) et un entrepôt (heap)

```
STACK (Bureau - Rapide)          HEAP (Entrepôt - Plus lent)
┌─────────────────┐             ┌─────────────────────┐
│ Variables       │             │ Objets volumineux   │
│ locales         │             │ ou à durée de vie   │
│ (int, bool,     │             │ indéterminée        │
│  small structs) │             │ (slices, maps,      │
│                 │             │  grandes structs)   │
│ Accès: ~1ns     │             │ Accès: ~10ns        │
│ Nettoyage: auto │             │ Nettoyage: GC       │
└─────────────────┘             └─────────────────────┘
```

### Escape Analysis : Savoir quand Go met sur heap vs stack

**Outil diagnostic :**
```bash
go build -gcflags="-m" main.go
```

**Exemples pratiques :**

```go
// ✅ Reste sur STACK - variable locale simple
func stackExample() {
    x := 42  // Pas d'escape
    fmt.Println(x)
}

// ❌ Échappe au HEAP - pointeur retourné
func heapExample() *int {
    x := 42
    return &x  // ESCAPE: &x escapes to heap
}

// ✅ Reste sur STACK - slice petite et connue
func stackSlice() {
    data := make([]int, 10)  // Petite taille connue
    processLocally(data)
}

// ❌ Échappe au HEAP - taille dynamique ou retournée
func heapSlice(n int) []int {
    return make([]int, n)  // ESCAPE: taille variable
}
```

**Règles d'escape :**
- Retourner un pointeur → HEAP
- Assigner à une interface → HEAP  
- Slice/map de taille variable → HEAP
- Fermeture (closure) capturant des variables → HEAP

### Taille des types et alignement mémoire

```go
// Coûts mémoire par type (architecture 64-bit)
var sizes = map[string]int{
    "bool":     1,  // mais aligné sur 8 bytes
    "int8":     1,
    "int16":    2,
    "int32":    4,
    "int64":    8,
    "int":      8,  // 64-bit systems
    "string":   16, // pointeur(8) + longueur(8)
    "slice":    24, // pointeur(8) + len(8) + cap(8)
    "map":      8,  // pointeur vers structure interne
    "chan":     8,  // pointeur
    "interface": 16, // type(8) + valeur(8)
}

// ✅ Struct bien alignée (40 bytes)
type OptimalStruct struct {
    ID       int64   // 8 bytes
    Value    int64   // 8 bytes  
    Name     string  // 16 bytes
    Active   bool    // 1 byte + 7 padding = 8 bytes
}

// ❌ Struct mal alignée (48 bytes à cause du padding)
type BadStruct struct {
    Active   bool    // 1 byte + 7 padding
    ID       int64   // 8 bytes
    Name     string  // 16 bytes
    Value    int64   // 8 bytes
}
```

---

## 2. Slices et Maps

### Préallocation : make([]T, 0, capacity)

**Le problème des append dynamiques :**

```go
// ❌ Croissance dynamique - coûteux
func badAppend() []string {
    var result []string  // capacity = 0

    for i := 0; i < 1000; i++ {
        result = append(result, fmt.Sprintf("item%d", i))
        // À chaque dépassement de capacité :
        // 1. Allouer nouveau slice (2x la taille)
        // 2. Copier ancien → nouveau (memcpy)
        // 3. Marquer ancien pour GC
    }
    return result
}

// ✅ Préallocation - une seule allocation
func goodPrealloc() []string {
    result := make([]string, 0, 1000)  // Capacité connue

    for i := 0; i < 1000; i++ {
        result = append(result, fmt.Sprintf("item%d", i))
        // Pas de réallocation, juste assignation
    }
    return result
}
```

**Séquence de croissance Go pour 50 éléments :**
```
append 1  → cap=1   (1 allocation)
append 2  → cap=2   (2 allocations + 1 copie) 
append 3  → cap=4   (3 allocations + 2 copies)
append 5  → cap=8   (4 allocations + 3 copies)
append 9  → cap=16  (5 allocations + 4 copies)
append 17 → cap=32  (6 allocations + 5 copies)
append 33 → cap=64  (7 allocations + 6 copies)
─────────────────────────────────────────────
TOTAL: 7 allocations + 6 copies complètes

Avec make([]T, 0, 50): 1 allocation + 0 copie
```

### Réinitialisation : slice = slice[:0]

```go
// ❌ Mauvaise réutilisation - nouvelle allocation
func badReuse() {
    for i := 0; i < 100; i++ {
        data := make([]string, 0, 50)  // 100 allocations !
        // ... utiliser data
    }
}

// ✅ Bonne réutilisation - même mémoire 
func goodReuse() {
    data := make([]string, 0, 50)  // 1 seule allocation

    for i := 0; i < 100; i++ {
        data = data[:0]  // Réinitialise len=0, garde cap=50
        // ... réutiliser data
    }
}
```

### Maps : préallocation et optimisations

```go
// ❌ Map non préallouée
func badMap() map[string]int {
    m := make(map[string]int)  // Taille par défaut

    for i := 0; i < 1000; i++ {
        m[fmt.Sprintf("key%d", i)] = i
        // Resize interne à 8, 16, 32, 64... buckets
    }
    return m
}

// ✅ Map préallouée
func goodMap() map[string]int {
    m := make(map[string]int, 1000)  // Évite les rehash

    for i := 0; i < 1000; i++ {
        m[fmt.Sprintf("key%d", i)] = i
    }
    return m
}

// ✅ Set optimisé avec struct{} (0 bytes)
func optimizedSet() map[string]struct{} {
    set := make(map[string]struct{}, 100)

    set["item1"] = struct{}{}  // Valeur 0 bytes
    set["item2"] = struct{}{}

    // Vérifier présence
    if _, exists := set["item1"]; exists {
        // Item trouvé
    }
    return set
}
```

---

## 3. Strings et I/O

### strings.Builder pour concaténation

**Le problème de l'opérateur + :**

```go
// ❌ Opérateur + en boucle - O(n²) complexity
func badConcat(items []string) string {
    result := ""
    for _, item := range items {
        result += item + ","  // Nouvelle allocation à chaque +
        // Chaque + crée un nouveau string et copie tout
    }
    return result
}

// ✅ strings.Builder - O(n) complexity  
func goodConcat(items []string) string {
    var builder strings.Builder
    builder.Grow(len(items) * 10)  // Préallocation estimée

    for _, item := range items {
        builder.WriteString(item)
        builder.WriteByte(',')
    }
    return builder.String()
}
```

**Benchmark pour 1000 strings :**
```
BenchmarkBadConcat    1000    1500000 ns/op    500000 B/op    1000 allocs/op
BenchmarkGoodConcat   5000     300000 ns/op     10240 B/op       5 allocs/op
                                 ↑ 5x plus rapide      ↑ 50x moins d'allocations
```

### bufio.Scanner vs os.ReadFile

```go
// ✅ Petits fichiers : os.ReadFile + string operations
func readSmallFile(filename string) ([]string, error) {
    data, err := os.ReadFile(filename)  // Charge tout en mémoire
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(data), "\n")
    return lines, nil
}

// ✅ Gros fichiers : bufio.Scanner (streaming)
func readLargeFile(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {  // Lit ligne par ligne
        lines = append(lines, scanner.Text())
    }

    return lines, scanner.Err()
}
```

### Éviter conversions []byte ↔ string

```go
// ❌ Conversions multiples
func badConversions(data []byte) []string {
    text := string(data)              // Conversion 1
    lines := strings.Split(text, "\n") // string operations

    var result []string
    for _, line := range lines {
        trimmed := strings.TrimSpace(line)  // Plus de string ops
        if trimmed != "" {
            result = append(result, trimmed)
        }
    }
    return result
}

// ✅ Operations sur []byte directement  
func goodConversions(data []byte) []string {
    var result []string

    for len(data) > 0 {
        // Trouve \n sans conversion string
        i := bytes.IndexByte(data, '\n')
        if i == -1 {
            i = len(data)
        }

        line := bytes.TrimSpace(data[:i])
        if len(line) > 0 {
            result = append(result, string(line))  // Conversion unique
        }

        data = data[i+1:]
    }
    return result
}
```

---

# Partie II : Concepts Avancés Go 1.24+

## 6. strings.SplitSeq - Concept de Base

### Qu'est-ce qu'un itérateur lazy

**Analogie :** Imagine une machine distributrice de bonbons

```
SPLIT (Machine qui donne tout d'un coup)    |  SPLITSEQ (Machine à pièce)
┌─────────────────────────────────────┐    |  ┌────────────────────────────┐
│ Input: "a,b,c,d,e"                  │    |  │ Input: "a,b,c,d,e"         │
│          ↓                          │    |  │          ↓                 │
│ Alloue: []string{"a","b","c","d","e"}│    |  │ for item := range SplitSeq │
│ (tout en mémoire immédiatement)      │    |  │   yield "a"  ← première    │
│                                     │    |  │   yield "b"  ← suivante    │
│ Coût: 5 strings + 1 slice          │    |  │   yield "c"  ← à la demande│
│ Mémoire: ~120 bytes                 │    |  │                            │
└─────────────────────────────────────┘    |  │ Coût: 1 string à la fois  │
                                           |  │ Mémoire: ~24 bytes         │
                                           |  └────────────────────────────┘
```

### Split vs SplitSeq : différences fondamentales

```go
// strings.Split - Eager (tout d'un coup)
func demoSplit() {
    data := "ligne1\nligne2\nligne3\nligne4\nligne5"

    lines := strings.Split(data, "\n")  // ← Allocation du slice complet
    // Mémoire: []string{} + 5 sous-strings = ~120 bytes

    for _, line := range lines {
        if strings.Contains(line, "3") {
            fmt.Println("Trouvé:", line)
            return  // Mais on a quand même alloué lines[4] et lines[5]
        }
    }
}

// strings.SplitSeq - Lazy (à la demande) - Go 1.24+
func demoSplitSeq() {
    data := "ligne1\nligne2\nligne3\nligne4\nligne5"

    // Pas d'allocation ici, juste une fonction
    for line := range strings.SplitSeq(data, "\n") {
        if strings.Contains(line, "3") {
            fmt.Println("Trouvé:", line)
            return  // ligne4 et ligne5 ne sont jamais créées !
        }
    }
}
```

### Comment ça marche (iter.Seq[string])

```go
// Signature de strings.SplitSeq (Go 1.24)
func SplitSeq(s, sep string) iter.Seq[string] {
    return func(yield func(string) bool) {
        // yield est appelé pour chaque segment trouvé
        for {
            if i := strings.Index(s, sep); i >= 0 {
                if !yield(s[:i]) {  // Retourne le segment
                    return  // Arrêt si yield retourne false
                }
                s = s[i+len(sep):]
            } else {
                yield(s)  // Dernier segment
                return
            }
        }
    }
}
```

---

## 7. Quand Utiliser SplitSeq

### ✅ Cas 1 : Boucle une seule fois

```go
// ✅ SplitSeq optimal - itération unique
func parseConfigSplitSeq(data string) map[string]string {
    config := make(map[string]string)

    for line := range strings.SplitSeq(data, "\n") {
        line = strings.TrimSpace(line)
        if line == "" || line[0] == '#' {
            continue  // Skip vides et commentaires
        }

        if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
            config[parts[0]] = parts[1]
        }
    }
    return config
}

// ❌ Split inutile - allocation du slice complet non utilisé
func parseConfigSplit(data string) map[string]string {
    config := make(map[string]string)
    lines := strings.Split(data, "\n")  // Alloue tout même si break early

    for _, line := range lines {
        // ... même logique
    }
    return config
}
```

### ✅ Cas 2 : Early exit (recherche)

```go
// ✅ SplitSeq avec early exit - optimal
func findInLogSplitSeq(logData string, pattern string) (string, bool) {
    for line := range strings.SplitSeq(logData, "\n") {
        if strings.Contains(line, pattern) {
            return line, true  // Arrêt dès trouvé
        }
    }
    return "", false
}

// ❌ Split gaspille - parse toutes les lignes même après match
func findInLogSplit(logData string, pattern string) (string, bool) {
    lines := strings.Split(logData, "\n")  // Alloue 100% du contenu

    for _, line := range lines {
        if strings.Contains(line, pattern) {
            return line, true  // Mais 90% des lignes sont gaspillées
        }
    }
    return "", false
}
```

### ❌ Besoin du slice complet

```go
// ❌ Ne pas utiliser SplitSeq ici
func needFullSliceBad() {
    data := readCSVFile()

    // ERREUR: Impossible avec SplitSeq
    // lines := strings.SplitSeq(data, "\n")  
    // fmt.Println("Nombre de lignes:", len(lines))  ← len() impossible
    // lastLine := lines[len(lines)-1]              ← index impossible
}

// ✅ Split obligatoire pour accès indices/taille
func needFullSliceGood() {
    data := readCSVFile()

    lines := strings.Split(data, "\n")  // Split requis
    fmt.Println("Nombre de lignes:", len(lines))

    if len(lines) > 0 {
        lastLine := lines[len(lines)-1]
        fmt.Println("Dernière ligne:", lastLine)
    }

    // Accès aléatoire
    if len(lines) > 10 {
        middleLine := lines[len(lines)/2]
        fmt.Println("Ligne du milieu:", middleLine)
    }
}
```

---

## 8. Gains de Performance

### Benchmarks mesurés

```go
func BenchmarkSplit(b *testing.B) {
    data := strings.Repeat("ligne de test\n", 1000)  // 1000 lignes

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        lines := strings.Split(data, "\n")

        for _, line := range lines {
            if strings.Contains(line, "500") {  // Trouve vers la moitié
                break
            }
        }
    }
}

func BenchmarkSplitSeq(b *testing.B) {
    data := strings.Repeat("ligne de test\n", 1000)  // 1000 lignes

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for line := range strings.SplitSeq(data, "\n") {
            if strings.Contains(line, "500") {  // Trouve vers la moitié
                break
            }
        }
    }
}
```

**Résultats typiques :**
```
BenchmarkSplit      5000    240000 ns/op    85000 B/op     1001 allocs/op
BenchmarkSplitSeq   6500    190000 ns/op    35000 B/op      501 allocs/op
                            ↑ -21% temps    ↑ -59% mémoire  ↑ -50% allocs
```

### Mécanisme de substring (pas de copie)

```go
// Exemple d'optimisation substring
original := "ligne1\nligne2\nligne3"

// strings.Split crée des COPIES
lines := strings.Split(original, "\n")
// lines[0] = "ligne1"  ← nouvelle allocation
// lines[1] = "ligne2"  ← nouvelle allocation  
// lines[2] = "ligne3"  ← nouvelle allocation

// strings.SplitSeq utilise des SUBSTRINGS (même mémoire)
for line := range strings.SplitSeq(original, "\n") {
    // line pointe vers original[0:6], original[7:13], etc.
    // Pas de copie, juste des pointeurs + longueurs
}
```

**Visualisation mémoire :**
```
Mémoire originale: "ligne1\nligne2\nligne3"
                    ├─────┤ ├─────┤ ├─────┤
                        │       │       │
Split (copies):     "ligne1" "ligne2" "ligne3"  ← 3 allocations

SplitSeq (refs):    same memory, different views ← 0 allocations
```

---

## 9. Limitations et Trade-offs

### Bug connu : escape analysis sous-optimale (Go 1.24)

```go
// Problème actuel Go 1.24 - closure échappe au heap
func demonstrateBug() {
    data := "a,b,c,d,e"

    // Cette closure échappe au heap même si pas nécessaire
    iter := strings.SplitSeq(data, ",")

    // Workaround temporaire : forcer inline
    func() {
        for item := range iter {
            processItem(item)
        }
    }()
}
```

### Coût de closure (~40 bytes)

```go
// Overhead de la closure SplitSeq
type SplitSeqClosure struct {
    s   string  // 16 bytes (string header)
    sep string  // 16 bytes (string header)  
    pos int     // 8 bytes (position courante)
    // Total: ~40 bytes par SplitSeq créé
}

// Pour de TRÈS petits datasets, Split peut être plus efficace
func microOptimization() {
    tiny := "a,b"  // 2 éléments seulement

    // Split: 24 bytes (slice header) + 32 bytes (2 strings) = 56 bytes
    parts := strings.Split(tiny, ",")

    // SplitSeq: 40 bytes (closure) + overhead itération = ~45 bytes
    // Mais différence négligeable en pratique
}
```

---

# Partie III : Patterns Essentiels

## 11. Gestion d'Erreurs Moderne

### Error wrapping avec fmt.Errorf(..., %w, err)

**Principe central :** Créer une chaîne d'erreurs avec contexte du bas niveau au haut niveau.

```go
// Couche 1 : Bas niveau (OS/système)
func readSysFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        // ✅ Wrap avec le chemin exact pour debugging
        return nil, fmt.Errorf("lecture fichier %q: %w", path, err)
    }
    return data, nil
}

// Couche 2 : Logique métier (parsing GPU)  
func readUevent(cardPath string) (UeventInfo, error) {
    data, err := readSysFile(filepath.Join(cardPath, "uevent"))
    if err != nil {
        // ✅ Wrap avec le contexte métier
        return UeventInfo{}, fmt.Errorf("lecture uevent pour %q: %w", cardPath, err)
    }

    info := parseUevent(data)
    if info.PCISlot == "" {
        // ✅ Erreur métier avec contexte
        return UeventInfo{}, fmt.Errorf("PCI slot manquant dans %q", cardPath)
    }

    return info, nil
}

// Couche 3 : API utilisateur
func GetGPUInfo(cardID string) (*GPUInfo, error) {
    cardPath := fmt.Sprintf("/sys/class/drm/card%s", cardID)

    uevent, err := readUevent(cardPath)
    if err != nil {
        // ✅ Wrap avec contexte utilisateur final
        return nil, fmt.Errorf("impossible de lire GPU %q: %w", cardID, err)
    }

    // ... continuer processing
    return gpu, nil
}
```

### Chaînes d'erreurs - exemple complet

```go
// Appel qui échoue
gpu, err := GetGPUInfo("1")
if err != nil {
    fmt.Printf("Erreur: %v\n", err)

    // Affiche la chaîne complète:
    // impossible de lire GPU "1": lecture uevent pour "/sys/class/drm/card1": 
    // lecture fichier "/sys/class/drm/card1/uevent": open /sys/class/drm/card1/uevent: 
    // no such file or directory
}

// Unwrapping pour accéder aux erreurs sous-jacentes
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Printf("Problème de chemin: %s\n", pathErr.Path)
}

// Vérifier type d'erreur spécifique
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("Le fichier n'existe pas")
}
```

### Quand arrêter la propagation

```go
// ✅ Arrêter et logger au niveau approprié
func ProcessAllGPUs() error {
    cards, err := listGPUCards()
    if err != nil {
        // Erreur critique - on ne peut pas continuer
        return fmt.Errorf("impossible de lister les GPUs: %w", err)
    }

    var errors []error
    for _, cardID := range cards {
        gpu, err := GetGPUInfo(cardID)
        if err != nil {
            // ✅ Log et continue - erreur non-critique
            log.Printf("GPU %s ignoré: %v", cardID, err)
            errors = append(errors, err)
            continue
        }

        processGPU(gpu)
    }

    // Retourner erreur composite si nécessaire
    if len(errors) == len(cards) {
        return fmt.Errorf("aucun GPU lisible: %d erreurs", len(errors))
    }

    return nil
}
```

---

## 12. Pointeurs et Performance

### Retour par valeur vs pointeur

**Règle fondamentale :** La taille compte plus que le type.

```go
// Struct petite (≤ 32 bytes) → retourner par VALEUR
type SmallStruct struct {
    ID    int64   // 8 bytes
    Value int64   // 8 bytes
    Flags uint32  // 4 bytes + 4 padding = 8 bytes
    // Total: 24 bytes → OK pour valeur
}

func ProcessSmall() SmallStruct {  // ✅ Par valeur
    return SmallStruct{ID: 1, Value: 42, Flags: 0xFF}
}

// Struct volumineuse (> 32 bytes) → retourner par POINTEUR  
type LargeStruct struct {
    ID          int64     // 8 bytes
    Name        string    // 16 bytes
    Description string    // 16 bytes  
    Tags        []string  // 24 bytes
    Properties  map[string]interface{}  // 8 bytes
    CreatedAt   time.Time // 24 bytes
    // Total: 96 bytes → Trop gros, utiliser pointeur
}

func ProcessLarge() *LargeStruct {  // ✅ Par pointeur
    return &LargeStruct{
        ID:   1,
        Name: "Example",
        // ... init
    }
}
```

### Visualisation mémoire stack/heap

```go
// Scénario : fonction qui crée et retourne une struct
func createData() (SmallStruct, *LargeStruct) {
    // small allouée sur STACK de createData
    small := SmallStruct{ID: 1}  

    // large allouée sur STACK de createData
    large := LargeStruct{ID: 2}

    // Retour par valeur → COPIE 24 bytes sur stack de l'appelant
    // Retour par pointeur → COPIE 8 bytes (adresse), large échappe au HEAP
    return small, &large
}

func caller() {
    small, large := createData()

    // Mémoire finale:
    // small: 24 bytes sur STACK de caller (rapide)
    // large: 96 bytes sur HEAP + 8 bytes pointeur sur STACK (plus lent)
}
```

### Piège courant : struct vs *struct en paramètres

```go
// ❌ Mauvais : grande struct par valeur en paramètre
func ProcessByValue(data LargeStruct) {  // Copie 96 bytes !
    // Modification locale, ne change pas l'original
    data.Name = "Modified"
}

// ✅ Bon : pointeur pour éviter la copie
func ProcessByPointer(data *LargeStruct) {  // Copie 8 bytes seulement
    // Modification de l'original
    data.Name = "Modified"
}

// ✅ Alternative : méthode avec receiver par pointeur
func (l *LargeStruct) Process() {  // Receiver par pointeur si struct > 32 bytes
    l.Name = "Modified"
}
```

---

## 13. Préallocation Optimale

### Pourquoi make([]string, 0, 50) au lieu de var

**Le coût caché des réallocations :**

```go
// ❌ var result []string → capacité 0, réallocations multiples
func parseProcCPUInfoBad(data string) []CPUInfo {
    var cpus []CPUInfo  // cap=0, len=0

    lines := strings.Split(data, "\n")  // ~4000 lignes

    for _, line := range lines {
        if strings.HasPrefix(line, "processor") {
            // À chaque append: vérifier capacité, possiblement réallouer
            cpus = append(cpus, parseCPU(line))
        }
    }
    return cpus
}

// ✅ Estimation et préallocation
func parseProcCPUInfoGood(data string) []CPUInfo {
    // Estimation: 1 CPU par 100 lignes dans /proc/cpuinfo
    estimatedCPUs := len(strings.Split(data, "\n")) / 100
    if estimatedCPUs < 4 { estimatedCPUs = 4 }    // Minimum raisonnable
    if estimatedCPUs > 128 { estimatedCPUs = 128 } // Maximum réaliste

    cpus := make([]CPUInfo, 0, estimatedCPUs)  // Une seule allocation

    for line := range strings.SplitSeq(data, "\n") {
        if strings.HasPrefix(line, "processor") {
            cpus = append(cpus, parseCPU(line))
        }
    }
    return cpus
}
```

### Mécanisme de croissance des slices Go

```go
// Politique de croissance interne de Go
func growSlice(oldCap, minCap int) int {
    if minCap > oldCap*2 {
        return minCap
    }

    if oldCap < 1024 {
        return oldCap * 2  // Double jusqu'à 1024
    } else {
        // Croissance de 25% au-dessus de 1024
        newCap := oldCap + oldCap/4
        for newCap < minCap {
            newCap += newCap / 4
        }
        return newCap
    }
}

// Exemple concret pour atteindre 100 éléments
var trace []string
fmt.Printf("cap=%d ", cap(trace))      // 0

trace = append(trace, "1")
fmt.Printf("cap=%d ", cap(trace))      // 1

for i := 2; i <= 100; i++ {
    oldCap := cap(trace)
    trace = append(trace, fmt.Sprintf("%d", i))
    if cap(trace) != oldCap {
        fmt.Printf("→ cap=%d ", cap(trace))
    }
}
// Sortie: cap=0 cap=1 → cap=2 → cap=4 → cap=8 → cap=16 → cap=32 → cap=64 → cap=128
// 7 réallocations pour 100 éléments
```

---

## 14. Maps et Clés Composites

### Maps : principe du dictionnaire

**Analogie visuelle :** Un dictionnaire français-anglais

```
┌─────────────────────┬─────────────────────┐
│       CLÉ           │       VALEUR        │
├─────────────────────┼─────────────────────┤
│ "chat"              │ "cat"               │
│ "chien"             │ "dog"               │  
│ "ordinateur"        │ "computer"          │
│ "performance"       │ "performance"       │
└─────────────────────┴─────────────────────┘

En Go:
ages := map[string]int{
    "Alice": 30,    // clé: "Alice", valeur: 30
    "Bob":   25,    // clé: "Bob",   valeur: 25  
}

fmt.Println(ages["Alice"])  // Résultat: 30
```

### Structs comme clés de map

**Cas d'usage :** Identifier uniquement un cache CPU

```go
// Problème : les caches L3 sont partagés entre cores
// /sys/devices/system/cpu/cpu0/cache/index3/shared_cpu_list: "0-7"
// /sys/devices/system/cpu/cpu1/cache/index3/shared_cpu_list: "0-7"  ← MÊME cache !

// ✅ Solution : struct comme clé composite
type CacheKey struct {
    Level          int     // 1, 2, 3  
    SharedCPUList  string  // "0-7", "8-15", etc.
}

func deduplicateCaches(cpuCaches []CacheInfo) map[CacheKey]CacheInfo {
    uniqueCaches := make(map[CacheKey]CacheInfo)

    for _, cache := range cpuCaches {
        key := CacheKey{
            Level:         cache.Level,
            SharedCPUList: cache.SharedCPUList,
        }

        // Chaque combinaison unique (level, shared_cpu_list) = 1 entrée
        uniqueCaches[key] = cache
    }

    return uniqueCaches
    // Résultat: map avec 3 entrées au lieu de 24 (8 cores × 3 levels)
}
```

### Pourquoi string et pas []int pour clés

```go
// ❌ []int comme clé - NE COMPILE PAS
func badMapKey() {
    // ERREUR: invalid map key type []int
    // m := make(map[[]int]string)  
}

// ✅ string comme clé - FONCTIONNE
func goodMapKey() {
    m := make(map[string]string)  // ✅ string est comparable
    m["0-7"] = "L3 Cache Group 1"
    m["8-15"] = "L3 Cache Group 2"
}

// ✅ Alternative avec struct contenant des types comparables
type CPURange struct {
    Start int  // comparable
    End   int  // comparable
}

func structMapKey() {
    m := make(map[CPURange]string)  // ✅ struct avec champs comparables
    m[CPURange{0, 7}] = "L3 Cache Group 1"
    m[CPURange{8, 15}] = "L3 Cache Group 2"
}
```

**Types utilisables comme clés de map :**
- ✅ Types de base : `bool`, `int*`, `uint*`, `float*`, `string`
- ✅ Pointeurs : `*T`
- ✅ Arrays : `[N]T` où T est comparable  
- ✅ Structs : tous les champs doivent être comparables
- ✅ Interfaces : si la valeur sous-jacente est comparable
- ❌ Slices : `[]T`
- ❌ Maps : `map[K]V`
- ❌ Functions : `func(...) ...`

---

# Partie IV : Diagnostic et Amélioration

## 15. Diagnostic des Concepts Manquants

### 🔴 Niveau Critique (7h) - À maîtriser immédiatement

**1. Gestion d'erreurs moderne (2h)**
```go
// ❌ TON CODE INITIAL - masque les erreurs
func readLspci(pciSlot string) (string, string) {
    out, err := exec.Command("lspci", "-mm", "-nn", "-D").Output()
    if err != nil {
        return "Unknown", "Unknown"  // ← Perte de contexte !
    }
    // ...
}

// ✅ VERSION CORRIGÉE - erreur traceable
func readLspci(pciSlot string) (PCIDevice, error) {
    out, err := exec.Command("lspci", "-mm", "-nn", "-D", "-s", pciSlot).Output()
    if err != nil {
        return PCIDevice{}, fmt.Errorf("échec lspci pour slot %q: %w", pciSlot, err)
    }
    // ...
}
```

**2. Validation d'entrées (1.5h)**  
```go
// ✅ Regex pour valider PCI slot format
var pciSlotRegex = regexp.MustCompile(`^[0-9a-fA-F]{4}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}\.[0-7]$`)

func validatePCISlot(slot string) error {
    if !pciSlotRegex.MatchString(slot) {
        return fmt.Errorf("format PCI slot invalide %q, attendu: XXXX:XX:XX.X", slot)
    }
    return nil
}
```

**3. Types personnalisés (2h)**
```go
// ✅ Types métier au lieu de strings magiques
type Vendor string
const (
    VendorAMD    Vendor = "AMD"
    VendorNVIDIA Vendor = "NVIDIA" 
    VendorIntel  Vendor = "Intel"
)

type PCIDevice struct {
    Slot   string
    Vendor Vendor  // ← Type fort au lieu de string
    Model  string
}
```

**4. SplitSeq pour parsing (1.5h)**
```go
// ✅ Utiliser SplitSeq pour parsing /sys
for line := range strings.SplitSeq(string(data), "\n") {
    line = strings.TrimSpace(line)
    if line == "" { continue }

    if key, value, found := strings.Cut(line, "="); found {
        processKeyValue(key, value)
    }
}
```

### 🟠 Niveau Important (6h) - Performance et idiomes

**1. Préallocation systématique (1.5h)**
```go
// Estimer la capacité selon le contexte
func estimateCapacity(source string) int {
    switch {
    case strings.Contains(source, "/proc/cpuinfo"):
        return 50  // ~50 lignes par CPU
    case strings.Contains(source, "lspci"):
        return 100 // ~100 périphériques max
    default:
        return 16  // Valeur par défaut raisonnable
    }
}
```

**2. Retour d'erreur vs valeurs par défaut (2h)**
```go
// ❌ Masquer avec valeurs par défaut
func getGPUMemory() int64 { return 8000000000 } // "8GB par défaut"

// ✅ Erreur explicite  
func getGPUMemory() (int64, error) {
    // tentative de lecture...
    return 0, fmt.Errorf("impossible de lire la mémoire GPU")
}
```

### 🟡 Niveau Intermédiaire (4.5h) - Maintenabilité

**Documentation, testabilité, séparation des responsabilités**

---

## 16. Plan d'Apprentissage Structuré

### Planning sur 3 semaines (17h30 total)

**Semaine 1 : Concepts critiques (7h)**
- Lundi : Error wrapping (2h)
  - Pratique : Refactoriser ton code GPU avec fmt.Errorf(%w)
  - Validation : Tracer une erreur sur 3 niveaux
- Mercredi : Validation entrées (1.5h)  
  - Pratique : Regex PCI slots, validation chemins /sys
  - Validation : 0 panic sur entrées malformées
- Vendredi : Types + SplitSeq (3.5h)
  - Pratique : Enum Vendor, refactor avec SplitSeq
  - Validation : Benchmark avant/après

**Semaine 2 : Performance et idiomes (6h)**  
- Lundi : Préallocation (2h)
  - Pratique : make([]T, 0, cap) sur tous tes parsers
  - Validation : -50% allocations au benchmark
- Mercredi : Structs vs maps (2h)
  - Pratique : CacheKey, CPUInfo bien structurés
  - Validation : 0 fonction avec >3 retours
- Vendredi : Profiling (2h)
  - Pratique : go test -bench + pprof
  - Validation : Identifier le bottleneck #1

**Semaine 3 : Maintenabilité (4.5h)**
- Lundi : Documentation (1.5h)
- Mercredi : Tests unitaires (2h) 
- Vendredi : Refactoring final (1h)

### Ressources par priorité

**🔴 Critiques (consulter cette semaine)**
- Go Error Handling Best Practices
- OWASP Go Security Guide  
- Go 1.24 Release Notes (SplitSeq)

**🟠 Importantes (mois prochain)**
- Effective Go (official)
- Google Go Style Guide
- Go Memory Model

**🟡 Complémentaires (quand tu as le temps)**
- Go Doc Comments Guide
- Advanced Go Patterns

---

## 17. Auto-Diagnostic Personnel

### ✅ Tes Forces Identifiées

**Pipeline de Données - Pattern maîtrisé à 90%**
```go
// Ton style naturel - très bon !
func tonPattern() {
    // 1. Path construction
    path := filepath.Join("/sys", "devices", "cpu0")

    // 2. Glob pour lister
    files, _ := filepath.Glob(path + "/*")

    // 3. Read file content  
    data, _ := os.ReadFile(files[0])

    // 4. Trim et clean
    content := strings.TrimSpace(string(data))

    // 5. Convert types
    value, _ := strconv.Atoi(content)

    // 6. Loop processing
    for _, item := range items {
        process(item)
    }
}
```

**Tu maîtrises parfaitement :**
- `filepath.Join`, `filepath.Glob`
- `os.ReadFile`, `strings.TrimSpace`  
- `strconv.Atoi`, `strconv.ParseInt`
- Boucles `for range` simples

### 🔴 Faiblesses à Corriger

**1. Structs - Règle des ≥3 valeurs liées**
```go
// ❌ TON CODE - trop de retours individuels  
func getCPUInfo() (string, string, int, bool, error) {
    return vendor, model, cores, hyperthreading, nil
}

// ✅ VERSION STRUCTURÉE
type CPUInfo struct {
    Vendor         string
    Model          string  
    Cores          int
    Hyperthreading bool
}

func getCPUInfo() (CPUInfo, error) {
    return CPUInfo{...}, nil
}
```

**2. Maps - Règle du "déjà vu?"**
```go
// ❌ TON STYLE - slice avec recherche O(n)
func isDuplicateUSB(usbID string, seen []string) bool {
    for _, id := range seen {  // ← O(n) à chaque appel
        if id == usbID {
            return true
        }
    }
    return false
}

// ✅ VERSION OPTIMISÉE - map O(1)  
func checkUSBDuplicates(devices []USBDevice) []USBDevice {
    seen := make(map[string]bool)
    unique := make([]USBDevice, 0, len(devices))

    for _, device := range devices {
        if !seen[device.ID] {  // ← O(1) lookup
            seen[device.ID] = true
            unique = append(unique, device)
        }
    }
    return unique
}
```

### 📋 Plan d'Amélioration Personnalisé

**Semaine 1 : Structs**
- **Exercice pratique :** Dans ton prochain projet, dès que tu as ≥3 variables liées, créer une struct
- **Validation :** 0 fonction avec >3 paramètres de retour
- **Exemple :** `type GPUInfo struct { Vendor, Model, Memory }`

**Semaine 2 : Maps + Structs**  
- **Exercice pratique :** Combinaisons struct en tant que clé de map
- **Validation :** Utiliser map dès que tu vérifies "déjà vu?"  
- **Exemple :** `map[CacheKey]CacheInfo` pour déduplication

**Semaine 3 : Conditions imbriquées**
- **Exercice pratique :** Pattern `if + continue` pour filtrage
- **Validation :** Éviter else imbriqués >2 niveaux
- **Exemple :** Parser avec early continue sur lignes vides

**Stratégie générale :** 1 struct OU 1 map par projet pour t'habituer progressivement.

---

# Index des Concepts

**A-C**
- Arrays vs Slices → [Slices et Maps](#2-slices-et-maps)
- Benchmarking → [Profiling et Mesure](#5-profiling-et-mesure)  
- Capacity estimation → [Préallocation Optimale](#13-préallocation-optimale)
- Closure overhead → [Limitations SplitSeq](#9-limitations-et-trade-offs)

**D-H**  
- Error wrapping → [Gestion d'Erreurs](#11-gestion-derreurs-moderne)
- Escape analysis → [Gestion Mémoire](#1-gestion-mémoire)
- Heap vs Stack → [Gestion Mémoire](#1-gestion-mémoire)

**I-P**
- Itérateurs lazy → [strings.SplitSeq](#6-stringssplitseq---concept-de-base)
- Maps, clés composites → [Maps et Clés Composites](#14-maps-et-clés-composites)  
- Performance benchmarks → [Gains de Performance](#8-gains-de-performance)
- Pointeurs → [Pointeurs et Performance](#12-pointeurs-et-performance)
- Préallocation → [Préallocation Optimale](#13-préallocation-optimale)

**S-Z**
- SplitSeq → [Concepts Avancés Go 1.24+](#partie-ii--concepts-avancés-go-124)
- strings.Builder → [Strings et I/O](#3-strings-et-io)
- Structs comme clés → [Maps et Clés Composites](#14-maps-et-clés-composites)

---

# Références et Ressources

**Documentation officielle Go 1.24+**
- [Go 1.24 Release Notes](https://golang.org/doc/go1.24) - SplitSeq et itérateurs
- [Effective Go](https://golang.org/doc/effective_go.html) - Styles et idiomes  
- [Go Memory Model](https://golang.org/ref/mem) - Concurrence et mémoire

**Performance et Profiling**  
- [Go pprof Guide](https://blog.golang.org/pprof) - CPU et memory profiling
- [Go Benchmarking](https://pkg.go.dev/testing#hdr-Benchmarks) - Mesurer les performances

**Sécurité et Bonnes Pratiques**
- [Go Security Guide](https://github.com/Checkmarx/Go-SCP) - Pratiques sécurisées
- [Google Go Style Guide](https://google.github.io/styleguide/go/) - Standards industrie

**Ce guide en version PDF**
- Généré le 27 octobre 2025 à partir de 13 PDFs consolidés
- Version Markdown source disponible pour modifications
- Mis à jour avec Go 1.25+ et retour d'expérience personnalisé
