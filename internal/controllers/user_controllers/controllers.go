package user_controllers

import (
	"back_api/internal/controllers"
	"back_api/internal/dto"
	"back_api/internal/interfaces/middelware"
	"back_api/internal/service/user_service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type UserController struct {
	logger  *slog.Logger
	service user_service.IUserService
	jwt     *middelware.Jwt
}

func NewUserController(logger *slog.Logger, service user_service.IUserService, jwt *middelware.Jwt) *UserController {
	return &UserController{logger: logger, service: service, jwt: jwt}
}

func (uc *UserController) AuthUser(c *gin.Context) {
	// Извлекаем параметры из запроса
	data := map[string]string{
		"id":         c.Query("id"),
		"first_name": c.Query("first_name"),
		"photo_url":  c.Query("photo_url"),
		"username":   c.Query("username"),
		"auth_date":  c.Query("auth_date"),
		"hash":       c.Query("hash"),
	}
	uc.logger.Info("data", data)
	// Проверяем, что хэш присутствует
	if data["hash"] == "" {
		uc.logger.Warn("Telegram hash is missing")
		controllers.NewResponse(false, "Telegram hash is missing", "", c, http.StatusBadRequest)
		return
	}

	if ok, err := uc.service.ValidateTelegramHash(data["hash"], data); !ok || err != nil {
		if err != nil {
			uc.logger.Error("Error validating telegram hash", err)
			controllers.NewResponse(false, "Error validate telegram hash", "", c, http.StatusUnauthorized)
			return
		}
		uc.logger.Warn("Error validating telegram hash")
		controllers.NewResponse(false, "Telegram hash invalid", "", c, http.StatusUnauthorized)
		return
	}

	user, err := uc.service.GetUserForTelegram(data["id"])
	if err != nil {
		uc.logger.Error(err.Error())
		controllers.NewResponse(false, "Internal server error", []string{}, c, http.StatusInternalServerError)
		return
	}

	if user == nil {
		user := &dto.UserDTO{
			TelegramId: data["id"],
			Name:       data["first_name"],
			Role:       "user",
			Phone:      "",
		}

		userDomain, err := uc.service.CreateUser(user)
		if err != nil {
			uc.logger.Error("Error create user", err.Error())
			controllers.NewResponse(false, "Internal server error", []string{}, c, http.StatusInternalServerError)
			return
		}

		jwt, err := uc.jwt.GenerateToken(userDomain.Name, userDomain.UserUUID)
		if err != nil {
			uc.logger.Error(err.Error())
			controllers.NewResponse(false, "Internal server error", []string{}, c, http.StatusInternalServerError)
			return
		}

		controllers.NewResponse(true, "", map[string]string{"token": jwt}, c, http.StatusOK)
		return
	} else {
		userDomain, err := uc.service.GetUserForTelegram(data["id"])
		if err != nil {
			uc.logger.Error(err.Error())
			controllers.NewResponse(false, "Internal server error", []string{}, c, http.StatusInternalServerError)
			return
		}
		jwt, err := uc.jwt.GenerateToken(userDomain.Name, userDomain.UserUUID)
		if err != nil {
			uc.logger.Error(err.Error())
			controllers.NewResponse(false, "Internal server error", []string{}, c, http.StatusInternalServerError)
			return
		}
		controllers.NewResponse(true, "", map[string]string{"token": jwt}, c, http.StatusOK)
		return
	}
}

func (uc *UserController) GetUser(userUuid string, c *gin.Context) {
	user, err := uc.service.GetUserForUuid(userUuid)
	if err != nil {
		uc.logger.Error(err.Error())
		controllers.NewResponse(false, "internal server error", []string{}, c, http.StatusInternalServerError)
	}
	controllers.NewResponse(true, "", user, c, http.StatusOK)
}
