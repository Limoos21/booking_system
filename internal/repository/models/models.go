package models

import (
	"back_api/internal/domain"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"time"
)

type User struct {
	gorm.Model
	UserUUID   string `gorm:"type:varchar(36);primaryKey"`
	TelegramID string `gorm:"type:varchar(255);not null"`
	Name       string `gorm:"type:varchar(255);not null"`
	Role       string `gorm:"type:varchar(50);not null"`
}

func ConvertUserMTODomain(u *User) (*domain.User, error) {
	domainUser := &domain.User{
		UserUuid:   u.UserUUID,
		TelegramId: u.TelegramID,
		Name:       u.Name,
		Role:       u.Role,
	}
	return domainUser, nil
}

func ConvertDomainUserTOMUser(u *domain.User) (*User, error) {
	modelUser := &User{
		UserUUID:   u.UserUuid,
		Name:       u.Name,
		Role:       u.Role,
		TelegramID: u.TelegramId,
	}
	return modelUser, nil
}

type Restaurant struct {
	gorm.Model
	UUIDRestaurant string    `gorm:"type:varchar(36);primaryKey"`
	Name           string    `gorm:"type:varchar(255);not null"`
	Address        string    `gorm:"type:varchar(255);not null"`
	Open           time.Time `gorm:"type:time;not null"`
	Close          time.Time `gorm:"type:time;not null"`
}

type Table struct {
	gorm.Model
	UUIDTable      string     `gorm:"type:varchar(36);primaryKey"`
	UUIDRestaurant string     `gorm:"type:varchar(36);not null"`
	TableID        string     `gorm:"type:varchar(255);not null"`
	PositionX      int        `gorm:"not null"`
	PositionY      int        `gorm:"not null"`
	MaxUser        int        `gorm:"not null"`
	Restaurant     Restaurant `gorm:"foreignKey:UUIDRestaurant;references:UUIDRestaurant"`
}

type Booking struct {
	gorm.Model
	UUIDBooking string    `gorm:"type:varchar(36);primaryKey"`
	DateStart   time.Time `gorm:"type:date;not null"`
	DateEnd     time.Time `gorm:"type:date;not null"`
	Comment     string    `gorm:"type:varchar(50)"`
	NumGuests   int       `gorm:"not null"`
	Status      string    `gorm:"type:varchar(50);not null"`
	UserUUID    string    `gorm:"type:varchar(36);not null"`
	User        User      `gorm:"foreignKey:UserUUID;references:UserUUID;constraint:OnDelete:CASCADE;"`
	Contact     string    `gorm:"type:varchar(50);not null"`
}

type BookingTable struct {
	gorm.Model
	UUIDTable   string  `gorm:"type:varchar(36);not null"`
	UUIDBooking string  `gorm:"type:varchar(36);not null"`
	Table       Table   `gorm:"foreignKey:UUIDTable;references:UUIDTable"`
	Booking     Booking `gorm:"foreignKey:UUIDBooking;references:UUIDBooking"`
}

func Migrate(db *gorm.DB, logger *slog.Logger) {
	err := db.AutoMigrate(&BookingTable{}, &Table{}, &Restaurant{}, &Booking{}, &User{})
	if err != nil {
		log.Fatalf("Ошибка при создании таблиц: %v", err)
	}
}
