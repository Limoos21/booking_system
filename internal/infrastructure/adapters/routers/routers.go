package routers

import (
	"booking_system/cmd/providers/middelware"
	"booking_system/internal/app/ports"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Router struct {
	logger      *slog.Logger
	controllers ports.IController
}

func New(server *gin.Engine, logger *slog.Logger, controllers ports.IController, jwt *middelware.Jwt) {
	rout := Router{
		logger:      logger,
		controllers: controllers,
	}
	r := server.Group("api/v1")

	// Роуты, связанные с аутентификацией и пользователем
	r.GET("/auth/telegram", rout.Auth)
	r.PATCH("/me", rout.UpdateInfo)
	r.GET("/me", rout.GetUser)
	// Роуты, связанные с бронированиями пользователя
	r.GET("/booking/me", jwt.JwtMiddleware(), rout.GetUserBooking)
	r.GET("/booking/me/:date", jwt.JwtMiddleware(), rout.GetUserBooking)

	// Роуты для работы с конкретными бронированиями
	r.GET("/booking/:id", rout.GetBooking)
	r.PATCH("/booking/:id", jwt.JwtMiddleware(), rout.UpdateBooking)
	r.PATCH("/booking/:id/:status", jwt.JwtMiddleware(), rout.UpdateStatus)

	// Роуты для работы с бронированиями в ресторанах
	r.POST("/:restaurantId/booking", jwt.JwtMiddleware(), rout.CreateBooking)
	r.GET("/:restaurantId/bookings/:date", rout.GetBookingDate)

}

func (r Router) UpdateStatus(c *gin.Context) {
	r.controllers.UpdateBooking(c)
}

func (r Router) Auth(c *gin.Context) {
	r.controllers.Authorize(c)
}

func (r Router) GetBooking(c *gin.Context) {

}

func (r Router) GetUserBooking(c *gin.Context) {
	date := c.Param("date")
	if date == "" {
		r.controllers.GetUserBookings(c)
	} else {
		r.controllers.GetUserBookingsDate(c)
	}

}

func (r Router) CreateBooking(c *gin.Context) {
	r.controllers.CreateBooking(c)
}

func (r Router) GetBookingDate(c *gin.Context) {
	r.controllers.GetBookingsDate(c)
}

func (r Router) UpdateBooking(c *gin.Context) {
	r.controllers.UpdateBooking(c)
}
