package main

import (
	"mynet/actor"
	"mynet/common"
	"mynet/server/world"
	"mynet/server/world/chat"
	"mynet/server/world/cmd"
	"mynet/server/world/data"
	"mynet/server/world/mail"
	"mynet/server/world/player"
	"mynet/server/world/social"
	"mynet/server/world/toprank"
)

func InitMgr(serverName string) {
	//一些共有数据量初始化
	common.Init()
	if serverName == "account" {
	} else if serverName == "netgate" {
	} else if serverName == "world" {
		cmd.Init()
		data.InitRepository()
		player.MGR.Init(1000)
		chat.MGR.Init(1000)
		mail.MGR.Init(1000)
		toprank.MGR().Init(1000)
		player.SIMPLEMGR.Init(1000)
		social.MGR().Init(1000)
		actor.MGR.InitActorHandle(world.SERVER.GetCluster())
	}
}

//程序退出后执行
func ExitMgr(serverName string) {
	if serverName == "account" {
	} else if serverName == "netgate" {
	} else if serverName == "world" {
	}
}
