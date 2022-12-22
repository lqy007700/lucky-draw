package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"math/rand"
	"sync"
	"time"
)

/**
抢红包 channel 版本
*/

type lotteryController struct {
	Ctx iris.Context
}

type packageMap struct {
	m map[uint32][]uint
	sync.Mutex
}

var packageList = packageMap{m: make(map[uint32][]uint)}

type task struct {
	id       uint32
	callback chan uint
}

// 单任务
//var chTasks = make(chan task)
// 多任务
const taskNum = 16

var chTasks = make([]chan task, taskNum)

// 红包集合 红包id:红包列表
//var packageList = make(map[uint32][]uint)
//var packageList sync.Map

func (p *packageMap) set(key uint32, val []uint) {
	p.Lock()
	defer p.Unlock()
	p.m[key] = val
}

func (p *packageMap) get(key uint32) []uint {
	p.Lock()
	defer p.Unlock()
	v, ok := p.m[key]
	if !ok {
		return nil
	}
	return v
}

func (p *packageMap) del(key uint32) {
	p.Lock()
	defer p.Unlock()
	delete(p.m, key)
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	for i := 0; i < taskNum; i++ {
		chTasks[i] = make(chan task)
		go fetchPackageListMoney(i)
	}

	app := newApp()
	app.Run(iris.Addr(":8080"))
}

//Get 获取红包信息
func (l *lotteryController) Get() map[uint32][2]int {
	res := make(map[uint32][2]int)

	for u, uints := range packageList.m {
		money := 0
		for _, u2 := range uints {
			money += int(u2)
		}
		res[u] = [2]int{len(uints), money}
	}
	//
	//packageList.Range(func(key, value interface{}) bool {
	//	money := 0
	//	id := key.(uint32)
	//	list := value.([]uint)
	//	for _, i2 := range list {
	//		money += int(i2)
	//	}
	//	res[id] = [2]int{len(list), money}
	//	return true
	//})
	return res
}

//GetSet 发红包
func (l *lotteryController) GetSet() string {
	uid, uidErr := l.Ctx.URLParamInt("uid")
	num, numErr := l.Ctx.URLParamInt("num")
	money, moneyErr := l.Ctx.URLParamInt("money")
	if uidErr != nil || numErr != nil || moneyErr != nil {
		return "错误"
	}

	// 转换单位为分
	totalMoney := money * 100
	if uid < 1 || totalMoney < num || num < 1 {
		return "参数错误"
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 最大红包占比
	rate := float64(num / totalMoney)
	rmax := 0.3
	if rate > 0.7 {
		rmax = 0.2
	} else if rate > 0.5 {
		rmax = 0.5
	} else if rate > 0.3 {
		rmax = 0.3
	} else {
		rmax = 0.6
	}
	list := make([]uint, num)
	leftNum := num
	leftMoney := totalMoney

	for leftNum > 0 {
		// 最后一个红包
		if leftNum == 1 {
			list[num-1] = uint(leftMoney)
			break
		}

		// 剩余最小金额和剩余红包数量相同
		if leftMoney == leftNum {
			for i := num - leftNum; i < num; i++ {
				list[i] = 1
			}
			break
		}

		/**
		随机红包
		1/ 保留剩余数量最小红包
		2/ *最大红包占比
		3/ 计算随机数
		*/
		rMoney := int((float64(leftMoney-leftNum) - 1) * rmax)
		m := r.Intn(rMoney)
		list[num-leftNum] = uint(m)
		leftMoney -= m
		leftNum--
	}

	id := r.Uint32()
	packageList.set(id, list)
	log.Printf("发红包记录:%v", list)
	return fmt.Sprintf("/get?id=%d", id)
}

//GetGet 领取红包
func (l *lotteryController) GetGet() string {
	id, err := l.Ctx.URLParamInt("id")
	if err != nil {
		return err.Error()
	}

	list := packageList.get(uint32(id))
	if list == nil {
		return "红包不存在"
	}

	// 构造抢红包任务
	callback := make(chan uint)
	t := task{id: uint32(id), callback: callback}
	// 投递任务，接受任务处理结果
	chTask := chTasks[id%taskNum]
	chTask <- t
	money := <-callback

	if money <= 0 {
		return "未抢到"
	}
	log.Printf("抢红包记录 id:%d, 金额:%d", id, money)
	return fmt.Sprintf("恭喜你抢到红包%d", money)
}

//func fetchPackageListMoney(chTask chan task) {
func fetchPackageListMoney(taskId int) {
	for {
		t := <-chTasks[taskId]
		id := t.id

		list := packageList.get(id)
		if list != nil {
			num := len(list)

			// 随机红包id
			idx := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)

			// 红包金额
			money := list[idx]
			if len(list) > 1 {
				if idx == len(list)-1 {
					list = list[:idx]
				} else if idx == 0 {
					list = list[1:]
				} else {
					list = append(list[:idx], list[idx+1:]...)
				}
				packageList.set(id, list)
			} else {
				packageList.del(id)
			}
			t.callback <- money
		} else {
			t.callback <- 0
		}
	}
}
