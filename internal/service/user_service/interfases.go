package user_service

import (
	"back_api/internal/dto"
)

type IUserService interface {
	IfUserExists(userId string) (bool, error)
	CreateUser(user *dto.UserDTO) (*dto.UserDTO, error)
	ValidateTelegramHash(telegramHash string, data map[string]string) (bool, error)
	GetUserForTelegram(telegramID string) (*dto.UserDTO, error)
	GetUserForUuid(uuid string) (*dto.UserDTO, error)
}
