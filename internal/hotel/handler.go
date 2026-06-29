package hotel

import (
	"errors"
	"net/http"

	"github.com/aliasadiwastaken/hotel-booking-system/internal/api"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for the hotel domain.
type Handler struct {
	service *Service
}

// NewHandler returns a new Handler backed by the given service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /hotels.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name    string `json:"name"    validate:"required"`
		Address string `json:"address" validate:"required"`
		City    string `json:"city"    validate:"required"`
		Country string `json:"country" validate:"required"`
	}

	if err := api.Decode(r, &body); err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if msg := api.Validate(body); msg != "" {
		api.WriteError(w, http.StatusBadRequest, msg)
		return
	}

	hotel, err := h.service.Create(r.Context(), CreateInput{
		Name:    body.Name,
		Address: body.Address,
		City:    body.City,
		Country: body.Country,
	})
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "failed to create hotel")
		return
	}

	api.WriteJSON(w, http.StatusCreated, hotel)
}

// GetByID handles GET /hotels/{id}.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid hotel id")
		return
	}

	hotel, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		api.WriteError(w, http.StatusNotFound, "hotel not found")
		return
	}
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "failed to get hotel")
		return
	}

	api.WriteJSON(w, http.StatusOK, hotel)
}

// List handles GET /hotels.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	hotels, err := h.service.List(r.Context())
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "failed to list hotels")
		return
	}

	api.WriteJSON(w, http.StatusOK, hotels)
}
