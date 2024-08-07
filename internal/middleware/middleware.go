package middleware

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jha-captech/golambdaotel/internal/telemetry"
)

type (
	HandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	Middleware  func(next HandlerFunc) HandlerFunc
)

func AddToHandler(handler HandlerFunc, middlewares ...Middleware) HandlerFunc {
	mds := func(next HandlerFunc) HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next
	}
	return mds(handler)
}

func Logger(logger *slog.Logger, telemeter telemetry.Telemeter) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			ctx, span := telemeter.Tracer("middleware").Start(ctx, request.HTTPMethod+" "+request.Path)
			defer span.End()
			event, err := next(ctx, request)
			logger.InfoContext(ctx, "Request details", "status", event.StatusCode, "path", request.HTTPMethod+" "+request.Path)
			return event, err
		}
	}
}
