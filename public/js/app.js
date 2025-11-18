// ============= API CLIENT WITH NEW ENDPOINTS =============
class MaestroAPI {
    constructor(baseURL = 'http://localhost:8080') {
        this.baseURL = baseURL;
    }

    async request(endpoint, options = {}) {
        try {
            const response = await fetch(`${this.baseURL}${endpoint}`, {
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers
                },
                ...options
            });

            if (!response.ok) throw new Error(`API Error: ${response.statusText}`);
            return await response.json();
        } catch (error) {
            console.error('API Request failed:', error);
            throw error;
        }
    }

    getExercises(page = 1, pageSize = 10) {
        return this.request(`/api/exercises?page=${page}&page_size=${pageSize}`);
    }

    rateExercise(exerciseId, rating) {
        return this.request('/api/rate', {
            method: 'POST',
            body: JSON.stringify({
                exercise_id: exerciseId,
                rating: rating
            })
        });
    }

    getStats() {
        return this.request('/api/stats');
    }

    // Advanced stats endpoints
    getAdvancedStats() {
        return this.request('/api/stats/advanced');
    }

    getStruggling() {
        return this.request('/api/stats/struggling');
    }

    getMastered() {
        return this.request('/api/stats/mastered');
    }

    getNeedsPractice() {
        return this.request('/api/stats/needs-practice');
    }

    getInsights() {
        return this.request('/api/stats/insights');
    }

    getDomainAnalysis(domain) {
        return this.request(`/api/stats/domain/${domain}`);
    }

    // Calendar endpoints
    getCalendarUpcoming(days = 30) {
        return this.request(`/api/calendar/upcoming?days=${days}`);
    }

    getCalendarWeek() {
        return this.request('/api/calendar/week');
    }

    getCalendarMonth() {
        return this.request('/api/calendar/month');
    }
}

// ============= APP STATE WITH ADVANCED STATS =============
class MaestroApp {
    constructor() {
        this.api = new MaestroAPI();
        this.allExercises = [];
        this.filteredExercises = [];
        this.selectedExercise = null;
        this.currentPage = 'exercises';
        this.currentStatus = 'all';
        this.currentDomain = 'all';
        this.pageNum = 1;
        this.pageSize = 10;
        this.stats = {};
        this.advancedStats = {};
        this.reviewDates = {};
        this.completedSteps = {};
        this.init();
    }

    async init() {
        await this.loadExercises();
        await this.updateStats();
        this.setupEventListeners();
        this.render();
    }

    async loadExercises() {
        try {
            const response = await this.api.getExercises(1, 1000);
            this.allExercises = response.exercises || [];
            this.reviewDates = response.review_dates || {};
        } catch (error) {
            console.error('Erreur chargement:', error);
        }
    }

    async updateStats() {
        try {
            this.stats = await this.api.getStats();

            // Load advanced stats if available
            try {
                this.advancedStats = await this.api.getAdvancedStats();
            } catch (e) {
                console.warn('Advanced stats endpoint not available yet');
            }

            this.renderNavStats();
            this.renderStatsPage();
        } catch (error) {
            console.error('Erreur stats:', error);
        }
    }

    setupEventListeners() {
        document.querySelectorAll('[data-status]').forEach(btn => {
            btn.addEventListener('click', (e) => {
                document.querySelectorAll('[data-status]').forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
                this.currentStatus = e.target.dataset.status;
                this.pageNum = 1;
                this.applyFilters();
            });
        });

        document.querySelectorAll('[data-domain]').forEach(btn => {
            btn.addEventListener('click', (e) => {
                document.querySelectorAll('[data-domain]').forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
                this.currentDomain = e.target.dataset.domain;
                this.pageNum = 1;
                this.applyFilters();
            });
        });
    }

    goToPage(page) {
        document.querySelectorAll('.nav-link').forEach(btn => btn.classList.remove('active'));
        event.target.classList.add('active');

        document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
        document.getElementById(page).classList.add('active');

        this.currentPage = page;

        if (page === 'stats') {
            this.updateStats();
        } else if (page === 'calendar') {
            this.renderCalendar();
        }
    }

    switchStatsTab(tab) {
        this.currentStatsTab = tab;

        document.querySelectorAll('.stats-tab').forEach(t => t.classList.remove('active'));
        document.querySelectorAll('.stats-content').forEach(c => c.classList.remove('active'));

        event.target.classList.add('active');
        document.getElementById(tab + '-stats').classList.add('active');
    }

