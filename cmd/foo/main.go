package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/dgyoshi/jaeger-example/internal/pkg/echo"
	"github.com/dgyoshi/jaeger-example/internal/pkg/logger"
	pb "github.com/dgyoshi/jaeger-example/internal/pkg/proto/echo"
	"github.com/dgyoshi/jaeger-example/internal/pkg/tracer"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
)

type Config struct {
	ServerProtocol string `env:"FOO_PROTOCOL,required"`
	ServerPort     string `env:"FOO_PORT,required"`
	BarHost        string `env:"BAR_HOST,required"`
	BarPort        string `env:"BAR_PORT,required"`
}

func main() {
	ctx := context.Background()
	logger.Info(ctx, "Start Foo Main")
	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		panic(err)
	}

	// initialize tracer
	tracer, closer, err := tracer.NewTracer("Foo Service")
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	// initialize grpc server
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// add opentracing stream interceptor to chain
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// add opentracing unary interceptor to chain
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
		)))

	barClient := echo.NewClient(conf.BarHost, conf.BarPort, tracer)
	if err != nil {
		panic(err)
	}

	pb.RegisterEchoServer(s, &server{BarClient: barClient})

	listen, err := net.Listen(conf.ServerProtocol, fmt.Sprintf(":%s", conf.ServerPort))
	if err != nil {
		panic(err.Error())
	}

	go func() {
		logger.Info(ctx, "Foo Server start")
		if err := s.Serve(listen); err != nil {
			panic(err.Error())
		}
	}()

	logger.Info(ctx, "Waiting Foo Server to be available")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for s := range c {
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return

		case syscall.SIGHUP:
		default:
			return
		}
	}

	return
}

type server struct {
	EchoService *echo.Service
	BarClient   *echo.EchoServiceClient
}

func (s *server) Echo(ctx context.Context, req *pb.EchoMsg) (*pb.Reply, error) {
	logger.Info(ctx, "foo.server.Echo")

	msg := req.GetMsg()

	msg = s.EchoService.EchoDelay(ctx, msg)

	res, err := s.BarClient.Echo(ctx, msg)
	if err != nil {
		fmt.Println(err)
		return &pb.Reply{}, err
	}

	sp := opentracing.SpanFromContext(ctx)
	spc := sp.Context().(jaeger.SpanContext)

	return &pb.Reply{
		Msg:     res,
		TraceId: spc.TraceID().String(),
	}, nil
}
