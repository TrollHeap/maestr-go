// ============= MAESTRO APP v3.2 - FULLY CONNECTED =============

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
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`API Error: ${response.statusText} - ${errorText}`);
            }
            return await response.json();
        } catch (error) {
            console.error('API Request failed:', error);
            throw error;
        }
    }

    getExercises() {
        return this.request('/api/exercises');
    }

    updateExerciseSteps(exerciseId, completedSteps) {
        return this.request(`/api/exercises/${exerciseId}/steps`, {
            method: 'PUT',
            body: JSON.stringify({ completed_steps: completedSteps })
        });
    }

    toggleExerciseCompletion(exerciseId, completed) {
        return this.request(`/api/exercises/${exerciseId}/completion`, {
            method: 'PUT',
            body: JSON.stringify({ completed: completed })
        });
    }

    reviewExercise(exerciseId, rating) {
        return this.request(`/api/exercises/${exerciseId}/review`, {
            method: 'POST',
            body: JSON.stringify({ rating: rating })
        });
    }

    getStats() {
        return this.request('/api/stats');
    }

    getPlannerToday() {
        return this.request('/api/planner/today');
    }

    getPlannerWeek() {
        return this.request('/api/planner/week');
    }

    createSession(session) {
        return this.request('/api/planner/sessions', {
            method: 'POST',
            body: JSON.stringify(session)
        });
    }

    updateSession(sessionId, session) {
        return this.request(`/api/planner/sessions/${sessionId}`, {
            method: 'PUT',
            body: JSON.stringify(session)
        });
    }

    deleteSession(sessionId) {
        return this.request(`/api/planner/sessions/${sessionId}`, {
            method: 'DELETE'
        });
    }
}

