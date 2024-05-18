package main

import (
	"tcpserver/znet"
)


func main()  {
	s := znet.NewServer("zinx")
	s.Serve()
}

