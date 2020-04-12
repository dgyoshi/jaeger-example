package tracing

import (
	"context"

	"github.com/dgyoshi/jaeger-example/internal/pkg/log"
	"google.golang.org/grpc"
)

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		log.Infof(ctx, "grpc.client")

		err := invoker(ctx, method, req, reply, cc, opts...)

		return err

	}
}
