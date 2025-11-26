<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

## 📸 Visualisation Mentale

Tu as **3 jours** pour implémenter **Phase 1 : Sessions Adaptatives**.

C'est faisable avec une approche **MVP minimal** → Tests → Itération.

***

## 🎯 PLAN 3 JOURS : Session Adaptatives

### **JOUR 1 : Structure de Données + Routes**

**Matin (3-4h) :**

```go
// 1. Ajouter modèles
// internal/models/session.go (NEW)

type EnergyLevel int
const (
    Low    EnergyLevel = 1 // 5-15min
    Medium EnergyLevel = 2 // 20-30min
    High   EnergyLevel = 3 // 45-60min
)

type SessionMode string
const (
    MicroSession    SessionMode = "micro"    // 10min
    StandardSession SessionMode = "standard" // 25min
    DeepSession     SessionMode = "deep"     // 50min
)

type AdaptiveSession struct {
    ID               string            `json:"id"`
    UserID           string            `json:"user_id,omitempty"`
    Mode             SessionMode       `json:"mode"`
    EnergyLevel      EnergyLevel       `json:"energy_level"`
    EstimatedTime    time.Duration     `json:"estimated_time"`
    Exercises        []int             `json:"exercises"`     // IDs
    BreakSchedule    []time.Duration   `json:"break_schedule"`
    StartedAt        time.Time         `json:"started_at"`
    Status           string            `json:"status"` // "pending", "active", "completed"
}
```

**Après-midi (2-3h) :**

```go
// 2. Handler builder
// internal/handlers/sessions.go (NEW)

func HandleSessionBuilder(w http.ResponseWriter, r *http.Request) {
    // GET /session/builder
    // Affiche : Sélection énergie + Prévisualisation
    Tmpl.ExecuteTemplate(w, "session-builder", nil)
}

func HandleStartSession(w http.ResponseWriter, r *http.Request) {
    // POST /session/start?energy=medium
    energy := r.URL.Query().Get("energy")
    
    // Parser energy
    var energyLevel EnergyLevel
    switch energy {
    case "low":
        energyLevel = Low
    case "medium":
        energyLevel = Medium
    case "high":
        energyLevel = High
    default:
        http.Error(w, "Invalid energy level", http.StatusBadRequest)
        return
    }
    
    // Construire session
    session := BuildAdaptiveSession(energyLevel)
    
    // Stocker en mémoire (ou JSON file pour MVP)
    sessions[session.ID] = session
    
    // Rediriger
    http.Redirect(w, r, "/session/"+session.ID, http.StatusSeeOther)
}

func HandleCurrentSession(w http.ResponseWriter, r *http.Request) {
    // GET /session/{id}
    sessionID := r.PathValue("id")
    session, exists := sessions[sessionID]
    if !exists {
        http.NotFound(w, r)
        return
    }
    
    data := map[string]any{
        "Session":  session,
        "Exercise": getExerciseByID(session.Exercises[0]),
    }
    
    Tmpl.ExecuteTemplate(w, "session-current", data)
}
```

**Soir (1h) :**

```go
// 3. Logique BuildAdaptiveSession
// internal/store/sessions.go (NEW)

func BuildAdaptiveSession(energy EnergyLevel) *AdaptiveSession {
    session := &AdaptiveSession{
        ID:            generateSessionID(),
        EnergyLevel:   energy,
        StartedAt:     time.Now(),
        Status:        "pending",
    }
    
    switch energy {
    case Low:
        session.Mode = MicroSession
        session.EstimatedTime = 10 * time.Minute
        session.Exercises = pickDueExercises(1)  // 1 exo urgent
        session.BreakSchedule = []time.Duration{} // Aucune pause
        
    case Medium:
        session.Mode = StandardSession
        session.EstimatedTime = 25 * time.Minute
        session.Exercises = pickDueExercises(2)  // 2 exos liés
        session.BreakSchedule = []time.Duration{
            25 * time.Minute,  // Pause après
        }
        
    case High:
        session.Mode = DeepSession
        session.EstimatedTime = 50 * time.Minute
        session.Exercises = pickDueExercises(3)  // 3 exos progressifs
        session.BreakSchedule = []time.Duration{
            25 * time.Minute,  // Pause 1
            50 * time.Minute,  // Pause 2
        }
    }
    
    return session
}

func pickDueExercises(count int) []int {
    // Sélectionne les N exercices les plus urgents
    // (NextReviewAt < now, ou mieux matching)
    // Pour MVP : sélection simple aléatoire parmi Done
    
    var result []int
    for _, ex := range store.GetAll() {
        if ex.Done && len(result) < count {
            result = append(result, ex.ID)
        }
    }
    return result
}
```


