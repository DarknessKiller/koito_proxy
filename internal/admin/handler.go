package admin

import (
	"database/sql"
	"koito_proxy/internal/model"
	"koito_proxy/internal/repository"
	"koito_proxy/internal/response"
	"koito_proxy/internal/rules"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type Handler struct {
	repo   repository.Repository[model.Rule]
	engine *rules.RuleEngine
}

func NewHandler(repo repository.Repository[model.Rule], engine *rules.RuleEngine) *Handler {
	return &Handler{repo: repo, engine: engine}
}

type RuleRequest struct {
	MatchTrackName      *string  `json:"match_track_name,omitempty"`
	MatchArtistName     *string  `json:"match_artist_name,omitempty"`
	MatchReleaseName    *string  `json:"match_release_name,omitempty"`
	MatchArtistNames    []string `json:"match_artist_names,omitempty"`
	MatchDurationBucket *int32   `json:"match_duration_bucket,omitempty"`
	MatchMBID           *string  `json:"match_mbid,omitempty"`
	ReplaceTrackName    *string  `json:"replace_track_name,omitempty"`
	ReplaceArtistName   *string  `json:"replace_artist_name,omitempty"`
	ReplaceReleaseName  *string  `json:"replace_release_name,omitempty"`
	ReplaceArtistNames  []string `json:"replace_artist_names,omitempty"`
	Enabled             *bool    `json:"enabled,omitempty"`
}

type RuleResponse struct {
	ID                  string   `json:"id"`
	MatchTrackName      *string  `json:"match_track_name,omitempty"`
	MatchArtistName     *string  `json:"match_artist_name,omitempty"`
	MatchReleaseName    *string  `json:"match_release_name,omitempty"`
	MatchArtistNames    []string `json:"match_artist_names,omitempty"`
	MatchDurationBucket *int32   `json:"match_duration_bucket,omitempty"`
	MatchMBID           *string  `json:"match_mbid,omitempty"`
	ReplaceTrackName    *string  `json:"replace_track_name,omitempty"`
	ReplaceArtistName   *string  `json:"replace_artist_name,omitempty"`
	ReplaceReleaseName  *string  `json:"replace_release_name,omitempty"`
	ReplaceArtistNames  []string `json:"replace_artist_names,omitempty"`
	Enabled             bool     `json:"enabled"`
	Priority            int      `json:"priority"`
	Valid               bool     `json:"valid"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
}

func (h *Handler) ListRules(c *gin.Context) {
	rules, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		slog.Error("failed to list rules", "error", err)
		response.RespondInternalError(c)
		return
	}

	resp := make([]RuleResponse, 0, len(rules))
	for _, r := range rules {
		resp = append(resp, ruleToResponse(r))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetRule(c *gin.Context) {
	id := c.Param("id")
	rule, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to get rule", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	c.JSON(http.StatusOK, ruleToResponse(*rule))
}

func (h *Handler) CreateRule(c *gin.Context) {
	var req RuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind create rule request", "error", err)
		response.RespondBadRequest(c, response.ErrInvalidRequest)
		return
	}

	rule := requestToRule(req)

	if err := h.repo.Create(c.Request.Context(), rule); err != nil {
		slog.Error("failed to create rule", "error", err)
		response.RespondInternalError(c)
		return
	}

	h.engine.Add(*rule)

	slog.Info("rule created", "id", rule.ID.String())
	c.JSON(http.StatusCreated, ruleToResponse(*rule))
}

func (h *Handler) UpdateRule(c *gin.Context) {
	id := c.Param("id")

	existing, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to get rule for update", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	var req RuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind update rule request", "error", err)
		response.RespondBadRequest(c, response.ErrInvalidRequest)
		return
	}

	updated := applyRequestToRule(existing, &req)

	if err := h.repo.Update(c.Request.Context(), id, updated); err != nil {
		slog.Error("failed to update rule", "id", id, "error", err)
		response.RespondInternalError(c)
		return
	}

	parsedID, err := ksuid.Parse(id)
	if err == nil {
		slog.Error("failed to parse rule id", "id", id, "error", err)
		h.engine.Remove(parsedID)
	}
	h.engine.Add(*updated)

	slog.Info("rule updated", "id", id)
	c.JSON(http.StatusOK, ruleToResponse(*updated))
}

func (h *Handler) DeleteRule(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		slog.Error("failed to delete rule", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		slog.Error("failed to delete rule", "id", id, "error", err)
		response.RespondInternalError(c)
		return
	}

	parsedID, err := ksuid.Parse(id)
	if err != nil {
		slog.Error("failed to parse rule id", "id", id, "error", err)
		response.RespondInternalError(c)
		return
	}

	h.engine.Remove(parsedID)

	slog.Info("rule deleted", "id", id)
	c.Status(http.StatusNoContent)
}

func ruleToResponse(r model.Rule) RuleResponse {
	resp := RuleResponse{
		ID:                 r.ID.String(),
		MatchArtistNames:   r.MatchArtistNames,
		ReplaceArtistNames: r.ReplaceArtistNames,
		Enabled:            *r.Enabled,
		CreatedAt:          r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if r.MatchDurationBucket != 0 {
		v := r.MatchDurationBucket
		resp.MatchDurationBucket = &v
	}
	if r.MatchTrackName.Valid {
		resp.MatchTrackName = &r.MatchTrackName.String
	}
	if r.MatchArtistName.Valid {
		resp.MatchArtistName = &r.MatchArtistName.String
	}
	if r.MatchReleaseName.Valid {
		resp.MatchReleaseName = &r.MatchReleaseName.String
	}
	if r.MatchMBID.Valid {
		resp.MatchMBID = &r.MatchMBID.String
	}
	if r.ReplaceTrackName.Valid {
		resp.ReplaceTrackName = &r.ReplaceTrackName.String
	}
	if r.ReplaceArtistName.Valid {
		resp.ReplaceArtistName = &r.ReplaceArtistName.String
	}
	if r.ReplaceReleaseName.Valid {
		resp.ReplaceReleaseName = &r.ReplaceReleaseName.String
	}

	compiled := rules.CompileRule(r)
	resp.Priority = compiled.Priority
	resp.Valid = compiled.Valid

	return resp
}

func requestToRule(req RuleRequest) *model.Rule {
	return &model.Rule{
		MatchTrackName:      nullStr(req.MatchTrackName),
		MatchArtistName:     nullStr(req.MatchArtistName),
		MatchReleaseName:    nullStr(req.MatchReleaseName),
		MatchArtistNames:    req.MatchArtistNames,
		MatchDurationBucket: int32Val(req.MatchDurationBucket),
		MatchMBID:           nullStr(req.MatchMBID),
		ReplaceTrackName:    nullStr(req.ReplaceTrackName),
		ReplaceArtistName:   nullStr(req.ReplaceArtistName),
		ReplaceReleaseName:  nullStr(req.ReplaceReleaseName),
		ReplaceArtistNames:  req.ReplaceArtistNames,
		Enabled:             req.Enabled,
	}
}

func applyRequestToRule(existing *model.Rule, req *RuleRequest) *model.Rule {
	existing.MatchTrackName = nullStr(req.MatchTrackName)
	existing.MatchArtistName = nullStr(req.MatchArtistName)
	existing.MatchReleaseName = nullStr(req.MatchReleaseName)
	existing.MatchArtistNames = req.MatchArtistNames
	existing.MatchDurationBucket = int32Val(req.MatchDurationBucket)
	existing.MatchMBID = nullStr(req.MatchMBID)
	existing.ReplaceTrackName = nullStr(req.ReplaceTrackName)
	existing.ReplaceArtistName = nullStr(req.ReplaceArtistName)
	existing.ReplaceReleaseName = nullStr(req.ReplaceReleaseName)
	existing.ReplaceArtistNames = req.ReplaceArtistNames
	existing.Enabled = req.Enabled
	return existing
}

func nullStr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	trimmed := strings.TrimSpace(*s)
	return sql.NullString{String: trimmed, Valid: trimmed != ""}
}

func int32Val(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}
