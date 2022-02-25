package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func (a *App) InitializeRouters() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	a.Router = gin.Default()
	a.Router.Use(cors())
	a.Router.Use(gin.Logger())
	a.Router.Use(gin.Recovery())
	a.Router.GET("/", indexpage)
	a.Router.GET("/luck", nonepage)
	a.Router.Static("/luck/static", "./static")
	a.Router.Static("/luck/static1", "./luckview/build")
	a.Router.StaticFile("/luck/login", "./view/login.html")
	a.Router.StaticFile("/luck/admin", "./view/admin.html")
	a.Router.LoadHTMLFiles("view/luck.html")
	a.Router.LoadHTMLFiles("./luckview/build/index.html")
	a.Router.GET("/luck/:shortlink", a.luckpage)
	a.Router.POST("/luck/prizes", a.getprizes)
	a.Router.POST("/luck/login/verify", a.loginVerify)
	a.Router.POST("/luck/admin/genlink", a.AuthMiddleware, a.genlinkapi)
	a.Router.POST("/luck/admin/shortlinks", a.AuthMiddleware, a.shortlinks)
	a.Router.POST("/luck/:shortlink/goodluck", a.goodluck)
	a.Router.POST("/luck/:shortlink/info", a.shortlinkpage)
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

func (a *App) luckpage(ctx *gin.Context) {
	eid := ctx.Param("shortlink")
	if len(eid) != 8 {
		ctx.JSON(400, gin.H{
			"code":    404,
			"message": "404頁面",
		})
		return
	}
	b, err := a.RedisCli.Exists(fmt.Sprintf(ShortlinkKey, eid))
	if err != nil {
		ctx.JSON(500, ServerError(err))
		return
	}
	if !b {
		ctx.JSON(400, gin.H{
			"code":    404,
			"message": "404頁面",
		})
		return
	}
	// info, e := a.RedisCli.GetShortlinkInfo(eid)
	// if e != nil {
	// 	ctx.JSON(e.Status(), e.Error())
	// 	return
	// }
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		// "title": "Main website",
	})
}
