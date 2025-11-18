// ============= API CLIENT =============
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

            if (!response.ok) {
                throw new Error(`API Error: ${response.statusText}`);
            }

            return await response.json();
        } catch (error) {
            console.error('API Request failed:', error);
            throw error;
        }
    }

    getExercises(page = 1, pageSize = 10) {
        return this.request(`/api/exercises?page=${page}&page_size=${pageSize}`);
    }

    getRecommended() {
        return this.request('/api/recommended');
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
}

// ============= APP STATE =============
class MaestroApp {
    constructor() {
        this.api = new MaestroAPI();
        this.exercises = [];
        this.allExercises = [];
        this.selectedExercise = null;
        this.currentFilter = 'all';
        this.currentPage = 1;
        this.pageSize = 10;
        this.stats = {};
        this.reviewDates = {};
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
            const response = await this.api.getExercises(this.currentPage, this.pageSize);
            this.exercises = response.exercises || [];
            this.allExercises = response.exercises || [];
            this.reviewDates = response.review_dates || {};
            console.log('✅ Exercices chargés:', this.exercises.length);
        } catch (error) {
            this.showError('Erreur chargement exercices: ' + error.message);
            console.error(error);
        }
    }

    async updateStats() {
        try {
            this.stats = await this.api.getStats();
            this.renderStats();
        } catch (error) {
            console.error('Erreur stats:', error);
        }
    }

    setupEventListeners() {
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
                this.currentFilter = e.target.dataset.domain;
                this.currentPage = 1;
                this.renderExerciseList();
            });
        });
    }

    getFilteredExercises() {
        if (this.currentFilter === 'all') {
            return this.exercises;
        }
        return this.exercises.filter(ex => ex.domain === this.currentFilter);
    }

    renderPagination() {
        const pagination = document.getElementById('pagination');
        const totalPages = Math.ceil(this.allExercises.length / this.pageSize);

        if (totalPages <= 1) {
            pagination.innerHTML = '';
            return;
        }

        let html = '';
        for (let i = 1; i <= Math.min(totalPages, 5); i++) {
            html += `
                        <button class="pagination-btn ${i === this.currentPage ? 'active' : ''}"
                                onclick="app.goToPage(${i})">
                            ${i}
                        </button>
                    `;
        }
        if (totalPages > 5) {
            html += `<button class="pagination-btn" disabled>...</button>`;
        }
        pagination.innerHTML = html;
    }

    goToPage(page) {
        this.currentPage = page;
        this.loadExercises();
        this.renderExerciseList();
        this.renderPagination();
    }

    renderExerciseList() {
        const list = document.getElementById('exerciseList');
        const filtered = this.getFilteredExercises();

        list.innerHTML = filtered.map(ex => {
            const isSelected = this.selectedExercise?.id === ex.id;
            const nextReviewText = this.reviewDates[ex.id] || 'Nouveau';

            return `
                        <div class="exercise-item ${isSelected ? 'selected' : ''}" 
                             data-id="${ex.id}">
                            <div class="exercise-title">${ex.title}</div>
                            <div class="exercise-meta">
                                <span style="background: #5D4E60; padding: 2px 8px; border-radius: 3px;">${ex.domain}</span>
                                <span style="background: #2d4a6e; padding: 2px 8px; border-radius: 3px;">D${ex.difficulty}</span>
                            </div>
                            <div class="next-review-badge">
                                ⏱ ${nextReviewText}
                            </div>
                        </div>
                    `;
        }).join('');

        this.renderPagination();

        document.querySelectorAll('.exercise-item').forEach(item => {
            item.addEventListener('click', () => {
                const id = item.dataset.id;
                this.selectedExercise = [...this.allExercises].find(ex => ex.id === id);
                this.renderExerciseDetail();
                this.renderExerciseList();
            });
        });
    }

    renderExerciseDetail() {
        const detail = document.getElementById('exerciseDetail');

        if (!this.selectedExercise) {
            detail.innerHTML = '<div class="loading">Sélectionnez un exercice</div>';
            return;
        }

        const ex = this.selectedExercise;
        const nextReviewText = this.reviewDates[ex.id] || 'Nouveau';

        detail.innerHTML = `
                    <div class="detail-header">
                        <h3>${ex.title}</h3>
                        <p class="detail-description">${ex.description}</p>
                    </div>

                    <div class="next-review-info">
                        <strong>Prochaine révision:</strong> ${nextReviewText}
                    </div>

                    <div class="steps">
                        <strong>Étapes à suivre:</strong>
                        ${ex.steps.map((step, i) => `
                            <div class="step">
                                <div class="step-number">${i + 1}</div>
                                <div class="step-text">${step}</div>
                            </div>
                        `).join('')}
                    </div>

                    <pre>${ex.content}</pre>

                    ${ex.last_reviewed ? `
                        <div style="margin-top: 20px; padding: 15px; background: #0f1429; border-radius: 4px;">
                            <strong>Progression:</strong><br>
                            ✓ Complété: ${ex.completed ? 'Oui' : 'Non'}<br>
                            📊 EF: ${ex.ease_factor.toFixed(2)}<br>
                            ⏱ Interval: ${ex.interval_days} jour(s)
                        </div>
                    ` : ''}

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

    renderStats() {
        const totalCompleted = this.allExercises.filter(ex => ex.completed).length;

        let totalMastery = 0;
        let count = 0;
        this.allExercises.forEach(ex => {
            if (ex.completed) {
                const mastery = Math.round(((ex.ease_factor - 1.3) / (2.5 - 1.3)) * 100);
                totalMastery += mastery;
                count++;
            }
        });
        const avgMastery = count > 0 ? Math.round(totalMastery / count) : 0;

        document.getElementById('completed').textContent = totalCompleted;
        document.getElementById('mastery').textContent = avgMastery;

        // Afficher le streak
        if (this.stats.streak) {
            const streakDisplay = this.stats.streak.display || '';
            document.getElementById('streak').textContent = `${streakDisplay || '○'} ${this.stats.streak.current || 0}`;
        }
    }

    async rate(exerciseId, rating) {
        try {
            const response = await this.api.rateExercise(exerciseId, rating);

            const exercise = this.allExercises.find(ex => ex.id === exerciseId);
            if (exercise && response.exercise) {
                Object.assign(exercise, response.exercise);
            }

            this.showSuccess(response.message);
            this.updateStats();
            this.renderExerciseDetail();
            this.renderExerciseList();
        } catch (error) {
            this.showError('Erreur notation: ' + error.message);
        }
    }

    showError(message) {
        const detail = document.getElementById('exerciseDetail');
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error';
        errorDiv.textContent = message;
        detail.insertBefore(errorDiv, detail.firstChild);
        setTimeout(() => errorDiv.remove(), 5000);
    }

    showSuccess(message) {
        const detail = document.getElementById('exerciseDetail');
        const successDiv = document.createElement('div');
        successDiv.className = 'success-message';
        successDiv.textContent = message;
        detail.insertBefore(successDiv, detail.firstChild);
        setTimeout(() => successDiv.remove(), 3000);
    }

    render() {
        this.renderExerciseList();
    }
}

// ============= INIT =============
let app;
document.addEventListener('DOMContentLoaded', () => {
    app = new MaestroApp();
});
