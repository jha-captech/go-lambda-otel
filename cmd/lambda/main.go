package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jha-captech/golambdaotel/internal/handlers"
	"github.com/jha-captech/golambdaotel/internal/middleware"
	"github.com/jha-captech/golambdaotel/internal/services"
	"github.com/jha-captech/golambdaotel/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "My lambda name")
	defer os.Unsetenv("AWS_LAMBDA_FUNCTION_NAME")

	logger := slog.Default()
	logger.InfoContext(ctx, "running main")

	telemeter, shutdown, err := telemetry.NewTelemeter(ctx, "golambdaotel")
	defer shutdown(ctx)
	if err != nil {
		return fmt.Errorf("[in main.run] setup telemeter client: %w", err)
	}

	service := services.NewService(telemeter)

	handler := handlers.HandlerSample(logger, telemeter, service)

	handler = middleware.AddToHandler(
		handler,
		middleware.Logger(logger, telemeter),
	)

	InstrumentedHandler := otellambda.InstrumentHandler(
		handler,
		otellambda.WithTracerProvider(telemeter.TracerProvider),
		otellambda.WithPropagator(telemeter.TextMapPropagator),
		otellambda.WithFlusher(telemeter.TracerProvider),
	)

	lambda.StartWithOptions(
		InstrumentedHandler,
		lambda.WithContext(ctx),
		lambda.WithEnableSIGTERM(func() {
			fmt.Println("calling shutdown")
			if err = shutdown(ctx); err != nil {
				logger.Error("Error shutting down telemeter", "err", err)
			}
		}),
	)

	return nil
}
