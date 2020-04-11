package echo

import (
	"context"
	"fmt"

	"github.com/dgyoshi/jaeger-example/internal/pkg/interceptor"
	"github.com/dgyoshi/jaeger-example/internal/pkg/logger"
	pb "github.com/dgyoshi/jaeger-example/internal/pkg/proto/echo"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type EchoServiceClient struct {
	client pb.EchoClient
}

func NewClient(host, port string, tracer opentracing.Tracer) *EchoServiceClient {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
			interceptor.OpenTracingUnaryClientInterceptor(),
		)))

	if err != nil {
		panic(err.Error())
	}

	client := pb.NewEchoClient(conn)

	return &EchoServiceClient{
		client: client,
	}
}

func (c *EchoServiceClient) Echo(ctx context.Context, msg string) (string, error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "Call Echo")
	defer sp.Finish()
	logger.Info(ctx, "echo.client.Echo")

	res, err := c.client.Echo(ctx, &pb.EchoMsg{
		Msg: msg,
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return res.GetMsg(), nil

}
