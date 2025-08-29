package service

import (
	"context"
	"fmt"

	"github.com/aliffatulmf/mkp-eticket-service/internal/model"
	"github.com/aliffatulmf/mkp-eticket-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	FindByUsername(ctx context.Context, username string) (*model.Admin, error)
	FindByID(ctx context.Context, id int) (*model.Admin, error)
	Authenticate(ctx context.Context, req *model.AdminLoginRequest) (*model.Admin, error)
	Create(ctx context.Context, admin *model.Admin) error
	Delete(ctx context.Context, id int) error
}

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) AdminService {
	return &adminService{repo: repo}
}

func (s *adminService) FindByUsername(ctx context.Context, username string) (*model.Admin, error) {
	return s.repo.FindByUsername(ctx, username)
}

func (s *adminService) FindByID(ctx context.Context, id int) (*model.Admin, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *adminService) Authenticate(ctx context.Context, req *model.AdminLoginRequest) (*model.Admin, error) {
	admin, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	return admin, nil
}

func (s *adminService) Create(ctx context.Context, admin *model.Admin) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin.Password = string(hashed)

	return s.repo.Create(ctx, admin)
}

func (s *adminService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
