package dto

type UserDTO struct {
	UserUUID   string `json:"user_uuid"`
	TelegramId string `json:"telegramId"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Phone      string `json:"phone"`
}
