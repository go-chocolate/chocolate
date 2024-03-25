package chocogin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-chocolate/chocolate/pkg/chocolate/chocohttp/chocomux"
)

func Gin(h chocomux.Handler) gin.HandlerFunc {
	return func(c *gin.Context) { h(WithGin(c)) }
}

func FromGin(h gin.HandlerFunc) chocomux.Handler {
	return func(ctx chocomux.Context) {
		ginCtx, _ := gin.CreateTestContext(ctx.Writer())
		ginCtx.Request = ctx.Request()
		h(ginCtx)
	}
}
