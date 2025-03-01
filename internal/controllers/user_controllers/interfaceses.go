package user_controllers

import (
	"github.com/gin-gonic/gin"
)

type IUserController interface {
	AuthUser(c *gin.Context)
	GetUser(userUuid string, c *gin.Context)
}
