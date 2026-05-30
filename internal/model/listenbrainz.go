package model

type ListenBrainzTrackMetaData struct {
	TrackName      string                     `json:"track_name"`
	ArtistName     string                     `json:"artist_name"`
	ReleaseName    string                     `json:"release_name,omitempty"`
	AdditionalInfo ListenBrainzAdditionalInfo `json:"additional_info"`
}

type ListenBrainzAdditionalInfo struct {
	CustomTrackID int64             `json:"custom_track_id"`
	ArtistIds     []int64           `json:"artist_ids"`
	ArtistNames   []string          `json:"artist_names"`
	MusicbrainzID int64             `json:"musicbrainz_id"`
	ListenCount   int64             `json:"listen_count"`
	Duration      int64             `json:"duration"`
	AlbumID       int64             `json:"album_id"`
	FirstListen   int64             `json:"first_listen"`
	AllTimeRank   int64             `json:"all_time_rank"`
	Image         ListenBrainzImage `json:"image"`
}

type ListenBrainzImage struct {
	XS     string `json:"xs"`
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
	XL     string `json:"xl"`
}

type ListenBrainzPayload struct {
	ListenedAt    int64                     `json:"listened_at"`
	TrackMetaData ListenBrainzTrackMetaData `json:"track_metadata"`
}

type ListenBrainzSubmitRequest struct {
	ListenType string                `json:"listen_type"`
	Payload    []ListenBrainzPayload `json:"payload"`
}
