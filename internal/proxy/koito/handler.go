package koito

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"koito_proxy/internal/config"
	"koito_proxy/internal/model"
	"koito_proxy/internal/repository"
	"koito_proxy/internal/response"
	"koito_proxy/internal/rules"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine *rules.RuleEngine
	rule   repository.Repository[model.Rule]
	config *config.Config
}

func NewHandler(e *rules.RuleEngine, rule repository.Repository[model.Rule], cfg *config.Config) *Handler {
	return &Handler{
		engine: e,
		rule:   rule,
		config: cfg,
	}
}

type mergeRequest struct {
	MergeFromID int64 `json:"merge_from_id"`
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

func (h *Handler) InterceptMerge(c *gin.Context) {
	entity := c.Param("entity")
	targetID := c.Param("id")

	var req mergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to parse merge request", "error", err, "path", c.Request.URL.Path)
		response.RespondBadRequest(c, response.ErrInvalidRequest)
		return
	}

	modifiedBytes, err := json.Marshal(req)
	if err != nil {
		slog.Error("failed to marshal merge request", "error", err, "entity", entity, "target_id", targetID)
		response.RespondInternalError(c)
		return
	}

	targetURL, err := newPathBuilder().MergeEntity().URLWithParams(h.config.UpstreamURL, map[string]string{"entity": entity, "id": targetID})
	if err != nil {
		slog.Error("failed to build merge target URL", "error", err, "entity", entity, "target_id", targetID)
		response.RespondInternalError(c)
		return
	}

	proxyReq, err := http.NewRequestWithContext(c, c.Request.Method, targetURL.String(), bytes.NewReader(modifiedBytes))
	if err != nil {
		slog.Error("failed to create proxy request", "error", err, "entity", entity, "target_id", targetID, "method", c.Request.Method)
		response.RespondInternalError(c)
		return
	}

	session, err := c.Cookie("koito_session")
	if err == nil {
		proxyReq.AddCookie(&http.Cookie{
			Name:  "koito_session",
			Value: session,
		})
	}

	if h.engine != nil {
		if err := h.addMergeRule(c.Request.Context(), entity, targetID, req.MergeFromID); err != nil {
			slog.Error("koito merge rule add failed", "entity", entity, "target_id", targetID, "merge_from_id", req.MergeFromID, "error", err)
		}
	}

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		slog.Error("failed to execute upstream proxy request", "error", err, "entity", entity, "target_id", targetID, "method", c.Request.Method)
		response.RespondBadGateway(c)
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func (h *Handler) addMergeRule(ctx context.Context, entity, targetID string, sourceID int64) error {
	slog.Info("koito merge request detected", "entity", entity, "target_id", targetID, "source_id", sourceID)

	switch entity {
	case "track":
		source, err := h.fetchTrack(ctx, sourceID)
		if err != nil {
			return fmt.Errorf("fetch source track: %w", err)
		}

		sourceAlbum, err := h.fetchAlbum(ctx, source.AlbumID)
		if err != nil {
			return fmt.Errorf("fetch release: %w", err)
		}

		target, err := h.fetchTrackString(ctx, targetID)
		if err != nil {
			return fmt.Errorf("fetch target track: %w", err)
		}

		targetAlbum, err := h.fetchAlbum(ctx, target.AlbumID)
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

		if err := h.rule.Create(ctx, &rule); err != nil {
			return err
		}
		h.engine.Add(rule)
		return nil
	case "artist":
		source, err := h.fetchArtist(ctx, sourceID)
		if err != nil {
			return fmt.Errorf("fetch source artist: %w", err)
		}
		target, err := h.fetchArtistString(ctx, targetID)
		if err != nil {
			return fmt.Errorf("fetch target artist: %w", err)
		}

		rule := model.Rule{
			MatchArtistName:   newNullString(source.Name),
			ReplaceArtistName: newNullString(target.Name),
			Enabled:           new(true),
		}

		if err := h.rule.Create(ctx, &rule); err != nil {
			return err
		}
		h.engine.Add(rule)
		return nil
	default:
		return nil
	}
}

func (h *Handler) fetchTrack(ctx context.Context, id int64) (*koitoTrack, error) {
	return h.fetchTrackString(ctx, fmt.Sprintf("%d", id))
}

func (h *Handler) fetchTrackString(ctx context.Context, id string) (*koitoTrack, error) {
	pathBuilder := newPathBuilder()
	api := pathBuilder.Track(id)
	body, err := h.fetchUpstreamAPI(ctx, "GET", api, nil)
	if err != nil {
		return nil, err
	}

	var track koitoTrack
	if err := json.Unmarshal(body, &track); err != nil {
		return nil, err
	}

	if track.Title == "" || track.Artists[0].Name == "" {
		return nil, fmt.Errorf("unexpected track payload: missing track or artist name")
	}

	return &track, nil
}

func (h *Handler) fetchArtist(ctx context.Context, id int64) (*koitoArtist, error) {
	return h.fetchArtistString(ctx, fmt.Sprintf("%d", id))
}

func (h *Handler) fetchArtistString(ctx context.Context, id string) (*koitoArtist, error) {
	pathBuilder := newPathBuilder()
	api := pathBuilder.Artist(id)
	body, err := h.fetchUpstreamAPI(ctx, "GET", api, nil)
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

func (h *Handler) fetchAlbum(ctx context.Context, id int64) (*koitoAlbum, error) {
	return h.fetchAlbumString(ctx, fmt.Sprintf("%d", id))
}

func (h *Handler) fetchAlbumString(ctx context.Context, id string) (*koitoAlbum, error) {
	pathBuilder := newPathBuilder()
	api := pathBuilder.Album(id)
	body, err := h.fetchUpstreamAPI(ctx, "GET", api, nil)
	if err != nil {
		return nil, err
	}

	var album koitoAlbum
	if err := json.Unmarshal(body, &album); err != nil {
		return nil, err
	}

	return &album, nil
}

func (h *Handler) fetchUpstreamAPI(ctx context.Context, method string, api APIPath, body []byte) ([]byte, error) {
	target, err := api.URL(h.config.UpstreamURL)
	if err != nil {
		return nil, err
	}
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
