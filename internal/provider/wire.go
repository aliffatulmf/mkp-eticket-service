//go:build wireinject
// +build wireinject

package provider

import (
	"github.com/aliffatulmf/mkp-eticket-service/internal/auth"
	"github.com/aliffatulmf/mkp-eticket-service/internal/handler"
	"github.com/aliffatulmf/mkp-eticket-service/internal/repository"
	"github.com/aliffatulmf/mkp-eticket-service/internal/service"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTerminalHandler(db *pgxpool.Pool) handler.TerminalHandler {
	wire.Build(
		repository.NewTerminalRepository,
		service.NewTerminalService,
		handler.NewTerminalHandler,
	)
	return nil
}

func NewAuthHandler(db *pgxpool.Pool, jwtSecret string) handler.AuthHandler {
	wire.Build(
		repository.NewAdminRepository,
		service.NewAdminService,
		auth.NewService,
		handler.NewAuthHandler,
	)
	return nil
}

func NewAdminHandler(db *pgxpool.Pool) handler.AdminHandler {
	wire.Build(
		repository.NewAdminRepository,
		service.NewAdminService,
		handler.NewAdminHandler,
	)
	return nil
}
