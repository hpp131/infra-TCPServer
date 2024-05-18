package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp4", ":8999")
	if err != nil {
		panic("net.Dial tcp server failed...")
	}
	defer conn.Close()
	for {
		// 持续向TCP Server发送数据
		time.Sleep(3 * time.Second)
		req, err := conn.Write([]byte("hello, I'm client"))
		if err != nil {
			fmt.Println("client write content failed...")
			return
		}
		fmt.Printf("client write succ, totally write %d bytes\n", req)
		recvBuf := make([]byte,10)
		resp, err := conn.Read(recvBuf)
		if err != nil {
			fmt.Println("client recv failed")
		}else {
			fmt.Printf("client recv content is %s, length is %d\n", string(recvBuf), resp)
		}
	}

}
