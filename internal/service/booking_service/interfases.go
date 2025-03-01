package booking_service

import (
	"back_api/internal/dto"
)

type IBookingService interface {
	CreateBooking(booking *dto.BookingDTO) (*dto.BookingDTO, error)
}
