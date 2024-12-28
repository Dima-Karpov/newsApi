package router

import (
	"github.com/gin-gonic/gin"
	"newsApi/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()

	news := router.Group("/news")
	{
		news.GET("/", h.getNews)
	}

	return router
}