***

### **JOUR 2 : Templates + UI**

**Matin (4-5h) :**

```html
<!-- templates/session-builder.html (NEW) -->
{{define "session-builder"}}
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Nouvelle Session - Maestro</title>
    <link rel="stylesheet" href="/public/css/list-exercice.css">
    <style>
        .energy-selector {
            max-width: 600px;
            margin: 60px auto;
        }
        
        .energy-option {
            display: flex;
            align-items: center;
            padding: 20px;
            margin: 15px 0;
            border: 2px solid var(--border-terminal);
            background: #0d0d0d;
            cursor: pointer;
            transition: all 0.2s;
        }
        
        .energy-option:hover {
            border-color: var(--fg-terminal);
            box-shadow: 0 0 10px var(--shadow-glow);
        }
        
        .energy-option input[type="radio"] {
            width: 24px;
            height: 24px;
            margin-right: 20px;
            cursor: pointer;
        }
        
        .energy-option.low {
            --color: #ff6600;
        }
        
        .energy-option.medium {
            --color: #ffff00;
        }
        
        .energy-option.high {
            --color: #66ff66;
        }
        
        .energy-option .icon {
            font-size: 28px;
            margin-right: 20px;
            min-width: 40px;
        }
        
        .energy-info {
            flex: 1;
        }
        
        .energy-title {
            font-size: 14px;
            font-weight: bold;
            text-transform: uppercase;
            color: var(--color, var(--fg-terminal));
            margin-bottom: 5px;
        }
        
        .energy-subtitle {
            font-size: 11px;
            opacity: 0.7;
        }
        
        .preview-box {
            border: 1px solid var(--border-terminal);
            padding: 15px;
            background: #111;
            margin: 30px 0;
            min-height: 120px;
        }
        
        .preview-title {
            font-size: 12px;
            text-transform: uppercase;
            margin-bottom: 10px;
            color: var(--fg-terminal);
        }
        
        .preview-content {
            font-size: 11px;
            opacity: 0.8;
            line-height: 1.8;
        }
        
        .btn-start {
            display: block;
            width: 100%;
            padding: 15px;
            margin-top: 30px;
            background: #0a0a0a;
            color: var(--fg-terminal);
            border: 2px solid var(--fg-terminal);
            font-family: var(--mono-font);
            font-size: 14px;
            text-transform: uppercase;
            cursor: pointer;
            transition: all 0.2s;
        }
        
        .btn-start:hover {
            background: var(--fg-terminal);
            color: #0a0a0a;
            box-shadow: 0 0 20px var(--shadow-glow);
        }
    </style>
</head>
<body>
    <div class="scan-line"></div>
    
    <nav>
        <a href="/">[HOME]</a>
        <a href="/exercises">[EXERCISES]</a>
        <a href="/session/builder">[← BACK]</a>
    </nav>
    
    <div class="terminal-header">
        <div class="terminal-title">
            ╔═══════════════════════════════════╗<br>
            ║   NOUVELLE SESSION                ║<br>
            ╚═══════════════════════════════════╝
        </div>
        <div class="terminal-subtitle">
            &gt; Comment te sens-tu maintenant ? | STATUS: <span class="blink">ONLINE</span>
        </div>
    </div>
    
    <div style="max-width: 700px; margin: 0 auto; padding: 20px;">
        <form method="POST" action="/session/start" class="energy-selector">
            <!-- Option 1 : Faible -->
            <label class="energy-option low">
                <input type="radio" name="energy" value="low" required>
                <div class="icon">🔋</div>
                <div class="energy-info">
                    <div class="energy-title">Faible - 5-15min</div>
                    <div class="energy-subtitle">Juste me réveiller</div>
                </div>
            </label>
            
            <!-- Option 2 : Moyenne -->
            <label class="energy-option medium">
                <input type="radio" name="energy" value="medium" checked>
                <div class="icon">🔋🔋</div>
                <div class="energy-info">
                    <div class="energy-title">Moyenne - 20-30min</div>
                    <div class="energy-subtitle">Concentré mais pas marathon</div>
                </div>
            </label>
            
            <!-- Option 3 : Haute -->
            <label class="energy-option high">
                <input type="radio" name="energy" value="high" required>
                <div class="icon">🔋🔋🔋</div>
                <div class="energy-info">
                    <div class="energy-title">Haute - 45-60min</div>
                    <div class="energy-subtitle">Je suis en feu, allons-y !</div>
                </div>
            </label>
            
            <!-- Prévisualisation -->
            <div class="preview-box">
                <div class="preview-title">▶ Session proposée :</div>
                <div class="preview-content" id="preview">
                    <div>• Mode : Micro (10min)</div>
                    <div>• Exercices : 1</div>
                    <div>• Pauses : 0</div>
                </div>
            </div>
            
            <button type="submit" class="btn-start">
                [▶ DÉMARRER MAINTENANT]
            </button>
        </form>
    </div>
    
    <div class="terminal-footer">
        ┌────────────────────────────────────┐<br>
        │ © 2025 GO v2 TERMINAL | ALL SYSTEMS NOMINAL │<br>
        └────────────────────────────────────┘
    </div>
    
    <script>
        // Update preview au changement d'énergie
        document.querySelectorAll('input[name="energy"]').forEach(radio => {
            radio.addEventListener('change', function() {
                const preview = document.getElementById('preview');
                
                switch(this.value) {
                case 'low':
                    preview.innerHTML = `
                        <div>• Mode : Micro (10min)</div>
                        <div>• Exercices : 1</div>
                        <div>• Pauses : 0</div>
                    `;
                    break;
                case 'medium':
                    preview.innerHTML = `
                        <div>• Mode : Standard (25min)</div>
                        <div>• Exercices : 2</div>
                        <div>• Pauses : 1 (5min après ex.1)</div>
                    `;
                    break;
                case 'high':
                    preview.innerHTML = `
                        <div>• Mode : Deep (50min)</div>
                        <div>• Exercices : 3</div>
                        <div>• Pauses : 2 (25min, 50min)</div>
                    `;
                    break;
                }
            });
        });
    </script>
</body>
</html>
{{end}}
```

