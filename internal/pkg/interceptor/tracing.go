package interceptor

import (
	"context"

	"github.com/dgyoshi/jaeger-example/internal/pkg/logger"
	"google.golang.org/grpc"
)

func OpenTracingUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		logger.Info(ctx, "grpc.client")

		err := invoker(ctx, method, req, reply, cc, opts...)

		return err

	}
}
