package controllers

import "github.com/gin-gonic/gin"

type userResponse struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type responseTableAvaible struct {
	ID           string  `json:"id"`
	RestaurantID string  `json:"restaurant_id"`
	TableNumber  int     `json:"table_number"`
	Capacity     int     `json:"capacity"`
	PositionX    float64 `json:"position_x"`
	PositionY    float64 `json:"position_y"`
	PositionZ    float64 `json:"position_z"`
	IsAvaible    bool    `json:"is_available"`
}

type anyResponse struct {
	Status string       `json:"status"`
	Data   *interface{} `json:"data,omitempty"` // Используем указатель и omitempty
	Errors *interface{} `json:"errors,omitempty"`
	Meta   *interface{} `json:"meta,omitempty"`
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

	response := anyResponse{
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
