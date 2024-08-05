package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type HandlerFunc = func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
