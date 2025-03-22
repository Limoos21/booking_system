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
	r.POST("/auth/test", rout.AuthTest)
	// Роуты, связанные с аутентификацией и пользователем
	r.GET("/auth/telegram", rout.Auth)
	r.PATCH("/me", jwt.JwtMiddleware(), rout.UpdateUserInfo)
	r.GET("/me", jwt.JwtMiddleware(), rout.GetUser)
	// Роуты, связанные с бронированиями пользователя
	r.GET("/booking/me", jwt.JwtMiddleware(), rout.GetUserBooking)
	r.GET("/booking/me/:date", jwt.JwtMiddleware(), rout.GetUserBooking)

	// Роуты для работы с конкретными бронированиями
	r.GET("/booking/:id", rout.GetBooking)
	r.PATCH("/booking/:id", jwt.JwtMiddleware(), rout.UpdateBooking)
	r.PATCH("/booking/:id/:status", jwt.JwtMiddleware(), rout.UpdateStatus)

	// Роуты для работы с бронированиями в ресторанах
	r.POST("/:restaurantId/booking", jwt.JwtMiddleware(), rout.CreateBooking)
	r.GET("/:restaurantId/bookings/:date", rout.GetUserBookingsDate)

}

func (r Router) AuthTest(c *gin.Context) {
	r.controllers.AuthTest(c)
}

func (r Router) GetUser(c *gin.Context) {
	r.controllers.GetUserInfo(c)
}

func (r Router) UpdateUserInfo(c *gin.Context) {
	r.controllers.UpdateUserInfo(c)
}

func (r Router) UpdateStatus(c *gin.Context) {
	r.controllers.UpdateBooking(c)
}

func (r Router) Auth(c *gin.Context) {
	r.controllers.Authorize(c)
}

func (r Router) GetBooking(c *gin.Context) {
	r.controllers.GetBooking(c)
}

func (r Router) GetUserBookingsDate(c *gin.Context) {
	r.controllers.GetUserBookingsDate(c)
}

func (r Router) GetUserBooking(c *gin.Context) {
	r.controllers.GetUserBookings(c)

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
