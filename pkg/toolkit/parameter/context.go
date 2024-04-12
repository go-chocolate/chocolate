package parameter

import "context"

type parameterContextKey struct{}

var _parameterContextKey = &parameterContextKey{}

func (p Parameters) WithContext(ctx context.Context) context.Context {
	if val := ctx.Value(_parameterContextKey); val != nil {
		val.(Parameters).Append(p)
		return ctx
	} else {
		return context.WithValue(ctx, _parameterContextKey, p)
	}
}

func FromContext(ctx context.Context) Parameters {
	if val := ctx.Value(_parameterContextKey); val != nil {
		return val.(Parameters)
	}
	return Parameters{}
}