**Après-midi (3-4h) :**

```html
<!-- templates/session-current.html (NEW) -->
{{define "session-current"}}
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Session en Cours - Maestro</title>
    <link rel="stylesheet" href="/public/css/list-exercice.css">
    <style>
        .session-container {
            max-width: 900px;
            margin: 0 auto;
            padding: 20px;
        }
        
        .session-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
            padding: 15px;
            border: 1px solid var(--border-terminal);
            background: #0d0d0d;
        }
        
        .session-progress {
            flex: 1;
        }
        
        .progress-bar {
            width: 100%;
            height: 20px;
            border: 1px solid var(--border-terminal);
            background: #111;
            overflow: hidden;
            margin: 10px 0;
        }
        
        .progress-fill {
            height: 100%;
            background: repeating-linear-gradient(
                90deg,
                var(--fg-terminal) 0px,
                var(--fg-terminal) 2px,
                transparent 2px,
                transparent 4px
            );
            width: 20%;
            transition: width 0.3s ease;
        }
        
        .session-timer {
            font-size: 28px;
            font-weight: bold;
            color: var(--fg-terminal);
            text-shadow: 0 0 10px var(--shadow-glow);
            text-align: right;
        }
        
        .exercise-content {
            border: 2px solid var(--border-terminal);
            padding: 30px;
            background: #0d0d0d;
            margin-bottom: 20px;
            box-shadow: 0 0 30px var(--shadow-glow);
        }
        
        .session-controls {
            display: flex;
            gap: 10px;
            justify-content: center;
        }
        
        .btn-control {
            padding: 12px 20px;
            border: 1px solid var(--border-terminal);
            background: #0a0a0a;
            color: var(--fg-terminal);
            font-family: var(--mono-font);
            font-size: 12px;
            text-transform: uppercase;
            cursor: pointer;
            transition: all 0.2s;
        }
        
        .btn-control:hover {
            background: #1a1a1a;
            box-shadow: 0 0 10px var(--shadow-glow);
        }
    </style>
</head>
<body>
    <div class="scan-line"></div>
    
    <nav>
        <a href="/">[HOME]</a>
        <a href="/exercises">[EXERCISES]</a>
        <a href="/session/builder">[← NEW SESSION]</a>
    </nav>
    
    <div class="terminal-header">
        <div class="terminal-title">
            ╔═══════════════════════════════════╗<br>
            ║   SESSION {{.Session.Mode | upper}}  ║<br>
            ╚═══════════════════════════════════╝
        </div>
        <div class="terminal-subtitle">
            &gt; {{len .Session.Exercises}}/{{len .Session.Exercises}} exercices | STATUS: <span class="blink">ACTIVE</span>
        </div>
    </div>
    
    <div class="session-container">
        <!-- Progress Header -->
        <div class="session-header">
            <div class="session-progress">
                <div style="font-size: 12px; text-transform: uppercase; margin-bottom: 5px;">
                    Progression
                </div>
                <div class="progress-bar">
                    <div class="progress-fill" id="progressFill"></div>
                </div>
                <div style="font-size: 10px; opacity: 0.7;">
                    1/{{len .Session.Exercises}} exercices
                </div>
            </div>
            <div class="session-timer" id="sessionTimer">
                {{.Session.EstimatedTime.Minutes | int}}:00
            </div>
        </div>
        
        <!-- Exercise Content -->
        <div class="exercise-content">
            <!-- Insert exercise detail here -->
            {{template "exercise-detail" .Exercise}}
        </div>
        
        <!-- Controls -->
        <div class="session-controls">
            <button class="btn-control" onclick="pauseSession()">
                [⏸ PAUSE]
            </button>
            <button class="btn-control" onclick="completeExercise()">
                [✓ TERMINÉ]
            </button>
            <button class="btn-control" onclick="stopSession()">
                [⏹ STOP SESSION]
            </button>
        </div>
    </div>
    
    <script>
        // Timer simulation (MVP - pas de WebSocket)
        let secondsRemaining = {{.Session.EstimatedTime.Seconds | int}};
        
        setInterval(() => {
            if (secondsRemaining > 0) {
                secondsRemaining--;
                const mins = Math.floor(secondsRemaining / 60);
                const secs = secondsRemaining % 60;
                document.getElementById('sessionTimer').textContent = 
                    `${mins}:${secs.toString().padStart(2, '0')}`;
            }
        }, 1000);
        
        function pauseSession() {
            alert('⏸ Session en pause');
            // TODO: Implémenter pause
        }
        
        function completeExercise() {
            // Marquer exercice comme complété
            window.location.href = '/session/{{.Session.ID}}/complete';
        }
        
        function stopSession() {
            if (confirm('⏹ Arrêter la session ?')) {
                window.location.href = '/session/{{.Session.ID}}/stop';
            }
        }
    </script>
</body>
</html>
{{end}}
```


