package controllers

import (
	"booking_system/cmd/providers/middelware"
	"booking_system/internal/app/ports"
	"booking_system/internal/dto"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
	logger  *slog.Logger
	useCase ports.IUseCase
	jwt     *middelware.Jwt
}

func New(logger *slog.Logger, useCase ports.IUseCase, jwt *middelware.Jwt) *Controller {
	return &Controller{
		logger:  logger,
		useCase: useCase,
		jwt:     jwt,
	}
}

func (c *Controller) Authorize(context *gin.Context) {
	data := map[string]string{
		"id":         context.Query("id"),
		"first_name": context.Query("first_name"),
		"photo_url":  context.Query("photo_url"),
		"username":   context.Query("username"),
		"auth_date":  context.Query("auth_date"),
		"hash":       context.Query("hash"),
	}
	c.logger.Info("data", data)

	if data["hash"] == "" {
		c.logger.Warn("Telegram hash is missing")
		response(false, nil, "Telegram hash is missing", nil, context, http.StatusBadRequest)
		return
	}

	if ok, err := c.useCase.ValidateTelegramHash(data["hash"], data); !ok || err != nil {
		if err != nil {
			c.logger.Error("Error validating telegram hash", err)
			response(false, nil, "Error validating telegram hash", nil, context, http.StatusBadRequest)
			return
		}
		c.logger.Warn("Telegram hash is unvalidated")
		response(false, nil, "Telegram hash is unvalidated", nil, context, http.StatusBadRequest)
		return
	}
	telegramId, err := strconv.Atoi(context.Query("id"))
	if err != nil {
		c.logger.Warn("Telegram id is missing")
		response(false, nil, "Telegram id is missing", nil, context, http.StatusBadRequest)
		return
	}
	dtoUser := dto.UserDTO{
		Name:       data["username"],
		TelegramID: int64(telegramId),
	}
	createUser, token, err := c.useCase.AuthUser(dtoUser)
	if err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err, nil, context, http.StatusBadRequest)
		return
	}
	userResponses := &userResponse{
		Name: createUser.Name,
		Id:   createUser.ID,
	}
	meta := map[string]string{
		"token": token,
	}
	response(true, userResponses, nil, meta, context, http.StatusOK)
}

func (c *Controller) GetUserBookings(context *gin.Context) {
	userUUid, ok := context.Get("userUuid")
	if !ok {
		c.logger.Warn("User uuid is missing")
		response(false, nil, "User uuid is missing", "My be jwt token missing?", context, http.StatusBadRequest)
		return
	}
	bookings, err := c.useCase.GetUserReservations(userUUid.(string))
	c.logger.Debug("bookings", bookings)
	if err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err, nil, context, http.StatusBadRequest)
		return
	}
	response(true, bookings, nil, nil, context, http.StatusOK)
	return

}

func (c *Controller) UpdateBooking(context *gin.Context) {
	reservationID := context.Param("id")
	userUUID, ok := context.Get("userUuid")
	status := context.Param("status")

	if !ok {
		c.logger.Warn("User uuid is missing")
		response(false, nil, "User uuid is missing", "", context, http.StatusBadRequest)
		return
	}

	if reservationID == "" {
		c.logger.Warn("reservation id is missing")
		response(false, nil, "reservation id is missing", nil, context, http.StatusBadRequest)
		return
	}
	reservationDto, err := c.useCase.GetReservationForId(reservationID)
	if err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err.Error(), nil, context, http.StatusBadRequest)
		return
	}
	if userUUID.(string) != reservationDto.UserID {
		c.logger.Warn("User uuid is unvalidated")
		response(false, nil, "User uuid is unvalidated", "", context, http.StatusBadRequest)
		return
	}
	if status == "canceled" {
		updateReservation := dto.ReservationDTO{
			ID:           reservationDto.ID,
			UserID:       reservationDto.UserID,
			RestaurantID: reservationDto.RestaurantID,
			StartTime:    reservationDto.StartTime,
			EndTime:      reservationDto.EndTime,
			Status:       status,
			Table:        reservationDto.Table,
			Contacts: dto.ContactsDTO{
				Name:  reservationDto.Contacts.Name,
				Phone: reservationDto.Contacts.Phone,
			},
			Capacity: reservationDto.Capacity,
		}
		ok, err = c.useCase.UpdateReservation(updateReservation)
		if err != nil {
			c.logger.Error(err.Error())
			response(false, nil, err.Error(), nil, context, http.StatusBadRequest)
			return
		}
		if !ok {
			c.logger.Warn("Update reservation fail")
			response(false, nil, "Update reservation fail", nil, context, http.StatusBadRequest)
			return
		}
		response(true, "Update reservation success", nil, nil, context, http.StatusOK)
		return
	}
	var data updateReservationRequest
	if err := context.ShouldBindJSON(&data); err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err, nil, context, http.StatusBadRequest)
		return
	}

	c.logger.Debug("reservation", reservationDto)
	updateReservation := dto.ReservationDTO{
		ID:           reservationDto.ID,
		UserID:       reservationDto.UserID,
		RestaurantID: reservationDto.RestaurantID,
		StartTime:    reservationDto.StartTime,
		EndTime:      reservationDto.EndTime,
		Status:       reservationDto.Status,
		Table:        reservationDto.Table,
		Contacts: dto.ContactsDTO{
			Name:  data.Contacts.Name,
			Phone: data.Contacts.Phone,
		},
		Capacity: data.Capacity,
	}
	ok, err = c.useCase.UpdateReservation(updateReservation)
	if err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err.Error(), nil, context, http.StatusBadRequest)
		return
	}
	if !ok {
		c.logger.Warn("Update reservation fail")
		response(false, nil, "Update reservation fail", nil, context, http.StatusBadRequest)
		return
	}
	response(true, "Update reservation success", nil, nil, context, http.StatusOK)
}

