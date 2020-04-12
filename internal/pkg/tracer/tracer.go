package tracer

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	zaplogger "github.com/uber/jaeger-client-go/log/zap"
	"go.uber.org/zap"
)

func NewTracer(service string, logger *zap.Logger) (opentracing.Tracer, io.Closer, error) {
	cfg, _ := config.FromEnv()
	l := zaplogger.NewLogger(logger)

	return cfg.New(service, config.Logger(l))
}
