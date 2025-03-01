package api_server

import (
	"back_api/internal/controllers/booking_controllers"
	"back_api/internal/dto"
	"back_api/internal/interfaces/middelware"
	"net/http"

	"github.com/gin-gonic/gin"
	"log/slog"
)

type ApiRouterBooking struct {
	controller booking_controllers.IBookingControllers
	logger     *slog.Logger
	jwt        *middelware.Jwt
}

func NewApiRouter(r *gin.RouterGroup, logger *slog.Logger, controller booking_controllers.IBookingControllers, jwt *middelware.Jwt) {
	apiRouter := &ApiRouterBooking{controller: controller, logger: logger, jwt: jwt}
	r.POST("/booking/", jwt.JwtMiddleware(), apiRouter.createBooking)

}

func (a *ApiRouterBooking) createBooking(c *gin.Context) {
	booking := &dto.BookingDTO{}
	if err := c.ShouldBindJSON(&booking); err != nil {
		a.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, "invalid data")
		return
	}
	a.controller.CreateBooking(booking, c)
}
