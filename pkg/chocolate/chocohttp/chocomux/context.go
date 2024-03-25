package chocomux

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/creasty/defaults"

	"github.com/go-chocolate/chocolate/pkg/chocolate/errorx"
)

type Context interface {
	Writer() http.ResponseWriter
	Request() *http.Request
	Bind(data any) error
	Error(err error)
	Ok(body []byte)
	JSON(status int, body any)
	OkJSON(body any)
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

func (c *httpContext) Bind(data any) error {
	if err := defaults.Set(data); err != nil {
		return err
	}
	request := c.Request()
	contentType := request.Header.Get("Content-Type")
	var err error
	switch contentType {
	case "application/json":
		err = json.NewDecoder(request.Body).Decode(data)
	case "application/x-www-form-urlencoded":
		if err = request.ParseForm(); err != nil {
			return err
		}
		err = UnmarshalForm(request.Form, data)
	case "form-data":
		if request.Form == nil {
			if err = request.ParseMultipartForm(32 << 20); err != nil {
				return err
			}
		}
		err = UnmarshalMultipartForm(request.MultipartForm, data)
	default:
		err = UnmarshalForm(request.URL.Query(), data)
	}
	if err != nil {
		return err
	}
	if v, ok := data.(Validator); ok {
		return v.Validate(c.Context)
	}
	return nil
}

func (c *httpContext) Error(err error) {
	if e, ok := err.(*errorx.Error); ok {
		c.JSON(http.StatusInternalServerError, e)
	} else {
		c.w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *httpContext) Ok(body []byte) {
	c.w.WriteHeader(http.StatusOK)
	_, _ = c.Writer().Write(body)
}

func (c *httpContext) JSON(status int, body any) {
	c.w.WriteHeader(status)
	c.w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(c.w).Encode(body)
}

func (c *httpContext) OkJSON(body any) {
	c.JSON(http.StatusOK, body)
}
