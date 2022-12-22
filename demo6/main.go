package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"
)

/**
大转盘类
*/

// Prate 中奖概率
type Prate struct {
	Rate  int    // 万分之
	Total int    // 总量 0为无限
	Left  *int32 // 剩余
	CodeA int    // 起始编码
	CodeB int    // 终止编码
}

// 奖品列表
var prizeList = []string{
	"一等奖",
	//"二等奖",
	//"三等奖",
}

var left int32 = 1000

// 奖品中奖概率对应
var prizeRateList = []Prate{
	{100, 1000, &left, 0, 100},
	//{2, 2, 2, 1, 3},
	//{5, 3, 3, 4, 9},
}

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

//Get 转盘信息
func (l *lotteryController) Get() string {
	l.Ctx.Header("Content-Type", "text/html")
	return fmt.Sprintf("奖品列表 <br/> %s", strings.Join(prizeList, "<br/>\n"))
}

//GetPrize 发红包
func (l *lotteryController) GetPrize() string {

	// 随机数
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	code := r.Intn(100)

	//log.Printf("中奖编码为%d", code)

	var myPrize string
	var prizeRate *Prate

	for i, s := range prizeList {
		rate := &prizeRateList[i]

		// 奖品有限量并且剩余为0
		if rate.Total != 0 && *rate.Left <= 0 {
			//log.Printf("奖品%s已经发完，跳过 %+v", s, rate)
			continue
		}

		if code >= rate.CodeA && code < rate.CodeB {
			myPrize = s
			prizeRate = rate
			break
		}
	}

	if myPrize == "" {
		//log.Printf("奖品%s剩余信息%+v", myPrize, prizeRate)
		return "未中奖"
	}

	// 发奖
	if prizeRate.Total == 0 {
		return myPrize
	}
	log.Printf("奖品%s", myPrize)
	left := atomic.AddInt32(prizeRate.Left, -1)
	if left >= 0 {
		return myPrize
	}
	return "未中奖"
}