package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/jha-captech/golambdaotel/internal/telemetry"

	"github.com/aws/aws-lambda-go/events"
)

func HandlerSample(logger *slog.Logger, telemeter *telemetry.Telemeter) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		logger.Info("running HandlerSample")
		ctx, span := telemeter.Tracer.Start(ctx, "HandlerSample")
		defer span.End()

		// ctx, span := tracer.Start(ctx, "HandlerSample")

		logger.Info("In HandlerSample")

		return returnJSON(http.StatusOK, map[string]any{
			"message": "Hello World!",
		})
	}
}