    applyFilters() {
        let filtered = this.allExercises;

        if (this.currentStatus === 'pending') {
            filtered = filtered.filter(ex => !ex.completed);
        } else if (this.currentStatus === 'completed') {
            filtered = filtered.filter(ex => ex.completed);
        }

        if (this.currentDomain !== 'all') {
            filtered = filtered.filter(ex => ex.domain === this.currentDomain);
        }

        this.filteredExercises = filtered;
        this.renderExerciseList();
    }

    getPaginatedExercises() {
        const start = (this.pageNum - 1) * this.pageSize;
        return this.filteredExercises.slice(start, start + this.pageSize);
    }

    getTotalPages() {
        return Math.ceil(this.filteredExercises.length / this.pageSize);
    }

    nextPage() {
        if (this.pageNum < this.getTotalPages()) {
            this.pageNum++;
            this.renderExerciseList();
        }
    }

    previousPage() {
        if (this.pageNum > 1) {
            this.pageNum--;
            this.renderExerciseList();
        }
    }

    renderExerciseList() {
        const list = document.getElementById('exerciseList');
        const paginated = this.getPaginatedExercises();

        if (paginated.length === 0) {
            list.innerHTML = '<div style="text-align: center; color: #888; padding: 40px;">Aucun exercice trouvé</div>';
            document.getElementById('pageInfo').textContent = '';
            return;
        }

        list.innerHTML = paginated.map(ex => {
            const isSelected = this.selectedExercise?.id === ex.id;
            const nextReviewText = this.reviewDates[ex.id] || 'Nouveau';

            return `
                <div class="exercise-item ${isSelected ? 'selected' : ''}" data-id="${ex.id}">
                    <div class="exercise-title">${ex.completed ? '✓' : '○'} ${ex.title}</div>
                    <div class="exercise-meta">
                        <span style="background: #5D4E60; padding: 2px 8px; border-radius: 3px;">${ex.domain}</span>
                        <span style="background: #2d4a6e; padding: 2px 8px; border-radius: 3px;">D${ex.difficulty}</span>
                        <span style="color: #90EE90;">⏱ ${nextReviewText}</span>
                    </div>
                </div>
            `;
        }).join('');

        document.querySelectorAll('.exercise-item').forEach(item => {
            item.addEventListener('click', () => {
                const id = item.dataset.id;
                this.selectedExercise = this.allExercises.find(ex => ex.id === id);
                this.renderExerciseDetail();
                this.renderExerciseList();
            });
        });

        const total = this.getTotalPages();
        document.getElementById('pageInfo').textContent = `Page ${this.pageNum}/${total} (${this.filteredExercises.length} total)`;
        document.getElementById('prevBtn').disabled = this.pageNum <= 1;
        document.getElementById('nextBtn').disabled = this.pageNum >= total;
    }

    toggleStep(exerciseId, stepIndex) {
        if (!this.completedSteps[exerciseId]) {
            this.completedSteps[exerciseId] = [];
        }

        const idx = this.completedSteps[exerciseId].indexOf(stepIndex);
        if (idx > -1) {
            this.completedSteps[exerciseId].splice(idx, 1);
        } else {
            this.completedSteps[exerciseId].push(stepIndex);
        }

        this.renderExerciseDetail();
    }

    toggleCompletion() {
        if (!this.selectedExercise) return;

        const ex = this.selectedExercise;
        const newStatus = ex.completed ? 2 : 3;
        this.rate(ex.id, newStatus);
    }

