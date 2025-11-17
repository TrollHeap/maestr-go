// Data Management (Using in-memory storage instead of localStorage)
let appData = {
  exercises: [
    {
      id: 'go-001',
      title: 'Goroutines Basics',
      description: 'Learn how goroutines work in Go - lightweight threads managed by Go runtime',
      domain: 'golang',
      difficulty: 1,
      steps: [
        'Create a simple goroutine',
        'Use sync.WaitGroup to coordinate',
        'Understand goroutine scheduling'
      ],
      completed: false,
      last_reviewed: null,
      ease_factor: 2.5,
      interval_days: 0,
      repetitions: 0,
      completed_steps: []
    },
    {
      id: 'go-002',
      title: 'Channels and Communication',
      description: 'Master channel patterns for safe goroutine communication',
      domain: 'golang',
      difficulty: 2,
      steps: [
        'Create channel types',
        'Send and receive operations',
        'Implement producer-consumer pattern'
      ],
      completed: true,
      last_reviewed: '2025-11-15',
      ease_factor: 2.6,
      interval_days: 3,
      repetitions: 2,
      completed_steps: [0, 1, 2]
    },
    {
      id: 'linux-001',
      title: 'Tmux Window Management',
      description: 'Master tmux for efficient terminal multiplexing',
      domain: 'linux',
      difficulty: 1,
      steps: [
        'Create and name sessions',
        'Split windows and panes',
        'Navigate between windows'
      ],
      completed: false,
      last_reviewed: null,
      ease_factor: 2.5,
      interval_days: 0,
      repetitions: 0,
      completed_steps: []
    },
    {
      id: 'linux-002',
      title: 'Shell Scripting Fundamentals',
      description: 'Learn bash scripting for system automation',
      domain: 'linux',
      difficulty: 2,
      steps: [
        'Variables and conditionals',
        'Loops and functions',
        'Error handling with exit codes'
      ],
      completed: true,
      last_reviewed: '2025-11-10',
      ease_factor: 2.7,
      interval_days: 5,
      repetitions: 3,
      completed_steps: [0, 1, 2]
    },
    {
      id: 'arch-001',
      title: 'Understanding Memory Hierarchy',
      description: 'Visualize how CPU caches, RAM, and storage interact',
      domain: 'architecture',
      difficulty: 2,
      steps: [
        'Study cache levels (L1, L2, L3)',
        'Understand memory access patterns',
        'Optimize for cache locality'
      ],
      completed: false,
      last_reviewed: null,
      ease_factor: 2.5,
      interval_days: 0,
      repetitions: 0,
      completed_steps: []
    },
    {
      id: 'arch-002',
      title: 'Process vs Thread: Mental Model',
      description: 'Build clear mental model of process isolation and thread sharing',
      domain: 'architecture',
      difficulty: 1,
      steps: [
        'Draw process memory layout',
        'Understand thread memory sharing',
        'Visualize context switching'
      ],
      completed: true,
      last_reviewed: '2025-11-12',
      ease_factor: 2.8,
      interval_days: 4,
      repetitions: 2,
      completed_steps: [0, 1, 2]
    },
    {
      id: 'go-003',
      title: 'Interfaces and Polymorphism',
      description: 'Write flexible, reusable Go code with interfaces',
      domain: 'golang',
      difficulty: 2,
      steps: [
        'Define interface contracts',
        'Implement interfaces implicitly',
        'Use interface{} for flexibility'
      ],
      completed: false,
      last_reviewed: null,
      ease_factor: 2.5,
      interval_days: 0,
      repetitions: 0,
      completed_steps: []
    },
    {
      id: 'linux-003',
      title: 'File Permissions and Ownership',
      description: 'Master Linux file permission system',
      domain: 'linux',
      difficulty: 1,
      steps: [
        'Understand read/write/execute',
        'User/group/other model',
        'chmod and chown commands'
      ],
      completed: false,
      last_reviewed: null,
      ease_factor: 2.5,
      interval_days: 0,
      repetitions: 0,
      completed_steps: []
    },
    {
      id: 'arch-003',
      title: 'Virtual Memory System',
      description: 'Understand paging, segmentation, and translation lookaside buffers',
      domain: 'architecture',
      difficulty: 3,
      steps: [
        'Understand virtual address space',
        'Learn paging mechanism',
        'Study TLB optimization'
      ],
      completed: false,
      last_reviewed: null,
      ease_factor: 2.5,
      interval_days: 0,
      repetitions: 0,
      completed_steps: []
    },
    {
      id: 'go-004',
      title: 'Error Handling Patterns',
      description: 'Go error handling best practices and custom error types',
      domain: 'golang',
      difficulty: 1,
      steps: [
        'Understand error interface',
        'Create custom error types',
        'Handle errors gracefully'
      ],
      completed: true,
      last_reviewed: '2025-11-16',
      ease_factor: 2.9,
      interval_days: 1,
      repetitions: 4,
      completed_steps: [0, 1, 2]
    }
  ],
  userStats: {
    current_streak: 5,
    total_completed: 4,
    total_reviews: 9,
    domains: {
      golang: { completed: 2, total: 4, mastery: 50 },
      linux: { completed: 1, total: 3, mastery: 33 },
      architecture: { completed: 1, total: 3, mastery: 33 }
    }
  }
};

