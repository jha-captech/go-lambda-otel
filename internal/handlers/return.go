package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func returnJSON(statusCode int, data any) (events.APIGatewayProxyResponse, error) {
	JSONData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(JSONData),
	}, err
}
