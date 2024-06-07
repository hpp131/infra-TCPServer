package main

import (
	"tcpserver/ziface"
	"tcpserver/znet"
) 

func PingFunc(request ziface.IRequest)  {
	request.GetConnection().SendMsg([]byte("Ping... Ping..."), 1)
}

func AuthFunc(request ziface.IRequest)  {
	// fmt.Println("Auth... Auth...")
	request.GetConnection().SendMsg([]byte("Auth... Auth..."), 1)
}


func main()  {
	server := znet.NewServer("[zinx testing]")
	server.Use(AuthFunc)
	server.AddRouterSlice(1, PingFunc)
	server.Serve()
}