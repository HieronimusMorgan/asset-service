package response

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Response struct {
	Status    int         `json:"status"`          // HTTP status code
	Message   string      `json:"message"`         // Descriptive message
	Timestamp string      `json:"timestamp"`       // Descriptive message
	Data      interface{} `json:"data,omitempty"`  // Any additional data
	Error     interface{} `json:"error,omitempty"` // Error details (if any)
}

func SendResponse(c *gin.Context, status int, message string, data interface{}, err interface{}) {
	c.JSON(status, Response{
		Status:    status,
		Timestamp: time.Now().In(time.FixedZone("GMT+7", 7*3600)).Format(time.DateTime),
		Message:   message,
		Data:      data,
		Error:     err,
	})
}
