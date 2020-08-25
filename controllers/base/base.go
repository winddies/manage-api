package base

import (
	"winddies/manage-api/code"
	"winddies/manage-api/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Base struct{}

const (
	reqidKey = "X-Reqid"
)

func (base *Base) Logger(ctx *gin.Context) *log.Entry {
	reqId := utils.GetReqId(ctx)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true, // 神奇的东西，如果不设置发现并不能显示全时间
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger := log.WithFields(log.Fields{"request_id": reqId, "user_ip": ctx.ClientIP()})
	return logger
}

func (base *Base) Send(ctx *gin.Context, code code.Code, data interface{}) {
	utils.SendResponse(ctx, code, data)
}
