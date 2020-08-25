package utils

import (
	codes "winddies/manage-api/code"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func SendResponse(ctx *gin.Context, code codes.Code, data interface{}) {
	ctx.Header("Content-Type", "application/json; charset=UTF-8")
	if code == codes.OK {
		ctx.JSON(code.Int(), data)
		return
	}

	httpCode := code
	if httpCode.Int() >= 1000 {
		httpCode = codes.ResultError
	}

	ctx.JSON(httpCode.Int(), &ErrorBody{
		Code:  code.Int(),
		Error: code.Error(),
	})
}
