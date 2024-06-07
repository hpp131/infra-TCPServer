package main

import (
	"fmt"
	"net"
	"tcpserver/znet"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	
	dp := znet.NewDataPackage()
	for {
		// sendmsg
		msg := znet.NewMessage([]byte("ZinxPing"), 1)
		sendMsg, _ := dp.Pack(msg)
		_, err = conn.Write(sendMsg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}
		// recvmsg
		buf := make([]byte, 100)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
		}
		fmt.Println(string(buf))
		time.Sleep(1*time.Second)
	}
}