<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# 🎨 REFACTORING MAESTRO → RETROWAVE TERMINAL

## 📸 VISUALISATION MENTALE

```
🌊 RETROWAVE TERMINAL VISION
┌────────────────────────────────────────────────────┐
│ AVANT: 26 fichiers templ éparpillés               │
│ APRÈS: Architecture modulaire 3-tiers             │
├────────────────────────────────────────────────────┤
│                                                    │
│  📁 internal/views/                               │
│  ├─ 🎨 ui/          → Design system atoms         │
│  ├─ 🧩 components/  → Business composants         │
│  ├─ 📄 pages/       → Vues complètes              │
│  └─ 🎭 layouts/     → Wrappers HTML               │
│                                                    │
│  📁 internal/views/logic/                         │
│  ├─ styles.go       → Classes CSS computées       │
│  ├─ builders.go     → URL/Query helpers           │
│  ├─ formatters.go   → Date/Number formats         │
│  └─ validators.go   → Règles validation           │
│                                                    │
│  🎨 RETROWAVE PALETTE                             │
│  ├─ Cyan néon     #00E5FF (primary)              │
│  ├─ Magenta néon  #FF10F0 (accent)               │
│  ├─ Violet profond #1A0033 (bg)                  │
│  └─ Grille CRT    rgba(0,229,255,0.1)            │
└────────────────────────────────────────────────────┘
```


## 🧠 ANCRAGE MNÉMOTECHNIQUE

**ATOMIC = Architecture Totalement Optimisée Modulaire Isolée Composable**

- **A**toms → Boutons, inputs, badges (ui/)
- **T**okens → Variables CSS centralisées
- **O**rchestration → Logic/ gère calculs
- **M**olecules → Composants métier (components/)
- **I**solation → Zéro duplication
- **C**omposition → Tout est réutilisable


## 📑 NOUVELLE ARCHITECTURE

### Structure Cible

```bash
internal/views/
├── ui/                    # 🎨 DESIGN SYSTEM ATOMS
│   ├── Button.templ       # Bouton universel
│   ├── Badge.templ        # Badge unifié
│   ├── Card.templ         # Card container
│   ├── Input.templ        # Input + validation
│   ├── Select.templ       # Dropdown select
│   ├── Progress.templ     # Barre progression
│   ├── Icon.templ         # Icônes SVG inline
│   └── Spinner.templ      # Loading state
│
├── components/            # 🧩 BUSINESS COMPONENTS
│   ├── ExerciseCard.templ    # Uses ui/Card + ui/Badge
│   ├── FilterBar.templ       # Uses ui/Select + ui/Input
│   ├── ReviewPanel.templ     # Uses ui/Button + ui/Card
│   ├── PlannerView.templ     # Unifié (day/week/month)
│   └── StepsManager.templ    # Uses ui/Progress + ui/Icon
│
├── pages/                 # 📄 FULL PAGES
│   ├── Dashboard.templ
│   ├── ExerciseList.templ
│   ├── ExerciseDetail.templ
│   ├── PlannerPage.templ
│   └── SessionBuilder.templ
│
├── layouts/               # 🎭 HTML WRAPPERS
│   ├── Base.templ         # Layout principal
│   └── Empty.templ        # Pour modals/fragments
│
└── logic/                 # 🎯 PURE GO LOGIC
    ├── styles.go          # Compute CSS classes
    ├── builders.go        # URL/Query builders
    ├── formatters.go      # Date/number/text
    ├── validators.go      # Input validation
    └── constants.go       # Design tokens Go

internal/views/tokens/     # 🎨 CSS DESIGN SYSTEM
└── retrowave.css          # Variables CSS centralisées
```


## 🎨 DESIGN SYSTEM RETROWAVE

### `internal/views/tokens/retrowave.css`

