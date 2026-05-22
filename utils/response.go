package utils

import (

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool `json:"success" example:"true"`
	Data    any  `json:"data"`
}

type ErrorResponse struct {
	Success bool      `json:"success" example:"false"`
	Error   *APIError `json:"error"`
}

type APIError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func JSONSuccess(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func JSONError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

