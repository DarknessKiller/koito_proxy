package model

import "database/sql"

type Rule struct {
	BaseModel
	MatchTrackName      sql.NullString `gorm:"type:varchar(255)"`
	MatchArtistName     sql.NullString `gorm:"type:varchar(255)"`
	MatchReleaseName    sql.NullString `gorm:"type:varchar(255)"`
	MatchArtistNames    []string       `gorm:"serializer:json"`
	MatchDurationBucket int32          `gorm:"nullable"`
	MatchMBID           sql.NullString `gorm:"type:varchar(36)"`
	ReplaceTrackName    sql.NullString `gorm:"type:varchar(255)"`
	ReplaceArtistName   sql.NullString `gorm:"type:varchar(255)"`
	ReplaceReleaseName  sql.NullString `gorm:"type:varchar(255)"`
	ReplaceArtistNames  []string       `gorm:"serializer:json"`
	Enabled             *bool          `gorm:"not null;default:true"`
	Timestamps
}
