package model

type LbzListenType string

const (
	ListenTypeSingle     LbzListenType = "single"
	ListenTypePlayingNow LbzListenType = "playing_now"
	ListenTypeImport     LbzListenType = "import"
)

type ListenBrainzSubmitRequest struct {
	ListenType LbzListenType               `json:"listen_type,omitempty"`
	Payload    []ListenBrainzSubmitPayload `json:"payload,omitempty"`
}

type ListenBrainzSubmitPayload struct {
	ListenedAt    int64                     `json:"listened_at,omitempty"`
	TrackMetaData ListenBrainzTrackMetaData `json:"track_metadata"`
}

type ListenBrainzTrackMetaData struct {
	ArtistName     string                     `json:"artist_name"` // required
	TrackName      string                     `json:"track_name"`  // required
	ReleaseName    string                     `json:"release_name,omitempty"`
	MBIDMapping    ListenBrainzMBIDMapping    `json:"mbid_mapping,omitempty"`
	AdditionalInfo ListenBrainzAdditionalInfo `json:"additional_info,omitempty"`
}
type ListenBrainzArtist struct {
	ArtistMBID string `json:"artist_mbid,omitempty"`
	ArtistName string `json:"artist_credit_name,omitempty"`
}
type ListenBrainzMBIDMapping struct {
	ReleaseMBID   string               `json:"release_mbid,omitempty"`
	RecordingMBID string               `json:"recording_mbid,omitempty"`
	ArtistMBIDs   []string             `json:"artist_mbids,omitempty"`
	Artists       []ListenBrainzArtist `json:"artists,omitempty"`
}

type ListenBrainzAdditionalInfo struct {
	MediaPlayer             string   `json:"media_player,omitempty"`
	SubmissionClient        string   `json:"submission_client,omitempty"`
	SubmissionClientVersion string   `json:"submission_client_version,omitempty"`
	ReleaseMBID             string   `json:"release_mbid,omitempty"`
	ReleaseGroupMBID        string   `json:"release_group_mbid,omitempty"`
	ArtistMBIDs             []string `json:"artist_mbids,omitempty"`
	ArtistNames             []string `json:"artist_names,omitempty"`
	RecordingMBID           string   `json:"recording_mbid,omitempty"`
	DurationMs              int32    `json:"duration_ms,omitempty"`
	Duration                int32    `json:"duration,omitempty"`
	Tags                    []string `json:"tags,omitempty"`
	AlbumArtist             string   `json:"albumartist,omitempty"`
}
