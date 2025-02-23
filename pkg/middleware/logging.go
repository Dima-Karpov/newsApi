package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

// RequestLoggerMiddleware логирует все входящие HTTP-запросы
func RequestLoggerMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID := c.GetHeader("X-Request-ID")

		// Выполняем обработку запроса
		c.Next()

		// Логируем после обработки
		log.WithFields(logrus.Fields{
			"trace_id": traceID,
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"query":    c.Request.URL.RawQuery,
			"status":   c.Writer.Status(),
			"time":     start,
		}).Info("Incoming request")
	}
}
