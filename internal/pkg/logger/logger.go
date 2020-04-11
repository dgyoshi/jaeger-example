package logger

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

func init() {
	var err error
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)
	myConfig := zap.Config{
		Level:    level,
		Encoding: "console",
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

	log, err = myConfig.Build()

	if err != nil {
		panic(err)
	}
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	tracingFields := extractTracingInfo(ctx)

	f := make([]zap.Field, len(tracingFields)+len(fields))

	i := 0
	for n := 0; n < len(tracingFields); n++ {
		f[i] = tracingFields[n]
		i++
	}

	for n := 0; n < len(fields); n++ {
		f[i] = fields[n]
		i++
	}

	log.WithOptions(zap.AddCallerSkip(2))
	log.Info(msg, f...)
}

func extractTracingInfo(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 3)
	if ctx == nil {
		return []zap.Field{}
	}

	sp := opentracing.SpanFromContext(ctx)
	if sp == nil {
		return []zap.Field{}
	}
	spc := sp.Context().(jaeger.SpanContext)

	fields[0] = zap.String("TraceID", spc.TraceID().String())
	fields[1] = zap.String("ParentID", spc.ParentID().String())
	fields[2] = zap.String("SpanID", spc.SpanID().String())
	return fields
}
