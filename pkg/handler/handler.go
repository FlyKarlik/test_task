package handler

import (
	"mongo_db/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	auth := router.Group("/auth")
	{
		auth.POST("/getTokens", h.GetTokens)
		auth.POST("/refreshTokens", h.RefreshTokens)
	}

	return router
}
