package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/jha-captech/golambdaotel/internal/handlers"
	"github.com/jha-captech/golambdaotel/internal/telemetry"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	logger := slog.Default()
	logger.Info("running main")
	telemeter, shutdown, err := telemetry.SetupOTELSDK(ctx, "github.com/jha-captech/golambdaotel")
	defer shutdown(ctx)
	if err != nil {
		return fmt.Errorf("[in main.run] setup telemeter client: %w", err)
	}

	handler := handlers.HandlerSample(logger, telemeter)

	lambda.StartWithOptions(
		handler,
		lambda.WithContext(ctx),
		lambda.WithEnableSIGTERM(func() {
			if err = shutdown(ctx); err != nil {
				logger.Error("Error shutting down telemeter", "err", err)
			}
		}),
	)

	return nil
}
