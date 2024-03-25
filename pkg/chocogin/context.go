package chocogin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinContext struct {
	ginCtx *gin.Context
	context.Context
}

func (ctx *GinContext) Writer() http.ResponseWriter {
	return ctx.ginCtx.Writer
}
func (ctx *GinContext) Request() *http.Request {
	return ctx.ginCtx.Request
}

func WithGin(ctx *gin.Context) *GinContext {
	return &GinContext{ctx, ctx.Request.Context()}
}
