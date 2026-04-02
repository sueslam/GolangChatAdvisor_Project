package api

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// API Response Helper, helps Lambda API return JSON consistently

// Function builds an api response by taking in status code and payload
// Returns APIGatewayProxyResponse and errors
func jsonResponse(statusCode int, payload any) (events.APIGatewayProxyResponse, error) {
	//Converts Go to Json
	body, err := json.Marshal(payload)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"failed to marshal response"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func errorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return jsonResponse(statusCode, map[string]string{
		"error": message,
	})
}
