package rules

import (
	"context"
	"database/sql"
	"log/slog"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Load(ctx context.Context) ([]Rule, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			id,
			match_track_name,
			match_artist_name,
			match_release_name,
			replace_track_name,
			replace_artist_name,
			replace_release_name
		FROM overwrite_rules
		WHERE enabled = 1
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Rule

	for rows.Next() {
		var r Rule
		if err := rows.Scan(
			&r.ID,
			&r.MatchTrackName,
			&r.MatchArtistName,
			&r.MatchReleaseName,
			&r.ReplaceTrackName,
			&r.ReplaceArtistName,
			&r.ReplaceReleaseName,
		); err != nil {
			return nil, err
		}

		out = append(out, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Store) Add(ctx context.Context, r Rule) error {
	existing, err := s.Load(ctx)
	if err != nil {
		return err
	}

	for _, ex := range existing {
		if ruleEqual(ex, r) {
			slog.Warn("duplicate rule detected; skipping insert")
			return nil
		}
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO overwrite_rules (
			match_track_name,
			match_artist_name,
			match_release_name,
			replace_track_name,
			replace_artist_name,
			replace_release_name,
			enabled
		) VALUES (?, ?, ?, ?, ?, ?, 1)
	`, nullableString(r.MatchTrackName), nullableString(r.MatchArtistName), nullableString(r.MatchReleaseName), nullableString(r.ReplaceTrackName), nullableString(r.ReplaceArtistName), nullableString(r.ReplaceReleaseName))

	return err
}

func ruleEqual(a, b Rule) bool {
	return nullStringEqual(a.MatchTrackName, b.MatchTrackName) &&
		nullStringEqual(a.MatchArtistName, b.MatchArtistName) &&
		nullStringEqual(a.MatchReleaseName, b.MatchReleaseName) &&
		nullStringEqual(a.ReplaceTrackName, b.ReplaceTrackName) &&
		nullStringEqual(a.ReplaceArtistName, b.ReplaceArtistName) &&
		nullStringEqual(a.ReplaceReleaseName, b.ReplaceReleaseName)
}

func nullStringEqual(a, b sql.NullString) bool {
	if a.Valid != b.Valid {
		return false
	}
	if !a.Valid {
		return true
	}
	return a.String == b.String
}

func nullableString(s sql.NullString) interface{} {
	if s.Valid {
		return s.String
	}
	return nil
}
