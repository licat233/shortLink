package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/licat233/goutil/readfile"
)

type App struct {
	Config      *Config
	Router      *gin.Engine
	Middlewares *Middleware
	RedisCli    *RedisCli
}

func (a *App) AuthMiddleware(ctx *gin.Context) {
	if getting(ctx) == nil {
		return
	}
	ctx.Next()
}

func (a *App) InitializeRouters() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	a.Router = gin.Default()
	a.Router.Use(gin.Logger())
	a.Router.Use(gin.Recovery())
	a.Router.GET("/", indexpage)
	a.Router.GET("/luck", nonepage)
	a.Router.Static("/luck/static", "./static")
	a.Router.StaticFile("/luck/login", "./view/login.html")
	a.Router.StaticFile("/luck/admin", "./view/admin.html")
	a.Router.GET("/luck/prizes", a.getprizes)
	a.Router.POST("/luck/login/verify", a.loginVerify)
	a.Router.POST("/luck/admin/genlink", a.AuthMiddleware, a.genlinkapi)
	a.Router.POST("/luck/admin/shortlinks", a.AuthMiddleware, a.shortlinks)
	a.Router.GET("/luck/:shortlink", a.shortlinkpage)
}

func (a *App) Initialize() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	readfile.YamlConfig("./config.yaml", &a.Config, func(err error) {
		if err != nil {
			panic(err)
		}
	})
	a.Middlewares = &Middleware{}
	a.RedisCli = &RedisCli{}
	a.RedisCli.InitializeRedis()
	a.InitializeRouters()
}

func (a *App) Run() {
	log.Fatal(a.Router.Run(a.Config.Addr))
}

func indexpage(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "首页",
	})
}
func nonepage(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "404頁面",
	})
}
