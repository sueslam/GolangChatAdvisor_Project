package store

import (
	"context"

	"GolangChatAdvisor_Project/internal/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

//Middle layer between app and DB
//Database Access Layer: How to save and fetch ai advisors from DynamoDB

// Contains connection to DynamoDB and which table to use
type AdvisorRepository struct {
	client    *dynamodb.Client
	tableName string
}

// Constructor to create new connection and get table
func NewAdvisorRepository(client *dynamodb.Client, tableName string) *AdvisorRepository {
	return &AdvisorRepository{
		client:    client,
		tableName: tableName,
	}
}

// Save an advisor into Dynamo DB
func (r *AdvisorRepository) CreateAdvisor(ctx context.Context, advisor models.Advisor) error {
	//Convert to DynamoDB map
	item, err := attributevalue.MarshalMap(advisor)
	if err != nil {
		return err
	}

	//Save in DynamoDB
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	return err
}

// Fetch an advisor by querying ID from DB and convert back into Go object
func (r *AdvisorRepository) GetAdvisorByID(ctx context.Context, id string) (*models.Advisor, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, nil
	}

	var advisor models.Advisor
	if err := attributevalue.UnmarshalMap(out.Item, &advisor); err != nil {
		return nil, err
	}

	return &advisor, nil
}
