package echo

import (
	"context"
	"math/rand"
	"time"

	"github.com/dgyoshi/jaeger-example/internal/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

type Service struct{}

func (s *Service) EchoDelay(ctx context.Context, msg string) string {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "EchoDelay")
	defer sp.Finish()
	logger.Info(ctx, "echo.service.EchoDelay")

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(3)
	time.Sleep(time.Duration(n) * time.Second)
	return msg
}