```css
/* ═══════════════════════════════════════════════════
   🌊 RETROWAVE TERMINAL DESIGN SYSTEM
   Inspired by: Tron, Blade Runner, VHS aesthetics
   ═══════════════════════════════════════════════════ */

:root {
  /* ═══ PRIMARY COLORS ═══ */
  --retro-cyan: #00E5FF;           /* Cyan néon principal */
  --retro-magenta: #FF10F0;        /* Magenta accent */
  --retro-purple: #BD00FF;         /* Violet vif */
  --retro-pink: #FF006E;           /* Rose fluo */
  
  /* ═══ BACKGROUNDS ═══ */
  --retro-void: #0A0015;           /* Noir spatial profond */
  --retro-dark: #1A0033;           /* Violet très sombre */
  --retro-surface: #2A0A4D;        /* Violet surface */
  --retro-elevated: #3D1566;       /* Violet élevé */
  
  /* ═══ TEXT COLORS ═══ */
  --retro-text-primary: #E0F7FF;   /* Cyan très clair */
  --retro-text-secondary: #A8C5D1; /* Cyan atténué */
  --retro-text-muted: #6B7C8A;     /* Gris bleuté */
  
  /* ═══ NEON GLOWS ═══ */
  --glow-cyan: 0 0 10px rgba(0, 229, 255, 0.5),
               0 0 20px rgba(0, 229, 255, 0.3),
               0 0 30px rgba(0, 229, 255, 0.1);
  
  --glow-magenta: 0 0 10px rgba(255, 16, 240, 0.5),
                  0 0 20px rgba(255, 16, 240, 0.3),
                  0 0 30px rgba(255, 16, 240, 0.1);
  
  --glow-purple: 0 0 10px rgba(189, 0, 255, 0.5),
                 0 0 20px rgba(189, 0, 255, 0.3);
  
  /* ═══ GRID CRT ═══ */
  --grid-color: rgba(0, 229, 255, 0.05);
  --scanline-color: rgba(255, 16, 240, 0.02);
  
  /* ═══ BORDERS ═══ */
  --border-neon: 1px solid var(--retro-cyan);
  --border-subtle: 1px solid rgba(0, 229, 255, 0.2);
  --border-muted: 1px solid rgba(0, 229, 255, 0.1);
  
  /* ═══ SPACING (8px base) ═══ */
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-6: 24px;
  --space-8: 32px;
  --space-12: 48px;
  
  /* ═══ TYPOGRAPHY ═══ */
  --font-mono: 'Berkeley Mono', 'Courier New', monospace;
  --font-display: 'Orbitron', 'Impact', sans-serif;
  
  --text-xs: 11px;
  --text-sm: 13px;
  --text-base: 15px;
  --text-lg: 18px;
  --text-xl: 24px;
  --text-2xl: 32px;
  
  /* ═══ ANIMATIONS ═══ */
  --transition-fast: 150ms cubic-bezier(0.4, 0, 0.2, 1);
  --transition-base: 250ms cubic-bezier(0.4, 0, 0.2, 1);
  --transition-slow: 400ms cubic-bezier(0.4, 0, 0.2, 1);
}

/* ═══════════════════════════════════════════════════
   🎬 GLOBAL EFFECTS
   ═══════════════════════════════════════════════════ */

body {
  background: var(--retro-void);
  color: var(--retro-text-primary);
  font-family: var(--font-mono);
  font-size: var(--text-base);
  position: relative;
  overflow-x: hidden;
}

/* Grid CRT effect */
body::before {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: 
    repeating-linear-gradient(
      0deg,
      var(--grid-color) 0px,
      transparent 1px,
      transparent 2px,
      var(--grid-color) 3px
    ),
    repeating-linear-gradient(
      90deg,
      var(--grid-color) 0px,
      transparent 1px,
      transparent 2px,
      var(--grid-color) 3px
    );
  background-size: 3px 3px;
  pointer-events: none;
  z-index: 9999;
  opacity: 0.3;
}

/* Scanlines effect */
body::after {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: repeating-linear-gradient(
    0deg,
    var(--scanline-color) 0px,
    transparent 1px,
    transparent 2px
  );
  pointer-events: none;
  z-index: 9998;
  animation: scanlines 8s linear infinite;
}

@keyframes scanlines {
  0% { transform: translateY(0); }
  100% { transform: translateY(100px); }
}

/* ═══════════════════════════════════════════════════
   🎨 UTILITY CLASSES
   ═══════════════════════════════════════════════════ */

.neon-text-cyan {
  color: var(--retro-cyan);
  text-shadow: var(--glow-cyan);
}

.neon-text-magenta {
  color: var(--retro-magenta);
  text-shadow: var(--glow-magenta);
}

.neon-border {
  border: var(--border-neon);
  box-shadow: var(--glow-cyan);
}

.neon-border-magenta {
  border: 1px solid var(--retro-magenta);
  box-shadow: var(--glow-magenta);
}

.retro-card {
  background: var(--retro-surface);
  border: var(--border-subtle);
  border-radius: 2px; /* Sharp edges retro style */
  position: relative;
}

.retro-card::before {
  content: '';
  position: absolute;
  inset: -1px;
  border-radius: 2px;
  padding: 1px;
  background: linear-gradient(135deg, var(--retro-cyan), var(--retro-magenta));
  -webkit-mask: 
    linear-gradient(#fff 0 0) content-box, 
    linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  opacity: 0;
  transition: opacity var(--transition-base);
}

.retro-card:hover::before {
  opacity: 0.5;
}

/* Terminal prompt effect */
.terminal-prompt::before {
  content: '>';
  color: var(--retro-cyan);
  margin-right: var(--space-2);
  text-shadow: var(--glow-cyan);
}

/* Blink cursor */
@keyframes blink {
  0%, 49% { opacity: 1; }
  50%, 100% { opacity: 0; }
}

.cursor-blink::after {
  content: '▮';
  color: var(--retro-cyan);
  animation: blink 1s infinite;
  margin-left: 2px;
}
```


