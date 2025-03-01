package api_server

import (
	"back_api/internal/controllers/user_controllers"
	"back_api/internal/interfaces/middelware"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type ApiRoutersUser struct {
	logger      *slog.Logger
	controllers user_controllers.IUserController
	jwt         *middelware.Jwt
}

func NewApiRoutersUser(r *gin.RouterGroup, logger *slog.Logger, controllers user_controllers.IUserController, jwt *middelware.Jwt) {
	api := &ApiRoutersUser{logger: logger, controllers: controllers, jwt: jwt}
	r.GET("/auth/telgram/", api.AuthTelegram)
	r.GET("/me", api.jwt.JwtMiddleware(), api.GetUser)

}

func (api *ApiRoutersUser) GetUser(c *gin.Context) {
	claims := c.MustGet("claims").(*middelware.Claims)
	if claims == nil {
		api.logger.Warn("claims is nil")
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	api.controllers.GetUser(claims.UserUuid, c)
	return
}

func (api *ApiRoutersUser) AuthTelegram(c *gin.Context) {
	api.controllers.AuthUser(c)
	return
}
