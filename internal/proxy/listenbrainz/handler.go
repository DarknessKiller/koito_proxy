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
	"net/url"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine *rules.Engine
	config *config.Config
}

func NewHandler(e *rules.Engine, cfg *config.Config) *Handler {
	return &Handler{
		engine: e,
		config: cfg,
	}
}

func (h *Handler) targetURL(apiPath APIPath) (*url.URL, error) {
	return apiPath.URL(h.config.UpstreamURL)
}

func (h *Handler) InterceptSubmitListen(c *gin.Context) {
	var req model.ListenBrainzSubmitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.engine != nil {
		for i := range req.Payload {
			h.engine.Apply(&req.Payload[i].TrackMetaData)
		}
	}

	modified, err := json.Marshal(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pathBuilder := newPathBuilder()
	targetURL, err := h.targetURL(pathBuilder.SubmitListen())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	proxyReq, err := http.NewRequestWithContext(
		c.Request.Context(),
		c.Request.Method,
		targetURL.String(),
		bytes.NewReader(modified),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for k, v := range c.Request.Header {
		for _, vv := range v {
			proxyReq.Header.Add(k, vv)
		}
	}

	proxyReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	slog.Info("koito submit listen intercepted", "original_body", c.Request.Body, "modified_body", string(modified))
	respBody, _ := io.ReadAll(resp.Body)

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
