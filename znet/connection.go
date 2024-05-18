package znet

import (
	"fmt"
	"net"
	"tcpserver/ziface"
)

type Connection struct {
	Conn        *net.TCPConn
	ConnID      uint32
	IsClosed    bool
	// 处理该connection的功能函数
	Handle      ziface.HandleFunc
	ExitBufChan chan bool
}

func  NewConnection(conn *net.TCPConn, connID uint32, handle ziface.HandleFunc) *Connection {
	return &Connection{
		ExitBufChan: make(chan bool),
		Conn: conn,
		IsClosed: false,
		ConnID: connID,
		Handle: handle,
	}
}

// implement ziface.IConnection

// 从conn读数据
func (c *Connection) startReader()  {
	fmt.Println("Read Goroutine is running")
	defer fmt.Printf("Terninating connectin with remoteaddr: %s\n", c.GetRemoteAddr().String())
	defer c.Stop()
	for {
		readBuf := make([]byte, 512)
		contentBytes, err := c.Conn.Read(readBuf)
		if err != nil {
			fmt.Println("recv buf err", err)
			c.ExitBufChan <- true
			continue
		}
		// 调用与当前connection绑定的handle
		err = c.Handle(c.Conn, readBuf[:contentBytes], contentBytes)
		if err != nil {
			fmt.Println("ConnID", c.ConnID, "handle exec error")
			c.ExitBufChan <- true
			return
		}
		
	}
}


func (c *Connection) Start() {
	go c.startReader()
	// 1. 读协程退出则该函数退出
	// 2. 当调用c.Stop()时，会向c.ExitBufChan发送信号，从而退出该函数
	for {
		select {
		case <-c.ExitBufChan:
			return
		}
	}
}

func (c *Connection) Stop()  {
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	if err := c.Conn.Close();err != nil{
		fmt.Println("ConnID", c.ConnID, "close error")
		return
	}
	c.ExitBufChan <- true
	// channel资源回收
	close(c.ExitBufChan)
}

// 获取FD
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}


// type Addr interface {
// 	Network() string // name of the network (for example, "tcp", "udp")
// 	String() string  // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
// }
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}