package repository

import (
	"back_api/internal/domain"
	"back_api/internal/repository/models"
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

type PostgresRepository struct {
	Database *gorm.DB
	logger   *slog.Logger
}

func (p *PostgresRepository) CreateBooking(booking *domain.Booking, table []string) (*domain.Booking, error) {
	panic("implement me")
	return nil, nil

}
func (p *PostgresRepository) GetAvailbleTable(tableUUid string, dateS, dateE time.Time) (*domain.Table, error) {
	panic("implement me")
	return nil, nil
}

func (p *PostgresRepository) GetUserForUuid(uuid string) (*domain.User, error) {
	var user models.User
	result := p.Database.First(&user, "user_uuid = ?", uuid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			p.logger.Warn("User not found for UUID: %s", uuid)
			return nil, nil
		}
		p.logger.Error("Error fetching user for UUID: %s, %v", uuid, result.Error)
		return nil, result.Error
	}
	domainUser, err := models.ConvertUserMTODomain(&user)
	if err != nil {
		p.logger.Error("Error converting user to domain: %v", err)
		return nil, err
	}
	return domainUser, nil
}

func (p *PostgresRepository) CreateUser(user *domain.User) (*domain.User, error) {
	modelsUser, err := models.ConvertDomainUserTOMUser(user)
	if err != nil {
		return nil, err
	}
	result := p.Database.Create(&modelsUser)
	if result.Error != nil {
		p.logger.Error("Error creating user: %v", result.Error)
		return nil, result.Error
	}
	domainUser, err := models.ConvertUserMTODomain(modelsUser)
	if err != nil {
		return nil, err
	}

	return domainUser, nil
}

func (p *PostgresRepository) GetUserForTelegram(telegramId string) (*domain.User, error) {
	var user models.User
	result := p.Database.First(&user, "telegram_id = ?", telegramId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			p.logger.Warn("User not found for Telegram ID: %s", telegramId)
			return nil, nil
		}
		p.logger.Error("Error fetching user for Telegram ID: %s, %v", telegramId, result.Error)
		return nil, result.Error
	}
	domainUser, err := models.ConvertUserMTODomain(&user)
	if err != nil {
		p.logger.Error("Error converting user to domain: %v", err)
		return nil, err
	}
	return domainUser, nil
}

func (p *PostgresRepository) IfUserExists(userUuid string) (bool, error) {
	var count int64
	result := p.Database.Model(&models.User{}).Where("user_uuid = ?", userUuid).Count(&count)
	if result.Error != nil {
		p.logger.Error("Error checking user existence for UUID: %s, %v", userUuid, result.Error)
		return false, result.Error
	}
	return count > 0, nil
}

func NewPostgresRepository(dsn string, logger *slog.Logger) *PostgresRepository {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database: %v", err)
		return nil
	}
	return &PostgresRepository{Database: db, logger: logger}
}
