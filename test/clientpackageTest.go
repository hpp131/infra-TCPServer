package main

import (
	"fmt"
	"net"
	"tcpserver/znet"
)

func main()  {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		panic(err)
	}
	msg1 := znet.NewMessage([]byte("hello"), 1, 5)
	msg2 := znet.NewMessage([]byte("world"), 2, 5)
	dp := znet.NewDataPackage()
	sendMsg1, _ := dp.Pack(msg1)
	sendMsg2, _ := dp.Pack(msg2)
	sendMsg1 = append(sendMsg1, sendMsg2...)
	content, err := conn.Write(sendMsg1)
	if err != nil {
		fmt.Println("Data Write error", err)
	}
	fmt.Println("Data Write Successfully, length is %d\n", content)
	// 阻塞执行，在server端观察状态
	select {}
}