// Current state
let currentView = 'dashboard';
let currentExercise = null;
let sessionTimer = null;
let sessionTimeRemaining = 15 * 60; // 15 minutes in seconds
let selectedExerciseIndex = 0;

// Utility functions
function getDifficultyClass(difficulty) {
  if (difficulty === 1) return 'easy';
  if (difficulty === 2) return 'medium';
  return 'hard';
}

function getDifficultySymbol(difficulty) {
  if (difficulty === 1) return '●';
  if (difficulty === 2) return '●●';
  return '●●●';
}

function getExerciseStatus(exercise) {
  if (exercise.completed) {
    const dueDate = getDueDate(exercise);
    if (dueDate && new Date(dueDate) <= new Date()) {
      return { symbol: '⏱', text: 'Due for review', class: 'due' };
    }
    return { symbol: '✓', text: 'Completed', class: 'complete' };
  }
  if (exercise.completed_steps && exercise.completed_steps.length > 0) {
    return { symbol: '◐', text: 'In progress', class: 'progress' };
  }
  return { symbol: '○', text: 'Not started', class: 'progress' };
}

function getDueDate(exercise) {
  if (!exercise.last_reviewed) return null;
  const lastReview = new Date(exercise.last_reviewed);
  lastReview.setDate(lastReview.getDate() + exercise.interval_days);
  return lastReview;
}

function getDueDateText(exercise) {
  const dueDate = getDueDate(exercise);
  if (!dueDate) return 'Not reviewed yet';
  
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  dueDate.setHours(0, 0, 0, 0);
  
  const diffTime = dueDate - today;
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  
  if (diffDays < 0) return `Overdue by ${Math.abs(diffDays)} days`;
  if (diffDays === 0) return 'Due today';
  if (diffDays === 1) return 'Due tomorrow';
  return `Due in ${diffDays} days`;
}

function getRecommendedExercises() {
  // Get exercises that are due or not started
  const due = appData.exercises.filter(ex => {
    if (!ex.completed) return true;
    const dueDate = getDueDate(ex);
    return dueDate && new Date(dueDate) <= new Date();
  });
  
  // Sort by priority: easier first, then by due date
  return due.sort((a, b) => a.difficulty - b.difficulty).slice(0, 3);
}

function formatTime(seconds) {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
}

function updateStreak() {
  const streak = appData.userStats.current_streak;
  document.getElementById('streak-value').textContent = `${streak} days`;
  document.getElementById('streak-visual').textContent = '✓'.repeat(Math.min(streak, 10));
  document.getElementById('stats-streak').textContent = streak;
}

function updateDashboard() {
  updateStreak();
  
  // Update today's progress
  const completed = appData.userStats.total_completed;
  const total = appData.exercises.length;
  document.getElementById('today-progress').textContent = `${completed}/${total} exercises`;
  const progressPercent = (completed / total) * 100;
  document.getElementById('daily-progress-bar').style.width = `${progressPercent}%`;
  
  // Update recommended exercise
  const recommended = getRecommendedExercises()[0];
  if (recommended) {
    const recCard = document.getElementById('recommended-exercise');
    const status = getExerciseStatus(recommended);
    recCard.innerHTML = `
      <div class="exercise-header">
        <span class="exercise-domain ${recommended.domain}">${recommended.domain}</span>
        <span class="exercise-difficulty ${getDifficultyClass(recommended.difficulty)}">${getDifficultySymbol(recommended.difficulty)}</span>
        <span class="exercise-title">${recommended.title}</span>
      </div>
      <div class="exercise-desc">${recommended.description}</div>
      <div class="exercise-status">Status: ${status.symbol} ${status.text}</div>
    `;
  }
  
  // Update mastery tree
  Object.keys(appData.userStats.domains).forEach(domain => {
    const stats = appData.userStats.domains[domain];
    const mastery = (stats.completed / stats.total) * 100;
    // Update is already in HTML with inline styles
  });
}

