package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func ResponseHandler(c *gin.Context, success bool, code int, message string, payload any) {
	if success {
		log.Printf("[Success] %d - %s", code, message)
		c.JSON(code, APIResponse{
			Success: true,
			Code:    code,
			Message: message,
			Data:    payload,
		})
	} else {
		if err, ok := payload.(error); ok {
			log.Printf("[ERROR] %d - %s: %v", code, message, err)
			payload = err.Error()
		} else {
			log.Printf("[ERROR] %d - %s: %v", code, message, payload)
		}

		c.JSON(code, APIResponse{
			Success: false,
			Code:    code,
			Message: message,
			Error:   payload.(string),
		})
	}
}
