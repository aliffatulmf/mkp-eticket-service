package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aliffatulmf/mkp-eticket-service/internal/auth"
	"github.com/aliffatulmf/mkp-eticket-service/internal/model"
	"github.com/aliffatulmf/mkp-eticket-service/internal/service"
	"github.com/aliffatulmf/mkp-eticket-service/internal/validator"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	adminService service.AdminService
	jwtService   *auth.Service
}

func NewAuthHandler(adminService service.AdminService, jwtService *auth.Service) AuthHandler {
	return &authHandler{
		adminService: adminService,
		jwtService:   jwtService,
	}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	admin, err := h.adminService.Authenticate(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.jwtService.GenerateToken(admin.Username, "admin")
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(admin.Username)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(15 * time.Minute).Unix()

	response := model.AdminLoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Username:     admin.Username,
		Role:         "admin",
		ExpiresAt:    expiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	refreshClaims, err := h.jwtService.VerifyToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	admin, err := h.adminService.FindByUsername(r.Context(), refreshClaims.Username)
	if err != nil {
		http.Error(w, "Admin not found", http.StatusUnauthorized)
		return
	}

	token, err := h.jwtService.GenerateToken(admin.Username, "admin")
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(15 * time.Minute).Unix()

	response := model.RefreshTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
