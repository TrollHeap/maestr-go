// assets/js/app.js - Application Logic

class RetroApp {
  constructor() {
    this.terminal = null;
    this.init();
  }

  init() {
    console.log('🟢 Retro App initialized');
    
    // HTMX event listeners
    this.setupHTMX();
    
    // Terminal setup
    this.setupTerminal();
    
    // Global event listeners
    this.setupGlobalEvents();
  }

  setupHTMX() {
    // On avant chaque requête HTMX
    document.addEventListener('htmx:beforeRequest', (event) => {
      console.log('📤 HTMX Request:', event.detail.xhr.responseURL);
    });

    // Après succès
    document.addEventListener('htmx:afterSwap', (event) => {
      console.log('✅ HTMX Swap successful');
      this.addGlowEffect(event.detail.target);
    });

    // Erreur
    document.addEventListener('htmx:responseError', (event) => {
      console.error('❌ HTMX Error:', event.detail.xhr.status);
    });
  }

  setupTerminal() {
    const terminalOutput = document.getElementById('terminal-output');
    const terminalForm = document.querySelector('form[hx-post]');

    if (terminalForm) {
      terminalForm.addEventListener('htmx:afterRequest', (event) => {
        // Auto-scroll vers le bas
        if (terminalOutput) {
          terminalOutput.scrollTop = 0;
        }
      });
    }

    // Focus sur l'input au chargement
    const terminalInput = document.querySelector('.terminal-input');
    if (terminalInput) {
      terminalInput.focus();
    }
  }

  setupGlobalEvents() {
    // Hover glow sur les boutons
    document.querySelectorAll('.btn-retro').forEach(btn => {
      btn.addEventListener('mouseenter', () => {
        btn.style.textShadow = '0 0 15px var(--retro-shadow)';
      });
      btn.addEventListener('mouseleave', () => {
        btn.style.textShadow = 'var(--text-glow)';
      });
    });

    // Glow sur les stat cards
    document.querySelectorAll('.stat-card').forEach(card => {
      card.addEventListener('mouseenter', () => {
        card.style.boxShadow = '0 0 20px var(--retro-shadow)';
      });
      card.addEventListener('mouseleave', () => {
        card.style.boxShadow = '0 0 8px var(--retro-shadow)';
      });
    });

    // Keyboard shortcuts
    document.addEventListener('keydown', (e) => {
      // Ctrl/Cmd + K = focus terminal
      if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        const input = document.querySelector('.terminal-input');
        if (input) input.focus();
      }

      // Ctrl/Cmd + L = clear (pas vraiment pour sécurité)
      if ((e.ctrlKey || e.metaKey) && e.key === 'l') {
        e.preventDefault();
        console.log('Clear request intercepted');
      }
    });
  }

  addGlowEffect(element) {
    element.style.animation = 'fadeIn 0.3s ease-out';
  }
}

// Statistiques live (optionnel - update toutes les 30s)
class StatsUpdater {
  constructor() {
    this.interval = null;
    this.init();
  }

  init() {
    // Optionnel: ne pas activer par défaut
    // this.start();
  }

  start() {
    this.interval = setInterval(() => {
      this.updateStats();
    }, 30000); // Chaque 30 secondes
  }

  stop() {
    if (this.interval) clearInterval(this.interval);
  }

  async updateStats() {
    try {
      const response = await fetch('/api/stats');
      const data = await response.json();
      console.log('📊 Stats updated:', data);
      // Tu peux mettre à jour les éléments DOM ici
    } catch (error) {
      console.error('Stats update failed:', error);
    }
  }
}

// Terminal interactif (optionnel - augmente les capacités)
class Terminal {
  constructor() {
    this.history = [];
    this.historyIndex = -1;
    this.setupTerminal();
  }

  setupTerminal() {
    const input = document.querySelector('.terminal-input');
    if (input) {
      // Arrow up/down: history navigation
      input.addEventListener('keydown', (e) => {
        if (e.key === 'ArrowUp') {
          e.preventDefault();
          this.showHistoryPrev();
        } else if (e.key === 'ArrowDown') {
          e.preventDefault();
          this.showHistoryNext();
        }
      });
    }
  }

  addToHistory(command) {
    if (command.trim()) {
      this.history.push(command);
      this.historyIndex = -1;
    }
  }

  showHistoryPrev() {
    if (this.history.length === 0) return;
    if (this.historyIndex < this.history.length - 1) {
      this.historyIndex++;
      this.setInputValue(this.history[this.history.length - 1 - this.historyIndex]);
    }
  }

  showHistoryNext() {
    if (this.historyIndex > 0) {
      this.historyIndex--;
      this.setInputValue(this.history[this.history.length - 1 - this.historyIndex]);
    } else if (this.historyIndex === 0) {
      this.historyIndex = -1;
      this.setInputValue('');
    }
  }

  setInputValue(value) {
    const input = document.querySelector('.terminal-input');
    if (input) {
      input.value = value;
      input.focus();
    }
  }
}

// Notification système (optionnel)
class Notifier {
  static success(message) {
    console.log('✅', message);
    this.show(message, 'success');
  }

  static error(message) {
    console.error('❌', message);
    this.show(message, 'error');
  }

  static info(message) {
    console.log('ℹ️', message);
    this.show(message, 'info');
  }

  static show(message, type) {
    // TODO: Implémenter un système de notification visuel
  }
}

// Init au chargement du DOM
document.addEventListener('DOMContentLoaded', () => {
  const app = new RetroApp();
  const terminal = new Terminal();
  // const stats = new StatsUpdater();

  // Expose globalement pour debugging
  window.retroApp = app;
  window.terminal = terminal;
  window.notifier = Notifier;

  console.log('🎮 Retro Terminal ready!');
  console.log('Shortcuts: Ctrl/Cmd+K to focus terminal');
});

// Service Worker basic (optional)
if ('serviceWorker' in navigator) {
  // navigator.serviceWorker.register('/sw.js').catch(() => {});
}
