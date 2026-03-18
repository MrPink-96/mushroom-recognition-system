package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func TimeoutMiddleware(timeOut time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, timeOut)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()

	}

}
