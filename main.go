package main

import (
	"flag"
	"net/http"
	"time"
	"winddies/manage-api/global"
	"winddies/manage-api/middlewares"
	"winddies/manage-api/models"
	"winddies/manage-api/routes"

	"github.com/gin-gonic/gin"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "get config file")
}

func main() {
	flag.Parse()
	global.Init(configPath)
	models.RedisInit()
	models.MongoInit()

	gin.SetMode(getGinMode())
	app := gin.New()
	app.Use(middlewares.Logger(), gin.Recovery())
	routes.InitRoutes(app)

	s := &http.Server{
		Addr:           global.Conf.Port,
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func getGinMode() string {
	switch global.Conf.Mode {
	case global.DevMode:
		return gin.DebugMode
	case global.TestMode:
		return gin.TestMode
	case global.ProdMode:
		return gin.ReleaseMode
	default:
		return gin.DebugMode
	}
}
