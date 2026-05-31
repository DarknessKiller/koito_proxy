package listenbrainz

import (
	"bytes"
	"encoding/json"
	"io"
	"koito_proxy/internal/config"
	"koito_proxy/internal/model"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.ruleEngine != nil {
		for i := range req.Payload {
			h.ruleEngine.Apply(&req.Payload[i].TrackMetaData)
		}
	}

	modifiedBytes, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetURL, err := newPathBuilder().SubmitListen().URL(h.config.UpstreamURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	proxyReq, err := http.NewRequestWithContext(
		c.Request.Context(),
		c.Request.Method,
		targetURL.String(),
		bytes.NewReader(modifiedBytes),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	proxyReq.Header.Set("Authorization", c.GetHeader("Authorization"))

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("koito submit listen intercepted",
		"body", string(modifiedBytes),
	)

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
