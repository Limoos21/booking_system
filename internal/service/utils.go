package service

import (
	"back_api/internal/domain"
	"back_api/internal/dto"
	"github.com/google/uuid"
	"time"
)

func GenerateUuid() (string, error) {
	uuId, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uuId.String(), nil
}

func ConvertToDTOUser(u domain.User) (*dto.UserDTO, error) {
	newUser := &dto.UserDTO{
		UserUUID:   u.UserUuid,
		TelegramId: u.TelegramId,
		Name:       u.Name,
		Role:       u.Role,
		Phone:      u.Phone,
	}
	return newUser, nil
}

func ConvertToDomainUser(u dto.UserDTO) (*domain.User, error) {
	newUser := &domain.User{
		TelegramId: u.TelegramId,
		Name:       u.Name,
		Role:       u.Role,
		Phone:      u.Phone,
	}
	return newUser, nil
}

func ConvertTODomainBooking(b dto.BookingDTO) (*domain.Booking, []string, error) {
	dateStart, err := time.Parse("2006-01-02 15:04:05", b.DateStart)
	if err != nil {
		return nil, nil, err
	}
	dateEnd, err := time.Parse("2006-01-02 15:04:05", b.DateEnd)
	if err != nil {
		return nil, nil, err
	}
	newBooking := &domain.Booking{
		BookingUuid: b.BookingUuid,
		DateStart:   dateStart,
		DateEnd:     dateEnd,
		Comment:     b.Comment,
		NumGuests:   b.NumGuests,
		Status:      b.Status,
		UserUuid:    b.UserUuid,
	}
	return newBooking, b.Table, nil

}

func ConvertToDTOBoking(b domain.Booking, table []string) (*dto.BookingDTO, error) {
	dateStart := b.DateStart.Format("2006-01-02 15:04:05")
	dateEnd := b.DateEnd.Format("2006-01-02 15:04:05")
	newBooking := &dto.BookingDTO{
		BookingUuid: b.BookingUuid,
		DateStart:   dateStart,
		DateEnd:     dateEnd,
		Comment:     b.Comment,
		NumGuests:   b.NumGuests,
		Status:      b.Status,
		UserUuid:    b.UserUuid,
		Table:       table,
	}
	return newBooking, nil
}
