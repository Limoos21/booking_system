package booking_service

import (
	"back_api/internal/dto"
	"back_api/internal/repository"
	"back_api/internal/service"
	"errors"
	"log/slog"
)

type BookingService struct {
	logger     *slog.Logger
	repository repository.IRepository
}

func NewBookingService(logger *slog.Logger, repository repository.IRepository) *BookingService {
	return &BookingService{logger, repository}
}

func (b *BookingService) CreateBooking(booking *dto.BookingDTO) (*dto.BookingDTO, error) {
	domainBooking, table, err := service.ConvertTODomainBooking(*booking)
	if err != nil {
		return nil, err
	}
	if len(table) == 0 {
		b.logger.Warn("Count table is zero")
		return nil, errors.New("count table cant be is zero")
	}
	for _, val := range table {
		avaibleT, err := b.repository.GetAvailbleTable(val, domainBooking.DateStart, domainBooking.DateEnd)
		if err != nil || avaibleT == nil {
			if err != nil {
				b.logger.Error("Error check Table", "err", err)
				return nil, err
			}
			b.logger.Info("Table is busy", "table", val)
			return nil, nil
		}
	}
	newBooking, err := b.repository.CreateBooking(domainBooking, table)
	if err != nil {
		b.logger.Error("Error create Booking", "err", err)
		return nil, err
	}
	if newBooking == nil {
		b.logger.Warn("Create Booking", "err", err)
		return nil, nil
	}
	dtoBooking, err := service.ConvertToDTOBoking(*newBooking, table)
	if err != nil {
		return nil, err
	}
	return dtoBooking, nil

}
