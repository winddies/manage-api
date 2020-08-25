package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

const reqidKey = "X-ReqId"

func GenerateReqId(ctx *gin.Context) {
	reqId := xid.New()
	ctx.Request.Header.Set(reqidKey, reqId.String())
}

func GetReqId(ctx *gin.Context) string {
	return ctx.Request.Header.Get(reqidKey)
}
