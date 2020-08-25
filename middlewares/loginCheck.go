package middlewares

import (
	"winddies/manage-api/controllers/base"
	"winddies/manage-api/global"
	"winddies/manage-api/models"

	"github.com/gin-gonic/gin"
)

func LoginRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var base *base.Base
		logger := base.Logger((ctx))
		sessionKey, err := ctx.Cookie(global.Conf.Session.Name)
		if err != nil {
			logger.Errorf("get session name failed: %s", err)
			ctx.Abort()
			return
		}

		result, err := models.RedisDb.Get(ctx, sessionKey+global.Conf.Session.Secret).Result()
		// json.Unmarshal([]byte(result), map[])
		logger.Infof(result)
		if err != nil {
			logger.Error("get userInfo failed")
			ctx.Abort()
			return
		}

		ctx.Set(global.Conf.User, result)
		ctx.Next()
	}
}
