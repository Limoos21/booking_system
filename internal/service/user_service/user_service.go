package user_service

import (
	"back_api/internal/dto"
	"back_api/internal/repository"
	"back_api/internal/service"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"time"
)

type UserService struct {
	logger     *slog.Logger
	repository repository.IRepository
	tokenBot   string
}

func NewUserService(logger *slog.Logger, repository repository.IRepository, tokenBot string) *UserService {
	return &UserService{logger: logger, repository: repository, tokenBot: tokenBot}
}

func (u UserService) ValidateTelegramHash(telegramHash string, data map[string]string) (bool, error) {
	// Проверяем, что все необходимые поля присутствуют
	requiredFields := []string{"id", "first_name", "auth_date", "hash"}
	for _, field := range requiredFields {
		if _, ok := data[field]; !ok {
			return false, fmt.Errorf("missing required field: %s", field)
		}
	}

	// Проверяем, что auth_date не устарел
	authDate, err := strconv.ParseInt(data["auth_date"], 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid auth_date: %v", err)
	}
	if time.Now().Unix()-authDate > 86400 { // 86400 секунд = 24 часа
		return false, fmt.Errorf("auth_date is too old")
	}

	// Создаем data-check-string
	keys := make([]string, 0, len(data))
	for k := range data {
		if k != "hash" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var dataCheckArr []string
	for _, k := range keys {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", k, data[k]))
	}
	dataCheckString := strings.Join(dataCheckArr, "\n")

	// Вычисляем HMAC-SHA256
	secretKey := sha256.Sum256([]byte(u.tokenBot)) // Вычисляем SHA256 от токена бота
	h := hmac.New(sha256.New, secretKey[:])        // Преобразуем [32]byte в []byte
	h.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	// Сравниваем хэши
	return strings.ToLower(expectedHash) == strings.ToLower(telegramHash), nil
}

func (u UserService) IfUserExists(userId string) (bool, error) {
	if ok, err := u.repository.IfUserExists(userId); err != nil || !ok {
		if err != nil {
			u.logger.Error(err.Error())
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (u UserService) CreateUser(user *dto.UserDTO) (*dto.UserDTO, error) {
	uuid, err := service.GenerateUuid()
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}
	userModel, err := service.ConvertToDomainUser(*user)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}
	userModel.UserUuid = uuid
	newUser, err := u.repository.CreateUser(userModel)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}
	dtoUser, err := service.ConvertToDTOUser(*newUser)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}
	return dtoUser, nil
}

func (u UserService) GetUserForTelegram(telegramID string) (*dto.UserDTO, error) {
	user, err := u.repository.GetUserForTelegram(telegramID)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	dtoUser, err := service.ConvertToDTOUser(*user)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}
	return dtoUser, nil
}

func (u UserService) GetUserForUuid(uuid string) (*dto.UserDTO, error) {
	user, err := u.repository.GetUserForUuid(uuid)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}
	dtoUser, err := service.ConvertToDTOUser(*user)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}
	return dtoUser, nil
}
