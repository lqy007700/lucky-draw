/**
微信摇一摇
/luck 抽奖接口
*/
package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// 奖品类型

const (
	giftTypeCoin      = iota // 虚拟币
	giftTypeCoupon           // 不同券
	giftTypeCouponFix        // 相同券
	giftTypeRealSmall        // 实物小
	giftTypeRealLarge        // 实物大
)

type gift struct {
	id       int      // id
	name     string   // 名称
	pic      string   // 图片
	link     string   // 连接
	gtype    int      // 类型
	data     string   // 奖品的数据
	dataList []string // 奖品数据集合
	total    int      // 总数量 0不限
	left     int      // 剩余
	inuse    bool     // 是否使用中
	rate     int      // 中奖概率 0-9999
	rateMin  int      // 中奖编码最小
	rateMax  int      // 中奖编码最大
}

const rateMax = 10000

var logger *log.Logger

var mu sync.Mutex

var giftList []*gift

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	initLog()
	initGift()
	return app
}

func initLog() {
	create, _ := os.Create("/Users/lqy007700/Data/log/lottery.log")
	logger = log.New(create, "", log.Ldate|log.Lmicroseconds)
}

func initGift() {
	giftList = make([]*gift, 5)
	g1 := &gift{
		id:       1,
		name:     "Iphone",
		pic:      "",
		link:     "www.iphone.com",
		gtype:    giftTypeRealLarge,
		data:     "",
		dataList: nil,
		total:    20000,
		left:     20000,
		inuse:    true,
		rate:     10000,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[0] = g1
	g2 := &gift{
		id:       2,
		name:     "BTC",
		pic:      "",
		link:     "www.btc.com",
		gtype:    giftTypeCoin,
		data:     "",
		dataList: nil,
		total:    1,
		left:     1,
		inuse:    false,
		rate:     1000,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[1] = g2

	g3 := &gift{
		id:       3,
		name:     "充电宝",
		pic:      "",
		link:     "www.chong.com",
		gtype:    giftTypeRealSmall,
		data:     "",
		dataList: nil,
		total:    2,
		left:     2,
		inuse:    false,
		rate:     3000,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[2] = g3

	g4 := &gift{
		id:       4,
		name:     "满200-50",
		pic:      "",
		link:     "www.chong.com",
		gtype:    giftTypeCouponFix,
		data:     "11112020",
		dataList: nil,
		total:    4,
		left:     4,
		inuse:    false,
		rate:     2000,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[3] = g4

	g5 := &gift{
		id:       5,
		name:     "无门槛50",
		pic:      "",
		link:     "www.chong.com",
		gtype:    giftTypeCoupon,
		data:     "",
		dataList: []string{"202201", "202202", "202203"},
		total:    6,
		left:     6,
		inuse:    false,
		rate:     6000,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[4] = g5
	rateStart := 0
	for _, g := range giftList {
		if !g.inuse {
			continue
		}
		g.rateMin = rateStart
		g.rateMax = rateStart + g.rate
		if g.rateMax >= rateMax {
			g.rateMax = rateMax
			rateStart = 0
		} else {
			rateStart += g.rate
		}
	}
}

func main() {
	mu = sync.Mutex{}
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

//Get 奖品数量的信息
func (l *lotteryController) Get() string {
	count := 0
	total := 0

	for _, g := range giftList {
		if g.inuse && (g.total == 0 || (g.total > 0 && g.left > 0)) {
			count++
			total += g.left
		}
	}

	return fmt.Sprintf("当前有效奖品总类数量：%d，限量奖品数量：%d \n", count, total)
}

func (l *lotteryController) GetLucky() map[string]interface{} {
	mu.Lock()
	defer mu.Unlock()

	code := luckyCode()
	ok := false
	result := make(map[string]interface{})
	result["success"] = ok

	for _, g := range giftList {
		if !g.inuse || (g.total > 0 && g.left <= 0) {
			continue
		}

		if g.rateMin <= code && code < g.rateMax {
			logger.Printf("中奖了，奖品信息 %v", g)
			if g.total != 0 {
				g.left--
				if g.left <= 0 {
					g.inuse = false
				}
			}
			result["success"] = true
			result["id"] = g.id
			result["name"] = g.name
			break
		}
	}
	return result
}

func luckyCode() int {
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Intn(rateMax)
	return code
}