## 📑 PATTERN CODE: ATOMIC COMPONENTS

### 1. `internal/views/ui/Button.templ`

```go
package ui

import (
    "maestro/internal/views/logic"
)

type ButtonVariant string

const (
    ButtonPrimary   ButtonVariant = "primary"
    ButtonSecondary ButtonVariant = "secondary"
    ButtonDanger    ButtonVariant = "danger"
    ButtonGhost     ButtonVariant = "ghost"
)

type ButtonSize string

const (
    ButtonSM ButtonSize = "sm"
    ButtonMD ButtonSize = "md"
    ButtonLG ButtonSize = "lg"
)

type ButtonProps struct {
    Label       string
    Icon        string
    Variant     ButtonVariant
    Size        ButtonSize
    HxPost      string
    HxGet       string
    HxTarget    string
    HxSwap      string
    HxIndicator string
    AriaLabel   string
    Disabled    bool
}

// Button - Bouton atomique universel
templ Button(props ButtonProps) {
    <button
        type="button"
        class={ logic.GetButtonClasses(props.Variant, props.Size) }
        if props.HxPost != "" {
            hx-post={ props.HxPost }
        }
        if props.HxGet != "" {
            hx-get={ props.HxGet }
        }
        if props.HxTarget != "" {
            hx-target={ props.HxTarget }
        }
        if props.HxSwap != "" {
            hx-swap={ props.HxSwap }
        } else {
            hx-swap="outerHTML"
        }
        if props.HxIndicator != "" {
            hx-indicator={ props.HxIndicator }
        }
        hx-on::before-request="this.disabled=true"
        hx-on::after-settle="this.disabled=false"
        aria-label={ props.AriaLabel }
        disabled?={ props.Disabled }>
        
        <span class="inline-flex items-center gap-2">
            if props.Icon != "" {
                <span class="neon-text-cyan" aria-hidden="true">{ props.Icon }</span>
            }
            <span class="terminal-prompt">{ props.Label }</span>
        </span>
        
        if props.HxIndicator != "" {
            <span class="htmx-indicator ml-2">
                @Spinner()
            </span>
        }
    </button>
}
```


### 2. `internal/views/ui/Badge.templ`

