// Package main is the entry point of the application
package main

import (
	"context"
	"log"

	"GolangChatAdvisor_Project/internal/ai"
	"GolangChatAdvisor_Project/internal/api"
	"GolangChatAdvisor_Project/internal/config"
	"GolangChatAdvisor_Project/internal/service"
	"GolangChatAdvisor_Project/internal/store"

	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	cfg := config.Load()

	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		log.Fatalf("failed to load aws config: %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(awsCfg)

	advisorRepo := store.NewAdvisorRepository(dynamoClient, cfg.AdvisorsTable)
	sessionRepo := store.NewSessionRepository(dynamoClient, cfg.SessionsTable)

	responder := ai.NewResponder()

	advisorService := service.NewAdvisorService(advisorRepo)
	sessionService := service.NewSessionService(advisorRepo, sessionRepo, responder)

	handler := api.NewHandler(advisorService, sessionService)

	lambda.Start(handler.HandleRequest)
}
