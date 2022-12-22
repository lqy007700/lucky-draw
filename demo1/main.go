/**
抽奖
 */
package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var userList []string

var lock sync.Mutex

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
	userList = []string{}
	lock = sync.Mutex{}

	app.Run(iris.Addr(":8080"))
}

func (c *lotteryController) Get() string {
	i := len(userList)
	return fmt.Sprintf("当前参与抽奖的用户数%d \n", i)
}

func (c *lotteryController) PostImport() string {
	value := c.Ctx.FormValue("users")
	split := strings.Split(value, ",")
	lock.Lock()
	defer lock.Unlock()
	count1 := len(split)

	for _, s := range split {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			userList = append(userList, s)
		}
	}

	count2 := len(userList)
	return fmt.Sprintf("当前总共参与人数%d，成功导入的用户数%d \n", count1, count2)
}

func (c *lotteryController) GetLucky() string {
	lock.Lock()
	defer lock.Unlock()

	count := len(userList)
	if count > 1 {
		seed := time.Now().UnixNano()
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		user := userList[index]
		userList = append(userList[0:index], user[index+1:])
		return fmt.Sprintf("当前中奖用户：%s,剩余用户数量：%d \n", user, count-1)
	} else if count == 1 {
		user := userList[0]
		userList = []string{}
		return fmt.Sprintf("当前中奖用户：%s,剩余用户数量：%d \n", user, count-1)
	} else {
		return "没有参与用户"
	}
}