```go
package ui

import "maestro/internal/views/logic"

type BadgeVariant string

const (
    BadgeInfo    BadgeVariant = "info"
    BadgeSuccess BadgeVariant = "success"
    BadgeWarning BadgeVariant = "warning"
    BadgeDanger  BadgeVariant = "danger"
    BadgeCyan    BadgeVariant = "cyan"
    BadgeMagenta BadgeVariant = "magenta"
)

// Badge - Badge atomique avec glow néon
templ Badge(text string, variant BadgeVariant) {
    <span class={ logic.GetBadgeClasses(variant) }>
        <span class="badge-dot"></span>
        <span class="badge-text">{ text }</span>
    </span>
}

// BadgeCount - Badge compteur avec animation
templ BadgeCount(count int, variant BadgeVariant) {
    <span class={ logic.GetBadgeClasses(variant) }>
        <span class="badge-count" data-count={ fmt.Sprint(count) }>
            { fmt.Sprint(count) }
        </span>
    </span>
}
```


### 3. `internal/views/ui/Card.templ`

```go
package ui

// CardProps - Props universelles card
type CardProps struct {
    ID          string
    Title       string
    Subtitle    string
    HeaderIcon  string
    Interactive bool
    Href        string
    HxBoost     bool
}

// Card - Container universel avec neon border
templ Card(props CardProps, content templ.Component) {
    if props.Interactive && props.Href != "" {
        <a
            href={ templ.URL(props.Href) }
            if props.HxBoost {
                hx-boost="true"
            }
            id={ props.ID }
            class="retro-card block p-4 transition-all hover:translate-y-[-2px]">
            @CardContent(props, content)
        </a>
    } else {
        <div
            id={ props.ID }
            class="retro-card p-4">
            @CardContent(props, content)
        </div>
    }
}

templ CardContent(props CardProps, content templ.Component) {
    if props.Title != "" {
        <div class="card-header mb-4 pb-3 border-b border-retro-cyan/20">
            <div class="flex items-center gap-3">
                if props.HeaderIcon != "" {
                    <span class="neon-text-cyan text-xl">{ props.HeaderIcon }</span>
                }
                <div class="flex-1">
                    <h3 class="neon-text-cyan text-lg font-bold">{ props.Title }</h3>
                    if props.Subtitle != "" {
                        <p class="text-retro-text-secondary text-sm mt-1">{ props.Subtitle }</p>
                    }
                </div>
            </div>
        </div>
    }
    
    <div class="card-body">
        @content
    </div>
}
```


### 4. `internal/views/logic/styles.go`

