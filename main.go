package main


import (
	"fmt"
	"tcpserver/ziface"
	"tcpserver/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

type HelloRouter struct {
	znet.BaseRouter
}

//Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	//回写数据
	err := request.GetConnection().SendMsg([]byte("ping...ping...ping"), request.GetMsgID())
	if err != nil {
		fmt.Println(err)
	}
}

func (h *HelloRouter) Handle(request  ziface.IRequest)  {
	fmt.Println("Call HelloRouter Handle")
	//先读取客户端的数据，再回写hello...hello...hello
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	//回写数据
	err := request.GetConnection().SendMsg([]byte("hello...hello...hello"), request.GetMsgID())
	if err != nil {
		fmt.Println(err)
	}
}

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
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//开启服务
	s.Serve()
}