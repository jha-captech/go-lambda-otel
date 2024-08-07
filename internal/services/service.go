package services

import (
	"context"
	"time"

	"github.com/jha-captech/golambdaotel/internal/telemetry"
)

type Service struct {
	telemeter telemetry.Telemeter
}

func NewService(telemeter telemetry.Telemeter) Service {
	return Service{telemeter: telemeter}
}

func (s Service) DoStuff(ctx context.Context) {
	ctx, span := s.telemeter.Tracer("service").Start(ctx, "DoStuff")
	defer span.End()

	time.Sleep(time.Second)
}
