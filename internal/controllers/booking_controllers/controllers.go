package booking_controllers

import (
	"back_api/internal/dto"
	"back_api/internal/interfaces/middelware"
	"back_api/internal/service/booking_service"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type BookingController struct {
	logger  *slog.Logger
	service booking_service.IBookingService
	jwt     *middelware.Jwt
}

func NewBookingService(logger *slog.Logger, service booking_service.IBookingService, jwt *middelware.Jwt) *BookingController {
	return &BookingController{
		logger:  logger,
		service: service,
		jwt:     jwt,
	}
}

func (b *BookingController) CreateBooking(booking *dto.BookingDTO, c *gin.Context) {
	//booking, err := b.service.
}
