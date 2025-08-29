package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aliffatulmf/mkp-eticket-service/internal/model"
	"github.com/aliffatulmf/mkp-eticket-service/internal/service"
	"github.com/aliffatulmf/mkp-eticket-service/internal/validator"
	"github.com/go-chi/chi/v5"
)

type AdminHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type adminHandler struct {
	service service.AdminService
}

func NewAdminHandler(service service.AdminService) AdminHandler {
	return &adminHandler{service: service}
}

func (h *adminHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	admin := &model.Admin{
		Username: req.Username,
		Password: req.Password,
	}

	err := h.service.Create(r.Context(), admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"data": admin,
	})
}

func (h *adminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
