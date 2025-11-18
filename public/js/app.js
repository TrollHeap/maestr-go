// ============= MAESTRO APP v3.1 WITH FILTERS & PLANNER =============

class MaestroAPI {
    constructor(baseURL = 'http://localhost:8080') {
        this.baseURL = baseURL;
    }

    async request(endpoint, options = {}) {
        try {
            const response = await fetch(`${this.baseURL}${endpoint}`, {
                headers: { 'Content-Type': 'application/json', ...options.headers },
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
            body: JSON.stringify({ exercise_id: exerciseId, rating: rating })
        });
    }

    getStats() {
        return this.request('/api/stats');
    }

    getPlannerToday() {
        return this.request('/api/planner/today');
    }

    getPlannerWeek(startDate) {
        const query = startDate ? `?start=${startDate}` : '';
        return this.request(`/api/planner/week${query}`);
    }

    getPlannerMonth(startDate) {
        const query = startDate ? `?start=${startDate}` : '';
        return this.request(`/api/planner/month${query}`);
    }

    createSession(date, timeSlot, exerciseIDs, duration) {
        return this.request('/api/planner/session', {
            method: 'POST',
            body: JSON.stringify({ date, time_slot: timeSlot, exercise_ids: exerciseIDs, duration })
        });
    }

    updateSession(sessionID, status, notes = '') {
        return this.request(`/api/planner/session/${sessionID}`, {
            method: 'PUT',
            body: JSON.stringify({ status, notes })
        });
    }

    deleteSession(sessionID) {
        return this.request(`/api/planner/session/${sessionID}`, { method: 'DELETE' });
    }
}

class MaestroApp {
    constructor() {
        this.api = new MaestroAPI();
        this.allExercises = [];
        this.filteredExercises = [];
        this.selectedExercise = null;
        this.currentPage = 'exercises';
        this.currentStatus = 'all';
        this.currentDomain = 'all';
        this.currentDateFilter = 'all'; // NOUVEAU
        this.pageNum = 1;
        this.pageSize = 10;
        this.stats = {};
        this.completedSteps = {};
        this.reviewDates = {};
        this.selectedTimeSlot = null;
        this.selectedExercisesForSession = [];

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
            this.renderNavStats();
        } catch (error) {
            console.error('Erreur stats:', error);
        }
    }

