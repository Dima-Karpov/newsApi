package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"newsApi/internal/service"
	"newsApi/pkg/middleware"

	_ "newsApi/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()
	logger := logrus.New()
	router.Use(middleware.RequestLoggerMiddleware(logger))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	news := router.Group("/news")
	{
		news.GET("/", h.getAllNews)
		news.GET("/:id", h.getNews)
	}

	return router
}
