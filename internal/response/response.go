package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

const (
	// Authentication errors
	ErrMissingKoitoSession = "missing_koito_session"
	ErrInvalidKoitoSession = "invalid_koito_session"
	ErrMissingAuthHeader   = "missing_authorization_header"
	ErrInvalidToken        = "invalid_token"

	ErrInvalidRequest      = "invalid_request"
	ErrInternalServer      = "internal_server_error"
	ErrUpstreamUnavailable = "upstream_service_unavailable"
)

func RespondError(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, ErrorResponse{Error: errorMessage})
}

func RespondUnauthorized(c *gin.Context, errorMessage string) {
	RespondError(c, http.StatusUnauthorized, errorMessage)
}

func RespondBadRequest(c *gin.Context, errorMessage string) {
	RespondError(c, http.StatusBadRequest, errorMessage)
}

func RespondInternalError(c *gin.Context) {
	RespondError(c, http.StatusInternalServerError, ErrInternalServer)
}

func RespondBadGateway(c *gin.Context) {
	RespondError(c, http.StatusBadGateway, ErrUpstreamUnavailable)
}
