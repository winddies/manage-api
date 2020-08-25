package middlewares

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"winddies/manage-api/global"
	"winddies/manage-api/utils"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

var LogClient *log.Logger

func Logger() gin.HandlerFunc {

	LogClient = log.New()

	LogClient.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true, // 神奇的东西，如果不设置发现并不能显示全时间
		TimestampFormat: "2006-01-02 15:04:05",
	})

	apiLogPath := "api.log"
	logWriter, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(apiLogPath),       // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		fmt.Println("err", err)
	}

	switch level := global.Conf.LogLevel; level {
	/*如果日志级别不是 debug 就不要打印日志到控制台了*/
	case "debug":
		LogClient.SetLevel(log.DebugLevel)
		LogClient.SetOutput(os.Stdout)
	case "info":
		setNull(LogClient)
		LogClient.SetLevel(log.InfoLevel)
	case "warn":
		setNull(LogClient)
		LogClient.SetLevel(log.WarnLevel)
	case "error":
		setNull(LogClient)
		LogClient.SetLevel(log.ErrorLevel)
	default:
		setNull(LogClient)
		LogClient.SetLevel(log.InfoLevel)
	}

	writeMap := lfshook.WriterMap{
		log.InfoLevel:  logWriter,
		log.FatalLevel: logWriter,
		log.ErrorLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &log.JSONFormatter{})
	log.AddHook(lfHook)

	return func(ctx *gin.Context) {
		utils.GenerateReqId(ctx)

		// 开始时间
		start := time.Now()
		LogClient.Infof(
			"[REQ_BEG] | %s | %s | %s",
			ctx.Request.Method,
			ctx.Request.URL.Path,
			ctx.ClientIP(),
		)
		// 处理请求
		ctx.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)

		path := ctx.Request.URL.Path

		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		LogClient.Infof(
			"[REQ_BEG] | %d | %v | %s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}

func setNull(Logger *log.Logger) {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	Logger.SetOutput(writer)
}
