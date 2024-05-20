package main

import (
	"tcpserver/znet"
)


func main()  {
	s := znet.NewServer("zinx")
	
	s.AddRouter(&znet.PingRouter{})
	s.Serve()
}