function renderExerciseList(filter = 'all') {
  const list = document.getElementById('exercise-list');
  const filtered = filter === 'all' 
    ? appData.exercises 
    : appData.exercises.filter(ex => ex.domain === filter);
  
  list.innerHTML = filtered.map((ex, index) => {
    const status = getExerciseStatus(ex);
    return `
      <div class="exercise-item" data-id="${ex.id}" data-index="${index}">
        <div class="exercise-header">
          <span class="exercise-domain ${ex.domain}">${ex.domain}</span>
          <span class="exercise-difficulty ${getDifficultyClass(ex.difficulty)}">${getDifficultySymbol(ex.difficulty)}</span>
          <span class="exercise-title">${ex.title}</span>
        </div>
        <div class="exercise-desc">${ex.description}</div>
        <div class="exercise-item-status">
          <span class="status-badge ${status.class}">${status.symbol} ${status.text}</span>
          ${ex.completed ? `<span>${getDueDateText(ex)}</span>` : ''}
        </div>
      </div>
    `;
  }).join('');
  
  // Add click handlers
  list.querySelectorAll('.exercise-item').forEach(item => {
    item.addEventListener('click', () => {
      const id = item.dataset.id;
      const exercise = appData.exercises.find(ex => ex.id === id);
      showExerciseView(exercise);
    });
  });
}

function showExerciseView(exercise) {
  currentExercise = exercise;
  currentView = 'exercise';
  switchView('exercise-view');
  
  // Update exercise details
  document.getElementById('detail-domain').textContent = exercise.domain;
  document.getElementById('detail-domain').className = `exercise-domain ${exercise.domain}`;
  document.getElementById('detail-difficulty').textContent = getDifficultySymbol(exercise.difficulty);
  document.getElementById('detail-difficulty').className = `exercise-difficulty ${getDifficultyClass(exercise.difficulty)}`;
  document.getElementById('detail-title').textContent = exercise.title;
  document.getElementById('detail-description').textContent = exercise.description;
  
  // Render steps
  const stepsContainer = document.getElementById('steps-container');
  stepsContainer.innerHTML = exercise.steps.map((step, index) => {
    const completed = exercise.completed_steps && exercise.completed_steps.includes(index);
    return `
      <div class="step-item ${completed ? 'completed' : ''}" data-index="${index}">
        <div class="step-checkbox">${completed ? '✓' : ''}</div>
        <div class="step-text">${step}</div>
      </div>
    `;
  }).join('');
  
  // Update progress
  updateStepProgress();
  
  // Add click handlers for steps
  stepsContainer.querySelectorAll('.step-item').forEach(item => {
    item.addEventListener('click', () => {
      toggleStep(parseInt(item.dataset.index));
    });
  });
  
  // Start session timer
  startSessionTimer();
}

function toggleStep(stepIndex) {
  if (!currentExercise.completed_steps) {
    currentExercise.completed_steps = [];
  }
  
  const index = currentExercise.completed_steps.indexOf(stepIndex);
  if (index > -1) {
    currentExercise.completed_steps.splice(index, 1);
  } else {
    currentExercise.completed_steps.push(stepIndex);
  }
  
  // Re-render steps
  const stepsContainer = document.getElementById('steps-container');
  const stepItems = stepsContainer.querySelectorAll('.step-item');
  stepItems[stepIndex].classList.toggle('completed');
  const checkbox = stepItems[stepIndex].querySelector('.step-checkbox');
  checkbox.textContent = currentExercise.completed_steps.includes(stepIndex) ? '✓' : '';
  
  updateStepProgress();
  updateEncouragement();
}

function updateStepProgress() {
  const completed = currentExercise.completed_steps ? currentExercise.completed_steps.length : 0;
  const total = currentExercise.steps.length;
  const percent = (completed / total) * 100;
  
  document.getElementById('step-progress-bar').style.width = `${percent}%`;
  document.getElementById('step-progress-text').textContent = `${completed}/${total} steps completed`;
}

function updateEncouragement() {
  const completed = currentExercise.completed_steps ? currentExercise.completed_steps.length : 0;
  const total = currentExercise.steps.length;
  const messages = [
    "💪 Let's do this! Start with step 1.",
    "✓ Good start! Keep the momentum!",
    "🔥 You're on fire! Almost there!",
    "🎉 Excellent work! One more step!",
    "🌟 Perfect! You've completed all steps!"
  ];
  
  const index = Math.min(completed, messages.length - 1);
  document.getElementById('encouragement').textContent = messages[index];
}

