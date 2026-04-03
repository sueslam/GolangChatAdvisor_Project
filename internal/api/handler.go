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
	advisorService *service.AdvisorService
	sessionService *service.SessionService
}

func NewHandler(
	advisorService *service.AdvisorService,
	sessionService *service.SessionService,
) *Handler {
	return &Handler{
		advisorService: advisorService,
		sessionService: sessionService,
	}
}

func (h *Handler) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := req.Path
	method := req.HTTPMethod

	switch {
	case method == "GET" && path == "/health":
		return jsonResponse(200, map[string]string{"status": "ok"})

	case method == "POST" && path == "/advisors":
		return h.createAdvisor(ctx, req)

	case method == "GET" && strings.HasPrefix(path, "/advisors/"):
		return h.getAdvisor(ctx, req)

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

func (h *Handler) createAdvisor(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body models.CreateAdvisorRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return errorResponse(400, "invalid request body")
	}

	advisor, err := h.advisorService.CreateAdvisor(ctx, body)
	if err != nil {
		return errorResponse(400, err.Error())
	}

	return jsonResponse(201, advisor)
}

func (h *Handler) getAdvisor(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := strings.TrimPrefix(req.Path, "/advisors/")

	advisor, err := h.advisorService.GetAdvisor(ctx, id)
	if err != nil {
		return errorResponse(400, err.Error())
	}
	if advisor == nil {
		return errorResponse(404, "advisor not found")
	}

	return jsonResponse(200, advisor)
}

func (h *Handler) createSession(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body models.CreateSessionRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return errorResponse(400, "invalid request body")
	}

	session, err := h.sessionService.StartSession(ctx, body)
	if err != nil {
		if err.Error() == "advisor not found" {
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
		if err.Error() == "session not found" || err.Error() == "advisor not found" {
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
