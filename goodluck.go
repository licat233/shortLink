package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

func (a *App) goodluck(ctx *gin.Context) {
	//檢測是否被抽過
	eid := ctx.Param("shortlink")
	if len(eid) != 8 {
		ctx.JSON(404, gin.H{
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
	if info.Status {
		ctx.JSON(200, gin.H{
			"code":    400,
			"message": "抽獎次數已用完",
			"data":    info.Prize,
			"line":    info.LineId,
		})
		return
	}
	prize := a.randomPrize()
	if prize == nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": "服務器配置錯誤,請聯繫站點管理員",
		})
		return
	}
	prize.Chance = 0
	info.Prize = prize
	info.Status = true
	info.Count += 1
	info.LuckDate = time.Now().String()
	infoStr, _ := json.Marshal(info)

	if e := a.RedisCli.Cli.Set(fmt.Sprintf(ShortlinkKey, eid), infoStr, 0).Err(); e != nil {
		fmt.Println(e)
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": e.Error(),
		})
		return
	}
	msg := "感謝參與"
	if prize.Win {
		msg = fmt.Sprintf("運氣爆表！抽中了【%s】", prize.Name)
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": msg,
		"data":    prize,
		"line":    info.LineId,
	})
}

func (a *App) randomPrize() *Prize {
	sort.Slice(a.Config.Prizes, func(i, j int) bool {
		return a.Config.Prizes[i].Chance < a.Config.Prizes[j].Chance
	})
	var allprob int32
	for _, v := range a.Config.Prizes {
		allprob += v.Chance
	}
	rand.Seed(time.Now().UnixNano())
	if allprob == 0 {
		return nil
	}
	random := rand.Int31n(allprob)
	for _, v := range a.Config.Prizes {
		if random < v.Chance {
			return v
		}
		random -= v.Chance
	}
	return nil
}
