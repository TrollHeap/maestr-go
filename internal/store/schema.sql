-- ============================================
-- MAESTRO GO v2 - DATABASE SCHEMA (Dates YYYYMMDD)
-- ============================================

CREATE TABLE IF NOT EXISTS exercises (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    domain TEXT NOT NULL,
    difficulty INTEGER CHECK(difficulty BETWEEN 1 AND 5),
    
    -- Contenu pédagogique
    content TEXT,
    mnemonic TEXT,
    conceptual_visuals TEXT,
    
    -- Steps (JSON arrays)
    steps TEXT,
    completed_steps TEXT,
    
    -- SRS (Spaced Repetition System) - DATES UNIQUEMENT
    done BOOLEAN DEFAULT 0,
    last_reviewed_date INTEGER,
    next_review_date INTEGER NOT NULL DEFAULT 0,
    ease_factor REAL DEFAULT 2.5 CHECK(ease_factor >= 1.3),
    interval_days INTEGER DEFAULT 1,
    repetitions INTEGER DEFAULT 0,
    skipped_count INTEGER DEFAULT 0,
    last_skipped_date INTEGER,
    
    -- Soft delete
    deleted BOOLEAN DEFAULT 0,
    deleted_at INTEGER,
    
    -- Metadata
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

-- Index optimisés
CREATE INDEX IF NOT EXISTS idx_next_review ON exercises(next_review_date, done, deleted);
CREATE INDEX IF NOT EXISTS idx_domain ON exercises(domain, deleted);
CREATE INDEX IF NOT EXISTS idx_done ON exercises(done) WHERE deleted = 0;
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_title ON exercises(title) WHERE deleted = 0;

-- ============================================
-- TABLE : SESSIONS
-- ============================================
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    started_at INTEGER NOT NULL,
    ended_at INTEGER,
    energy_level TEXT CHECK(energy_level IN ('low', 'medium', 'high')),
    mode TEXT CHECK(mode IN ('micro', 'standard', 'deep')),
    completed_count INTEGER DEFAULT 0,
    duration_min INTEGER,
    created_at INTEGER NOT NULL DEFAULT (strftime('%Y%m%d', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_session_date ON sessions(started_at);

-- ============================================
-- TABLE : SESSION_EXERCISES
-- ============================================
CREATE TABLE IF NOT EXISTS session_exercises (
    session_id INTEGER NOT NULL,
    exercise_id INTEGER NOT NULL,
    position INTEGER NOT NULL,
    completed BOOLEAN DEFAULT 0,
    quality INTEGER CHECK(quality BETWEEN 0 AND 3),
    reviewed_at INTEGER,
    PRIMARY KEY (session_id, exercise_id),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_session_ex ON session_exercises(session_id);

-- ============================================
-- TABLE : PROGRESS_LOG
-- ============================================
CREATE TABLE IF NOT EXISTS progress_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    exercise_id INTEGER NOT NULL,
    reviewed_at INTEGER NOT NULL,
    quality INTEGER NOT NULL CHECK(quality BETWEEN 0 AND 3),
    ease_factor REAL,
    interval_days INTEGER,
    repetitions INTEGER,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_progress_exercise ON progress_log(exercise_id, reviewed_at DESC);

-- ============================================
-- TABLE : ANALYTICS
-- ============================================
CREATE TABLE IF NOT EXISTS analytics (
    id INTEGER PRIMARY KEY CHECK(id = 1),
    avg_session_length_min REAL DEFAULT 0,
    best_time_slot TEXT,
    difficulty_success_rate TEXT DEFAULT '{"1":0.95,"2":0.85,"3":0.70,"4":0.50,"5":0.30}',
    current_streak INTEGER DEFAULT 0,
    longest_streak INTEGER DEFAULT 0,
    last_session_date INTEGER,
    total_sessions INTEGER DEFAULT 0,
    total_exercises_done INTEGER DEFAULT 0,
    updated_at INTEGER NOT NULL DEFAULT (strftime('%Y%m%d', 'now'))
);

INSERT OR IGNORE INTO analytics (id, updated_at) VALUES (1, strftime('%Y%m%d', 'now'));

-- ============================================
-- TABLE : SETTINGS
-- ============================================
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at INTEGER NOT NULL DEFAULT (strftime('%Y%m%d', 'now'))
);

INSERT OR IGNORE INTO settings (key, value) VALUES 
    ('theme', 'dark'),
    ('session_reminder', 'true'),
    ('default_energy', 'medium'),
    ('ascii_visuals_enabled', 'true');

-- ============================================
-- TRIGGERS
-- ============================================
CREATE TRIGGER IF NOT EXISTS update_exercise_timestamp
AFTER UPDATE ON exercises
FOR EACH ROW
BEGIN
    UPDATE exercises SET updated_at = strftime('%Y%m%d', 'now') WHERE id = OLD.id;
END;
