package middlewares

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"route256/loms/internal/metrics"
)

func Metrics(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)
	duration := time.Since(start)
	statusCode := status.Code(err)
	metrics.RequestCounterInc(info.FullMethod, strconv.FormatUint(uint64(statusCode), 10))
	metrics.RequestHandlerDuration(info.FullMethod, duration)
	return resp, err
}
