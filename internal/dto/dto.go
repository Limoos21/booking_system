package dto

import (
	"time"
)

// ReservationTableDTO — структура для передачи данных о связи бронирования и стола.
type ReservationTableDTO struct {
	ID            string `json:"id"`
	ReservationID string `json:"reservation_id"`
	TableID       string `json:"table_id"`
}

// UserDTO — структура для передачи данных о пользователе.
type UserDTO struct {
	ID         string `json:"id"`              // Уникальный идентификатор пользователя
	Name       string `json:"name"`            // Имя пользователя
	TelegramID int64  `json:"telegram_id"`     // Уникальный идентификатор Telegram
	Phone      string `json:"phone,omitempty"` // Номер телефона (опционально)
}

// RestaurantDTO — структура для передачи данных о ресторане.
type RestaurantDTO struct {
	ID      string `json:"id"`      // Уникальный идентификатор ресторана
	Name    string `json:"name"`    // Название ресторана
	Address string `json:"address"` // Адрес ресторана
	Phone   string `json:"phone"`   // Телефон ресторана
}

// TableDTO — структура для передачи данных о столе.
type TableDTO struct {
	ID           string  `json:"id"`
	RestaurantID string  `json:"restaurant_id"`
	TableNumber  int     `json:"table_number"`
	Capacity     int     `json:"capacity"`
	PositionX    float64 `json:"position_x"`
	PositionY    float64 `json:"position_y"`
	PositionZ    float64 `json:"position_z"`
}

// ReservationDTO — структура для передачи данных о бронировании.
type ReservationDTO struct {
	ID           string      `json:"id"`
	UserID       string      `json:"user_id"`
	RestaurantID string      `json:"restaurant_id"`
	StartTime    time.Time   `json:"start_time"`
	EndTime      time.Time   `json:"end_time"`
	Status       string      `json:"status"` // Статус брони (отменена, подтверждена, в ожидании подтверждения)
	Table        []TableDTO  `json:"table"`
	Contacts     ContactsDTO `json:"contacts"`
	Capacity     int         `json:"capacity"`
}

type ContactsDTO struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type AvaibleTableDTO struct {
	TableDTO
	IsAvaible bool `json:"is_avaible"`
}