***

### **JOUR 3 : Intégration + Tests**

**Matin (2-3h) :**

```go
// Ajouter routes dans handlers/routes.go

mux.HandleFunc("GET /session/builder", HandleSessionBuilder)
mux.HandleFunc("POST /session/start", HandleStartSession)
mux.HandleFunc("GET /session/{id}", HandleCurrentSession)
mux.HandleFunc("POST /session/{id}/complete", HandleCompleteSession)
mux.HandleFunc("POST /session/{id}/stop", HandleStopSession)

// Modifier dashboard pour inclure bouton Session

{{define "dashboard"}}
<div style="max-width: 900px; margin: 0 auto; padding: 20px;">
    <div class="terminal-header">
        <div class="terminal-title">
            ╔═══════════════════════════════════╗<br>
            ║   MAESTRO v2.0 DASHBOARD          ║<br>
            ╚═══════════════════════════════════╝
        </div>
    </div>
    
    <!-- PRIMARY ACTION : Session -->
    <div style="border: 2px solid #00ff00; padding: 30px; 
                background: #0d0d0d; text-align: center; 
                margin-bottom: 20px;">
        <h2 style="font-size: 18px; margin-bottom: 15px;">
            🎯 Prêt à apprendre ?
        </h2>
        <a href="/session/builder" 
           style="display: inline-block; padding: 15px 40px;
                  background: #0a0a0a; border: 2px solid #00ff00;
                  color: #00ff00; text-decoration: none;
                  font-family: var(--mono-font); font-size: 14px;
                  text-transform: uppercase; transition: all 0.2s;"
           onmouseover="this.style.background='#00ff00'; this.style.color='#0a0a0a';"
           onmouseout="this.style.background='#0a0a0a'; this.style.color='#00ff00';">
            [DÉMARRER SESSION]
        </a>
    </div>
    
    <!-- SECONDARY ACTIONS -->
    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 15px;">
        <a href="/exercises" style="padding: 15px; border: 1px solid var(--border-terminal); 
                                   background: #0d0d0d; text-decoration: none; 
                                   text-align: center; color: #00ff00;">
            [📚 Explorer Exercices]
        </a>
        <a href="/stats" style="padding: 15px; border: 1px solid var(--border-terminal); 
                               background: #0d0d0d; text-decoration: none; 
                               text-align: center; color: #00ff00;">
            [📊 Statistiques]
        </a>
    </div>
</div>
{{end}}
```

