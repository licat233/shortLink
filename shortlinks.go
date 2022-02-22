package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func (a *App) shortlinks(ctx *gin.Context) {
	res, err := a.RedisCli.getUrls()
	if err == redis.Nil {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "目前沒有數據",
		})
	} else if err != nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": "服務器出現錯誤",
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "請求成功",
			"total":   len(res),
			"data":    res,
		})
	}
}
