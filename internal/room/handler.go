package room

import (
	"errors"
	"net/http"

	"github.com/aliasadiwastaken/hotel-booking-system/internal/api"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for the room domain.
type Handler struct {
	service *Service
}

// NewHandler returns a new Handler backed by the given service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /hotels/{id}/rooms.
// The hotel ID is read from the URL path rather than the request body
// to enforce the ownership relationship at the routing level.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	hotelID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid hotel id")
		return
	}

	var body struct {
		Name          string  `json:"name"            validate:"required"`
		BedCount      int     `json:"bed_count"       validate:"required,gt=0"`
		Capacity      int     `json:"capacity"        validate:"required,gt=0"`
		Quantity      int     `json:"quantity"        validate:"required,gt=0"`
		PricePerNight float64 `json:"price_per_night" validate:"required,gt=0"`
	}

	if err := api.Decode(r, &body); err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if msg := api.Validate(body); msg != "" {
		api.WriteError(w, http.StatusBadRequest, msg)
		return
	}

	room, err := h.service.Create(r.Context(), CreateInput{
		HotelID:       hotelID,
		Name:          body.Name,
		BedCount:      body.BedCount,
		Capacity:      body.Capacity,
		Quantity:      body.Quantity,
		PricePerNight: body.PricePerNight,
	})
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "failed to create room")
		return
	}

	api.WriteJSON(w, http.StatusCreated, room)
}

// GetByID handles GET /rooms/{id}.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid room id")
		return
	}

	room, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		api.WriteError(w, http.StatusNotFound, "room not found")
		return
	}
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "failed to get room")
		return
	}

	api.WriteJSON(w, http.StatusOK, room)
}

// ListByHotelID handles GET /hotels/{id}/rooms.
func (h *Handler) ListByHotelID(w http.ResponseWriter, r *http.Request) {
	hotelID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid hotel id")
		return
	}

	rooms, err := h.service.ListByHotelID(r.Context(), hotelID)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "failed to list rooms")
		return
	}

	api.WriteJSON(w, http.StatusOK, rooms)
}
