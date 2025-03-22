package usecase

import (
	"booking_system/cmd/providers/middelware"
	"booking_system/internal/app/ports"
	"booking_system/internal/dto"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"time"
)

type UserService struct {
	storage  ports.IStorage
	logger   *slog.Logger
	tokenBot string
	jwt      *middelware.Jwt
}

func New(storage ports.IStorage, logger *slog.Logger, t string, jwt *middelware.Jwt) UserService {
	return UserService{
		storage:  storage,
		logger:   logger,
		tokenBot: t,
		jwt:      jwt,
	}
}

func (u UserService) AuthUser(dto dto.UserDTO) (dto.UserDTO, string, error) {
	ok, user, err := u.storage.CheckUserForTelegram(dto.TelegramID)
	if err != nil {
		return dto, "", err
	}
	if ok {
		token, err := u.jwt.GenerateToken(user.Name, user.ID)
		if err != nil {
			return dto, "", err
		}
		dtoUser := *fromUserDomain(&user)
		return dtoUser, token, nil
	} else {
		domainUser := toUserDomain(&dto)
		uuid := uuid.New().String()
		domainUser.ID = uuid
		newUser, err := u.storage.CreateUser(*domainUser)
		if err != nil {
			return dto, "", err
		}
		token, err := u.jwt.GenerateToken(newUser.Name, newUser.ID)
		if err != nil {
			return dto, "", err
		}
		dtoUser := *fromUserDomain(&newUser)
		return dtoUser, token, nil
	}
}

func (u UserService) GetReservationForDate(date *time.Time) ([]dto.ReservationDTO, error) {
	reservation, err := u.storage.GetReservationsForDate(*date)
	if err != nil {
		return nil, err
	}

	var reservations []dto.ReservationDTO
	for _, val := range reservation {
		tablesDomain, err := u.storage.GetTablesByReservationID(val.ID)
		if err != nil {
			return nil, err
		}

		dtoReservation := fromReservationDomain(val)
		for _, table := range tablesDomain {
			dtoReservation.Table = append(dtoReservation.Table, *fromTableDomain(&table))
		}
		reservations = append(reservations, *dtoReservation)
	}
	return reservations, nil

}

func (u UserService) CreateReservation(dtoReservation dto.ReservationDTO) (dto.ReservationDTO, error) {
	u.logger.Debug("Create Reservation len table first %v", len(dtoReservation.Table))
	domainReservation, tables := toReservationDomain(&dtoReservation)
	ok, err := domainReservation.CheckDate()
	if err != nil {
		return dtoReservation, err
	}
	if !ok {
		return dtoReservation, errors.New("invalid date")
	}
	if len(tables) > 4 {
		return dtoReservation, errors.New("too many tables")
	}

	_, tablesDomain, err := u.checkGuestCapacity(tables, domainReservation.Capacity)
	if err != nil {
		return dtoReservation, err
	}

	tableIds := map[string]string{}
	u.logger.Debug("len tableIds: %v", len(tables))
	u.logger.Info("Table", tables)
	for _, t := range tables {
		u.logger.Debug("LoopTable: %s", t.ID)
		TableOk, err := u.storage.IsTableAvailable(t.ID, domainReservation.StartTime, domainReservation.EndTime)
		if err != nil {
			return dtoReservation, err
		}
		if !TableOk {
			return dtoReservation, errors.New(fmt.Sprintf("table %s not available", t.ID))
		}
		tableIds[uuid.New().String()] = t.ID
	}

	domainReservation.ID = uuid.New().String()

	_, err = u.storage.CreateReservation(domainReservation, tableIds)
	if err != nil {
		u.logger.Error("Failed to create reservation: %v", err)
		return dtoReservation, err
	}
	dtoTables := make([]dto.TableDTO, 0, len(tables))
	for _, table := range tablesDomain {
		dtoTable := *fromTableDomain(table)
		dtoTables = append(dtoTables, dtoTable)
	}
	dtoReservationResult := fromReservationDomain(domainReservation)
	dtoReservationResult.Table = dtoTables
	return *dtoReservationResult, nil
}

