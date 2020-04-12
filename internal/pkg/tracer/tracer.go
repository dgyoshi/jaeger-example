package tracer

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	zaplogger "github.com/uber/jaeger-client-go/log/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewTracer(service string) (opentracing.Tracer, io.Closer, error) {

	//cfg := &config.Configuration{
	//	Sampler: &config.SamplerConfig{
	//		Type:  "const",
	//		Param: 1,
	//	},
	//	Reporter: &config.ReporterConfig{
	//		LogSpans: true,
	//	},
	//}
	cfg, _ := config.FromEnv()

	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)
	myConfig := zap.Config{
		Level:    level,
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	zl, _ := myConfig.Build()
	l := zaplogger.NewLogger(zl)

	// create tracer from config
	return cfg.New(service, config.Logger(l))
}
