package controllers

import "github.com/gin-gonic/gin"

type ResponseDTO struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"message"`
	Error   interface{} `json:"error"`
}

func NewResponse(success bool, error, data interface{}, c *gin.Context, statusCode int) {
	response := ResponseDTO{
		Success: success,
		Data:    data,
		Error:   error,
	}
	c.JSON(statusCode, response)
}
