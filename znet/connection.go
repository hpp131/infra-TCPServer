package znet

import (
	"errors"
	"fmt"
	"net"
	"tcpserver/ziface"

)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	IsClosed bool
	// 处理该connection的功能函数
	// Handle      ziface.HandleFunc
	// 使用Router处理业务，而不是将Handle固定在Connection中
	MsgHandler  ziface.IMsgHandler
	ExitBufChan chan bool
	// 将请求处理分为读写两个线程，使用该chan为两个goroutine提供通信
	MsgChan chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IMsgHandler) *Connection {
	return &Connection{
		ExitBufChan: make(chan bool),
		Conn:        conn,
		IsClosed:    false,
		ConnID:      connID,
		MsgHandler:  router,
		MsgChan:     make(chan []byte),
	}
}

// Implement ziface.IConnection

// 从conn读数据，即read goroutine
func (c *Connection) startReader() {
	fmt.Println("Read Goroutine is running")
	defer fmt.Printf("Terninating connectin with remoteaddr: %s\n", c.GetRemoteAddr().String())
	defer c.Stop()
	dp := NewDataPackage()
	for {
		headBuf := make([]byte, dp.GetHeadLen())
		_, err := c.Conn.Read(headBuf)
		if err != nil {
			fmt.Println("recv buf err", err)
			c.ExitBufChan <- true
			continue
		}

		// 解包操作
		// 获取head数据(msgLen, msgID)
		msg, err := dp.Unpack(headBuf)
		if err != nil {
			fmt.Println("TCP Unpack error:", err)
			c.ExitBufChan <- true
			continue
		}
		var body []byte
		if msg.GetDataLen() > 0 {
			body = make([]byte, msg.GetDataLen())
			_, err := c.Conn.Read(body)
			if err != nil {
				fmt.Println("Read Message Body error", err)
				c.ExitBufChan <- true
				continue
			}
		}
		msg.SetData(body)

		req := &Request{
			Conn: c,
			Data: msg,
		}
		go c.MsgHandler.DoMsgHandle(req)
	}
}

// 向conn写数据,即write goroutine
func (c *Connection) startWriter() {
	fmt.Println("Write Goroutine is Running...")
	for {
		select {
		case data := <- c.MsgChan:
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Println("Send data error", err)
				return
			}
		case <- c.ExitBufChan:
			fmt.Println("Connection has already been closed, Exiting to Write Goroutine...")
			return
		}
	}
}

func (c *Connection) Start() {
	// 读协程负责读取请求数据并将请求数据拆包
	go c.startReader()
	go c.startWriter()
	// 1. 读协程退出则该函数退出
	// 2. 当调用c.Stop()时，会向c.ExitBufChan发送信号，从而退出该函数
	for {
		select {
		case <-c.ExitBufChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	if err := c.Conn.Close(); err != nil {
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

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(data []byte, id uint32) error {
	// Message -> []bytes
	msg := NewMessage(data, id)
	dp := NewDataPackage()
	res, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("SendMsg Pack error", err)
		return errors.New("SendMsg Pack error")
	}
	// 获取到TLV数据后再通过c.MsgChan发送给Write goroutine
	c.MsgChan <- res
	return nil
}
