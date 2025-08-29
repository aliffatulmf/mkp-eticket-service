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

var adminSet = wire.NewSet(
	repository.NewAdminRepository,
	service.NewAdminService,
)

func NewTerminalHandler(db *pgxpool.Pool) handler.TerminalHandler {
	wire.Build(
		repository.NewTerminalRepository,
		service.NewTerminalService,
		handler.NewTerminalHandler,
	)
	return nil
}

func NewAuthHandler(db *pgxpool.Pool, jwt auth.JWTService) handler.AuthHandler {
	wire.Build(
		adminSet,
		handler.NewAuthHandler,
	)
	return nil
}

func NewAdminHandler(db *pgxpool.Pool) handler.AdminHandler {
	wire.Build(
		adminSet,
		handler.NewAdminHandler,
	)
	return nil
}
