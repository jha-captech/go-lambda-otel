package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jha-captech/golambdaotel/internal/services"
	"github.com/jha-captech/golambdaotel/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
)

func HandlerSample(logger *slog.Logger, telemeter telemetry.Telemeter, service services.Service) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		logger.InfoContext(ctx, "running HandlerSample")
		ctx, span := telemeter.Tracer("handlers").Start(ctx, "HandlerSample")
		defer span.End()

		logger.InfoContext(ctx, "In HandlerSample")

		span.SetAttributes(attribute.String("stringAttr", "hi!"))

		service.DoStuff(ctx)

		return returnJSON(http.StatusOK, map[string]any{
			"message": "Hello World!",
		})
	}
}