```go
package logic

import (
    "fmt"
    "strings"
)

// ═══════════════════════════════════════════════════
// 🎨 CSS CLASSES COMPUTATION
// Centralise toute la logique de classes CSS
// ═══════════════════════════════════════════════════

// GetButtonClasses - Compute button classes
func GetButtonClasses(variant string, size string) string {
    base := []string{
        "inline-flex",
        "items-center",
        "justify-center",
        "gap-2",
        "font-mono",
        "font-semibold",
        "uppercase",
        "tracking-wider",
        "transition-all",
        "disabled:opacity-50",
        "disabled:cursor-not-allowed",
        "relative",
        "overflow-hidden",
    }
    
    // Size variants
    sizeClasses := map[string][]string{
        "sm": {"px-3", "py-1.5", "text-xs"},
        "md": {"px-4", "py-2", "text-sm"},
        "lg": {"px-6", "py-3", "text-base"},
    }
    
    // Variant styles
    variantClasses := map[string][]string{
        "primary": {
            "bg-retro-cyan/10",
            "border",
            "border-retro-cyan",
            "text-retro-cyan",
            "hover:bg-retro-cyan/20",
            "hover:shadow-[0_0_20px_rgba(0,229,255,0.5)]",
        },
        "secondary": {
            "bg-retro-surface",
            "border",
            "border-retro-cyan/30",
            "text-retro-text-primary",
            "hover:border-retro-cyan",
            "hover:shadow-[0_0_15px_rgba(0,229,255,0.3)]",
        },
        "danger": {
            "bg-retro-pink/10",
            "border",
            "border-retro-pink",
            "text-retro-pink",
            "hover:bg-retro-pink/20",
            "hover:shadow-[0_0_20px_rgba(255,0,110,0.5)]",
        },
        "ghost": {
            "text-retro-cyan",
            "hover:text-retro-magenta",
            "hover:bg-retro-surface",
        },
    }
    
    classes := append(base, sizeClasses[size]...)
    classes = append(classes, variantClasses[variant]...)
    
    return strings.Join(classes, " ")
}

// GetBadgeClasses - Compute badge classes
func GetBadgeClasses(variant string) string {
    base := []string{
        "inline-flex",
        "items-center",
        "gap-1.5",
        "px-2",
        "py-0.5",
        "rounded-sm",
        "text-xs",
        "font-mono",
        "uppercase",
        "tracking-wider",
        "border",
    }
    
    variantClasses := map[string][]string{
        "cyan": {
            "bg-retro-cyan/10",
            "border-retro-cyan",
            "text-retro-cyan",
            "shadow-[0_0_10px_rgba(0,229,255,0.3)]",
        },
        "magenta": {
            "bg-retro-magenta/10",
            "border-retro-magenta",
            "text-retro-magenta",
            "shadow-[0_0_10px_rgba(255,16,240,0.3)]",
        },
        "success": {
            "bg-retro-purple/10",
            "border-retro-purple",
            "text-retro-purple",
        },
        "warning": {
            "bg-retro-pink/10",
            "border-retro-pink",
            "text-retro-pink",
        },
        "info": {
            "bg-retro-surface",
            "border-retro-cyan/30",
            "text-retro-text-secondary",
        },
    }
    
    classes := append(base, variantClasses[variant]...)
    return strings.Join(classes, " ")
}

// GetCardClasses - Compute card classes based on state
func GetCardClasses(hasCount bool, isToday bool, isInteractive bool) string {
    base := []string{
        "retro-card",
        "p-4",
        "transition-all",
        "duration-250",
    }
    
    if isInteractive {
        base = append(base, "cursor-pointer", "hover:translate-y-[-2px]")
    }
    
    if isToday {
        base = append(base, "neon-border", "shadow-[0_0_30px_rgba(0,229,255,0.4)]")
    }
    
    if hasCount {
        base = append(base, "border-retro-purple")
    } else {
        base = append(base, "border-retro-cyan/20", "opacity-60")
    }
    
    return strings.Join(base, " ")
}

// GetProgressClasses - Compute progress bar classes
func GetProgressClasses(percent int) string {
    base := []string{
        "h-2",
        "rounded-none",
        "relative",
        "overflow-hidden",
        "transition-all",
        "duration-500",
    }
    
    // Gradient based on progress
    if percent < 30 {
        base = append(base, "bg-gradient-to-r", "from-retro-pink", "to-retro-magenta")
    } else if percent < 70 {
        base = append(base, "bg-gradient-to-r", "from-retro-magenta", "to-retro-purple")
    } else {
        base = append(base, "bg-gradient-to-r", "from-retro-cyan", "to-retro-purple")
    }
    
    return strings.Join(base, " ")
}
```


### 5. `internal/views/logic/builders.go`

