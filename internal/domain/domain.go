package domain

import (
	"errors"
	"time"
)

// User представляет пользователя системы.
type User struct {
	ID         string // Уникальный идентификатор пользователя
	Name       string // Имя пользователя
	TelegramID int64  // Уникальный идентификатор Telegram
	Phone      string // Номер телефона (опционально)
}

// Restaurant представляет ресторан.
type Restaurant struct {
	ID      string // Уникальный идентификатор ресторана
	Name    string // Название ресторана
	Address string // Адрес ресторана
	Phone   string // Телефон ресторана
}

// Table представляет столик в ресторане.
type Table struct {
	ID           string
	RestaurantID string
	TableNumber  int
	Capacity     int
	PositionX    float64
	PositionY    float64
	PositionZ    float64
}

// Reservation представляет бронь столика.
type Reservation struct {
	ID           string
	UserID       string
	RestaurantID string
	StartTime    time.Time
	EndTime      time.Time
	Status       string // Статус брони (отменена, подтверждена, в ожидании подтверждения)
	Capacity     int
	Contacts     Contacts
}

type Contacts struct {
	Name  string
	Phone string
}

func (rv Reservation) CheckDate() (bool, error) {
	now := time.Now().Local()

	if rv.StartTime.Before(now) {
		return false, errors.New("StartTime должна быть позже или равна текущему времени")
	}
	duration := rv.EndTime.Sub(rv.StartTime)
	if duration > 2*time.Hour {
		return false, errors.New("разница между StartTime и EndTime не должна превышать 2 часа")
	}

	return true, nil
}

// ReservationTable представляет связь между бронированием и столиком.
type ReservationTable struct {
	ID            string
	ReservationID string
	TableID       string
}

type TableAvailability struct {
	Table
	IsAvailable bool `json:"is_available"`
}
