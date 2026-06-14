package admin

import (
	"koito_proxy/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(group *gin.RouterGroup, ruleSvc service.RuleService, auth gin.HandlerFunc) {
	group.Use(auth)

	h := NewHandler(ruleSvc)

	group.GET("/check", h.CheckAuth)
	group.GET("/rules", h.ListRules)
	group.GET("/rules/:id", h.GetRule)
	group.POST("/rules", h.CreateRule)
	group.PUT("/rules/:id", h.UpdateRule)
	group.DELETE("/rules/:id", h.DeleteRule)
	group.GET("/ui", UIHandler)
}
