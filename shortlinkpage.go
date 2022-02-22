package main

import "github.com/gin-gonic/gin"

func (a *App) shortlinkpage(ctx *gin.Context) {
	shortlink := ctx.Param("shortlink")
	info, err := a.RedisCli.GetShortlinkInfo(shortlink)
	if err != nil {
		ctx.JSON(200, err)
		return
	}
	ctx.JSON(200, gin.H{
		"message":   "短链接页面",
		"shortlink": info,
	})
}
