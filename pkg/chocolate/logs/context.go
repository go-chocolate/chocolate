package logs

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Context interface {
	context.Context

	logrus.Ext1FieldLogger
}

var _ logrus.Ext1FieldLogger = (Context)(nil)

type loggerContext struct {
	context.Context

	logrus.Ext1FieldLogger
}

func WithContext(ctx context.Context) Context {
	if c, ok := ctx.(Context); ok {
		return c
	}
	return &loggerContext{
		Context:         ctx,
		Ext1FieldLogger: logrus.WithContext(ctx),
	}
}
