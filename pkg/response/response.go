package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	InternalCode int         `json:"internal_code"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}

type ErrorResponseData struct {
	InternalCode int         `json:"internal_code"`
	Err          string      `json:"error"`
	Detail       interface{} `json:"detail"`
}

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		InternalCode: code,
		Message:      msg[code],
		Data:         data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	if message == "" {
		message = msg[code]
	}
	c.JSON(http.StatusInternalServerError, ErrorResponseData{
		InternalCode: code,
		Err:          message,
		Detail:       nil,
	})
}