function startSessionTimer() {
  sessionTimeRemaining = 15 * 60;
  updateSessionTimer();
  
  if (sessionTimer) clearInterval(sessionTimer);
  sessionTimer = setInterval(() => {
    sessionTimeRemaining--;
    updateSessionTimer();
    
    if (sessionTimeRemaining <= 0) {
      clearInterval(sessionTimer);
      alert('⏰ Session complete! Great work!');
    }
  }, 1000);
}

function updateSessionTimer() {
  document.getElementById('session-time').textContent = formatTime(sessionTimeRemaining);
}

function rateExercise(rating) {
  // SM-2 Algorithm
  const exercise = currentExercise;
  exercise.repetitions++;
  
  if (rating >= 3) {
    if (exercise.repetitions === 1) {
      exercise.interval_days = 1;
    } else if (exercise.repetitions === 2) {
      exercise.interval_days = 6;
    } else {
      exercise.interval_days = Math.round(exercise.interval_days * exercise.ease_factor);
    }
    
    exercise.ease_factor = exercise.ease_factor + (0.1 - (5 - rating) * (0.08 + (5 - rating) * 0.02));
  } else {
    exercise.repetitions = 0;
    exercise.interval_days = 1;
  }
  
  exercise.ease_factor = Math.max(1.3, exercise.ease_factor);
  exercise.last_reviewed = new Date().toISOString().split('T')[0];
  
  // Mark as completed if all steps done
  if (exercise.completed_steps && exercise.completed_steps.length === exercise.steps.length) {
    if (!exercise.completed) {
      exercise.completed = true;
      appData.userStats.total_completed++;
      appData.userStats.domains[exercise.domain].completed++;
      appData.userStats.domains[exercise.domain].mastery = 
        Math.round((appData.userStats.domains[exercise.domain].completed / appData.userStats.domains[exercise.domain].total) * 100);
    }
  }
  
  appData.userStats.total_reviews++;
  
  // Show success message
  const messages = [
    "Keep practicing! You'll get there! 💪",
    "Good effort! Review again soon. 📚",
    "Great job! See you next time! ⭐",
    "Perfect! You've mastered this! 🎉"
  ];
  alert(messages[rating - 1]);
  
  // Return to dashboard
  switchView('dashboard-view');
  updateDashboard();
}

function renderStats() {
  document.getElementById('total-completed').textContent = appData.userStats.total_completed;
  document.getElementById('total-reviews').textContent = appData.userStats.total_reviews;
  
  const totalExercises = appData.exercises.length;
  const overallPercent = Math.round((appData.userStats.total_completed / totalExercises) * 100);
  document.getElementById('overall-progress').textContent = `${overallPercent}%`;
  
  // Render review schedule
  const reviewList = document.getElementById('review-list');
  const dueExercises = appData.exercises
    .filter(ex => ex.completed)
    .map(ex => ({ ...ex, dueDate: getDueDate(ex), dueText: getDueDateText(ex) }))
    .sort((a, b) => new Date(a.dueDate) - new Date(b.dueDate))
    .slice(0, 5);
  
  reviewList.innerHTML = dueExercises.map(ex => `
    <div class="review-item">
      <div class="review-info">
        <div class="exercise-title">${ex.title}</div>
        <div class="exercise-desc">${ex.domain}</div>
      </div>
      <div class="review-due">${ex.dueText}</div>
    </div>
  `).join('');
  
  // Render domain breakdown
  const domainBreakdown = document.getElementById('domain-breakdown');
  domainBreakdown.innerHTML = Object.keys(appData.userStats.domains).map(domain => {
    const stats = appData.userStats.domains[domain];
    return `
      <div class="domain-stat">
        <div class="domain-stat-header">
          ${domain === 'golang' ? '🌿' : domain === 'linux' ? '🌳' : '🏛️'} ${domain}
        </div>
        <div class="mastery-bar-container">
          <div class="mastery-bar" style="width: ${stats.mastery}%"></div>
        </div>
        <div class="domain-stat-details">
          <span>Completed: ${stats.completed}/${stats.total}</span>
          <span>Mastery: ${stats.mastery}%</span>
        </div>
      </div>
    `;
  }).join('');
}

