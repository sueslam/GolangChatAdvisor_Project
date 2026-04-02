package models

/*
The structs define data models for advisors, sessions, and messages
along with request/response formats and map to both JSON APIs and DynamoDB storage
*/

// AI advisor attributes
type Advisor struct {
	ID        string `json:"id" dynamodbav:"id"`
	Name      string `json:"name" dynamodbav:"name"`
	Persona   string `json:"persona" dynamodbav:"persona"`
	Style     string `json:"style" dynamodbav:"style"`
	Greeting  string `json:"greeting" dynamodbav:"greeting"`
	CreatedAt string `json:"created_at" dynamodbav:"created_at"`
}

// Request sent by user to create avatar/advisor
// Backend generates ID and sets timestamp
type CreateAdvisorRequest struct {
	Name     string `json:"name"`
	Persona  string `json:"persona"`
	Style    string `json:"style"`
	Greeting string `json:"greeting"`
}

// Start a chat session
type CreateSessionRequest struct {
	AdvisorID string `json:"advisor_id"`
	UserID    string `json:"user_id"`
}

// User message content
type SendMessageRequest struct {
	Content string `json:"content"`
}

// Store session (one chat conversation) information/metadata about chat on DynamoDB
// PK is Partition Key and SK is Sort Key needed to organize data in DynamoDB
type SessionMeta struct {
	PK          string `json:"-" dynamodbav:"pk"`
	SK          string `json:"-" dynamodbav:"sk"`
	SessionID   string `json:"session_id" dynamodbav:"session_id"`
	CompanionID string `json:"advisor_id" dynamodbav:"advisor_id"`
	UserID      string `json:"user_id" dynamodbav:"user_id"`
	CreatedAt   string `json:"created_at" dynamodbav:"created_at"`
	ItemType    string `json:"item_type" dynamodbav:"item_type"`
}

// One message line inside the conversation
// Actual content of a message
// Role is whether it is user or ai advisor
type Message struct {
	PK        string `json:"-" dynamodbav:"pk"`
	SK        string `json:"-" dynamodbav:"sk"`
	SessionID string `json:"session_id" dynamodbav:"session_id"`
	Role      string `json:"role" dynamodbav:"role"`
	Content   string `json:"content" dynamodbav:"content"`
	Timestamp string `json:"timestamp" dynamodbav:"timestamp"`
	ItemType  string `json:"item_type" dynamodbav:"item_type"`
}

// What API returns
type SendMessageResponse struct {
	UserMessage Message `json:"user_message"`
	AIMessage   Message `json:"ai_message"`
}
