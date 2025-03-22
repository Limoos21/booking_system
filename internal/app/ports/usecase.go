package ports

import (
	"booking_system/internal/dto"
	"time"
)

type IUseCase interface {
	AuthUser(dto dto.UserDTO) (dto.UserDTO, string, error)
	GetReservationForDate(date *time.Time) ([]dto.ReservationDTO, error)
	CreateReservation(dto dto.ReservationDTO) (dto.ReservationDTO, error)
	GetUserReservations(userId string) ([]dto.ReservationDTO, error)
	GetUserReservationsDate(date *time.Time, userId string) ([]dto.ReservationDTO, error)
	ValidateTelegramHash(telegramHash string, data map[string]string) (bool, error)
	GetReservationForId(reservationId string) (dto.ReservationDTO, error)
	UpdateReservation(dto dto.ReservationDTO) (bool, error)
	GetTableForReservationDate(date time.Time, restaurantId string) ([]dto.AvaibleTableDTO, error)
}
