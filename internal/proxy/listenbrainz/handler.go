package listenbrainz

import (
	"bytes"
	"encoding/json"
	"io"
	"koito_proxy/internal/config"
	"koito_proxy/internal/model"
	"koito_proxy/internal/response"
	"koito_proxy/internal/rules"
	"log/slog"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ruleEngine *rules.RuleEngine
	config     *config.Config
}

func NewHandler(ruleEngine *rules.RuleEngine, cfg *config.Config) *Handler {
	return &Handler{
		ruleEngine: ruleEngine,
		config:     cfg,
	}
}

func (h *Handler) InterceptSubmitListen(c *gin.Context) {
	var req model.ListenBrainzSubmitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to parse listenbrainz submit request", "error", err, "path", c.Request.URL.Path)
		response.RespondBadRequest(c, response.ErrInvalidRequest)
		return
	}

	originalBody, err := json.Marshal(req)
	if err != nil {
		slog.Error("failed to marshal original listenbrainz request", "error", err)
		response.RespondInternalError(c)
		return
	}

	if h.ruleEngine != nil {
		for i := range req.Payload {
			h.ruleEngine.Apply(&req.Payload[i].TrackMetaData)
		}
	}

	modifiedBytes, err := json.Marshal(req)
	if err != nil {
		slog.Error("failed to marshal modified listenbrainz request", "error", err)
		response.RespondInternalError(c)
		return
	}

	targetURL, err := newPathBuilder().SubmitListen().URL(h.config.UpstreamURL)
	if err != nil {
		slog.Error("failed to build listenbrainz target URL", "error", err)
		response.RespondInternalError(c)
		return
	}

	proxyReq, err := http.NewRequestWithContext(
		c.Request.Context(),
		c.Request.Method,
		targetURL.String(),
		bytes.NewReader(modifiedBytes),
	)
	if err != nil {
		slog.Error("failed to create listenbrainz proxy request", "error", err, "method", c.Request.Method)
		response.RespondInternalError(c)
		return
	}

	proxyReq.Header.Set("Authorization", c.GetHeader("Authorization"))

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		slog.Error("failed to execute listenbrainz upstream proxy request", "error", err, "method", c.Request.Method)
		response.RespondBadGateway(c)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read listenbrainz proxy response body", "error", err, "status", resp.StatusCode)
		response.RespondInternalError(c)
		return
	}

	slog.Info("koito submit listen intercepted",
		"original_body", string(originalBody),
		"modified_body", string(modifiedBytes),
	)

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
