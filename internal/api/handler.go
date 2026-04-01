package api

import (
	"context"
	"encoding/json"
	"strings"

	"GolangChatAdvisor_Project/internal/models"
	"GolangChatAdvisor_Project/internal/service"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	companionService *service.CompanionService
	sessionService   *service.SessionService
}

func NewHandler(
	companionService *service.CompanionService,
	sessionService *service.SessionService,
) *Handler {
	return &Handler{
		companionService: companionService,
		sessionService:   sessionService,
	}
}

func (h *Handler) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := req.Path
	method := req.HTTPMethod

	switch {
	case method == "GET" && path == "/health":
		return jsonResponse(200, map[string]string{"status": "ok"})

	case method == "POST" && path == "/companions":
		return h.createCompanion(ctx, req)

	case method == "GET" && strings.HasPrefix(path, "/companions/"):
		return h.getCompanion(ctx, req)

	case method == "POST" && path == "/sessions":
		return h.createSession(ctx, req)

	case method == "POST" && strings.HasSuffix(path, "/messages"):
		return h.sendMessage(ctx, req)

	case method == "GET" && strings.HasSuffix(path, "/messages"):
		return h.getMessages(ctx, req)

	default:
		return errorResponse(404, "route not found")
	}
}

func (h *Handler) createCompanion(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body models.CreateAdvisorRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return errorResponse(400, "invalid request body")
	}

	companion, err := h.companionService.CreateCompanion(ctx, body)
	if err != nil {
		return errorResponse(400, err.Error())
	}

	return jsonResponse(201, companion)
}

func (h *Handler) getCompanion(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := strings.TrimPrefix(req.Path, "/companions/")

	companion, err := h.companionService.GetCompanion(ctx, id)
	if err != nil {
		return errorResponse(400, err.Error())
	}
	if companion == nil {
		return errorResponse(404, "companion not found")
	}

	return jsonResponse(200, companion)
}

func (h *Handler) createSession(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body models.CreateSessionRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return errorResponse(400, "invalid request body")
	}

	session, err := h.sessionService.StartSession(ctx, body)
	if err != nil {
		if err.Error() == "companion not found" {
			return errorResponse(404, err.Error())
		}
		return errorResponse(400, err.Error())
	}

	return jsonResponse(201, session)
}

func (h *Handler) sendMessage(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sessionID := extractSessionID(req.Path)
	if sessionID == "" {
		return errorResponse(400, "invalid session path")
	}

	var body models.SendMessageRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return errorResponse(400, "invalid request body")
	}

	resp, err := h.sessionService.SendMessage(ctx, sessionID, body)
	if err != nil {
		if err.Error() == "session not found" || err.Error() == "companion not found" {
			return errorResponse(404, err.Error())
		}
		return errorResponse(400, err.Error())
	}

	return jsonResponse(200, resp)
}

func (h *Handler) getMessages(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sessionID := extractSessionID(req.Path)
	if sessionID == "" {
		return errorResponse(400, "invalid session path")
	}

	messages, err := h.sessionService.GetMessages(ctx, sessionID)
	if err != nil {
		if err.Error() == "session not found" {
			return errorResponse(404, err.Error())
		}
		return errorResponse(400, err.Error())
	}

	return jsonResponse(200, messages)
}

func extractSessionID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 3 && parts[0] == "sessions" && parts[2] == "messages" {
		return parts[1]
	}
	return ""
}