func (c *Controller) CreateBooking(context *gin.Context) {
	var data reservationRequest
	restaurantId := context.Param("restaurantId")
	if restaurantId == "" {
		c.logger.Warn("restaurant id is missing")
		response(false, nil, "restaurant id is missing", nil, context, http.StatusBadRequest)
		return
	}
	if err := context.ShouldBind(&data); err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err, nil, context, http.StatusBadRequest)
		return
	}
	userUUid, ok := context.Get("userUuid")
	if !ok {
		c.logger.Warn("User uuid is missing")
		response(false, nil, "User uuid is missing", "", context, http.StatusBadRequest)
		return
	}

	tablesDto := make([]dto.TableDTO, 0, len(data.Table))
	for _, tableID := range data.Table {
		tableDto := dto.TableDTO{
			ID: tableID,
		}
		tablesDto = append(tablesDto, tableDto)
	}
	reservationDto := dto.ReservationDTO{
		UserID:       userUUid.(string),
		RestaurantID: restaurantId,
		StartTime:    data.DateStart,
		EndTime:      data.DateEnd,
		Status:       "wait",
		Table:        tablesDto,
		Capacity:     data.Capacity,
		Contacts: dto.ContactsDTO{
			Name:  data.Contacts.Name,
			Phone: data.Contacts.Phone,
		},
	}
	c.logger.Debug("lenTable controller %v", len(reservationDto.Table))
	createBooking, err := c.useCase.CreateReservation(reservationDto)
	if err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err.Error(), nil, context, http.StatusBadRequest)
		return
	}
	response(true, createBooking, nil, nil, context, http.StatusOK)
}

func (c *Controller) GetBookingsDate(context *gin.Context) {
	date := context.Param("date")
	restaurantId := context.Param("restaurantId")
	if restaurantId == "" {
		c.logger.Warn("restaurant id is missing")
		response(false, nil, "restaurant id is missing", "", context, http.StatusBadRequest)
		return
	}
	if date == "" {
		c.logger.Warn("date is missing")
		response(false, nil, "date is missing", "", context, http.StatusBadRequest)
		return
	}
	dateTime, err := time.Parse("2006-01-02T15:04:05", date)
	if err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err.Error(), nil, context, http.StatusBadRequest)
		return
	}
	avaibleTable, err := c.useCase.GetTableForReservationDate(dateTime, restaurantId)
	if err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err.Error(), nil, context, http.StatusBadRequest)
		return
	}
	responseDto := make([]responseTableAvaible, 0, len(avaibleTable))
	for _, table := range avaibleTable {
		a := responseTableAvaible{
			ID:           table.ID,
			RestaurantID: table.RestaurantID,
			TableNumber:  table.TableNumber,
			Capacity:     table.Capacity,
			IsAvaible:    table.IsAvaible,
			PositionZ:    table.PositionZ,
			PositionY:    table.PositionY,
			PositionX:    table.PositionX,
		}
		responseDto = append(responseDto, a)
	}
	response(true, responseDto, nil, nil, context, http.StatusOK)

}

func (c *Controller) GetUserBookingsDate(context *gin.Context) {
	date := context.Param("date")
	userUUID, ok := context.Get("userUuid")
	if !ok {
		c.logger.Warn("User uuid is missing")
	}
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err.Error(), nil, context, http.StatusBadRequest)
	}
	userBookings, err := c.useCase.GetUserReservationsDate(&dateTime, userUUID.(string))
	if err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err.Error(), nil, context, http.StatusBadRequest)
	}
	response(true, userBookings, nil, nil, context, http.StatusOK)
}