```go
package logic

import (
    "fmt"
    "net/url"
    "strings"
)

// ═══════════════════════════════════════════════════
// 🔗 URL & QUERY BUILDERS
// Centralise construction URLs + query params
// ═══════════════════════════════════════════════════

// BuildFilterURL - Construit URL avec params filtre
func BuildFilterURL(baseURL string, filters map[string]string) string {
    if len(filters) == 0 {
        return baseURL
    }
    
    params := url.Values{}
    for key, value := range filters {
        if value != "" && value != "all" {
            params.Add(key, value)
        }
    }
    
    if len(params) == 0 {
        return baseURL
    }
    
    return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// BuildReviewURL - Construit URL review avec contexte session
func BuildReviewURL(exerciseID, quality int, fromSession bool, sessionID string) string {
    base := fmt.Sprintf("/exercise/%d/review?quality=%d", exerciseID, quality)
    
    if fromSession && sessionID != "" {
        return fmt.Sprintf("%s&session=%s", base, sessionID)
    }
    
    return base
}

// BuildPlannerURL - Construit URL planner avec date
func BuildPlannerURL(view string, date string) string {
    return fmt.Sprintf("/planner/%s?date=%s", view, date)
}

// SanitizeURL - Nettoie URL pour sécurité
func SanitizeURL(rawURL string) string {
    // Remove dangerous characters
    cleaned := strings.ReplaceAll(rawURL, "<", "")
    cleaned = strings.ReplaceAll(cleaned, ">", "")
    cleaned = strings.ReplaceAll(cleaned, "\"", "")
    cleaned = strings.ReplaceAll(cleaned, "'", "")
    
    return cleaned
}

// SanitizeID - Nettoie ID HTML
func SanitizeID(id string) string {
    // Keep only alphanumeric and dashes
    var result strings.Builder
    for _, r := range id {
        if (r >= 'a' && r <= 'z') || 
           (r >= 'A' && r <= 'Z') || 
           (r >= '0' && r <= '9') || 
           r == '-' || r == '_' {
            result.WriteRune(r)
        }
    }
    return result.String()
}
```


### 6. Composant Business Refactoré

```go
// internal/views/components/ExerciseCard.templ
package components

import (
    "fmt"
    "maestro/internal/models"
    "maestro/internal/views/ui"
    "maestro/internal/views/logic"
)

// ExerciseCard - Utilise UNIQUEMENT des atoms UI
templ ExerciseCard(ex models.Exercise) {
    @ui.Card(ui.CardProps{
        ID: fmt.Sprintf("exercise-%d", ex.ID),
        Interactive: true,
        Href: fmt.Sprintf("/exercise/%d", ex.ID),
        HxBoost: true,
    }, ExerciseCardContent(ex))
}

templ ExerciseCardContent(ex models.Exercise) {
    <!-- Header avec titre + status -->
    <div class="flex items-start justify-between gap-3 mb-3">
        <div class="flex-1">
            <h3 class="neon-text-cyan text-sm font-bold mb-1">
                { ex.Title }
            </h3>
            if ex.Description != "" {
                <p class="text-retro-text-secondary text-xs">
                    { ex.Description }
                </p>
            }
        </div>
        
        @ui.Badge(
            logic.FormatStatus(ex.Done),
            logic.GetStatusVariant(ex.Done),
        )
    </div>
    
    <!-- Meta badges -->
    <div class="flex items-center gap-2 mb-3">
        @ui.Badge(ex.Domain, "cyan")
        @ui.Badge(fmt.Sprintf("D%d", ex.Difficulty), "magenta")
    </div>
    
    <!-- Progress -->
    <div class="mb-3">
        @ui.Progress(len(ex.CompletedSteps), len(ex.Steps))
    </div>
    
    <!-- Actions -->
    <div class="flex items-center justify-between gap-2 pt-3 border-t border-retro-cyan/20">
        @ui.Button(ui.ButtonProps{
            Label: "Éditer",
            Icon: "✏",
            Variant: "ghost",
            Size: "sm",
            HxGet: fmt.Sprintf("/exercise/%d/edit", ex.ID),
            HxTarget: "#main-content",
            AriaLabel: "Éditer l'exercice",
        })
        
        @ui.Button(ui.ButtonProps{
            Label: "Suppr",
            Icon: "×",
            Variant: "danger",
            Size: "sm",
            HxPost: fmt.Sprintf("/exercise/%d/delete", ex.ID),
            HxTarget: "#exercise-list",
            HxSwap: "innerHTML",
            AriaLabel: "Supprimer l'exercice",
        })
    </div>
}
```


## 🚀 ACTION FLASH

