package middlewares

import (
	"avito/internal/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)


func DurationMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func (c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		fmt.Printf("[%s] %s %s %d %s\n", start.Format(time.RFC3339), c.Request.Method, c.Request.URL.Path, statusCode, duration)
	}
}