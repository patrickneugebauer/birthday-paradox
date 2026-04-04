CREATE TABLE sessions (
    id INTEGER PRIMARY KEY,
    session_type TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    --
    exit_code INTEGER, -- NULL unhandled exception, 0 success, 1 handled exception
    finished_at TIMESTAMP,
    error TEXT
);

CREATE TABLE languages (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    session_id INTEGER NOT NULL REFERENCES sessions(id),
    created_at TIMESTAMP DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE TABLE dockerfile (
    id INTEGER PRIMARY KEY,
    filename TEXT NOT NULL UNIQUE,
    language_id INTEGER NOT NULL REFERENCES languages(id),
    session_id INTEGER NOT NULL REFERENCES sessions(id),
    created_at TIMESTAMP DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE TABLE images (
    id INTEGER PRIMARY KEY,
    tag TEXT NOT NULL,
    dockerfile_id INTEGER NOT NULL UNIQUE REFERENCES dockerfile(id),
    session_id INTEGER NOT NULL REFERENCES sessions(id),
    created_at TIMESTAMP DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE TABLE image_size (
    id INTEGER PRIMARY KEY,
    size TEXT NOT NULL,
    image_id INTEGER NOT NULL UNIQUE REFERENCES images(id),
    session_id INTEGER NOT NULL REFERENCES sessions(id),
    created_at TIMESTAMP DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
);

CREATE TABLE runs (
    id INTEGER PRIMARY KEY,
    tag TEXT NOT NULL,
    iterations INTEGER NOT NULL,
    sample_size INTEGER NOT NULL,
    percent REAL NOT NULL,
    seconds REAL NOT NULL,
    iterations_per_second INTEGER NOT NULL,
    image_id INTEGER NOT NULL REFERENCES images(id),
    session_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    -- columns come first
    UNIQUE(image_id, session_id)
);