    renderExerciseDetail() {
        const detail = document.getElementById('exerciseDetail');

        if (!this.selectedExercise) {
            detail.innerHTML = '<div style="text-align: center; color: #888;">Sélectionnez un exercice</div>';
            return;
        }

        const ex = this.selectedExercise;
        const nextReviewText = this.reviewDates[ex.id] || 'Nouveau';
        const stepsCompleted = this.completedSteps[ex.id] || [];
        const stepsProgress = `${stepsCompleted.length}/${ex.steps.length}`;

        detail.innerHTML = `
            <div class="detail-header">
                <h3>${ex.title}</h3>
                <p class="detail-description">${ex.description}</p>
            </div>

            <div class="next-review-info">
                <strong>🎯 Prochaine révision:</strong> ${nextReviewText}
            </div>

            <div class="steps-section">
                <div class="steps-header">
                    <strong>📋 Étapes à suivre</strong>
                    <span class="steps-progress">${stepsProgress}</span>
                </div>
                ${ex.steps.map((step, i) => {
            const isCompleted = stepsCompleted.includes(i);
            return `
                        <div class="step-item ${isCompleted ? 'completed' : ''}" 
                             onclick="app.toggleStep('${ex.id}', ${i})">
                            <div class="step-checkbox">${isCompleted ? '✓' : i + 1}</div>
                            <div class="step-text">${step}</div>
                        </div>
                    `;
        }).join('')}
            </div>

            <pre>${ex.content}</pre>

            <div class="completion-section">
                <div class="completion-status">
                    <div class="status-label">
                        <strong>📌 Statut d'Exercice</strong>
                        <div class="current ${ex.completed ? 'completed' : 'pending'}">
                            ${ex.completed ? '✓ Marqué comme complété' : '○ Non complété'}
                        </div>
                    </div>
                    <button class="toggle-btn ${ex.completed ? 'completed' : 'pending'}" 
                            onclick="app.toggleCompletion()">
                        ${ex.completed ? '✓ Marqué' : '○ À Faire'}
                    </button>
                </div>

                ${ex.last_reviewed ? `
                    <div style="font-size: 13px; color: #888;">
                        <strong>Progression:</strong><br>
                        ✓ Complété: ${ex.completed ? 'Oui' : 'Non'}<br>
                        📊 EF: ${ex.ease_factor.toFixed(2)}<br>
                        ⏱ Interval: ${ex.interval_days} jour(s)<br>
                        🔁 Révisions: ${ex.repetitions}
                    </div>
                ` : ''}
            </div>

            <div class="rating-section">
                <div class="rating-label">Comment avez-vous trouvé cet exercice?</div>
                <div class="rating-buttons">
                    <button class="rating-btn r1" onclick="app.rate('${ex.id}', 1)">1 - Oublié</button>
                    <button class="rating-btn r2" onclick="app.rate('${ex.id}', 2)">2 - Difficile</button>
                    <button class="rating-btn r3" onclick="app.rate('${ex.id}', 3)">3 - Normal</button>
                    <button class="rating-btn r4" onclick="app.rate('${ex.id}', 4)">4 - Facile</button>
                </div>
            </div>
        `;
    }

    renderNavStats() {
        const toDoCount = this.allExercises.filter(ex => !ex.completed).length;
        let totalMastery = 0, count = 0;
        this.allExercises.forEach(ex => {
            if (ex.completed) {
                const mastery = Math.round(((ex.ease_factor - 1.3) / (3.0 - 1.3)) * 100);
                totalMastery += mastery;
                count++;
            }
        });
        const avgMastery = count > 0 ? Math.round(totalMastery / count) : 0;

        document.getElementById('navToDo').textContent = toDoCount;
        document.getElementById('navMastery').textContent = avgMastery + '%';

        if (this.stats.streak) {
            const display = this.stats.streak.display || '';
            document.getElementById('navStreak').textContent = `${display || '○'} ${this.stats.streak.current || 0}`;
        }
    }

    renderStatsPage() {
        const completed = this.allExercises.filter(ex => ex.completed).length;
        const total = this.allExercises.length;

        let totalMastery = 0, count = 0;
        this.allExercises.forEach(ex => {
            if (ex.completed) {
                const mastery = Math.round(((ex.ease_factor - 1.3) / (3.0 - 1.3)) * 100);
                totalMastery += mastery;
                count++;
            }
        });
        const avgMastery = count > 0 ? Math.round(totalMastery / count) : 0;

        document.getElementById('completedCard').textContent = completed;
        document.getElementById('totalCard').textContent = total;
        document.getElementById('masteryCard').textContent = avgMastery + '%';
        document.getElementById('streakCard').textContent = this.stats.streak?.current || 0;
        document.getElementById('reviewsCard').textContent = this.stats.total_reviews || 0;

        const domains = {};
        this.allExercises.forEach(ex => {
            if (!domains[ex.domain]) {
                domains[ex.domain] = { completed: 0, total: 0, mastery: 0 };
            }
            domains[ex.domain].total++;
            if (ex.completed) {
                domains[ex.domain].completed++;
                const mastery = Math.round(((ex.ease_factor - 1.3) / (3.0 - 1.3)) * 100);
                domains[ex.domain].mastery = Math.max(domains[ex.domain].mastery, mastery);
            }
        });

        const strengths = Object.entries(domains)
            .filter(([, d]) => d.mastery >= 70)
            .sort((a, b) => b[1].mastery - a[1].mastery);

        const strengthsList = document.getElementById('strengthsList');
        if (strengths.length > 0) {
            strengthsList.innerHTML = strengths.map(([name, data]) => `
                <li>
                    <span class="item-label">${name}:</span>
                    <span class="item-value">${data.mastery}% maîtrisé</span>
                </li>
            `).join('');
        } else {
            strengthsList.innerHTML = '<li><span class="item-label">Continuez à pratiquer pour débloquer!</span></li>';
        }

        const weaknesses = Object.entries(domains)
            .filter(([, d]) => d.mastery < 70)
            .sort((a, b) => a[1].mastery - b[1].mastery);

        const weaknessList = document.getElementById('weaknessList');
        if (weaknesses.length > 0) {
            weaknessList.innerHTML = weaknesses.map(([name, data]) => `
                <li>
                    <span class="item-label">${name}:</span>
                    <span class="item-value">${data.completed}/${data.total} complétés</span>
                </li>
            `).join('');
        } else {
            weaknessList.innerHTML = '<li><span class="item-label">Excellent! Tous les domaines sont en bonne voie!</span></li>';
        }

        const domainStats = document.getElementById('domainStats');
        domainStats.innerHTML = Object.entries(domains).map(([name, data]) => {
            const percentage = Math.round((data.completed / data.total) * 100);
            return `
                <div class="domain-stat">
                    <div class="domain-name">
                        <strong>${name}</strong>
                        <span class="percentage">${percentage}%</span>
                    </div>
                    <div class="progress-bar">
                        <div class="progress-fill" style="width: ${percentage}%"></div>
                    </div>
                    <div style="font-size: 12px; color: #888; margin-top: 5px;">
                        ${data.completed}/${data.total} exercices · Mastery: ${data.mastery}%
                    </div>
                </div>
            `;
        }).join('');
    }

