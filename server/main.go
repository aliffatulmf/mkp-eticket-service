package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aliffatulmf/mkp-eticket-service/internal/auth"
	"github.com/aliffatulmf/mkp-eticket-service/internal/config"
	"github.com/aliffatulmf/mkp-eticket-service/internal/database"
	"github.com/aliffatulmf/mkp-eticket-service/internal/handler"
	"github.com/aliffatulmf/mkp-eticket-service/internal/middleware"
	"github.com/aliffatulmf/mkp-eticket-service/internal/provider"
	"github.com/aliffatulmf/mkp-eticket-service/internal/repository"
	"github.com/aliffatulmf/mkp-eticket-service/internal/service"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	cfg := config.Load()

	pool, err := database.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	defer pool.Close()

	jwtService := auth.NewService(cfg.JWTSecret)

	adminRepo := repository.NewAdminRepository(pool)
	adminService := service.NewAdminService(adminRepo)
	authHandler := handler.NewAuthHandler(adminService, jwtService)
	adminHandler := handler.NewAdminHandler(adminService)

	terminalHandler := provider.NewTerminalHandler(pool)

	r := chi.NewMux()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.RefreshToken)
		})

		r.Route("/admins", func(r chi.Router) {
			r.Post("/", adminHandler.Create)
			r.Delete("/{id}", adminHandler.Delete)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.AdminAuthMiddleware(jwtService))

			r.Route("/terminals", func(r chi.Router) {
				r.Get("/", terminalHandler.List)
				r.Get("/{id}", terminalHandler.FindByID)
				r.Post("/", terminalHandler.Create)
				r.Put("/{id}", terminalHandler.Update)
				r.Delete("/{id}", terminalHandler.Delete)
			})
		})
	})

	log.Printf("Server starting on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
