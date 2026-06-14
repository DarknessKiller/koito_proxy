package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"koito_proxy/internal/config"
	"koito_proxy/internal/model"
)

type KoitoService interface {
	AddMergeRule(ctx context.Context, entity, targetID string, sourceID int64) error
}

type koitoService struct {
	ruleService RuleService
	config      *config.Config
}

func NewKoitoService(rs RuleService, cfg *config.Config) KoitoService {
	return &koitoService{
		ruleService: rs,
		config:      cfg,
	}
}

type koitoTrack struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Artists []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"artists"`
	MusicbrainzID interface{} `json:"musicbrainz_id"`
	ListenCount   int64       `json:"listen_count"`
	Duration      int64       `json:"duration"`
	Image         struct {
		Xs     string `json:"xs"`
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
		Xl     string `json:"xl"`
	} `json:"image"`
	AlbumID      int64 `json:"album_id"`
	TimeListened int64 `json:"time_listened"`
	FirstListen  int64 `json:"first_listen"`
	AllTimeRank  int64 `json:"all_time_rank"`
}

type koitoArtist struct {
	ID            int64    `json:"id"`
	MusicbrainzID any      `json:"musicbrainz_id"`
	Name          string   `json:"name"`
	Aliases       []string `json:"aliases"`
	Image         struct {
		Xs     string `json:"xs"`
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
		Xl     string `json:"xl"`
	} `json:"image"`
	ListenCount  int64 `json:"listen_count"`
	TimeListened int64 `json:"time_listened"`
	FirstListen  int64 `json:"first_listen"`
	AllTimeRank  int64 `json:"all_time_rank"`
}

type koitoAlbum struct {
	ID            int64  `json:"id"`
	MusicbrainzID any    `json:"musicbrainz_id"`
	Title         string `json:"title"`
	Image         struct {
		Xs     string `json:"xs"`
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
		Xl     string `json:"xl"`
	} `json:"image"`
	Artists []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"artists"`
	IsVariousArtists bool  `json:"is_various_artists"`
	ListenCount      int64 `json:"listen_count"`
	TimeListened     int64 `json:"time_listened"`
	FirstListen      int64 `json:"first_listen"`
	AllTimeRank      int64 `json:"all_time_rank"`
}

func (s *koitoService) AddMergeRule(ctx context.Context, entity, targetID string, sourceID int64) error {
	slog.Info("koito merge request detected", "entity", entity, "target_id", targetID, "source_id", sourceID)

	switch entity {
	case "track":
		source, err := s.fetchTrack(ctx, sourceID)
		if err != nil {
			return fmt.Errorf("fetch source track: %w", err)
		}

		sourceAlbum, err := s.fetchAlbum(ctx, source.AlbumID)
		if err != nil {
			return fmt.Errorf("fetch release: %w", err)
		}

		target, err := s.fetchTrackString(ctx, targetID)
		if err != nil {
			return fmt.Errorf("fetch target track: %w", err)
		}

		targetAlbum, err := s.fetchAlbum(ctx, target.AlbumID)
		if err != nil {
			return fmt.Errorf("fetch release: %w", err)
		}

		if len(source.Artists) == 0 || len(target.Artists) == 0 {
			return fmt.Errorf("source or target artist is empty")
		}

		matchSourceArtists := []string{}
		for _, artist := range source.Artists {
			matchSourceArtists = append(matchSourceArtists, artist.Name)
		}

		matchReplaceArtists := []string{}
		for _, artist := range target.Artists {
			matchReplaceArtists = append(matchReplaceArtists, artist.Name)
		}

		rule := model.Rule{
			MatchTrackName:     newNullString(source.Title),
			MatchArtistName:    newNullString(matchSourceArtists[0]),
			MatchArtistNames:   matchSourceArtists,
			MatchReleaseName:   newNullString(sourceAlbum.Title),
			ReplaceTrackName:   newNullString(target.Title),
			ReplaceArtistName:  newNullString(matchReplaceArtists[0]),
			ReplaceArtistNames: matchReplaceArtists,
			ReplaceReleaseName: newNullString(targetAlbum.Title),
			Enabled:            new(true),
		}

		return s.ruleService.Create(ctx, &rule)

	case "artist":
		source, err := s.fetchArtist(ctx, sourceID)
		if err != nil {
			return fmt.Errorf("fetch source artist: %w", err)
		}
		target, err := s.fetchArtistString(ctx, targetID)
		if err != nil {
			return fmt.Errorf("fetch target artist: %w", err)
		}

		rule := model.Rule{
			MatchArtistName:   newNullString(source.Name),
			ReplaceArtistName: newNullString(target.Name),
			Enabled:           new(true),
		}

		return s.ruleService.Create(ctx, &rule)

	default:
		return nil
	}
}

func (s *koitoService) fetchTrack(ctx context.Context, id int64) (*koitoTrack, error) {
	return s.fetchTrackString(ctx, fmt.Sprintf("%d", id))
}

func (s *koitoService) fetchTrackString(ctx context.Context, id string) (*koitoTrack, error) {
	api := fmt.Sprintf("/apis/web/v1/track/%s", id)
	body, err := s.fetchUpstreamAPI(ctx, "GET", api, nil)
	if err != nil {
		return nil, err
	}

	var track koitoTrack
	if err := json.Unmarshal(body, &track); err != nil {
		return nil, err
	}

	if track.Title == "" || len(track.Artists) == 0 || track.Artists[0].Name == "" {
		return nil, fmt.Errorf("unexpected track payload: missing track or artist name")
	}

	return &track, nil
}

func (s *koitoService) fetchArtist(ctx context.Context, id int64) (*koitoArtist, error) {
	return s.fetchArtistString(ctx, fmt.Sprintf("%d", id))
}

func (s *koitoService) fetchArtistString(ctx context.Context, id string) (*koitoArtist, error) {
	api := fmt.Sprintf("/apis/web/v1/artist/%s", id)
	body, err := s.fetchUpstreamAPI(ctx, "GET", api, nil)
	if err != nil {
		return nil, err
	}

	var artist koitoArtist
	if err := json.Unmarshal(body, &artist); err != nil {
		return nil, err
	}

	if artist.Name == "" {
		return nil, fmt.Errorf("unexpected artist payload: missing artist name")
	}

	return &artist, nil
}

func (s *koitoService) fetchAlbum(ctx context.Context, id int64) (*koitoAlbum, error) {
	return s.fetchAlbumString(ctx, fmt.Sprintf("%d", id))
}

func (s *koitoService) fetchAlbumString(ctx context.Context, id string) (*koitoAlbum, error) {
	api := fmt.Sprintf("/apis/web/v1/album/%s", id)
	body, err := s.fetchUpstreamAPI(ctx, "GET", api, nil)
	if err != nil {
		return nil, err
	}

	var album koitoAlbum
	if err := json.Unmarshal(body, &album); err != nil {
		return nil, err
	}

	return &album, nil
}

func (s *koitoService) fetchUpstreamAPI(ctx context.Context, method string, apiPath string, body []byte) ([]byte, error) {
	base, err := url.Parse(s.config.UpstreamURL)
	if err != nil {
		return nil, err
	}
	rel, err := url.Parse(apiPath)
	if err != nil {
		return nil, err
	}
	target := base.ResolveReference(rel)

	proxyReq, err := http.NewRequestWithContext(ctx, method, target.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	proxyReq.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected upstream status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func newNullString(value string) sql.NullString {
	return sql.NullString{String: strings.TrimSpace(value), Valid: strings.TrimSpace(value) != ""}
}
