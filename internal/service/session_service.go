package service

import (
	"context"
	"errors"
	"time"

	"GolangChatAdvisor_Project/internal/ai"
	"GolangChatAdvisor_Project/internal/models"
	"GolangChatAdvisor_Project/internal/store"

	"github.com/google/uuid"
)

type SessionService struct {
	advisorRepo *store.AdvisorRepository
	sessionRepo *store.SessionRepository
	responder   *ai.Responder
}

func NewSessionService(
	advisorRepo *store.AdvisorRepository,
	sessionRepo *store.SessionRepository,
	responder *ai.Responder,
) *SessionService {
	return &SessionService{
		advisorRepo: advisorRepo,
		sessionRepo: sessionRepo,
		responder:   responder,
	}
}

func (s *SessionService) StartSession(ctx context.Context, req models.CreateSessionRequest) (*models.SessionMeta, error) {
	if req.AdvisorID == "" || req.UserID == "" {
		return nil, errors.New("advisor_id and user_id are required")
	}

	advisor, err := s.advisorRepo.GetAdvisorByID(ctx, req.AdvisorID)
	if err != nil {
		return nil, err
	}
	if advisor == nil {
		return nil, errors.New("advisor not found")
	}

	sessionID := "sess_" + uuid.NewString()
	now := time.Now().UTC().Format(time.RFC3339Nano)

	meta := models.SessionMeta{
		PK:        "SESSION#" + sessionID,
		SK:        "META",
		SessionID: sessionID,
		AdvisorID: req.AdvisorID,
		UserID:    req.UserID,
		CreatedAt: now,
		ItemType:  "session_meta",
	}

	if err := s.sessionRepo.CreateSessionMeta(ctx, meta); err != nil {
		return nil, err
	}

	greeting := models.Message{
		PK:        "SESSION#" + sessionID,
		SK:        "MSG#" + now,
		SessionID: sessionID,
		Role:      "assistant",
		Content:   advisor.Greeting,
		Timestamp: now,
		ItemType:  "message",
	}

	if err := s.sessionRepo.AddMessage(ctx, greeting); err != nil {
		return nil, err
	}

	return &meta, nil
}

func (s *SessionService) SendMessage(ctx context.Context, sessionID string, req models.SendMessageRequest) (*models.SendMessageResponse, error) {
	if sessionID == "" {
		return nil, errors.New("session id is required")
	}
	if req.Content == "" {
		return nil, errors.New("content is required")
	}

	meta, err := s.sessionRepo.GetSessionMeta(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, errors.New("session not found")
	}

	advisor, err := s.advisorRepo.GetAdvisorByID(ctx, meta.AdvisorID)
	if err != nil {
		return nil, err
	}
	if advisor == nil {
		return nil, errors.New("advisor not found")
	}

	userTS := time.Now().UTC().Format(time.RFC3339Nano)
	userMsg := models.Message{
		PK:        "SESSION#" + sessionID,
		SK:        "MSG#" + userTS,
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Content,
		Timestamp: userTS,
		ItemType:  "message",
	}

	if err := s.sessionRepo.AddMessage(ctx, userMsg); err != nil {
		return nil, err
	}

	history, err := s.sessionRepo.ListMessages(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	reply := s.responder.GenerateReply(*advisor, history, req.Content)

	assistantTS := time.Now().UTC().Add(1 * time.Millisecond).Format(time.RFC3339Nano)
	assistantMsg := models.Message{
		PK:        "SESSION#" + sessionID,
		SK:        "MSG#" + assistantTS,
		SessionID: sessionID,
		Role:      "assistant",
		Content:   reply,
		Timestamp: assistantTS,
		ItemType:  "message",
	}

	if err := s.sessionRepo.AddMessage(ctx, assistantMsg); err != nil {
		return nil, err
	}

	return &models.SendMessageResponse{
		UserMessage: userMsg,
		AIMessage:   assistantMsg,
	}, nil
}

func (s *SessionService) GetMessages(ctx context.Context, sessionID string) ([]models.Message, error) {
	if sessionID == "" {
		return nil, errors.New("session id is required")
	}

	meta, err := s.sessionRepo.GetSessionMeta(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, errors.New("session not found")
	}

	return s.sessionRepo.ListMessages(ctx, sessionID)
}