    setupEventListeners() {
        // Status filters
        document.querySelectorAll('[data-status]').forEach(btn => {
            btn.addEventListener('click', (e) => {
                document.querySelectorAll('[data-status]').forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
                this.currentStatus = e.target.dataset.status;
                this.pageNum = 1;
                this.applyFilters();
            });
        });

        // Domain filters
        document.querySelectorAll('[data-domain]').forEach(btn => {
            btn.addEventListener('click', (e) => {
                document.querySelectorAll('[data-domain]').forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
                this.currentDomain = e.target.dataset.domain;
                this.pageNum = 1;
                this.applyFilters();
            });
        });

        // NOUVEAU: Date filters
        document.querySelectorAll('[data-date]').forEach(btn => {
            btn.addEventListener('click', (e) => {
                document.querySelectorAll('[data-date]').forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
                this.currentDateFilter = e.target.dataset.date;
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

        if (page === 'planner') {
            this.renderPlannerToday();
        }
    }

    applyFilters() {
        let filtered = this.allExercises;

        // Status filter
        if (this.currentStatus === 'pending') {
            filtered = filtered.filter(ex => !ex.completed);
        } else if (this.currentStatus === 'completed') {
            filtered = filtered.filter(ex => ex.completed);
        }

        // Domain filter
        if (this.currentDomain !== 'all') {
            filtered = filtered.filter(ex => ex.domain === this.currentDomain);
        }

        // NOUVEAU: Date filter
        if (this.currentDateFilter !== 'all') {
            const now = new Date();
            const weekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
            const monthAgo = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);

            filtered = filtered.filter(ex => {
                if (!ex.last_reviewed) return false;
                const reviewDate = new Date(ex.last_reviewed);

                if (this.currentDateFilter === 'week') {
                    return reviewDate >= weekAgo;
                } else if (this.currentDateFilter === 'month') {
                    return reviewDate >= monthAgo;
                }
                return true;
            });
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
        document.getElementById('pageInfo').textContent = `Page ${this.pageNum}/${total}`;
        document.getElementById('prevBtn').disabled = this.pageNum <= 1;
        document.getElementById('nextBtn').disabled = this.pageNum >= total;
    }

    renderExerciseDetail() {
        const detail = document.getElementById('exerciseDetail');
        if (!this.selectedExercise) {
            detail.innerHTML = '<div style="text-align: center; color: #888;">Sélectionnez un exercice</div>';
            return;
        }

        const ex = this.selectedExercise;
        const nextReviewText = this.reviewDates[ex.id] || 'Nouveau';

        // Initialize completed steps for this exercise if not exists
        if (!this.completedSteps[ex.id]) {
            this.completedSteps[ex.id] = [];
        }

        // GÉNÉRER HTML DES STEPS AVEC CHECKBOXES
        const stepsHTML = ex.steps && ex.steps.length > 0 ? `
        <div class="steps-section">
            <div class="steps-header">
                <h4>📝 Étapes</h4>
                <span class="steps-progress">${this.completedSteps[ex.id].length}/${ex.steps.length} complétées</span>
            </div>
            <div class="steps-list">
                ${ex.steps.map((step, index) => {
            const isCompleted = this.completedSteps[ex.id].includes(index);
            return `
                        <div class="step-item ${isCompleted ? 'completed' : ''}" onclick="app.toggleStep('${ex.id}', ${index})">
                            <div class="step-checkbox">
                                ${isCompleted ? '✓' : (index + 1)}
                            </div>
                            <div class="step-text">${step}</div>
                        </div>
                    `;
        }).join('')}
            </div>
        </div>
    ` : '';

        // BOUTON MARQUER COMME TERMINÉ
        const allStepsCompleted = ex.steps && ex.steps.length > 0 &&
            this.completedSteps[ex.id].length === ex.steps.length;

        const completionButton = `
        <div class="completion-section">
            <div class="completion-status">
                <div class="status-label">
                    <strong>Statut de l'exercice</strong>
                    <span class="current ${ex.completed ? 'completed' : 'pending'}">
                        ${ex.completed ? '✓ Complété' : '○ En cours'}
                    </span>
                </div>
                <button 
                    class="toggle-btn ${ex.completed ? 'completed' : 'pending'} ${!allStepsCompleted && !ex.completed ? 'disabled' : ''}" 
                    onclick="app.toggleCompletion('${ex.id}')"
                    ${!allStepsCompleted && !ex.completed ? 'disabled' : ''}>
                    ${ex.completed ? '✓ Exercice Terminé' : allStepsCompleted ? '✓ Marquer comme terminé' : '○ Compléter d\'abord toutes les étapes'}
                </button>
            </div>
            ${!allStepsCompleted && !ex.completed ? `
                <div class="completion-hint">
                    💡 Complétez toutes les étapes ci-dessus pour débloquer la validation
                </div>
            ` : ''}
        </div>
    `;

        detail.innerHTML = `
        <div class="detail-header">
            <h3>${ex.title}</h3>
            <p class="detail-description">${ex.description}</p>
        </div>
        
        <div class="next-review-info">
            <strong>🎯 Prochaine révision:</strong> ${nextReviewText}
        </div>

        ${stepsHTML}
        ${completionButton}

        <div class="rating-section">
            <div class="rating-label">Comment avez-vous trouvé cet exercice?</div>
            <div class="rating-buttons">
                <button class="rating-btn r1" onclick="app.rate('${ex.id}', 1)">
                    <span class="rating-icon">😰</span>
                    <span class="rating-text">1 - Oublié</span>
                </button>
                <button class="rating-btn r2" onclick="app.rate('${ex.id}', 2)">
                    <span class="rating-icon">😓</span>
                    <span class="rating-text">2 - Difficile</span>
                </button>
                <button class="rating-btn r3" onclick="app.rate('${ex.id}', 3)">
                    <span class="rating-icon">😊</span>
                    <span class="rating-text">3 - Normal</span>
                </button>
                <button class="rating-btn r4" onclick="app.rate('${ex.id}', 4)">
                    <span class="rating-icon">😎</span>
                    <span class="rating-text">4 - Facile</span>
                </button>
            </div>
        </div>
    `;
    }

    // NOUVELLE MÉTHODE: Toggle step completion
    toggleStep(exerciseId, stepIndex) {
        if (!this.completedSteps[exerciseId]) {
            this.completedSteps[exerciseId] = [];
        }

        const index = this.completedSteps[exerciseId].indexOf(stepIndex);
        if (index > -1) {
            // Uncheck
            this.completedSteps[exerciseId].splice(index, 1);
        } else {
            // Check
            this.completedSteps[exerciseId].push(stepIndex);
        }

        // Re-render to update UI
        this.renderExerciseDetail();
    }

    // NOUVELLE MÉTHODE: Toggle exercise completion
    async toggleCompletion(exerciseId) {
        const exercise = this.allExercises.find(ex => ex.id === exerciseId);
        if (!exercise) return;

        // Si l'exercice n'est pas complété et que toutes les étapes ne sont pas cochées, on ne fait rien
        const allStepsCompleted = exercise.steps && exercise.steps.length > 0 &&
            this.completedSteps[exerciseId] &&
            this.completedSteps[exerciseId].length === exercise.steps.length;

        if (!exercise.completed && !allStepsCompleted) {
            return; // Bouton désactivé
        }

        // Toggle completion
        exercise.completed = !exercise.completed;

        // Si on décoche, reset les steps
        if (!exercise.completed) {
            this.completedSteps[exerciseId] = [];
        }

        // Re-render
        this.renderExerciseDetail();
        this.renderExerciseList();
        this.updateStats();
    }
    // ============= PLANNER METHODS =============

    switchPlannerView(view) {
        document.querySelectorAll('.planner-tab').forEach(t => t.classList.remove('active'));
        document.querySelectorAll('.planner-view').forEach(v => v.classList.remove('active'));

        event.target.classList.add('active');
        document.getElementById('planner-' + view).classList.add('active');

        if (view === 'today') {
            this.renderPlannerToday();
        } else if (view === 'week') {
            this.renderPlannerWeek();
        } else if (view === 'month') {
            this.renderPlannerMonth();
        }
    }

    async renderPlannerToday() {
        try {
            const plan = await this.api.getPlannerToday();

            ['morning', 'afternoon', 'evening'].forEach(slot => {
                const container = document.getElementById(slot + '-sessions');
                const sessions = plan.sessions?.filter(s => s.time_slot === slot) || [];

                if (sessions.length > 0) {
                    container.innerHTML = sessions.map(session => {
                        const exerciseTitles = session.exercise_ids.map(id => {
                            const ex = this.allExercises.find(e => e.id === id);
                            return ex ? ex.title : id;
                        }).join(', ');

                        return `
                            <div class="session-card ${session.status}">
                                <div class="session-header">
                                    <span>${session.exercise_ids.length} exercice(s): ${exerciseTitles}</span>
                                    <span style="background: #2a5a72; padding: 3px 10px; border-radius: 4px;">${session.duration} min</span>
                                </div>
                                <div class="session-actions">
                                    <button style="background: #90EE90; color: #000;" onclick="app.completeSession('${session.id}')">✓ Compléter</button>
                                    <button style="background: #FF6347; color: #fff;" onclick="app.deleteSession('${session.id}')">× Supprimer</button>
                                </div>
                            </div>
                        `;
                    }).join('');
                } else {
                    container.innerHTML = '<p style="text-align: center; color: #888;">Aucune session planifiée</p>';
                }
            });

            const completed = plan.sessions?.filter(s => s.status === 'completed').length || 0;
            const total = plan.sessions?.length || 0;
            const progress = total > 0 ? Math.round((completed / total) * 100) : 0;

            document.getElementById('plannerDailyStats').innerHTML = `
                <div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 15px;">
                    <div><span style="color: #888;">Sessions:</span> <strong style="color: #90EE90;">${completed}/${total}</strong></div>
                    <div><span style="color: #888;">Progression:</span> <strong style="color: #90EE90;">${progress}%</strong></div>
                    <div><span style="color: #888;">Temps:</span> <strong style="color: #90EE90;">${plan.total_minutes} min</strong></div>
                    <div><span style="color: #888;">Statut:</span> <strong style="color: #90EE90;">${completed === total && total > 0 ? '✓ Terminé' : '○ En cours'}</strong></div>
                </div>
            `;

        } catch (error) {
            console.error('Erreur planner today:', error);
        }
    }

    async renderPlannerWeek() {
        try {
            const plan = await this.api.getPlannerWeek();
            const weekGrid = document.getElementById('week-grid');

            if (!plan || !plan.days || plan.days.length === 0) {
                weekGrid.innerHTML = '<p style="text-align: center; color: #888;">Aucune session planifiée cette semaine</p>';
                return;
            }

            weekGrid.innerHTML = plan.days.map(day => {
                const date = new Date(day.date);
                const dayName = date.toLocaleDateString('fr-FR', { weekday: 'long' });
                const dateStr = date.toLocaleDateString('fr-FR', { day: 'numeric', month: 'long' });
                const progress = day.total > 0 ? Math.round((day.completed / day.total) * 100) : 0;

                const timeSlots = ['morning', 'afternoon', 'evening'];
                const slotEmojis = { morning: '🌅', afternoon: '🌞', evening: '🌙' };
                const slotNames = { morning: 'Matin', afternoon: 'Après-midi', evening: 'Soir' };

                const slotSummary = timeSlots.map(slot => {
                    const sessions = day.sessions?.filter(s => s.time_slot === slot) || [];
                    const count = sessions.length;
                    const minutes = sessions.reduce((sum, s) => sum + s.duration, 0);
                    return count > 0 ? `${slotEmojis[slot]} ${slotNames[slot]}: ${count} session(s) (${minutes} min)` : null;
                }).filter(Boolean).join('<br>');

                return `
                    <div class="week-day-card">
                        <div class="week-day-header">
                            <div>
                                <div class="week-day-title">${dayName}</div>
                                <div class="week-day-date">${dateStr}</div>
                            </div>
                            <div style="text-align: right;">
                                <strong style="color: #90EE90;">${day.completed}/${day.total} sessions</strong><br>
                                <span style="color: #888;">${day.total_minutes} min</span>
                            </div>
                        </div>
                        <div>${slotSummary || '<span style="color: #888;">Aucune session</span>'}</div>
                        <div class="week-day-progress">
                            <div class="week-day-progress-fill" style="width: ${progress}%"></div>
                        </div>
                    </div>
                `;
            }).join('');

            // Week summary
            weekGrid.innerHTML += `
                <div class="week-summary">
                    <h4 style="color: #90EE90; margin-bottom: 15px;">📊 Résumé de la semaine</h4>
                    <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px;">
                        <div><span style="color: #888;">Total:</span> <strong style="color: #90EE90;">${plan.completed}/${plan.total} sessions</strong></div>
                        <div><span style="color: #888;">Progression:</span> <strong style="color: #90EE90;">${plan.total > 0 ? Math.round((plan.completed / plan.total) * 100) : 0}%</strong></div>
                        <div><span style="color: #888;">Temps:</span> <strong style="color: #90EE90;">${plan.total_minutes} min</strong></div>
                    </div>
                </div>
            `;
        } catch (error) {
            console.error('Erreur planner week:', error);
        }
    }

    async renderPlannerMonth() {
        try {
            const response = await this.api.getPlannerMonth();
            const monthOverview = document.getElementById('month-overview');

            if (!response || !response.weeks || response.weeks.length === 0) {
                monthOverview.innerHTML = '<p style="text-align: center; color: #888;">Aucune session planifiée ce mois</p>';
                return;
            }

            const now = new Date();
            const currentWeekStart = this.getMonday(now).toISOString().split('T')[0];

            monthOverview.innerHTML = response.weeks.map((week, index) => {
                const isCurrent = week.start_date === currentWeekStart;
                const progress = week.total > 0 ? Math.round((week.completed / week.total) * 100) : 0;
                const blocks = Array(10).fill(0).map((_, i) => i < Math.floor(progress / 10));

                const startDate = new Date(week.start_date);
                const endDate = new Date(week.end_date);
                const weekLabel = `${startDate.getDate()}-${endDate.getDate()} ${startDate.toLocaleDateString('fr-FR', { month: 'short' })}`;

                return `
                    <div class="month-week-card ${isCurrent ? 'current' : ''}">
                        <div class="month-week-header">
                            <div class="month-week-title">Semaine ${index + 1} (${weekLabel}) ${isCurrent ? '← ACTUELLE' : ''}</div>
                            <div style="text-align: right;">
                                <strong style="color: #90EE90;">${progress}%</strong> (${week.completed}/${week.total})
                            </div>
                        </div>
                        <div class="progress-bar-visual">
                            ${blocks.map(filled => `<div class="progress-block ${filled ? 'filled' : ''}"></div>`).join('')}
                        </div>
                        <div style="color: #888; font-size: 14px;">⏱️ ${week.total_minutes} minutes</div>
                    </div>
                `;
            }).join('');

            // Month summary
            const totalCompleted = response.weeks.reduce((sum, w) => sum + w.completed, 0);
            const totalSessions = response.weeks.reduce((sum, w) => sum + w.total, 0);
            const totalMinutes = response.weeks.reduce((sum, w) => sum + w.total_minutes, 0);

            monthOverview.innerHTML += `
                <div class="week-summary">
                    <h4 style="color: #90EE90; margin-bottom: 15px;">📊 Résumé du mois</h4>
                    <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px;">
                        <div><span style="color: #888;">Total:</span> <strong style="color: #90EE90;">${totalCompleted}/${totalSessions} sessions</strong></div>
                        <div><span style="color: #888;">Taux:</span> <strong style="color: #90EE90;">${totalSessions > 0 ? Math.round((totalCompleted / totalSessions) * 100) : 0}%</strong></div>
                        <div><span style="color: #888;">Temps:</span> <strong style="color: #90EE90;">${totalMinutes} min</strong></div>
                    </div>
                </div>
            `;
        } catch (error) {
            console.error('Erreur planner month:', error);
        }
    }

    getMonday(date) {
        const d = new Date(date);
        const day = d.getDay();
        const diff = d.getDate() - day + (day === 0 ? -6 : 1);
        return new Date(d.setDate(diff));
    }

    addSession(timeSlot) {
        this.selectedTimeSlot = timeSlot;
        this.selectedExercisesForSession = [];
        this.showExerciseSelector();
    }

    showExerciseSelector() {
        const modal = document.getElementById('exerciseSelectorModal');
        const list = document.getElementById('exerciseSelectorList');

        list.innerHTML = this.allExercises.map(ex => `
            <div class="exercise-selector-item" data-id="${ex.id}" onclick="app.toggleExerciseSelection('${ex.id}')">
                <div style="font-weight: bold; margin-bottom: 5px;">${ex.title}</div>
                <div style="display: flex; gap: 10px; font-size: 12px; color: #888;">
                    <span>${ex.domain}</span>
                    <span>D${ex.difficulty}</span>
                    <span>${ex.completed ? '✓ Complété' : '○ À faire'}</span>
                </div>
            </div>
        `).join('');

        modal.classList.add('active');
    }

    toggleExerciseSelection(exerciseId) {
        const item = event.currentTarget;
        item.classList.toggle('selected');

        const index = this.selectedExercisesForSession.indexOf(exerciseId);
        if (index > -1) {
            this.selectedExercisesForSession.splice(index, 1);
        } else {
            this.selectedExercisesForSession.push(exerciseId);
        }
    }

    closeExerciseSelector() {
        document.getElementById('exerciseSelectorModal').classList.remove('active');
        this.selectedExercisesForSession = [];
    }

    async confirmExerciseSelection() {
        if (this.selectedExercisesForSession.length === 0) {
            alert('Veuillez sélectionner au moins un exercice');
            return;
        }

        const duration = this.selectedExercisesForSession.length * 30;
        const today = new Date().toISOString().split('T')[0];

        try {
            await this.api.createSession(today, this.selectedTimeSlot, this.selectedExercisesForSession, duration);
            this.closeExerciseSelector();
            this.renderPlannerToday();
        } catch (error) {
            console.error('Erreur create session:', error);
            alert('Erreur lors de la création de la session');
        }
    }

    async completeSession(sessionID) {
        try {
            await this.api.updateSession(sessionID, 'completed', '');
            this.renderPlannerToday();
        } catch (error) {
            console.error('Erreur complete session:', error);
        }
    }

    async deleteSession(sessionID) {
        if (!confirm('Supprimer cette session?')) return;

        try {
            await this.api.deleteSession(sessionID);
            this.renderPlannerToday();
        } catch (error) {
            console.error('Erreur delete session:', error);
        }
    }

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