    // ============= ADVANCED STATS RENDERING =============

    async renderAdvancedStats() {
        try {
            const struggling = await this.api.getStruggling().catch(() => []);
            const mastered = await this.api.getMastered().catch(() => []);
            const needsPractice = await this.api.getNeedsPractice().catch(() => []);
            const insights = await this.api.getInsights().catch(() => ({}));

            this.renderInsights(insights);
            this.renderAnalysisExerciseList('struggling', struggling);
            this.renderAnalysisExerciseList('mastered', mastered);
            this.renderAnalysisExerciseList('needs-practice', needsPractice);
            await this.renderDomainAnalysis();
        } catch (error) {
            console.error('Erreur chargement stats avancées:', error);
        }
    }

    renderInsights(insights) {
        document.getElementById('strongestDomain').textContent = insights.strongest_domain || '-';
        document.getElementById('weakestDomain').textContent = insights.weakest_domain || '-';
        document.getElementById('successRate').textContent = Math.round(insights.success_rate || 0) + '%';
        document.getElementById('overdueCount').textContent = insights.overdue || 0;
        document.getElementById('recommendFocus').innerHTML = '<strong>📌 Recommandation:</strong> ' + (insights.recommend_focus || 'À déterminer');
    }

    renderAnalysisExerciseList(category, exercises) {
        const container = document.getElementById(category + 'List');

        if (!exercises || exercises.length === 0) {
            container.innerHTML = '<p style="text-align: center; color: #888;">Aucun exercice dans cette catégorie</p>';
            return;
        }

        container.innerHTML = exercises.map(ex => {
            const level = ex.level_detected || 'unknown';
            const levelColor = {
                'very_easy': '#90EE90',
                'easy': '#90EE90',
                'medium': '#FFD700',
                'hard': '#FF6347',
                'very_hard': '#FF0000'
            }[level] || '#888';

            return `
                <div class="exercise-analysis-card" style="border-left: 4px solid ${levelColor}">
                    <div class="exercise-analysis-header">
                        <strong>${ex.title}</strong>
                        <span style="background: ${levelColor}; color: #000; padding: 4px 8px; border-radius: 3px; font-size: 12px;">
                            ${level.replace('_', ' ').toUpperCase()}
                        </span>
                    </div>
                    <div class="exercise-analysis-meta">
                        <div>EF: ${ex.ef ? ex.ef.toFixed(2) : 'N/A'}</div>
                        <div>Interval: ${ex.interval} j</div>
                        <div>Reps: ${ex.repetitions}</div>
                        <div>Mastery: ${ex.mastery || 0}%</div>
                        ${ex.is_overdue ? '<div style="color: #FF6347;">🔴 OVERDUE</div>' : ''}
                    </div>
                    <div class="exercise-analysis-confidence">
                        Confidence: <strong>${ex.confidence_level || 'Unknown'}</strong>
                    </div>
                    <div class="exercise-analysis-recommendation">
                        ${ex.recommendation}
                    </div>
                </div>
            `;
        }).join('');
    }

