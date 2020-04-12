package tracing

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (
		interface{},
		error,
	) {
		sp := opentracing.SpanFromContext(ctx)
		if sp != nil {
			spc := sp.Context().(jaeger.SpanContext)
			ctxzap.AddFields(ctx, zap.String("TraceID", spc.TraceID().String()))
			ctxzap.AddFields(ctx, zap.String("ParentID", spc.ParentID().String()))
			ctxzap.AddFields(ctx, zap.String("SpanID", spc.SpanID().String()))
		}

		return handler(ctx, req)
	}
}
