package user

import (
	"errors"
	"net/http"

	"github.com/aliasadiwastaken/hotel-booking-system/internal/api"
)

// Handler handles HTTP requests for the user domain.
type Handler struct {
	service *Service
}

// NewHandler returns a new Handler backed by the given service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /users.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name"  validate:"required"`
		Email     string `json:"email"      validate:"required,email"`
		Password  string `json:"password"   validate:"required,min=8"`
	}

	if err := api.Decode(r, &body); err != nil {
		api.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if msg := api.Validate(body); msg != "" {
		api.WriteError(w, http.StatusBadRequest, msg)
		return
	}

	user, err := h.service.Create(r.Context(), CreateInput{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
	})
	if errors.Is(err, ErrEmailTaken) {
		api.WriteError(w, http.StatusConflict, "email already in use")
		return
	}
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	api.WriteJSON(w, http.StatusCreated, user)
}
