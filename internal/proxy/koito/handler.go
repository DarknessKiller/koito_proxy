package koito

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	"koito_proxy/internal/config"
	"koito_proxy/internal/response"
	"koito_proxy/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	koitoService service.KoitoService
	config       *config.Config
}

func NewHandler(koitoSvc service.KoitoService, cfg *config.Config) *Handler {
	return &Handler{
		koitoService: koitoSvc,
		config:       cfg,
	}
}

type mergeRequest struct {
	MergeFromID int64 `json:"merge_from_id"`
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

	if err := h.koitoService.AddMergeRule(c.Request.Context(), entity, targetID, req.MergeFromID); err != nil {
		slog.Error("koito merge rule add failed", "entity", entity, "target_id", targetID, "merge_from_id", req.MergeFromID, "error", err)
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
