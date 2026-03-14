package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/Jeskay/micsvc/internal/metrics"
)

func NewMetricsUnary(grpcMetrics *metrics.GRPCMetrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		t := time.Now()

		resp, err = handler(ctx, req)

		method := info.FullMethod
		statusCode := status.Code(err).String()

		grpcMetrics.GRPCReuqestCounter.WithLabelValues(method, statusCode).Inc()
		grpcMetrics.GRPCRequestLatency.WithLabelValues(method).Observe(time.Since(t).Seconds())
		return
	}
}
