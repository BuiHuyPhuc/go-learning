package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	InternalCode int         `json:"internal_code"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, Response{
		InternalCode: code,
		Message:      msg[code],
		Data:         data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		InternalCode: code,
		Message:      message,
	})
}
