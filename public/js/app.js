// ============================================================
// MAESTRO - Ultra-Learning Frontend Logic
// ============================================================

const API_BASE = 'http://localhost:8080/api';

let appData = {
    exercises: [],
    userStats: {
        total_completed: 0,
        total_reviews: 0,
        domains: {}
    }
};

let currentView = 'dashboard';
let sessionTimer = null;
let sessionTimeRemaining = 0;

// ============================================================
// INITIALIZATION
// ============================================================

document.addEventListener('DOMContentLoaded', async () => {
    try {
        // Charger les données
        await loadExercises();
        await loadStats();
        
        // Cacher le loading, afficher le dashboard
        document.getElementById('loading-view').classList.remove('active');
        document.getElementById('dashboard-view').classList.add('active');
        
        // Mettre à jour l'interface
        updateDashboard();
        
        // Setup event listeners
        setupEventListeners();
    } catch (error) {
        console.error('Erreur lors du chargement:', error);
        alert('Erreur: ' + error.message);
    }
});

// ============================================================
// API CALLS
// ============================================================

async function loadExercises() {
    const response = await fetch(`${API_BASE}/exercises`);
    if (!response.ok) throw new Error('Impossible de charger les exercices');
    appData.exercises = await response.json() || [];
}

async function loadStats() {
    const response = await fetch(`${API_BASE}/stats`);
    if (!response.ok) throw new Error('Impossible de charger les stats');
    const stats = await response.json();
    appData.userStats = stats;
}

async function getRecommended() {
    const response = await fetch(`${API_BASE}/recommended`);
    if (!response.ok) throw new Error('Impossible de charger les recommandations');
    return await response.json();
}

