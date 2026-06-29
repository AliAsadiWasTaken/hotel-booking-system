package booking

import (
	"errors"
	"net/http"
	"time"

	"github.com/aliasadiwastaken/hotel-booking-system/internal/api"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for the booking domain.
type Handler struct {
	service *Service
}

// NewHandler returns a new Handler backed by the given service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /bookings.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RoomID   uuid.UUID `json:"room_id"   validate:"required"`
		UserID   uuid.UUID `json:"user_id"   validate:"required"`
		CheckIn  time.Time `json:"check_in"  validate:"required"`
		CheckOut time.Time `json:"check_out" validate:"required"`
	}

	if err := api.Decode(r, &body); err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if msg := api.Validate(body); msg != "" {
		api.WriteError(w, http.StatusBadRequest, msg)
		return
	}

	booking, err := h.service.Create(r.Context(), CreateInput{
		RoomID:   body.RoomID,
		UserID:   body.UserID,
		CheckIn:  body.CheckIn,
		CheckOut: body.CheckOut,
	})
	if errors.Is(err, ErrRoomNotFound) {
		api.WriteError(w, http.StatusNotFound, "room not found")
		return
	}
	if errors.Is(err, ErrRoomUnavailable) {
		api.WriteError(w, http.StatusConflict, "room is not available for the selected dates")
		return
	}
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "failed to create booking")
		return
	}

	api.WriteJSON(w, http.StatusCreated, booking)
}

// GetByID handles GET /bookings/{id}.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid booking id")
		return
	}

	booking, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		api.WriteError(w, http.StatusNotFound, "booking not found")
		return
	}
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "failed to get booking")
		return
	}

	api.WriteJSON(w, http.StatusOK, booking)
}

// Cancel handles DELETE /bookings/{id}.
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid booking id")
		return
	}

	booking, err := h.service.Cancel(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		api.WriteError(w, http.StatusNotFound, "booking not found")
		return
	}
	if errors.Is(err, ErrAlreadyCancelled) {
		api.WriteError(w, http.StatusConflict, "booking is already cancelled")
		return
	}
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "failed to cancel booking")
		return
	}

	api.WriteJSON(w, http.StatusOK, booking)
}
