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

// @Summary Get booking by ID
// @Description Get a specific booking by its ID
// @Tags booking
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} AnyResponse{data=userBookingResponse} "Successfully retrieved booking"
// @Failure 400 {object} AnyResponse{errors=string} "Bad request"
// @Failure 500 {object} AnyResponse{errors=string} "Internal server error"
// @Router /booking/{id} [get]
func (c *Controller) GetBooking(context *gin.Context) {
	reservationID := context.Param("id")
	if reservationID == "" {
		response(false, nil, "reservation id is missing", nil, context, http.StatusBadRequest)
		return
	}
	booking, err := c.useCase.GetReservationForId(reservationID)
	if err != nil {
		response(false, nil, "reservation id is invalid", nil, context, http.StatusBadRequest)
	}
	resp := convertTouserBookingResponse(booking)
	response(true, resp, nil, nil, context, http.StatusOK)

}

// @Summary Test authentication
// @Description Test authentication endpoint
// @Tags auth
// @Accept json
// @Produce json
// @Param user body userRequestRegisterTest true "User data"
// @Success 200 {object} AnyResponse{data=userResponse,meta=map[string]string} "Successfully authenticated"
// @Failure 400 {object} AnyResponse{errors=string} "Bad request"
// @Failure 500 {object} AnyResponse{errors=string} "Internal server error"
// @Router /auth/test [post]
func (c *Controller) AuthTest(context *gin.Context) {
	var data userRequestRegisterTest
	err := context.ShouldBind(&data)
	if err != nil {
		response(false, nil, err, nil, context, http.StatusBadRequest)
	}
	userDto := dto.UserDTO{
		Name:       data.Name,
		Phone:      data.Phone,
		TelegramID: data.Telegram,
	}
	createUser, token, err := c.useCase.AuthUser(userDto)
	if err != nil {
		response(false, nil, err, nil, context, http.StatusInternalServerError)
	}

	resp := convertUserResponse(createUser)
	dateToken := map[string]string{
		"token": token,
	}
	response(true, resp, nil, dateToken, context, http.StatusOK)

}

// @Summary Update user information
// @Description Update the authenticated user's information
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body updateUserData true "User data"
// @Success 200 {object} AnyResponse{data=string} "User updated successfully"
// @Failure 400 {object} AnyResponse{errors=string} "Bad request"
// @Failure 500 {object} AnyResponse{errors=string} "Internal server error"
// @Router /me [patch]
func (c *Controller) UpdateUserInfo(context *gin.Context) {
	var data updateUserData
	err := context.ShouldBind(&data)
	if err != nil {
		response(false, nil, err, nil, context, http.StatusBadRequest)
	}
	userId, ok := context.Get("userUuid")
	if !ok {
		response(false, nil, "UserId is missing", nil, context, http.StatusBadRequest)
	}
	dtoUser := dto.UserDTO{
		ID:    userId.(string),
		Name:  data.Name,
		Phone: data.Phone,
	}
	ok, err = c.useCase.UpdateUser(dtoUser, userId.(string))
	if err != nil {
		response(false, nil, err, nil, context, http.StatusInternalServerError)
		return
	}
	if !ok {
		response(false, nil, "User not found", nil, context, http.StatusNotFound)
		return
	}
	response(true, "User successfully updated", nil, nil, context, http.StatusOK)
	return
}

// @Summary Get user information
// @Description Get information about the authenticated user
// @Tags user
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} AnyResponse{data=userResponse} "Successfully retrieved user info"
// @Failure 400 {object} AnyResponse{errors=string} "Bad request"
// @Failure 500 {object} AnyResponse{errors=string} "Internal server error"
// @Router /me [get]
func (c *Controller) GetUserInfo(context *gin.Context) {
	userId, ok := context.Get("userUuid")
	if !ok {
		response(false, nil, "UserUUid is missin", nil, context, http.StatusBadRequest)
		return
	}
	user, err := c.useCase.GetUserForId(userId.(string))
	if err != nil {
		response(false, nil, err, nil, context, http.StatusInternalServerError)
		return
	}
	resp := convertUserResponse(user)
	response(true, resp, "User successfully retrieved", nil, context, http.StatusOK)
	return
}

// @Summary Authenticate via Telegram
// @Description Authenticate a user using Telegram credentials
// @Tags auth
// @Produce json
// @Param id query string true "Telegram user ID"
// @Param first_name query string true "User's first name"
// @Param photo_url query string false "User's photo URL"
// @Param username query string false "User's username"
// @Param auth_date query string true "Authentication date"
// @Param hash query string true "Telegram hash"
// @Success 200 {object} AnyResponse{data=userResponse,meta=map[string]string} "Successfully authenticated"
// @Failure 400 {object} AnyResponse{errors=string} "Bad request"
// @Failure 500 {object} AnyResponse{errors=string} "Internal server error"
// @Router /auth/telegram [get]
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
	userResponses := convertUserResponse(createUser)
	meta := map[string]string{
		"token": token,
	}
	response(true, userResponses, nil, meta, context, http.StatusOK)
}

