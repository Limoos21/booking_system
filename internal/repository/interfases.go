package repository

import (
	"back_api/internal/domain"
	"time"
)

type IRepository interface {
	GetUserForUuid(uuid string) (*domain.User, error)
	CreateUser(user *domain.User) (*domain.User, error)
	IfUserExists(userUuid string) (bool, error)
	GetUserForTelegram(telegramId string) (*domain.User, error)
	GetAvailbleTable(tableUUid string, dateS, dateE time.Time) (*domain.Table, error)
	CreateBooking(booking *domain.Booking, table []string) (*domain.Booking, error)
}
