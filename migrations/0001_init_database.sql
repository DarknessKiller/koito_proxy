-- +goose Up
CREATE TABLE IF NOT EXISTS overwrite_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    match_track_name TEXT NULL,
    match_artist_name TEXT NULL,
    match_release_name TEXT NULL,

    replace_track_name TEXT NULL,
    replace_artist_name TEXT NULL,
    replace_release_name TEXT NULL,

    enabled INTEGER NOT NULL DEFAULT 1,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE overwrite_rules;