package main

import (
	"fmt"
	"net"
	"tcpserver/znet"
)

func main()  {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8999")
	if err != nil {
		panic(err)
	}
	listenr, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listenr.AcceptTCP()
		if err != nil {
			fmt.Println("Listenner Accept error", err)
			break
		}
		var cid uint32
		cid++
		// lack
		c := znet.NewConnection(conn, cid, nil)
		dp := znet.NewDataPackage()
		for {
			headData := make([]byte, dp.GetHeadLen())
			_, err = c.Conn.Read(headData)
			headMsg, err := dp.Unpack(headData)
			if err != nil {
				fmt.Println("Unpack error", err)
				break
			}
			msg := headMsg.(*znet.Message)
			msgLen, _ := msg.GetDataLen(), msg.GetMsgID()
			if msgLen > 0 {
				msg.Data = make([]byte, msg.GetDataLen())
				if _, err := c.Conn.Read(msg.Data);err != nil {
					fmt.Println("Read request body error", err)
				}
				fmt.Printf("Read content from Message.Data, msgID is %d, length is %d, content is %v\n", msg.MsgID, msg.Len, string(msg.Data))
			}else {
				fmt.Println("Recv data from reqeust is empty")
			}
		}
	}
	
}