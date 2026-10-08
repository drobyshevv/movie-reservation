package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/handler/dto"
	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/go-playground/validator/v10"
)

type HallHandler struct {
	log       *slog.Logger
	service   HallService
	validator *validator.Validate
}

func NewHallHandler(log *slog.Logger, service HallService, validator *validator.Validate) *HallHandler {
	return &HallHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

type HallService interface {
	GetHalls() ([]models.Hall, error)
	GetHall(id int64) (*models.Hall, error)
	CreateHall(params models.CreateHallParams) (*models.Hall, error)
	UpdateHall(id int64, params models.UpdateHallParams) (*models.Hall, error)
	DeleteHall(id int64) error
}

func (h *HallHandler) GetHalls(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetHalls"
	log := h.log.With(
		slog.String("op", op),
	)

	halls, err := h.service.GetHalls()
	if err != nil {
		log.Error("failed to get halls", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := []dto.HallResponse{}

	for _, m := range halls {
		resp = append(resp, dto.HallResponse{
			ID:   m.ID,
			Name: m.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

func (h *HallHandler) GetHall(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetHall"
	log := h.log.With(
		slog.String("op", op),
	)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("failed to parse id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hall, err := h.service.GetHall(id)
	if err != nil {
		log.Error("failed to get hall", "err", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := dto.HallResponse{}
	resp.FromModel(hall)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

func (h *HallHandler) PostHall(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateHall"
	log := h.log.With(
		slog.String("op", op),
	)

	var req dto.CreateHallRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("failed to decode request body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Error("validation failed", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := req.ToModel()
	hall, err := h.service.CreateHall(*params)
	if err != nil {
		log.Error("failed to create hall", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.HallResponse{}
	resp.FromModel(hall)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

func (h *HallHandler) PutHall(w http.ResponseWriter, r *http.Request) {
	const op = "handler.UpdateHall"
	log := h.log.With(
		slog.String("op", op),
	)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("failed to parse id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req dto.UpdateHallRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("failed to decode request body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Error("validation failed", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := req.ToModel()
	hall, err := h.service.UpdateHall(id, *params)
	if err != nil {
		log.Error("failed to update hall", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.HallResponse{}
	resp.FromModel(hall)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

func (h *HallHandler) DeleteHall(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteHall"
	log := h.log.With(
		slog.String("op", op),
	)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("failed to parse id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.service.DeleteHall(id)
	if err != nil {
		log.Error("failed to delete hall", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
