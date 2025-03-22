package ports

import "github.com/gin-gonic/gin"

type IController interface {
	Authorize(*gin.Context)
	GetUserBookings(*gin.Context)
	UpdateBooking(*gin.Context)
	CreateBooking(*gin.Context)
	GetBookingsDate(*gin.Context)
	GetUserBookingsDate(*gin.Context)
}
