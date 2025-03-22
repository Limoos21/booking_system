package controllers

import (
	"booking_system/internal/dto"
	"github.com/gin-gonic/gin"
	"time"
)

type userResponse struct {
	Name       string `json:"name"`
	Id         string `json:"id"`
	Phone      string `json:"phone"`
	TelegramId int64  `json:"telegram_id,omitempty"`
}

func convertUserResponse(userDTO dto.UserDTO) userResponse {
	return userResponse{
		Name:       userDTO.Name,
		Id:         userDTO.ID,
		Phone:      userDTO.Phone,
		TelegramId: userDTO.TelegramID,
	}
}

type ResponseTableAvaible struct {
	ID           string  `json:"id"`
	RestaurantID string  `json:"restaurant_id"`
	TableNumber  int     `json:"table_number"`
	Capacity     int     `json:"capacity"`
	PositionX    float64 `json:"position_x"`
	PositionY    float64 `json:"position_y"`
	PositionZ    float64 `json:"position_z"`
	IsAvaible    bool    `json:"is_available,omitempty"`
}

type AnyResponse struct {
	Status string       `json:"status"`
	Data   *interface{} `json:"data,omitempty"` // Используем указатель и omitempty
	Errors *interface{} `json:"errors,omitempty"`
	Meta   *interface{} `json:"meta,omitempty"`
}

type userBookingResponse struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	RestaurantID string                 `json:"restaurant_id"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time"`
	Status       string                 `json:"status"`
	Table        []ResponseTableAvaible `json:"table"`
	Contacts     contactsResponse       `json:"contacts"`
	Capacity     int                    `json:"capacity"`
}

type contactsResponse struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func ConvertTableResponse(table dto.TableDTO) ResponseTableAvaible {
	return ResponseTableAvaible{
		ID:           table.ID,
		RestaurantID: table.RestaurantID,
		TableNumber:  table.TableNumber,
		Capacity:     table.Capacity,
		PositionX:    table.PositionX,
		PositionY:    table.PositionY,
		PositionZ:    table.PositionZ,
	}
}

func convertTouserBookingResponse(dto dto.ReservationDTO) userBookingResponse {
	listTable := make([]ResponseTableAvaible, 0, len(dto.Table))
	for _, table := range dto.Table {
		tableResponse := ConvertTableResponse(table)
		listTable = append(listTable, tableResponse)
	}
	response := userBookingResponse{
		ID:           dto.ID,
		UserID:       dto.UserID,
		RestaurantID: dto.RestaurantID,
		StartTime:    dto.StartTime,
		EndTime:      dto.EndTime,
		Status:       dto.Status,
		Table:        listTable,
		Contacts: contactsResponse{
			Name:  dto.Contacts.Name,
			Phone: dto.Contacts.Phone,
		},
		Capacity: dto.Capacity,
	}
	return response
}

func response(status bool,
	data interface{},
	errors interface{},
	meta interface{},
	c *gin.Context,
	httpStatusCode int) {
	var statusString string
	if status {
		statusString = "success"
	} else {
		statusString = "error"
	}

	response := AnyResponse{
		Status: statusString,
	}

	// Добавляем данные только если они не nil
	if data != nil {
		response.Data = &data
	}
	if errors != nil {
		response.Errors = &errors
	}
	if meta != nil {
		response.Meta = &meta
	}

	c.JSON(httpStatusCode, response)
}
