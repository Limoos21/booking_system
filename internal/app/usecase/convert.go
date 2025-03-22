package usecase

import (
	"booking_system/internal/domain"
	"booking_system/internal/dto"
)

// ToUserDomain преобразует структуру UserDTO в User.
func toUserDomain(dto *dto.UserDTO) *domain.User {
	return &domain.User{
		ID:         dto.ID,
		Name:       dto.Name,
		TelegramID: dto.TelegramID,
		Phone:      dto.Phone,
	}
}

// ToRestaurantDomain преобразует структуру RestaurantDTO в Restaurant.
func toRestaurantDomain(dto *dto.RestaurantDTO) *domain.Restaurant {
	return &domain.Restaurant{
		ID:      dto.ID,
		Name:    dto.Name,
		Address: dto.Address,
		Phone:   dto.Phone,
	}
}

// ToTableDomain преобразует структуру TableDTO в Table.
func toTableDomain(dto *dto.TableDTO) *domain.Table {
	return &domain.Table{
		ID:           dto.ID,
		RestaurantID: dto.RestaurantID,
		TableNumber:  dto.TableNumber,
		Capacity:     dto.Capacity,
		PositionX:    dto.PositionX,
		PositionY:    dto.PositionY,
		PositionZ:    dto.PositionZ,
	}
}

// ToReservationDomain преобразует структуру ReservationDTO в Reservation.
func toReservationDomain(dto *dto.ReservationDTO) (*domain.Reservation, []*domain.Table) {
	reservation := &domain.Reservation{
		ID:           dto.ID,
		UserID:       dto.UserID,
		RestaurantID: dto.RestaurantID,
		StartTime:    dto.StartTime,
		EndTime:      dto.EndTime,
		Status:       dto.Status,
		Contacts: domain.Contacts{
			Name:  dto.Contacts.Name,
			Phone: dto.Contacts.Phone,
		},
		Capacity: dto.Capacity,
	}
	tables := make([]*domain.Table, 0, len(dto.Table))
	for _, table := range dto.Table {
		domainTable := toTableDomain(&table)
		tables = append(tables, domainTable)
	}
	return reservation, tables

}

// ToReservationTableDomain преобразует структуру ReservationTableDTO в ReservationTable.
func toReservationTableDomain(dto *dto.ReservationTableDTO) *domain.ReservationTable {
	return &domain.ReservationTable{
		ID:            dto.ID,
		ReservationID: dto.ReservationID,
		TableID:       dto.TableID,
	}
}

// FromUserDomain преобразует структуру User в UserDTO.
func fromUserDomain(domain *domain.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:         domain.ID,
		Name:       domain.Name,
		TelegramID: domain.TelegramID,
		Phone:      domain.Phone,
	}
}

// FromRestaurantDomain преобразует структуру Restaurant в RestaurantDTO.
func fromRestaurantDomain(domain *domain.Restaurant) *dto.RestaurantDTO {
	return &dto.RestaurantDTO{
		ID:      domain.ID,
		Name:    domain.Name,
		Address: domain.Address,
		Phone:   domain.Phone,
	}
}

// FromTableDomain преобразует структуру Table в TableDTO.
func fromTableDomain(domain *domain.Table) *dto.TableDTO {
	return &dto.TableDTO{
		ID:           domain.ID,
		RestaurantID: domain.RestaurantID,
		TableNumber:  domain.TableNumber,
		Capacity:     domain.Capacity,
		PositionX:    domain.PositionX,
		PositionY:    domain.PositionY,
		PositionZ:    domain.PositionZ,
	}
}

// FromReservationDomain преобразует структуру Reservation в ReservationDTO.
func fromReservationDomain(domain *domain.Reservation) *dto.ReservationDTO {
	return &dto.ReservationDTO{
		ID:           domain.ID,
		UserID:       domain.UserID,
		RestaurantID: domain.RestaurantID,
		StartTime:    domain.StartTime,
		EndTime:      domain.EndTime,
		Status:       domain.Status,
		Capacity:     domain.Capacity,
		Contacts: dto.ContactsDTO{
			Name:  domain.Contacts.Name,
			Phone: domain.Contacts.Phone,
		},
		Table: make([]dto.TableDTO, 0),
	}
}

// FromReservationTableDomain преобразует структуру ReservationTable в ReservationTableDTO.
func fromReservationTableDomain(domain *domain.ReservationTable) *dto.ReservationTableDTO {
	return &dto.ReservationTableDTO{
		ID:            domain.ID,
		ReservationID: domain.ReservationID,
		TableID:       domain.TableID,
	}
}
