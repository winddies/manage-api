package routes

import (
	"winddies/manage-api/controllers/article"

	"winddies/manage-api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes(app *gin.Engine) {
	article := article.Create()
	routeAuth := app.Group("/api/manage/auth")
	routeAuth.Use(middlewares.LoginRequired())
	routeAuth.POST("/article/new", article.CreateNewArticle)
	routeGeneral := app.Group("/api/manage/general")
	routeGeneral.GET("/article/list", article.GetArticleSummaryList)
}
