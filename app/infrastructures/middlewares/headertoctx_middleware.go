package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
)

func HeaderToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header yang diperlukan
		source := c.GetHeader("X-Source")

		// Tambahkan ke context
		ctx := c.Request.Context()
		if source != "" {
			ctx = context.WithValue(ctx, "X-Source", source)
		}

		// Update request context
		c.Request = c.Request.WithContext(ctx)

		// Lanjutkan ke handler berikutnya
		c.Next()
	}
}
