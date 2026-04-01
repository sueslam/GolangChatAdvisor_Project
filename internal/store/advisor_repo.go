package store

import (
	"context"

	"GolangChatAdvisor_Project/internal/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AdvisorRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewAdvisorRepository(client *dynamodb.Client, tableName string) *AdvisorRepository {
	return &AdvisorRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *AdvisorRepository) CreateAdvisor(ctx context.Context, advisor models.Advisor) error {
	item, err := attributevalue.MarshalMap(advisor)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	return err
}

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
