package models

import (
	"booking_system/internal/domain"
	"time"
)

// ConvertUserToDomain конвертирует модель User в доменный объект User.
func ConvertUserToDomain(u *User) *domain.User {
	return &domain.User{
		ID:         u.ID,
		Name:       u.Name,
		TelegramID: u.TelegramID,
		Phone:      u.Phone,
	}
}

// ConvertRestaurantToDomain конвертирует модель Restaurant в доменный объект Restaurant.
func ConvertRestaurantToDomain(r *Restaurant) *domain.Restaurant {
	return &domain.Restaurant{
		ID:      r.ID,
		Name:    r.Name,
		Address: r.Address,
		Phone:   r.Phone,
	}
}

// ConvertTableToDomain конвертирует модель Table в доменный объект Table.
func ConvertTableToDomain(t *Table) *domain.Table {
	return &domain.Table{
		ID:           t.ID,
		RestaurantID: t.RestaurantID,
		TableNumber:  t.TableNumber,
		Capacity:     t.Capacity,
		PositionZ:    t.PositionZ,
		PositionX:    t.PositionX,
		PositionY:    t.PositionY,
	}
}

// ConvertReservationToDomain конвертирует модель Reservation в доменный объект Reservation.
func ConvertReservationToDomain(r *Reservation) *domain.Reservation {
	return &domain.Reservation{
		ID:           r.ID,
		UserID:       r.UserID,
		RestaurantID: r.RestaurantID,
		StartTime:    r.StartTime,
		EndTime:      r.EndTime,
		Status:       r.Status,
		Contacts: domain.Contacts{
			Name:  r.Contacts.Name,
			Phone: r.Contacts.Phone,
		},
		Capacity: r.Capacity,
	}
}

// ConvertReservationTableToDomain конвертирует модель ReservationTable в доменный объект ReservationTable.
func ConvertReservationTableToDomain(rt *ReservationTable) *domain.ReservationTable {
	return &domain.ReservationTable{
		ID:            rt.ID,
		ReservationID: rt.ReservationID,
		TableID:       rt.TableID,
	}
}

// ConvertUserToModel конвертирует доменный объект User в модель User.
func ConvertUserToModel(u *domain.User) *User {
	return &User{
		ID:         u.ID,
		Name:       u.Name,
		TelegramID: u.TelegramID,
		Phone:      u.Phone,
		CreatedAt:  time.Now(),
	}
}

// ConvertRestaurantToModel конвертирует доменный объект Restaurant в модель Restaurant.
func ConvertRestaurantToModel(r *domain.Restaurant) *Restaurant {
	return &Restaurant{
		ID:        r.ID,
		Name:      r.Name,
		Address:   r.Address,
		Phone:     r.Phone,
		CreatedAt: time.Now(),
	}
}

// ConvertTableToModel конвертирует доменный объект Table в модель Table.
func ConvertTableToModel(t *domain.Table) *Table {
	return &Table{
		ID:           t.ID,
		RestaurantID: t.RestaurantID,
		TableNumber:  t.TableNumber,
		Capacity:     t.Capacity,
		CreatedAt:    time.Now(),
		PositionZ:    t.PositionZ,
		PositionY:    t.PositionY,
		PositionX:    t.PositionX,
	}
}

// ConvertReservationToModel конвертирует доменный объект Reservation в модель Reservation.
func ConvertReservationToModel(r *domain.Reservation) *Reservation {
	return &Reservation{
		ID:           r.ID,
		UserID:       r.UserID,
		RestaurantID: r.RestaurantID,
		StartTime:    r.StartTime,
		EndTime:      r.EndTime,
		Status:       r.Status,
		CreatedAt:    time.Now(),
		Capacity:     r.Capacity,
		Contacts: Contact{
			Name:  r.Contacts.Name,
			Phone: r.Contacts.Phone,
		},
	}
}

// ConvertReservationTableToModel конвертирует доменный объект ReservationTable в модель ReservationTable.
func ConvertReservationTableToModel(rt *domain.ReservationTable) *ReservationTable {
	return &ReservationTable{
		ID:            rt.ID,
		ReservationID: rt.ReservationID,
		TableID:       rt.TableID,
		CreatedAt:     time.Now(),
	}
}
