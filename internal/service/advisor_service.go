package service

import (
	"context"
	"errors"
	"time"

	"GolangChatAdvisor_Project/internal/models"
	"GolangChatAdvisor_Project/internal/store"

	"github.com/google/uuid"
)

type CompanionService struct {
	repo *store.AdvisorRepository
}

func NewCompanionService(repo *store.AdvisorRepository) *CompanionService {
	return &CompanionService{repo: repo}
}

func (s *CompanionService) CreateCompanion(ctx context.Context, req models.CreateAdvisorRequest) (*models.Advisor, error) {
	if req.Name == "" || req.Persona == "" || req.Style == "" || req.Greeting == "" {
		return nil, errors.New("all companion fields are required")
	}

	companion := models.Advisor{
		ID:        "comp_" + uuid.NewString(),
		Name:      req.Name,
		Persona:   req.Persona,
		Style:     req.Style,
		Greeting:  req.Greeting,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if err := s.repo.CreateAdvisor(ctx, companion); err != nil {
		return nil, err
	}

	return &companion, nil
}

func (s *CompanionService) GetCompanion(ctx context.Context, id string) (*models.Advisor, error) {
	if id == "" {
		return nil, errors.New("companion id is required")
	}

	return s.repo.GetAdvisorByID(ctx, id)
}