function showQuickStart() {
  currentView = 'quick-start';
  switchView('quick-start-view');
  
  const recommended = getRecommendedExercises();
  const quickExercises = document.getElementById('quick-exercises');
  
  quickExercises.innerHTML = recommended.map((ex, index) => {
    const status = getExerciseStatus(ex);
    return `
      <div class="quick-exercise-card" data-id="${ex.id}">
        <div class="quick-number">${index + 1}</div>
        <div class="exercise-header">
          <span class="exercise-domain ${ex.domain}">${ex.domain}</span>
          <span class="exercise-difficulty ${getDifficultyClass(ex.difficulty)}">${getDifficultySymbol(ex.difficulty)}</span>
        </div>
        <div class="exercise-title" style="margin-top: 12px; margin-bottom: 8px;">${ex.title}</div>
        <div class="exercise-desc">${ex.description}</div>
        <div class="exercise-status" style="margin-top: 12px;">${status.symbol} ${status.text}</div>
      </div>
    `;
  }).join('');
  
  // Add click handlers
  quickExercises.querySelectorAll('.quick-exercise-card').forEach(card => {
    card.addEventListener('click', () => {
      const id = card.dataset.id;
      const exercise = appData.exercises.find(ex => ex.id === id);
      showExerciseView(exercise);
    });
  });
  
  // Start timer
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
}

function updateQuickTimer() {
  document.getElementById('quick-timer').textContent = formatTime(sessionTimeRemaining);
}

function switchView(viewId) {
  document.querySelectorAll('.view').forEach(view => view.classList.remove('active'));
  document.getElementById(viewId).classList.add('active');
  
  if (viewId === 'dashboard-view') {
    if (sessionTimer) clearInterval(sessionTimer);
    updateDashboard();
  }
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
  // Initialize dashboard
  updateDashboard();
  renderExerciseList();
  
  // Filter buttons
  document.querySelectorAll('.filter-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
      btn.classList.add('active');
      renderExerciseList(btn.dataset.domain);
    });
  });
  
  // Rating buttons
  document.querySelectorAll('.rating-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      const rating = parseInt(btn.dataset.rating);
      rateExercise(rating);
    });
  });
  
  // Keyboard navigation
  document.addEventListener('keydown', (e) => {
    const key = e.key.toLowerCase();
    
    // Global shortcuts
    if (key === 'escape') {
      if (currentView !== 'dashboard') {
        if (sessionTimer) clearInterval(sessionTimer);
        currentView = 'dashboard';
        switchView('dashboard-view');
      }
      return;
    }
    
    // Dashboard shortcuts
    if (currentView === 'dashboard' || document.getElementById('dashboard-view').classList.contains('active')) {
      if (key === 'q') {
        showQuickStart();
      } else if (key === 'b') {
        currentView = 'browse';
        switchView('browse-view');
      } else if (key === 's') {
        currentView = 'stats';
        switchView('stats-view');
        renderStats();
      }
    }
    
    // Quick start shortcuts
    if (currentView === 'quick-start') {
      if (key >= '1' && key <= '3') {
        const index = parseInt(key) - 1;
        const recommended = getRecommendedExercises();
        if (recommended[index]) {
          showExerciseView(recommended[index]);
        }
      }
    }
    
    // Exercise view shortcuts
    if (currentView === 'exercise') {
      if (key === ' ') {
        e.preventDefault();
        const firstIncomplete = currentExercise.steps.findIndex((_, index) => 
          !currentExercise.completed_steps || !currentExercise.completed_steps.includes(index)
        );
        if (firstIncomplete !== -1) {
          toggleStep(firstIncomplete);
        }
      } else if (key >= '1' && key <= '4') {
        const rating = parseInt(key);
        rateExercise(rating);
      }
    }
    
    // Browse view shortcuts
    if (currentView === 'browse') {
      const items = document.querySelectorAll('.exercise-item');
      if (key === 'j' || key === 'arrowdown') {
        e.preventDefault();
        selectedExerciseIndex = Math.min(selectedExerciseIndex + 1, items.length - 1);
        items[selectedExerciseIndex]?.scrollIntoView({ behavior: 'smooth', block: 'center' });
      } else if (key === 'k' || key === 'arrowup') {
        e.preventDefault();
        selectedExerciseIndex = Math.max(selectedExerciseIndex - 1, 0);
        items[selectedExerciseIndex]?.scrollIntoView({ behavior: 'smooth', block: 'center' });
      } else if (key === 'enter') {
        items[selectedExerciseIndex]?.click();
      }
    }
  });
});