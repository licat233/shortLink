package main

import "github.com/gin-gonic/gin"

func (a *App) shortlinkpage(ctx *gin.Context) {
	eid := ctx.Param("shortlink")
	if len(eid) != 8 {
		ctx.JSON(400, gin.H{
			"code":    404,
			"message": "404端口",
		})
		return
	}
	info, err := a.RedisCli.GetShortlinkInfo(eid)
	if err != nil {
		ctx.JSON(200, err)
		return
	}
	ctx.JSON(200, gin.H{
		"code":      200,
		"message":   "短链接信息",
		"shortlink": info,
	})
}