### Phase 1: Setup Design System (2h)

```bash
# 1. Créer structure folders
mkdir -p internal/views/{ui,logic,tokens}
mkdir -p internal/views/components/{exercises,planner,session}

# 2. Créer CSS retrowave
cat > internal/views/tokens/retrowave.css << 'EOF'
[Le CSS retrowave complet ci-dessus]
EOF

# 3. Créer logic/constants.go
cat > internal/views/logic/constants.go << 'EOF'
package logic

const (
    ColorCyan    = "#00E5FF"
    ColorMagenta = "#FF10F0"
    ColorPurple  = "#BD00FF"
    ColorPink    = "#FF006E"
    
    ColorVoid     = "#0A0015"
    ColorDark     = "#1A0033"
    ColorSurface  = "#2A0A4D"
    ColorElevated = "#3D1566"
)
EOF
```


### Phase 2: Créer Atoms UI (3h)

```bash
# Créer les 8 atoms essentiels
touch internal/views/ui/{Button,Badge,Card,Input,Select,Progress,Icon,Spinner}.templ

# Implémenter chacun avec le pattern montré
# Button.templ → Pattern complet ci-dessus
# Badge.templ → Pattern complet ci-dessus  
# Card.templ → Pattern complet ci-dessus
```


### Phase 3: Refactorer Components (4h)

```bash
# Refactorer chaque composant pour utiliser UI atoms
# ExerciseCard.templ → Uses ui.Card + ui.Badge + ui.Button
# FilterBar.templ → Uses ui.Select + ui.Input
# ReviewPanel.templ → Uses ui.Card + ui.Button
# PlannerView.templ → Uses ui.Card + ui.Badge
```


### Phase 4: Migrer Logic (2h)

```bash
# Extraire toute logique CSS dans logic/
touch internal/views/logic/{styles,builders,formatters,validators}.go

# Déplacer fonctions de utils/ vers logic/
# GetButtonClasses() → logic/styles.go
# BuildFilterURL() → logic/builders.go
# FormatDate() → logic/formatters.go
```


## ✅ CHECKLIST MIGRATION

```
Setup:
[ ] Créer structure folders ui/logic/tokens
[ ] Créer retrowave.css avec palette complète
[ ] Mettre à jour Base.templ pour inclure retrowave.css

Atoms UI (8 composants):
[ ] ui/Button.templ avec variants + sizes
[ ] ui/Badge.templ avec glow effects
[ ] ui/Card.templ avec neon borders
[ ] ui/Input.templ avec validation inline
[ ] ui/Select.templ avec custom dropdown
[ ] ui/Progress.templ avec gradient
[ ] ui/Icon.templ avec SVG inline
[ ] ui/Spinner.templ avec animation

Logic Go:
[ ] logic/styles.go avec toutes les fonctions CSS
[ ] logic/builders.go avec URL helpers
[ ] logic/formatters.go avec date/number
[ ] logic/validators.go avec règles validation
[ ] logic/constants.go avec design tokens

Refactor Components (15 → 8):
[ ] ExerciseCard → Uses ui atoms
[ ] FilterBar (merge Filter + FilterPill)
[ ] ReviewPanel → Uses ui atoms
[ ] PlannerView (merge Day/Week/Month)
[ ] StepsManager → Uses ui atoms
[ ] StatCard → Uses ui atoms
[ ] SessionBuilder → Uses ui atoms
[ ] FormValidator → Uses ui atoms

Test:
[ ] Vérifier tous les hx-get/hx-post fonctionnent
[ ] Vérifier responsive design OK
[ ] Vérifier effets néon sur hover
[ ] Vérifier grid CRT visible
```

Veux-tu que je génère **le code complet** pour:

1. **ui/Button.templ** + **ui/Card.templ** + **ui/Badge.templ** ?
2. **logic/styles.go** + **logic/builders.go** complets ?
3. **ExerciseCard.templ refactoré** utilisant les atoms ?

Ou préfères-tu un **guide CLI étape par étape** pour migrer progressivement ?

