package main

import (
	"fmt"
	"io"
	"net"
	"tcpserver/znet"
	"time"
)

func main()  {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		panic(err)
	}
	dp := znet.NewDataPackage()
	msg := znet.NewMessage([]byte("TCPServer V0.5 Client Test Message"), 1)
	for {
		cont, err := dp.Pack(msg)
		if err != nil {
			fmt.Println("Test Client Pack error", err)
		}
		conn.Write(cont)
		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = conn.Read(headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Msg: ID=", msg.MsgID, ", len=", msg.Len, ", data=", string(msg.Data))
		}

		time.Sleep(1*time.Second)
	}

	select {}
}
