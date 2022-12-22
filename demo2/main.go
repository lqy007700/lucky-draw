/**
1/ 即开即得
2/ 双色球自选
*/
package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"math/rand"
	"time"
)

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

func (l *lotteryController) Get() string {
	var prize string

	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Intn(10)
	switch {
	case code == 1:
		prize = "一等奖"
		break
	case code >= 2 && code <= 3:
		prize = "二等奖"
		break
	case code >= 4 && code <= 6:
		prize = "三等奖"
		break
	default:
		prize = "未中奖"
		break
	}
	return prize
}

func (l *lotteryController) GetPrize() string {
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed))
	prize := [7]int{}

	for i := 0; i < 6; i++ {
		prize[i] = code.Intn(33) + 1
	}

	prize[6] = code.Intn(16) + 1
	return fmt.Sprintf("开奖号码是：%v", prize)
}