    async renderDomainAnalysis() {
        const domains = new Set(this.allExercises.map(ex => ex.domain));
        const container = document.getElementById('domainAnalysis');

        let html = '';
        for (const domain of domains) {
            try {
                const analysis = await this.api.getDomainAnalysis(domain);
                html += `
                    <div class="domain-analysis-card">
                        <h4>${domain.toUpperCase()}</h4>
                        <div class="domain-analysis-grid">
                            <div>Total: ${analysis.total}</div>
                            <div>Complétés: ${analysis.completed}</div>
                            <div>Mastery: ${Math.round(analysis.mastery * 100)}%</div>
                            <div>Avg Reps: ${analysis.avg_repetitions.toFixed(1)}</div>
                        </div>
                        <div>
                            <strong>Struggling:</strong> ${analysis.struggling}
                            <strong>Mastered:</strong> ${analysis.mastered}
                        </div>
                        <div style="color: #90EE90; margin-top: 10px;">
                            ${analysis.recommendation}
                        </div>
                    </div>
                `;
            } catch (e) {
                console.warn(`Could not load domain analysis for ${domain}`);
            }
        }

        container.innerHTML = html;
    }

    // ============= CALENDAR RENDERING =============

    async renderCalendar() {
        try {
            const upcoming = await this.api.getCalendarUpcoming(30);

            // Render overdue
            const overdueList = document.getElementById('overdueList');
            if (upcoming.overdue && upcoming.overdue.length > 0) {
                overdueList.innerHTML = upcoming.overdue.map(ex => {
                    const daysLate = Math.abs(ex.days_until);
                    return `
                        <div class="calendar-item urgent">
                            <div class="calendar-item-header">
                                <strong>${ex.title}</strong>
                                <span class="domain-badge">${ex.domain}</span>
                            </div>
                            <div class="calendar-item-meta">
                                <span class="overdue-badge">-${daysLate}j</span>
                                <span>EF: ${ex.ef.toFixed(2)}</span>
                            </div>
                        </div>
                    `;
                }).join('');
            } else {
                overdueList.innerHTML = '<p class="empty-state">✅ Aucun exercice en retard!</p>';
            }

            // Render today
            const todayList = document.getElementById('todayList');
            if (upcoming.today && upcoming.today.length > 0) {
                todayList.innerHTML = upcoming.today.map(ex => `
                    <div class="calendar-item today">
                        <div class="calendar-item-header">
                            <strong>${ex.title}</strong>
                            <span class="domain-badge">${ex.domain}</span>
                        </div>
                    </div>
                `).join('');
            } else {
                todayList.innerHTML = '<p class="empty-state">😊 Rien pour aujourd\'hui</p>';
            }

            // Render upcoming
            const upcomingList = document.getElementById('upcomingList');
            if (upcoming.upcoming && Object.keys(upcoming.upcoming).length > 0) {
                const sortedDates = Object.keys(upcoming.upcoming).sort();

                upcomingList.innerHTML = sortedDates.map(dateKey => {
                    const exercises = upcoming.upcoming[dateKey];
                    const date = new Date(dateKey);
                    const dayName = date.toLocaleDateString('fr-FR', { weekday: 'long' });
                    const dateStr = date.toLocaleDateString('fr-FR', {
                        day: 'numeric',
                        month: 'long'
                    });

                    return `
                        <div class="calendar-day-group">
                            <div class="calendar-day-header">
                                <strong>${dayName}</strong> ${dateStr}
                            </div>
                            ${exercises.map(ex => `
                                <div class="calendar-item">
                                    <div class="calendar-item-header">
                                        <strong>${ex.title}</strong>
                                        <span class="domain-badge">${ex.domain}</span>
                                    </div>
                                    <div class="calendar-item-meta">
                                        <span class="days-badge">dans ${ex.days_until}j</span>
                                        <span>EF: ${ex.ef.toFixed(2)}</span>
                                    </div>
                                </div>
                            `).join('')}
                        </div>
                    `;
                }).join('');
            } else {
                upcomingList.innerHTML = '<p class="empty-state">Rien de prévu pour les 30 prochains jours</p>';
            }

        } catch (error) {
            console.error('Erreur calendar:', error);
            document.getElementById('overdueList').innerHTML =
                '<p class="empty-state" style="color: red;">Erreur de chargement</p>';
        }
    }

    // ============= RATING =============

    async rate(exerciseId, rating) {
        try {
            const response = await this.api.rateExercise(exerciseId, rating);
            const exercise = this.allExercises.find(ex => ex.id === exerciseId);
            if (exercise && response.exercise) {
                Object.assign(exercise, response.exercise);
            }

            this.updateStats();
            this.applyFilters();
            this.renderExerciseDetail();
        } catch (error) {
            console.error('Erreur notation:', error);
        }
    }

    render() {
        this.applyFilters();
    }
}

let app;
document.addEventListener('DOMContentLoaded', () => {
    app = new MaestroApp();
});