**Après-midi (2-3h) :**

```
TESTS MANUELS
─────────────

1. Ouvre http://localhost:8080
   ✓ Vois dashboard avec [DÉMARRER SESSION]

2. Clique [DÉMARRER SESSION]
   ✓ Affiche session-builder.html
   ✓ 3 radio buttons (Faible/Moyenne/Haute)
   ✓ Preview en temps réel change

3. Sélectionne "Moyenne", clique [▶ DÉMARRER MAINTENANT]
   ✓ POST /session/start?energy=medium
   ✓ Redirige vers /session/{id}

4. Affiche session-current.html
   ✓ Timer compte à rebours
   ✓ Progress bar visible
   ✓ Exercice affichée

5. Clique [✓ TERMINÉ]
   ✓ Session marquée complétée
   ✓ Exercice dans SRS

6. Clique [⏹ STOP SESSION]
   ✓ Confirmation
   ✓ Retour dashboard
```

**Soir (1h) :**

```
COMMIT CODE
───────────
git add .
git commit -m "feat: Phase 1 - Sessions adaptatives

- Modèles : EnergyLevel, SessionMode, AdaptiveSession
- Routes : /session/builder, /session/start, /session/{id}
- Templates : session-builder, session-current
- Logique : BuildAdaptiveSession avec 3 modes
- Dashboard : Nouveau bouton primaire [DÉMARRER SESSION]"
```


***

## 🎯 RÉSUMÉ 3 JOURS

| Jour | Durée | Tâche | Résultat |
| :-- | :-- | :-- | :-- |
| **J1** | 7h | Modèles + Handlers + Logique | Code backend 90% |
| **J2** | 7-8h | Templates + CSS + JS preview | UI frontend complète |
| **J3** | 5-6h | Intégration + Tests + Commit | Feature prête à tester |

**Total : ~20h de dev = 3 jours (6-7h/jour)**

***

## 🚀 APRÈS J3 : Phase 2 (Optionnel - J4)

```
J4 : ASCII Visuals + Testing

- Créer 3 exercices avec ASCII diagrams
- Tester avec vraies données
- Ajuster timings si nécessaire
```


***

## 🎯 CHECK-LIST JOUR 1

- [ ] `internal/models/session.go` créé
- [ ] `internal/handlers/sessions.go` créé
- [ ] `internal/store/sessions.go` créé
- [ ] Routes ajoutées
- [ ] Code compile sans erreur


## 🎯 CHECK-LIST JOUR 2

- [ ] `templates/session-builder.html` créé
- [ ] `templates/session-current.html` créé
- [ ] CSS applié
- [ ] Preview JS fonctionne
- [ ] Dashboard modifié


## 🎯 CHECK-LIST JOUR 3

- [ ] Tests manuels passent
- [ ] Aucun bug visible
- [ ] Code committé
- [ ] Prêt pour Phase 2

***

Tu veux que je te code **Jour 1 complet** en bloc ? Ou tu préfères aller step by step ? 🚀

