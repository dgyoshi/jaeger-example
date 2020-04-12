package echo

import (
	"context"
	"math/rand"
	"time"

	"github.com/dgyoshi/jaeger-example/internal/pkg/log"
	"github.com/opentracing/opentracing-go"
)

type Service struct{}

func (s *Service) EchoDelay(ctx context.Context, msg string) string {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "EchoDelay")
	defer sp.Finish()
	log.Infof(ctx, "echo.service.EchoDelay")

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(3)
	time.Sleep(time.Duration(n) * time.Second)
	return msg
}
