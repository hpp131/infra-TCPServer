package main


import (
	"fmt"
	"tcpserver/ziface"
	"tcpserver/znet"
)


func main() {
	// 创建一个server句柄
	s := znet.NewServer("tcpserver")
	s.SetStartHook(func(conn ziface.IConnection) {
		fmt.Println("Start Hook Worked...")
	})
	s.SetStopHook(func(conn ziface.IConnection) {
		fmt.Println("Stop Hook Worked...")
	})

	// 配置路由
	// TODO

	//开启服务
	s.Serve()
}