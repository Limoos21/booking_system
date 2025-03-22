package ports

import (
	"booking_system/internal/domain"
	"time"
)

type IStorage interface {
	GetTable(tableId string) (*domain.Table, error)
	GetTablesByReservationID(reservationID string) ([]domain.Table, error)
	IsTableAvailable(tableID string, startTime, endTime time.Time) (bool, error)
	CreateUser(user domain.User) (domain.User, error)
	//CheckUserForTelegram проверка существования пользвоателя по telegramId
	CheckUserForTelegram(telegramId int64) (bool, domain.User, error)
	// GetUserForId получение пользователя по Id
	GetUserForId(user domain.User) (*domain.User, error)
	// GetReservationForId получение резервации (бронирования) по Id
	GetReservationForId(id string) (*domain.Reservation, error)
	// CreateReservation создание резервации(бронирования), возвращает id созданной резервации
	CreateReservation(reservation *domain.Reservation, tableIDs map[string]string) (string, error)
	// GetUserReservationsUser получение всех резерваций пользователя
	GetUserReservationsUser(userId string) ([]*domain.Reservation, error)
	// GetUserReservationsUserForDate получение всех резерваций пользователя на указанную дату
	GetUserReservationsUserForDate(date time.Time, userID string) ([]*domain.Reservation, error)
	// GetReservationsForDate получение всех резерваций на указанную дату
	GetReservationsForDate(date time.Time) ([]*domain.Reservation, error)
	// UpdateReservation обноваление резервации
	UpdateReservation(reservation *domain.Reservation) (bool, error)
	GetTablesWithAvailability(restaurantID string, dateTime time.Time) ([]domain.TableAvailability, error)
}