func (u UserService) GetUserReservations(userId string) ([]dto.ReservationDTO, error) {
	reservations, err := u.storage.GetUserReservationsUser(userId)
	if err != nil {
		return nil, err
	}

	var reservationsDto []dto.ReservationDTO
	for _, val := range reservations {
		var tables []dto.TableDTO
		reservationDto := fromReservationDomain(val)
		table, err := u.storage.GetTablesByReservationID(val.ID)
		if err != nil {
			return nil, err
		}
		for _, t := range table {
			tables = append(tables, *fromTableDomain(&t))
		}
		reservationDto.Table = tables
		reservationsDto = append(reservationsDto, *reservationDto)
	}
	return reservationsDto, nil

}

func (u UserService) GetUserReservationsDate(date *time.Time, userId string) ([]dto.ReservationDTO, error) {
	reservations, err := u.storage.GetUserReservationsUserForDate(*date, userId)
	if err != nil {
		return nil, err
	}

	var reservationsDto []dto.ReservationDTO
	for _, val := range reservations {
		var tables []dto.TableDTO
		reservationDto := fromReservationDomain(val)
		table, err := u.storage.GetTablesByReservationID(val.ID)
		if err != nil {
			return nil, err
		}
		for _, t := range table {
			tables = append(tables, *fromTableDomain(&t))
		}
		reservationDto.Table = tables
		reservationsDto = append(reservationsDto, *reservationDto)
	}
	return reservationsDto, nil
}

func (u UserService) ValidateTelegramHash(telegramHash string, data map[string]string) (bool, error) {

	requiredFields := []string{"id", "first_name", "auth_date", "hash"}
	for _, field := range requiredFields {
		if _, ok := data[field]; !ok {
			return false, fmt.Errorf("missing required field: %s", field)
		}
	}

	authDate, err := strconv.ParseInt(data["auth_date"], 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid auth_date: %v", err)
	}
	if time.Now().Unix()-authDate > 86400 { // 86400 секунд = 24 часа
		return false, fmt.Errorf("auth_date is too old")
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		if k != "hash" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var dataCheckArr []string
	for _, k := range keys {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", k, data[k]))
	}
	dataCheckString := strings.Join(dataCheckArr, "\n")
	secretKey := sha256.Sum256([]byte(u.tokenBot))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	return strings.ToLower(expectedHash) == strings.ToLower(telegramHash), nil
}

func (u UserService) GetReservationForId(reservationId string) (dto.ReservationDTO, error) {
	domainReservation, err := u.storage.GetReservationForId(reservationId)
	if err != nil {
		return dto.ReservationDTO{}, err
	}
	u.logger.Debug("domainReservation: %v", domainReservation)
	if domainReservation == nil {
		return dto.ReservationDTO{}, errors.New("reservation not found")
	}
	return *fromReservationDomain(domainReservation), nil
}

func (u UserService) UpdateReservation(dto dto.ReservationDTO) (bool, error) {
	domainReservation, _ := toReservationDomain(&dto)
	table, err := u.storage.GetTablesByReservationID(domainReservation.ID)
	if err != nil {
		return false, err
	}

	_, _, err = u.checkGuestCapacityMax(table, domainReservation.Capacity)
	if err != nil {
		return false, err
	}
	ok, err := u.storage.UpdateReservation(domainReservation)
	if err != nil {
		return ok, err
	}
	if !ok {
		return ok, errors.New("reservation not found")
	}

	return ok, nil
}

func (u UserService) GetTableForReservationDate(date time.Time, restaurantId string) ([]dto.AvaibleTableDTO, error) {
	domainTables, err := u.storage.GetTablesWithAvailability(restaurantId, date)
	if err != nil {
		return nil, err
	}
	avaibleTablesDto := make([]dto.AvaibleTableDTO, 0, len(domainTables))
	for _, t := range domainTables {
		dtoTable := *fromTableDomain(&t.Table)
		avaibleTable := dto.AvaibleTableDTO{
			IsAvaible: t.IsAvailable,
			TableDTO:  dtoTable,
		}
		avaibleTablesDto = append(avaibleTablesDto, avaibleTable)
	}
	u.logger.Debug("domainTables: %v", avaibleTablesDto)
	return avaibleTablesDto, nil
}
