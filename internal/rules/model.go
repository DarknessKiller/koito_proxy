package rules

import "database/sql"

type Rule struct {
	ID int64

	MatchTrackName   sql.NullString
	MatchArtistName  sql.NullString
	MatchReleaseName sql.NullString

	ReplaceTrackName   sql.NullString
	ReplaceArtistName  sql.NullString
	ReplaceReleaseName sql.NullString
}
