package booking_controllers

import (
	"back_api/internal/dto"
	"github.com/gin-gonic/gin"
)

type IBookingControllers interface {
	CreateBooking(booking *dto.BookingDTO, c *gin.Context)
}
