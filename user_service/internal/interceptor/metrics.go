package interceptor

import (
	"context"
	"time"

	"github.com/Dokhoyan/go-messenger-microservices/user_service/internal/metric"
	"google.golang.org/grpc"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metric.IncRequestCounter()

	timeStart := time.Now()

	res, err := handler(ctx, req)
	diffTime := time.Since(timeStart).Seconds()

	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod)
		metric.HistogramResponseTimeObserve("error", diffTime)

		return nil, err
	}

	metric.IncResponseCounter("success", info.FullMethod)
	metric.HistogramResponseTimeObserve("success", diffTime)

	return res, nil
}