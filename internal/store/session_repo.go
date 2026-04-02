package store

import (
	"context"

	"GolangChatAdvisor_Project/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DB Layer for chat session and messages
// Save chat session/metadata, msgs/ individual texts and fetch chat session and messages
// Handles conversation data

// Contains which DB connection and table to use
type SessionRepository struct {
	client    *dynamodb.Client
	tableName string
}

// Constructor to create new repository object
func NewSessionRepository(client *dynamodb.Client, tableName string) *SessionRepository {
	return &SessionRepository{
		client:    client,
		tableName: tableName,
	}
}

// Save session metadata such as sessionID, userID, advisor, createTime etc. into DynamoDB
// Basically create a new chat session in the DB
func (r *SessionRepository) CreateSessionMeta(ctx context.Context, meta models.SessionMeta) error {
	item, err := attributevalue.MarshalMap(meta)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	return err
}

// Get session record for a particular ID
// PK: e.g. SESSION#22 and META identifies the session metadata
// all data for one session is grouped under the same pk and sk decides what kind of item it is
func (r *SessionRepository) GetSessionMeta(ctx context.Context, sessionID string) (*models.SessionMeta, error) {
	pk := "SESSION#" + sessionID

	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: "META"},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, nil
	}

	var meta models.SessionMeta
	if err := attributevalue.UnmarshalMap(out.Item, &meta); err != nil {
		return nil, err
	}

	return &meta, nil
}

// Save one chat msg into dynamoDB
func (r *SessionRepository) AddMessage(ctx context.Context, msg models.Message) error {
	item, err := attributevalue.MarshalMap(msg)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	return err
}

// Fetch ALL msgs from ONE session
// Get allitems where e.g. PK is SESSION#22 and SK starts with MSG#
// This way it ONLY fetches msgs and NOT the session metadata
// Return in ascending sortkey order (oldest msgs to newest)
func (r *SessionRepository) ListMessages(ctx context.Context, sessionID string) ([]models.Message, error) {
	pk := "SESSION#" + sessionID

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.tableName,
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :msgPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: pk},
			":msgPrefix": &types.AttributeValueMemberS{Value: "MSG#"},
		},
		ScanIndexForward: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	var messages []models.Message
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