class MaestroApp {
    constructor() {
        this.api = new MaestroAPI();
        this.allExercises = [];
        this.filteredExercises = [];
        this.filteredSelectorExercises = [];
        this.selectedExercise = null;
        this.currentPage = 'exercises';
        this.currentStatus = 'all';
        this.currentDomain = 'all';
        this.currentDateFilter = 'all';
        this.pageNum = 1;
        this.pageSize = 10;
        this.stats = {};
        this.completedSteps = {};
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
            const exercises = await this.api.getExercises();
            this.allExercises = Array.isArray(exercises) ? exercises : [];

            this.completedSteps = {};
            this.allExercises.forEach(ex => {
                if (ex.completed_steps && ex.completed_steps.length > 0) {
                    this.completedSteps[ex.id] = ex.completed_steps;
                } else {
                    this.completedSteps[ex.id] = [];
                }
            });
        } catch (error) {
            console.error('Erreur chargement:', error);
            this.allExercises = [];
        }
    }

    async updateStats() {
        try {
            this.stats = await this.api.getStats();
            this.renderNavStats();
        } catch (error) {
            console.error('Erreur stats:', error);
            this.stats = { total: 0, completed: 0, in_progress: 0, due_review: 0 };
        }
    }

    renderNavStats() {
        const statsContainer = document.getElementById('navStats');
        if (!statsContainer) return;

        const { total = 0, completed = 0, in_progress = 0, due_review = 0 } = this.stats;
        const completionRate = total > 0 ? Math.round((completed / total) * 100) : 0;

        statsContainer.innerHTML = `
            <div class="nav-stats-grid">
                <div class="stat-item">
                    <div class="stat-label">Total</div>
                    <div class="stat-value">${total}</div>
                </div>
                <div class="stat-item">
                    <div class="stat-label">Complétés</div>
                    <div class="stat-value" style="color: #90EE90;">${completed}</div>
                </div>
                <div class="stat-item">
                    <div class="stat-label">En cours</div>
                    <div class="stat-value" style="color: #FFA500;">${in_progress}</div>
                </div>
                <div class="stat-item">
                    <div class="stat-label">À réviser</div>
                    <div class="stat-value" style="color: #FF6B6B;">${due_review}</div>
                </div>
                <div class="stat-item stat-full">
                    <div class="stat-label">Taux de complétion</div>
                    <div class="stat-progress">
                        <div class="stat-progress-bar" style="width: ${completionRate}%"></div>
                    </div>
                    <div class="stat-value">${completionRate}%</div>
                </div>
            </div>
        `;
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
        document.querySelector(`[onclick="app.goToPage('${page}')"]`)?.classList.add('active');
        document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
        document.getElementById(page)?.classList.add('active');
        this.currentPage = page;

        if (page === 'planner') {
            this.renderPlannerToday();
        }
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
            const nextReview = this.getNextReviewText(ex);

            return `
                <div class="exercise-item ${isSelected ? 'selected' : ''}" data-id="${ex.id}">
                    <div class="exercise-title">${ex.completed ? '✓' : '○'} ${ex.title}</div>
                    <div class="exercise-meta">
                        <span style="background: #5D4E60; padding: 2px 8px; border-radius: 3px;">${ex.domain}</span>
                        <span style="background: #2d4a6e; padding: 2px 8px; border-radius: 3px;">D${ex.difficulty}</span>
                        <span style="color: #90EE90;">⏱ ${nextReview}</span>
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
        document.getElementById('pageInfo').textContent = `Page ${this.pageNum}/${total || 1}`;
        document.getElementById('prevBtn').disabled = this.pageNum <= 1;
        document.getElementById('nextBtn').disabled = this.pageNum >= total;
    }

    getNextReviewText(ex) {
        if (!ex.last_reviewed) return 'Nouveau';

        const lastReview = new Date(ex.last_reviewed);
        const nextReview = new Date(lastReview);
        nextReview.setDate(nextReview.getDate() + (ex.interval_days || 0));

        const now = new Date();
        const diffDays = Math.ceil((nextReview - now) / (1000 * 60 * 60 * 24));

        if (diffDays < 0) return 'À réviser maintenant';
        if (diffDays === 0) return 'Aujourd\'hui';
        if (diffDays === 1) return 'Demain';
        return `Dans ${diffDays} jours`;
    }

    renderExerciseDetail() {
        const detail = document.getElementById('exerciseDetail');
        if (!this.selectedExercise) {
            detail.innerHTML = '<div style="text-align: center; color: #888;">Sélectionnez un exercice</div>';
            return;
        }

        const ex = this.selectedExercise;
        const nextReview = this.getNextReviewText(ex);

        if (!this.completedSteps[ex.id]) {
            this.completedSteps[ex.id] = [];
        }

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
                <strong>🎯 Prochaine révision:</strong> ${nextReview}
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

    async toggleStep(exerciseId, stepIndex) {
        if (!this.completedSteps[exerciseId]) {
            this.completedSteps[exerciseId] = [];
        }

        const index = this.completedSteps[exerciseId].indexOf(stepIndex);
        if (index > -1) {
            this.completedSteps[exerciseId].splice(index, 1);
        } else {
            this.completedSteps[exerciseId].push(stepIndex);
        }

        try {
            await this.api.updateExerciseSteps(exerciseId, this.completedSteps[exerciseId]);
            const exercise = this.allExercises.find(ex => ex.id === exerciseId);
            if (exercise) {
                exercise.completed_steps = this.completedSteps[exerciseId];
            }
        } catch (error) {
            console.error('Erreur sauvegarde steps:', error);
            if (index > -1) {
                this.completedSteps[exerciseId].push(stepIndex);
            } else {
                const revertIndex = this.completedSteps[exerciseId].indexOf(stepIndex);
                if (revertIndex > -1) {
                    this.completedSteps[exerciseId].splice(revertIndex, 1);
                }
            }
            this.showNotification('Erreur lors de la sauvegarde', 'error');
        }

        this.renderExerciseDetail();
    }

    async toggleCompletion(exerciseId) {
        const exercise = this.allExercises.find(ex => ex.id === exerciseId);
        if (!exercise) return;

        const allStepsCompleted = exercise.steps && exercise.steps.length > 0 &&
            this.completedSteps[exerciseId] &&
            this.completedSteps[exerciseId].length === exercise.steps.length;

        if (!exercise.completed && !allStepsCompleted) {
            return;
        }

        const newCompleted = !exercise.completed;

        try {
            await this.api.toggleExerciseCompletion(exerciseId, newCompleted);
            exercise.completed = newCompleted;

            if (!exercise.completed) {
                this.completedSteps[exerciseId] = [];
                await this.api.updateExerciseSteps(exerciseId, []);
            }

            this.renderExerciseDetail();
            this.renderExerciseList();
            this.updateStats();
            this.showNotification(newCompleted ? 'Exercice complété!' : 'Exercice réouvert', 'success');
        } catch (error) {
            console.error('Erreur toggle completion:', error);
            this.showNotification('Erreur lors de la sauvegarde', 'error');
        }
    }

    async rate(exerciseId, rating) {
        try {
            const response = await this.api.reviewExercise(exerciseId, rating);
            const exercise = this.allExercises.find(ex => ex.id === exerciseId);
            if (exercise && response.exercise) {
                Object.assign(exercise, response.exercise);
            }
            await this.updateStats();
            this.applyFilters();
            this.renderExerciseDetail();
            this.showNotification('Révision enregistrée!', 'success');
        } catch (error) {
            console.error('Erreur notation:', error);
            this.showNotification('Erreur lors de la notation', 'error');
        }
    }

    // ============= PLANNER METHODS =============

    switchPlannerView(view) {
        document.querySelectorAll('.planner-tab').forEach(t => t.classList.remove('active'));
        document.querySelectorAll('.planner-view').forEach(v => v.classList.remove('active'));

        const clickedTab = Array.from(document.querySelectorAll('.planner-tab')).find(
            t => t.textContent.toLowerCase().includes(view)
        );
        if (clickedTab) clickedTab.classList.add('active');

        const viewElement = document.getElementById('planner-' + view);
        if (viewElement) viewElement.classList.add('active');

        if (view === 'today') {
            this.renderPlannerToday();
        } else if (view === 'week') {
            this.renderPlannerWeek();
        }
    }

    async renderPlannerToday() {
        try {
            const plan = await this.api.getPlannerToday();

            ['morning', 'afternoon', 'evening'].forEach(slot => {
                const container = document.getElementById(slot + '-sessions');
                if (!container) return;

                const sessions = plan.sessions?.filter(s => s.time_slot === slot) || [];

                if (sessions.length > 0) {
                    container.innerHTML = sessions.map(session => {
                        const exerciseTitles = session.exercise_ids?.map(id => {
                            const ex = this.allExercises.find(e => e.id === id);
                            return ex ? ex.title : id;
                        }).join(', ') || 'Aucun exercice';

                        return `
                            <div class="session-card ${session.status}">
                                <div class="session-header">
                                    <span>${session.exercise_ids?.length || 0} exercice(s): ${exerciseTitles}</span>
                                    <span style="background: #2a5a72; padding: 3px 10px; border-radius: 4px;">${session.duration} min</span>
                                </div>
                                <div class="session-actions">
                                    <button style="background: #90EE90; color: #000; padding: 8px 16px; border: none; border-radius: 6px; cursor: pointer;" onclick="app.completeSession('${session.id}')">✓ Compléter</button>
                                    <button style="background: #FF6347; color: #fff; padding: 8px 16px; border: none; border-radius: 6px; cursor: pointer;" onclick="app.deleteSession('${session.id}')">× Supprimer</button>
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

            const statsContainer = document.getElementById('plannerDailyStats');
            if (statsContainer) {
                statsContainer.innerHTML = `
                    <div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 15px;">
                        <div><span style="color: #888;">Sessions:</span> <strong style="color: #90EE90;">${completed}/${total}</strong></div>
                        <div><span style="color: #888;">Progression:</span> <strong style="color: #90EE90;">${progress}%</strong></div>
                        <div><span style="color: #888;">Temps:</span> <strong style="color: #90EE90;">${plan.total_minutes || 0} min</strong></div>
                        <div><span style="color: #888;">Statut:</span> <strong style="color: #90EE90;">${completed === total && total > 0 ? '✓ Terminé' : '○ En cours'}</strong></div>
                    </div>
                `;
            }
        } catch (error) {
            console.error('Erreur planner today:', error);
        }
    }

    async renderPlannerWeek() {
        try {
            const plan = await this.api.getPlannerWeek();
            const weekGrid = document.getElementById('week-grid');
            if (!weekGrid) return;

            if (!plan || !plan.days || plan.days.length === 0) {
                weekGrid.innerHTML = '<p style="text-align: center; color: #888;">Aucune session planifiée cette semaine</p>';
                return;
            }

            weekGrid.innerHTML = plan.days.map(day => {
                const date = new Date(day.date);
                const dayName = date.toLocaleDateString('fr-FR', { weekday: 'long' });
                const dateStr = date.toLocaleDateString('fr-FR', { day: 'numeric', month: 'long' });
                const progress = day.total > 0 ? Math.round((day.completed / day.total) * 100) : 0;

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
                        <div class="week-day-progress">
                            <div class="week-day-progress-fill" style="width: ${progress}%"></div>
                        </div>
                    </div>
                `;
            }).join('');
        } catch (error) {
            console.error('Erreur planner week:', error);
        }
    }

    addSession(timeSlot) {
        this.selectedTimeSlot = timeSlot;
        this.selectedExercisesForSession = [];
        this.showExerciseSelector();
    }

    showExerciseSelector() {
        const modal = document.getElementById('exerciseSelectorModal');
        const list = document.getElementById('exerciseSelectorList');
        if (!modal || !list) return;

        // ✅ FILTRER: Exercices non complétés + à réviser aujourd'hui
        const now = new Date();
        this.filteredSelectorExercises = this.allExercises.filter(ex => {
            if (ex.deleted) return false;

            // Non complétés
            if (!ex.completed) return true;

            // Révisions dues
            if (!ex.last_reviewed || !ex.interval_days) return false;

            const lastReview = new Date(ex.last_reviewed);
            const nextReview = new Date(lastReview);
            nextReview.setDate(nextReview.getDate() + ex.interval_days);

            return nextReview <= now;
        });

        if (this.filteredSelectorExercises.length === 0) {
            list.innerHTML = `
                <div style="text-align: center; padding: 40px; color: #888;">
                    <p style="font-size: 18px; margin-bottom: 10px;">🎉 Aucun exercice à faire!</p>
                    <p style="font-size: 14px;">Tous les exercices sont complétés ou pas encore à réviser.</p>
                </div>
            `;
            modal.classList.add('active');
            return;
        }

        this.renderFilteredSelectorList(this.filteredSelectorExercises);
        modal.classList.add('active');
    }

    filterExerciseSelector(searchTerm) {
        const filtered = this.filteredSelectorExercises.filter(ex =>
            ex.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
            ex.domain.toLowerCase().includes(searchTerm.toLowerCase())
        );
        this.renderFilteredSelectorList(filtered);
    }

    renderFilteredSelectorList(exercises) {
        const list = document.getElementById('exerciseSelectorList');

        list.innerHTML = exercises.map(ex => {
            const isDueReview = ex.completed && ex.last_reviewed;
            const statusBadge = isDueReview
                ? '<span class="status-badge due">À RÉVISER</span>'
                : '<span class="status-badge todo">À FAIRE</span>';

            return `
                <div class="exercise-selector-item" data-id="${ex.id}" onclick="app.toggleExerciseSelection('${ex.id}')">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 5px;">
                        <div style="font-weight: bold;">${ex.title}</div>
                        ${statusBadge}
                    </div>
                    <div style="display: flex; gap: 10px; font-size: 12px; color: #888;">
                        <span>${ex.domain}</span>
                        <span>D${ex.difficulty}</span>
                        <span>${ex.steps ? ex.steps.length + ' étapes' : 'Pas d\'étapes'}</span>
                    </div>
                </div>
            `;
        }).join('');
    }

    toggleExerciseSelection(exerciseId) {
        const items = document.querySelectorAll(`[data-id="${exerciseId}"]`);
        items.forEach(item => item.classList.toggle('selected'));

        const index = this.selectedExercisesForSession.indexOf(exerciseId);
        if (index > -1) {
            this.selectedExercisesForSession.splice(index, 1);
        } else {
            this.selectedExercisesForSession.push(exerciseId);
        }

        this.updateConfirmButton();
    }

    updateConfirmButton() {
        const confirmBtn = document.getElementById('confirmExerciseBtn');
        if (!confirmBtn) return;

        const count = this.selectedExercisesForSession.length;
        if (count === 0) {
            confirmBtn.textContent = 'Sélectionner des exercices';
            confirmBtn.disabled = true;
        } else {
            confirmBtn.textContent = `Créer session (${count} exercice${count > 1 ? 's' : ''})`;
            confirmBtn.disabled = false;
        }
    }

    closeExerciseSelector() {
        const modal = document.getElementById('exerciseSelectorModal');
        if (modal) modal.classList.remove('active');
        this.selectedExercisesForSession = [];
        this.updateConfirmButton();

        const searchInput = document.getElementById('exerciseSearch');
        if (searchInput) searchInput.value = '';
    }

    async confirmExerciseSelection() {
        if (this.selectedExercisesForSession.length === 0) {
            this.showNotification('Veuillez sélectionner au moins un exercice', 'error');
            return;
        }

        const duration = this.selectedExercisesForSession.length * 30;
        const today = new Date().toISOString().split('T')[0];

        const sessionData = {
            date: today,
            time_slot: this.selectedTimeSlot,
            exercise_ids: this.selectedExercisesForSession,
            duration: duration,
            status: "planned",
            notes: ""
        };

        console.log('📤 Création session:', sessionData);

        try {
            const response = await this.api.createSession(sessionData);
            console.log('✅ Session créée:', response);

            this.closeExerciseSelector();
            await this.renderPlannerToday();

            this.showNotification(`Session créée avec ${this.selectedExercisesForSession.length} exercice(s)`, 'success');
        } catch (error) {
            console.error('❌ Erreur create session:', error);
            this.showNotification('Erreur lors de la création de la session', 'error');
        }
    }

    async completeSession(sessionId) {
        try {
            await this.api.updateSession(sessionId, { status: 'completed' });
            await this.renderPlannerToday();
            this.showNotification('Session complétée!', 'success');
        } catch (error) {
            console.error('Erreur complete session:', error);
            this.showNotification('Erreur lors de la complétion', 'error');
        }
    }

    async deleteSession(sessionId) {
        if (!confirm('Supprimer cette session?')) return;

        try {
            await this.api.deleteSession(sessionId);
            await this.renderPlannerToday();
            this.showNotification('Session supprimée', 'success');
        } catch (error) {
            console.error('Erreur delete session:', error);
            this.showNotification('Erreur lors de la suppression', 'error');
        }
    }

    // ============= NOTIFICATION SYSTEM =============

    showNotification(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;
        notification.textContent = message;

        const container = document.getElementById('notificationContainer');
        if (container) {
            container.appendChild(notification);
        } else {
            document.body.appendChild(notification);
        }

        setTimeout(() => {
            notification.classList.add('fade-out');
            setTimeout(() => notification.remove(), 300);
        }, 3000);
    }

    render() {
        this.applyFilters();
    }
}

// ✅ INITIALISATION
let app;
document.addEventListener('DOMContentLoaded', () => {
    app = new MaestroApp();
});
