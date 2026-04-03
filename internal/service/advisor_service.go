package service

import (
	"context"
	"errors"
	"time"

	"GolangChatAdvisor_Project/internal/models"
	"GolangChatAdvisor_Project/internal/store"

	"github.com/google/uuid"
)

type AdvisorService struct {
	repo *store.AdvisorRepository
}

func NewAdvisorService(repo *store.AdvisorRepository) *AdvisorService {
	return &AdvisorService{repo: repo}
}

func (s *AdvisorService) CreateAdvisor(ctx context.Context, req models.CreateAdvisorRequest) (*models.Advisor, error) {
	if req.Name == "" || req.Persona == "" || req.Style == "" || req.Greeting == "" {
		return nil, errors.New("all advisor fields are required")
	}

	advisor := models.Advisor{
		ID:        "comp_" + uuid.NewString(),
		Name:      req.Name,
		Persona:   req.Persona,
		Style:     req.Style,
		Greeting:  req.Greeting,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if err := s.repo.CreateAdvisor(ctx, advisor); err != nil {
		return nil, err
	}

	return &advisor, nil
}

func (s *AdvisorService) GetAdvisor(ctx context.Context, id string) (*models.Advisor, error) {
	if id == "" {
		return nil, errors.New("Advisor id is required")
	}

	return s.repo.GetAdvisorByID(ctx, id)
}
