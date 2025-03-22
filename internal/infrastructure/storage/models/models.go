package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// User представляет модель пользователя.
type User struct {
	ID         string    `gorm:"primaryKey"`
	Name       string    `gorm:"size:255;not null"`
	TelegramID int64     `gorm:"unique;not null"`
	Phone      string    `gorm:"size:15"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// Restaurant представляет модель ресторана.
type Restaurant struct {
	ID        string    `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	Address   string    `gorm:"size:255;not null"`
	Phone     string    `gorm:"size:15"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Tables    []Table   `gorm:"foreignKey:RestaurantID"`
}

// Table представляет модель столика.
type Table struct {
	ID           string    `gorm:"primaryKey"`
	RestaurantID string    `gorm:"not null"`
	TableNumber  int       `gorm:"not null"`
	Capacity     int       `gorm:"not null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	PositionX    float64   `gorm:"not null"`
	PositionY    float64   `gorm:"not null"`
	PositionZ    float64   `gorm:"not null"`
}

// Reservation представляет модель бронирования.
type Reservation struct {
	ID           string     `gorm:"primaryKey"`
	UserID       string     `gorm:"not null"`
	RestaurantID string     `gorm:"not null"`
	StartTime    time.Time  `gorm:"not null"`
	EndTime      time.Time  `gorm:"not null"`
	Status       string     `gorm:"size:50;not null;check:status IN ('wait', 'sucess', 'canceled')"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	User         User       `gorm:"foreignKey:UserID"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID"`
	Capacity     int        `gorm:"not null"`
	Tables       []Table    `gorm:"many2many:reservation_tables;"`
	Contacts     Contact    `gorm:"type:jsonb"`
}

type Contact struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// Scan реализует интерфейс sql.Scanner
func (c *Contact) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Contact: expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, c)
}

// Value реализует интерфейс driver.Valuer
func (c Contact) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// ReservationTable представляет связь между бронированием и столиком.
type ReservationTable struct {
	ID            string    `gorm:"primaryKey"`
	ReservationID string    `gorm:"not null"`
	TableID       string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