// @Summary Get user bookings
// @Description Get all bookings for the authenticated user
// @Tags booking
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} controllers.AnyResponse{data=[]controllers.userBookingResponse} "Successfully retrieved bookings"
// @Failure 400 {object} controllers.AnyResponse{errors=string} "Bad request"
// @Failure 500 {object} controllers.AnyResponse{errors=string} "Internal server error"
// @Router /booking/me [get]
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
	resp := make([]userBookingResponse, 0, len(bookings))
	for _, booking := range bookings {
		b := convertTouserBookingResponse(booking)
		resp = append(resp, b)
	}
	response(true, resp, nil, nil, context, http.StatusOK)
	return

}

// @Summary Update booking
// @Description Update a specific booking by its ID
// @Tags booking
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Booking ID"
// @Param booking body updateReservationRequest true "Booking data"
// @Success 200 {object} AnyResponse{data=string} "Booking updated successfully"
// @Failure 400 {object} AnyResponse{errors=string} "Bad request"
// @Failure 500 {object} AnyResponse{errors=string} "Internal server error"
// @Router /booking/{id} [patch]

// @Summary Update booking status
// @Description Update the status of a specific booking by its ID
// @Tags booking
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Booking ID"
// @Param status path string true "New status"
// @Success 200 {object} AnyResponse{data=string} "Booking status updated successfully"
// @Failure 400 {object} AnyResponse{errors=string} "Bad request"
// @Failure 500 {object} AnyResponse{errors=string} "Internal server error"
// @Router /booking/{id}/{status} [patch]
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

// @Summary Create a new booking
// @Description Create a new booking for the specified restaurant
// @Tags booking
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param restaurantId path string true "Restaurant ID"
// @Param booking body reservationRequest true "Booking data"
// @Success 201 {object} AnyResponse{data=userBookingResponse} "Booking created successfully"
// @Failure 400 {object} AnyResponse{errors=string} "Bad request"
// @Failure 500 {object} AnyResponse{errors=string} "Internal server error"
// @Router /{restaurantId}/booking [post]
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
	resp := convertTouserBookingResponse(createBooking)
	response(true, resp, nil, nil, context, http.StatusOK)
}

// @Summary Get available bookings for a date
// @Description Get available bookings for a specific date and restaurant
// @Tags booking
// @Produce json
// @Param restaurantId path string true "Restaurant ID"
// @Param date path string true "Date in YYYY-MM-DD format"
// @Success 200 {object} AnyResponse{data=[]ResponseTableAvaible} "Successfully retrieved available bookings"
// @Failure 400 {object} AnyResponse{errors=string} "Bad request"
// @Failure 500 {object} AnyResponse{errors=string} "Internal server error"
// @Router /{restaurantId}/bookings/{date} [get]
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
	responseDto := make([]ResponseTableAvaible, 0, len(avaibleTable))
	for _, table := range avaibleTable {
		a := ResponseTableAvaible{
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

// @Summary Get user bookings for a date
// @Description Get all bookings for the authenticated user on a specific date
// @Tags booking
// @Produce json
// @Security ApiKeyAuth
// @Param date path string true "Date in YYYY-MM-DD format"
// @Success 200 {object} AnyResponse{data=[]userBookingResponse} "Successfully retrieved bookings"
// @Failure 400 {object} AnyResponse{errors=string} "Bad request"
// @Failure 500 {object} AnyResponse{errors=string} "Internal server error"
// @Router /booking/me/{date} [get]
func (c *Controller) GetUserBookingsDate(context *gin.Context) {
	date := context.Param("date")
	userUUID, ok := context.Get("userUuid")
	if !ok {
		c.logger.Warn("User uuid is missing")
		response(false, nil, "User uuid is missing", nil, context, http.StatusBadRequest)
	}
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err.Error(), nil, context, http.StatusBadRequest)
		return
	}
	userBookings, err := c.useCase.GetUserReservationsDate(&dateTime, userUUID.(string))
	if err != nil {
		c.logger.Error(err.Error())
		response(false, nil, err.Error(), nil, context, http.StatusBadRequest)
		return
	}
	resp := make([]userBookingResponse, 0, len(userBookings))
	for _, userBooking := range userBookings {
		a := convertTouserBookingResponse(userBooking)
		resp = append(resp, a)
	}
	response(true, resp, nil, nil, context, http.StatusOK)
	return
}
