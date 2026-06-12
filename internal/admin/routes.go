package admin

import (
	"koito_proxy/internal/model"
	"koito_proxy/internal/repository"
	"koito_proxy/internal/rules"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(group *gin.RouterGroup, repo repository.Repository[model.Rule], engine *rules.RuleEngine, auth gin.HandlerFunc) {
	group.Use(auth)

	h := NewHandler(repo, engine)

	group.GET("/check", h.CheckAuth)
	group.GET("/rules", h.ListRules)
	group.GET("/rules/:id", h.GetRule)
	group.POST("/rules", h.CreateRule)
	group.PUT("/rules/:id", h.UpdateRule)
	group.DELETE("/rules/:id", h.DeleteRule)
	group.GET("/ui", UIHandler)
}
