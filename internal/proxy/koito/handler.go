package koito

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"

	"koito_proxy/internal/config"
	"koito_proxy/internal/response"
	"koito_proxy/internal/service"

	"github.com/gin-gonic/gin"
)

var (
	validEntity = regexp.MustCompile(`^[a-z]+$`)
	validID     = regexp.MustCompile(`^\d+$`)
)

type Handler struct {
	koitoService *service.KoitoService
	config       *config.Config
	httpClient   *http.Client
}

func NewHandler(koitoService *service.KoitoService, cfg *config.Config, httpClient *http.Client) *Handler {
	return &Handler{
		koitoService: koitoService,
		config:       cfg,
		httpClient:   httpClient,
	}
}

type mergeRequest struct {
	MergeFromID int64 `json:"merge_from_id"`
}

func (h *Handler) InterceptMerge(c *gin.Context) {
	entity := c.Param("entity")
	targetID := c.Param("id")

	if !validEntity.MatchString(entity) || !validID.MatchString(targetID) {
		slog.Warn("invalid merge request parameters", "entity", entity, "target_id", targetID, "path", c.Request.URL.Path)
		response.RespondBadRequest(c, response.ErrInvalidRequest)
		return
	}

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

	targetURL, err := newAPIPathBuilder().MergeEntity().URLWithParams(
		h.config.UpstreamURL,
		map[string]string{"entity": entity, "id": targetID},
	)
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
	if err == nil && session != "" {
		if hasNonPrintable(session) {
			slog.Warn("koito_session cookie contains non-printable characters, rejecting", "path", c.Request.URL.Path)
		} else {
			proxyReq.AddCookie(&http.Cookie{
				Name:  "koito_session",
				Value: session,
			})
		}
	}

	if err := h.koitoService.CreateMergeRule(c.Request.Context(), entity, targetID, req.MergeFromID); err != nil {
		slog.Error("koito merge rule add failed", "entity", entity, "target_id", targetID, "merge_from_id", req.MergeFromID, "error", err)
		response.RespondInternalError(c)
		return
	}

	resp, err := h.httpClient.Do(proxyReq)
	if err != nil {
		slog.Error("failed to execute upstream proxy request", "error", err, "entity", entity, "target_id", targetID, "method", c.Request.Method)
		response.RespondBadGateway(c)
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func hasNonPrintable(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return true
		}
	}
	return false
}
