package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}

	res := Result{}
	res.Code = SuccessCode
	res.Message = GetMessage(SuccessCode)
	res.Data = data

	c.JSON(http.StatusOK, res)
}

func Failed(c *gin.Context, code int, message string) {
	res := Result{}
	res.Code = code
	res.Message = message
	res.Data = gin.H{}

	c.JSON(http.StatusOK, res)
}
