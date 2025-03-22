package storage

import (
	"booking_system/internal/domain"
	"booking_system/internal/infrastructure/storage/models"
	"errors"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type Storage struct {
	logger   *slog.Logger
	Database *gorm.DB
}

func New(logger *slog.Logger, db *gorm.DB) *Storage {
	return &Storage{
		logger:   logger,
		Database: db,
	}
}

// UpdateUser обновляет данные пользователя.
func (s *Storage) UpdateUser(user domain.User) (bool, error) {
	userModel := models.ConvertUserToModel(&user)
	result := s.Database.Save(userModel)
	if result.Error != nil {
		s.logger.Error("Failed to update user", "error", result.Error)
		return false, result.Error
	}
	s.logger.Info("User updated successfully", "userId", userModel.ID)
	return true, nil
}

// GetTablesWithAvailability возвращает все столы с пометками о их доступности на конкретную дату и время.
func (s *Storage) GetTablesWithAvailability(restaurantID string, dateTime time.Time) ([]domain.TableAvailability, error) {
	var tables []models.Table
	var reservations []models.Reservation

	// Получаем все столы для конкретного ресторана
	if err := s.Database.Where("restaurant_id = ?", restaurantID).Find(&tables).Error; err != nil {
		return nil, err
	}

	// Получаем все бронирования, которые пересекаются с указанной датой и временем
	// И загружаем связанные таблицы (Tables) для каждого бронирования
	if err := s.Database.
		Preload("Tables").
		Where("restaurant_id = ? AND start_time <= ? AND end_time >= ?", restaurantID, dateTime, dateTime).
		Find(&reservations).Error; err != nil {
		return nil, err
	}

	// Создаем мапу для быстрого поиска занятых столов
	occupiedTables := make(map[string]bool)
	for _, reservation := range reservations {
		for _, table := range reservation.Tables {
			occupiedTables[table.ID] = true
		}
	}

	// Формируем результат с пометками о доступности
	var result []domain.TableAvailability
	for _, table := range tables {
		isAvailable := !occupiedTables[table.ID]
		domainTable := *models.ConvertTableToDomain(&table)
		result = append(result, domain.TableAvailability{
			Table:       domainTable,
			IsAvailable: isAvailable,
		})
	}
	s.logger.Debug("Tables with availability: ", result)

	return result, nil
}

// GetTable возвращает столик по его ID.
func (s *Storage) GetTable(tableId string) (*domain.Table, error) {
	var table models.Table
	result := s.Database.First(&table, "id = ?", tableId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			s.logger.Info("Table not found", "tableId", tableId)
			return nil, errors.New("table not found")
		}
		s.logger.Error("Failed to get table", "error", result.Error)
		return nil, result.Error
	}
	return models.ConvertTableToDomain(&table), nil
}

// CheckUserForTelegram проверяет, существует ли пользователь с указанным Telegram ID.
func (s *Storage) CheckUserForTelegram(telegramId int64) (bool, domain.User, error) {
	var user models.User
	result := s.Database.First(&user, "telegram_id = ?", telegramId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			s.logger.Info("User not found", "telegramId", telegramId)
			return false, domain.User{}, nil
		}
		s.logger.Error("Failed to check user", "error", result.Error)
		return false, domain.User{}, result.Error
	}
	return true, *models.ConvertUserToDomain(&user), nil
}

// CreateUser создает нового пользователя.
func (s *Storage) CreateUser(user domain.User) (domain.User, error) {
	userModel := models.ConvertUserToModel(&user)
	result := s.Database.Create(userModel)
	if result.Error != nil {
		s.logger.Error("Failed to create user", "error", result.Error)
		return domain.User{}, result.Error
	}
	s.logger.Info("User created successfully", "userId", userModel.ID)
	return *models.ConvertUserToDomain(userModel), nil
}

// IsTableAvailable проверяет, свободен ли столик на указанное время.
func (s *Storage) IsTableAvailable(tableID string, startTime, endTime time.Time) (bool, error) {
	var count int64

	// Проверяем, есть ли бронирования, которые пересекаются с запрашиваемым временем
	err := s.Database.Model(&models.ReservationTable{}).
		Joins("JOIN reservations ON reservation_tables.reservation_id = reservations.id").
		Where("reservation_tables.table_id = ?", tableID).
		Where("(? <= reservations.end_time) AND (? >= reservations.start_time)", startTime, endTime).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	// Если count > 0, значит столик занят
	return count == 0, nil
}

