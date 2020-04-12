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
	"github.com/dgyoshi/jaeger-example/internal/pkg/log"
	pb "github.com/dgyoshi/jaeger-example/internal/pkg/proto/echo"
	"github.com/dgyoshi/jaeger-example/internal/pkg/tracer"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
)

type Config struct {
	Env            string `env:"ENV,required"`
	ServerProtocol string `env:"BAR_PROTOCOL,required"`
	ServerPort     string `env:"BAR_PORT,required"`
}

func main() {
	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		panic(err)
	}

	logger, err := log.New(conf.Env)
	if err != nil {
		panic(err)
	}

	// initialize tracer
	tracer, closer, err := tracer.NewTracer("Bar Service", logger)
	defer closer.Close()
	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)

	// initialize grpc server
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// add opentracing stream interceptor to chain
			grpc_zap.StreamServerInterceptor(logger),
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// add opentracing unary interceptor to chain
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
		)))

	pb.RegisterEchoServer(s, &server{})

	listen, err := net.Listen(conf.ServerProtocol, fmt.Sprintf(":%s", conf.ServerPort))
	if err != nil {
		panic(err.Error())
	}

	go func() {
		logger.Info("Bar Server start")
		if err := s.Serve(listen); err != nil {
			panic(err.Error())
		}
	}()

	logger.Info("Waiting Bar Server to be available")
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
}

func (s *server) Echo(ctx context.Context, req *pb.EchoMsg) (*pb.Reply, error) {
	log.Infof(ctx, "bar.server.Echo")

	s.EchoService.EchoDelay(ctx, req.GetMsg())

	sp := opentracing.SpanFromContext(ctx)
	spc := sp.Context().(jaeger.SpanContext)

	return &pb.Reply{
		Msg:     req.GetMsg(),
		TraceId: spc.TraceID().String(),
	}, nil
}
