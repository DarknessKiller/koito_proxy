package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContextKey string

const RequestIDContextKey ContextKey = "request_id"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set(string(RequestIDContextKey), requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}
