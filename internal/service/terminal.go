package service

import (
	"context"
	"time"

	"github.com/aliffatulmf/mkp-eticket-service/internal/model"
	"github.com/aliffatulmf/mkp-eticket-service/internal/repository"
	"github.com/google/uuid"
)

type TerminalService interface {
	List(ctx context.Context) ([]model.Terminal, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Terminal, error)
	Create(ctx context.Context, req *model.CreateTerminalRequest) (*model.Terminal, error)
	Update(ctx context.Context, id uuid.UUID, req *model.UpdateTerminalRequest) (*model.Terminal, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type terminalService struct {
	repo repository.TerminalRepository
}

func NewTerminalService(repo repository.TerminalRepository) TerminalService {
	return &terminalService{repo: repo}
}

func (s *terminalService) List(ctx context.Context) ([]model.Terminal, error) {
	return s.repo.List(ctx)
}

func (s *terminalService) FindByID(ctx context.Context, id uuid.UUID) (*model.Terminal, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *terminalService) Create(ctx context.Context, req *model.CreateTerminalRequest) (*model.Terminal, error) {
	terminal := &model.Terminal{
		ID:        uuid.New(),
		Code:      req.Code,
		Name:      req.Name,
		Address:   req.Address,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.Create(ctx, terminal)
	if err != nil {
		return nil, err
	}

	return terminal, nil
}

func (s *terminalService) Update(ctx context.Context, id uuid.UUID, req *model.UpdateTerminalRequest) (*model.Terminal, error) {
	terminal, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	terminal.Name = req.Name
	terminal.Address = req.Address
	terminal.IsActive = req.IsActive
	terminal.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, terminal)
	if err != nil {
		return nil, err
	}

	return terminal, nil
}

func (s *terminalService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
