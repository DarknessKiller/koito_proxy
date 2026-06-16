package listenbrainz

import (
	"bytes"
	"encoding/json"
	"koito_proxy/internal/config"
	"koito_proxy/internal/model"
	"koito_proxy/internal/response"
	"koito_proxy/internal/rules"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ruleEngine *rules.RuleEngine
	config     *config.Config
	httpClient *http.Client
}

func NewHandler(ruleEngine *rules.RuleEngine, cfg *config.Config, httpClient *http.Client) *Handler {
	return &Handler{
		ruleEngine: ruleEngine,
		config:     cfg,
		httpClient: httpClient,
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

	targetURL, err := newAPIPathBuilder().SubmitListen().URL(h.config.UpstreamURL)
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

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && !strings.HasPrefix(authHeader, "Bearer ") && !strings.HasPrefix(authHeader, "Token ") {
		slog.Warn("unexpected authorization header format, forwarding as-is", "path", c.Request.URL.Path)
		response.RespondUnauthorized(c, response.ErrMissingAuthHeader)
		return
	}
	proxyReq.Header.Set("Authorization", authHeader)

	resp, err := h.httpClient.Do(proxyReq)
	if err != nil {
		slog.Error("failed to execute listenbrainz upstream proxy request", "error", err, "method", c.Request.Method)
		response.RespondBadGateway(c)
		return
	}
	defer resp.Body.Close()

	slog.Debug("koito submit listen intercepted",
		"original_body", string(originalBody),
		"modified_body", string(modifiedBytes),
	)

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}
