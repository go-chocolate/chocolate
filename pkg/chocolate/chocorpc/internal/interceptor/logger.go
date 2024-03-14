package interceptor

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()
	defer logger(ctx, start, req, info, resp, err)
	return handler(ctx, req)

}

func logger(ctx context.Context, start time.Time, req any, info *grpc.UnaryServerInfo, resp any, err error) {
	log := logrus.WithContext(ctx).WithFields(map[string]interface{}{
		"method":   info.FullMethod,
		"duration": time.Since(start),
	})
	if err != nil {
		log.Error(err)
	} else {
		log.Info()
	}
}
