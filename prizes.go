package main

import "github.com/gin-gonic/gin"

type Prize struct {
	Id     int    `yaml:"Id"`     // 对应的前端产品列表的 index
	Name   string `yaml:"Name"`   // 礼品名称
	Image  string `yaml:"Image"`  // 礼品图片
	Chance int32  `yaml:"Chance"` // 运气值
	Win    bool   `yaml:"Win"`    //中獎了
}

func (a *App) getprizes(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "請求成功",
		"data":    a.Config.Prizes,
	})
}