async function rateExercise(exerciseId, rating) {
    const response = await fetch(`${API_BASE}/rate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ exercise_id: exerciseId, rating })
    });
    if (!response.ok) throw new Error('Impossible de sauvegarder la notation');
    return await response.json();
}

// ============================================================
// UI UPDATES
// ============================================================

function updateDashboard() {
    // Mettre à jour les stats
    document.getElementById('total-completed').textContent = appData.userStats.total_completed || 0;
    document.getElementById('total-reviews').textContent = appData.userStats.total_reviews || 0;
    
    const totalExercises = appData.exercises.length || 1;
    const progress = Math.round(((appData.userStats.total_completed || 0) / totalExercises) * 100);
    document.getElementById('overall-progress').textContent = progress + '%';
    document.getElementById('streak-days').textContent = '0'; // À implémenter
    
    // Charger les recommandés
    renderRecommended();
}

async function renderRecommended() {
    const recommended = await getRecommended();
    const container = document.getElementById('recommended-list');
    
    if (!recommended || recommended.length === 0) {
        container.innerHTML = '<p class="empty-state">Aucun exercice recommandé pour maintenant</p>';
        return;
    }
    
    container.innerHTML = recommended.map(ex => createExerciseCard(ex)).join('');
    
    // Ajouter les event listeners
    container.querySelectorAll('.exercise-card').forEach(card => {
        card.addEventListener('click', () => {
            const exId = card.dataset.id;
            const exercise = appData.exercises.find(e => e.id === exId);
            showExerciseView(exercise);
        });
    });
}

function createExerciseCard(ex) {
    const domainClass = getDomainClass(ex.domain);
    const diffClass = getDifficultyClass(ex.difficulty);
    const completed = ex.completed ? 'completed' : '';
    const isDue = isExerciseDue(ex);
    const dueClass = isDue && ex.completed ? 'due' : '';
    
    return `
        <div class="exercise-card ${completed} ${dueClass}" data-id="${ex.id}">
            <div class="exercise-header">
                <span class="exercise-domain ${ex.domain}">${formatDomain(ex.domain)}</span>
                <span class="exercise-difficulty ${diffClass}">D${ex.difficulty}</span>
            </div>
            <div class="exercise-title">${ex.title}</div>
            <div class="exercise-desc">${ex.description}</div>
            <div class="exercise-status">
                ${ex.completed ? '✓ Complété' : isDue ? '⚠ À réviser' : '○ Nouveau'}
            </div>
        </div>
    `;
}

function showExerciseView(exercise) {
    if (!exercise) return;
    
    const container = document.getElementById('exercise-content');
    document.getElementById('exercise-title-nav').textContent = exercise.title;
    
    const stepsHtml = exercise.steps.map((step, i) => `
        <div class="exercise-step">
            <input type="checkbox" id="step-${i}">
            <label for="step-${i}">${step}</label>
        </div>
    `).join('');
    
    container.innerHTML = `
        <section class="exercise-detail">
            <div class="exercise-meta">
                <span class="exercise-domain ${exercise.domain}">${formatDomain(exercise.domain)}</span>
                <span class="exercise-difficulty ${getDifficultyClass(exercise.difficulty)}">Difficulté ${exercise.difficulty}/3</span>
            </div>
            
            <h2>${exercise.title}</h2>
            <p class="exercise-description">${exercise.description}</p>
            
            <div class="exercise-steps">
                <h3>Étapes:</h3>
                ${stepsHtml}
            </div>
            
            <div class="exercise-content">
                <h3>Contenu:</h3>
                <pre><code>${escapeHtml(exercise.content)}</code></pre>
            </div>
            
            <div class="rating-section">
                <h3>Évaluation:</h3>
                <p>Comment trouvez-vous cet exercice ?</p>
                <div class="rating-buttons">
                    <button class="rating-btn" data-rating="1" title="Oublié">😟 Oublié</button>
                    <button class="rating-btn" data-rating="2" title="Difficile">😕 Difficile</button>
                    <button class="rating-btn" data-rating="3" title="Normal">😐 Normal</button>
                    <button class="rating-btn" data-rating="4" title="Facile">😊 Facile</button>
                </div>
            </div>
        </section>
    `;
    
    // Event listeners pour les ratings
    container.querySelectorAll('.rating-btn').forEach(btn => {
        btn.addEventListener('click', () => submitRating(exercise.id, parseInt(btn.dataset.rating)));
    });
    
    switchView('exercise-view');
}

async function submitRating(exerciseId, rating) {
    try {
        await rateExercise(exerciseId, rating);
        
        // Recharger les données et retourner au dashboard
        await loadExercises();
        await loadStats();
        updateDashboard();
        
        switchView('dashboard-view');
        alert('✅ Exercice évalué!');
    } catch (error) {
        alert('Erreur: ' + error.message);
    }
}

function renderStats() {
    document.getElementById('total-completed-stats').textContent = appData.userStats.total_completed || 0;
    document.getElementById('total-reviews-stats').textContent = appData.userStats.total_reviews || 0;
    
    // Domaines
    const domainContainer = document.getElementById('domain-breakdown');
    const domains = appData.userStats.domain_stats || {};
    
    domainContainer.innerHTML = Object.keys(domains).map(domain => {
        const stat = domains[domain];
        return `
            <div class="domain-stat">
                <div class="domain-stat-header">${formatDomain(domain)}</div>
                <div class="mastery-bar-container">
                    <div class="mastery-bar" style="width: ${(stat.mastery || 0)}%"></div>
                </div>
                <div class="domain-stat-details">
                    <span>Complétés: ${stat.completed}/${stat.total}</span>
                    <span>Maîtrise: ${stat.mastery || 0}%</span>
                </div>
            </div>
        `;
    }).join('');
}

function switchView(viewId) {
    // Masquer toutes les vues
    document.querySelectorAll('.view').forEach(v => v.classList.remove('active'));
    
    // Afficher la vue demandée
    document.getElementById(viewId).classList.add('active');
    
    currentView = viewId;
    
    // Mettre à jour les boutons de nav
    document.querySelectorAll('.nav-btn').forEach(btn => {
        btn.classList.remove('active');
        if (btn.dataset.view === viewId.replace('-view', '')) {
            btn.classList.add('active');
        }
    });
    
    // Actions spéciales
    if (viewId === 'stats-view') {
        renderStats();
    }
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

function formatDomain(domain) {
    const map = { golang: '🐹 Go', linux: '🐧 Linux', architecture: '🏗️ Architecture' };
    return map[domain] || domain;
}

function getDomainClass(domain) {
    return domain;
}

function getDifficultyClass(difficulty) {
    switch(difficulty) {
        case 1: return 'easy';
        case 2: return 'medium';
        case 3: return 'hard';
        default: return '';
    }
}

function isExerciseDue(exercise) {
    if (!exercise.last_reviewed) return false;
    const nextReview = new Date(exercise.last_reviewed);
    nextReview.setDate(nextReview.getDate() + exercise.interval_days);
    return new Date() > nextReview;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatTime(seconds) {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
}

// ============================================================
// EVENT LISTENERS
// ============================================================

function setupEventListeners() {
    // Navigation buttons
    document.querySelectorAll('.nav-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            const view = btn.dataset.view + '-view';
            switchView(view);
        });
    });
    
    // Back buttons
    document.querySelectorAll('.nav-back-btn').forEach(btn => {
        btn.addEventListener('click', () => switchView('dashboard-view'));
    });
    
    // Quick Start button
    document.getElementById('btn-quick-start')?.addEventListener('click', showQuickStart);
    
    // Browse button
    document.getElementById('btn-browse')?.addEventListener('click', () => {
        renderBrowse();
        switchView('browse-view');
    });
    
    // Skip session
    document.getElementById('btn-skip-session')?.addEventListener('click', () => {
        if (sessionTimer) clearInterval(sessionTimer);
        switchView('dashboard-view');
    });
}

async function showQuickStart() {
    const recommended = await getRecommended();
    const container = document.getElementById('quick-exercises');
    
    container.innerHTML = recommended.map((ex, i) => `
        <div class="quick-exercise-card" data-id="${ex.id}">
            <div class="quick-number">${i + 1}</div>
            <div>
                <div class="exercise-header" style="margin-bottom: 8px;">
                    <span class="exercise-domain ${ex.domain}">${formatDomain(ex.domain)}</span>
                    <span class="exercise-difficulty ${getDifficultyClass(ex.difficulty)}">D${ex.difficulty}</span>
                </div>
                <div class="exercise-title">${ex.title}</div>
                <div class="exercise-desc">${ex.description}</div>
            </div>
        </div>
    `).join('');
    
    // Ajouter les listeners
    container.querySelectorAll('.quick-exercise-card').forEach(card => {
        card.addEventListener('click', () => {
            const exId = card.dataset.id;
            const ex = appData.exercises.find(e => e.id === exId);
            showExerciseView(ex);
        });
    });
    
    // Démarrer le timer
    sessionTimeRemaining = 15 * 60;
    updateQuickTimer();
    
    if (sessionTimer) clearInterval(sessionTimer);
    sessionTimer = setInterval(() => {
        sessionTimeRemaining--;
        updateQuickTimer();
        if (sessionTimeRemaining <= 0) {
            clearInterval(sessionTimer);
        }
    }, 1000);
    
    switchView('quick-start-view');
}

function updateQuickTimer() {
    const timerEl = document.getElementById('quick-timer');
    if (timerEl) {
        timerEl.textContent = formatTime(sessionTimeRemaining);
    }
}

async function renderBrowse() {
    const container = document.getElementById('exercises-list');
    
    container.innerHTML = appData.exercises.map(ex => createExerciseCard(ex)).join('');
    
    // Ajouter les listeners
    container.querySelectorAll('.exercise-card').forEach(card => {
        card.addEventListener('click', () => {
            const exId = card.dataset.id;
            const exercise = appData.exercises.find(e => e.id === exId);
            showExerciseView(exercise);
        });
    });
    
    // Setup filters
    const searchInput = document.getElementById('search-input');
    const domainFilter = document.getElementById('domain-filter');
    
    [searchInput, domainFilter].forEach(el => {
        if (el) el.addEventListener('change', filterExercises);
    });
}

function filterExercises() {
    const search = document.getElementById('search-input')?.value.toLowerCase() || '';
    const domain = document.getElementById('domain-filter')?.value || '';
    
    const filtered = appData.exercises.filter(ex => {
        const matchSearch = ex.title.toLowerCase().includes(search) || ex.description.toLowerCase().includes(search);
        const matchDomain = !domain || ex.domain === domain;
        return matchSearch && matchDomain;
    });
    
    const container = document.getElementById('exercises-list');
    container.innerHTML = filtered.map(ex => createExerciseCard(ex)).join('');
    
    // Ré-ajouter les listeners
    container.querySelectorAll('.exercise-card').forEach(card => {
        card.addEventListener('click', () => {
            const exId = card.dataset.id;
            const exercise = appData.exercises.find(e => e.id === exId);
            showExerciseView(exercise);
        });
    });
}
