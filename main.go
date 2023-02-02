package main

import (
	"fmt"
	"log"
	"openwechat"
	"plug"
)

func main() {

	bot := openwechat.DefaultBot(openwechat.Desktop)

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()

	// 注册消息处理函数
	bot.MessageHandler = plug.Doomsday

	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			log.Printf("login error: %v \n", err)
			return
		}
	}

	// 获取登陆的用户
	// slef, err := bot.GetCurrentUser()
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	//获得群组
	groups, _ := self.Groups()
	go plug.Cron(groups, self)

	bot.Block()
}
