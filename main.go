package main

import (
	"fmt"
	"wangpan/dao/mysql"
	"wangpan/dao/redis"
	"wangpan/routes"
	"wangpan/settings"
)

func main() {
	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed,err:%v\n", err)
	}

	//2.初始化MySQL
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
	}
	defer mysql.Close()
	//3.初始化Redis
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed,err:%v\n", err)
	}
	defer redis.CLose()
	//4.注册路由
	r := routes.Setup()
	//5.启动服务
	r.Run()

}
