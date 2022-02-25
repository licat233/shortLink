package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type GenLinkReq struct {
	LineID string `json:"LineID"`
}

func (a *App) genlinkapi(ctx *gin.Context) {
	req := &GenLinkReq{}
	ctx.BindJSON(req)
	lineId := strings.TrimSpace(req.LineID)
	if len(lineId) == 0 {
		ctx.JSON(200, gin.H{
			"code":    400,
			"message": "The request is missing a parameter",
		})
		return
	}
	info, err := a.RedisCli.Shorten(lineId)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    err.Status(),
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "短链接生成端口",
		"data":    info,
	})
}
