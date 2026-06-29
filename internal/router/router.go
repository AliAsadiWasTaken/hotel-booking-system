package router

import (
	"net/http"

	"github.com/aliasadiwastaken/hotel-booking-system/internal/booking"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/hotel"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/room"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/user"
)

func New(
	hotelHandler *hotel.Handler,
	roomHandler *room.Handler,
	userHandler *user.Handler,
	bookingHandler *booking.Handler,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthCheck)

	// Hotels
	mux.HandleFunc("POST /hotels", hotelHandler.Create)
	mux.HandleFunc("GET /hotels", hotelHandler.List)
	mux.HandleFunc("GET /hotels/{id}", hotelHandler.GetByID)

	// Rooms — nested under hotels for creation and listing
	mux.HandleFunc("POST /hotels/{id}/rooms", roomHandler.Create)
	mux.HandleFunc("GET /hotels/{id}/rooms", roomHandler.ListByHotelID)
	mux.HandleFunc("GET /rooms/{id}", roomHandler.GetByID)

	// Users
	mux.HandleFunc("POST /users", userHandler.Create)

	// Bookings
	mux.HandleFunc("POST /bookings", bookingHandler.Create)
	mux.HandleFunc("GET /bookings/{id}", bookingHandler.GetByID)
	mux.HandleFunc("DELETE /bookings/{id}", bookingHandler.Cancel)

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