func (s *Storage) GetUserForId(userUid string) (*domain.User, error) {
	var dbUser models.User
	result := s.Database.First(&dbUser, userUid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Пользователь не найден
		}
		return nil, result.Error
	}
	return models.ConvertUserToDomain(&dbUser), nil
}

func (s *Storage) GetUserReservationsUser(userId string) ([]*domain.Reservation, error) {
	var dbReservations []models.Reservation
	result := s.Database.Preload("User").Preload("Restaurant").Preload("Tables").Where("user_id = ?", userId).Find(&dbReservations)
	if result.Error != nil {
		return nil, result.Error
	}

	var reservations []*domain.Reservation
	for _, dbReservation := range dbReservations {
		reservations = append(reservations, models.ConvertReservationToDomain(&dbReservation))
	}
	return reservations, nil
}

func (s *Storage) GetUserReservationsUserForDate(date time.Time, userID string) ([]*domain.Reservation, error) {
	var dbReservations []models.Reservation
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	result := s.Database.Preload("User").Preload("Restaurant").Preload("Tables").
		Where("start_time >= ? AND start_time < ? and user_id = ?", startOfDay, endOfDay, userID).
		Find(&dbReservations)
	if result.Error != nil {
		return nil, result.Error
	}

	var reservations []*domain.Reservation
	for _, dbReservation := range dbReservations {
		reservations = append(reservations, models.ConvertReservationToDomain(&dbReservation))
	}
	return reservations, nil
}

func (s *Storage) GetReservationsForDate(date time.Time) ([]*domain.Reservation, error) {
	var dbReservations []models.Reservation
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	result := s.Database.Preload("User").Preload("Restaurant").Preload("Tables").
		Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay).
		Find(&dbReservations)
	if result.Error != nil {
		return nil, result.Error
	}

	var reservations []*domain.Reservation
	for _, dbReservation := range dbReservations {
		reservations = append(reservations, models.ConvertReservationToDomain(&dbReservation))
	}
	return reservations, nil
}

func (s *Storage) UpdateReservation(reservation *domain.Reservation) (bool, error) {
	dbReservation := models.ConvertReservationToModel(reservation)
	result := s.Database.Save(dbReservation)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (s *Storage) GetReservationForId(id string) (*domain.Reservation, error) {
	var dbReservation models.Reservation

	result := s.Database.Preload("User").Preload("Restaurant").Preload("Tables").
		Where("id = ?", id).
		First(&dbReservation)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	s.logger.Debug("dbReservation: %v", dbReservation)
	return models.ConvertReservationToDomain(&dbReservation), nil
}

func (s *Storage) CreateReservation(reservation *domain.Reservation, tableIDs map[string]string) (string, error) {

	dbReservation := models.ConvertReservationToModel(reservation)

	tx := s.Database.Begin()
	if tx.Error != nil {
		return "", tx.Error
	}

	if err := tx.Create(dbReservation).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	for key, tableID := range tableIDs {
		reservationTable := models.ReservationTable{
			ReservationID: dbReservation.ID,
			TableID:       tableID,
			ID:            key,
		}

		if err := tx.Create(&reservationTable).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return dbReservation.ID, nil
}

func (s *Storage) GetTablesByReservationID(reservationID string) ([]domain.Table, error) {
	var tables []models.Table

	// Выполняем запрос к таблице ReservationTable и связываем ее с Table
	err := s.Database.
		Joins("JOIN reservation_tables ON reservation_tables.table_id = tables.id").
		Where("reservation_tables.reservation_id = ?", reservationID).
		Find(&tables).Error

	if err != nil {
		return nil, err
	}
	domainTables := make([]domain.Table, len(tables))
	for _, table := range tables {
		domainTable := models.ConvertTableToDomain(&table)
		domainTables = append(domainTables, *domainTable)
	}

	return domainTables, nil
}
