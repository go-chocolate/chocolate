package chocomux

import (
	"context"
	"net/http"
)

type Context interface {
	Writer() http.ResponseWriter
	Request() *http.Request
	context.Context
}

type httpContext struct {
	w http.ResponseWriter
	r *http.Request
	context.Context
}

func (c *httpContext) Writer() http.ResponseWriter { return c.w }
func (c *httpContext) Request() *http.Request      { return c.r }

func WithStd(w http.ResponseWriter, r *http.Request) Context {
	return &httpContext{w: w, r: r, Context: r.Context()}
}
