package znet

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"tcpserver/util"
	"tcpserver/ziface"

)

type Connection struct {
	Server   ziface.IServer
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
	// 与上面的chan作用相同，但是带有缓冲
	MsgBuffChan chan []byte
	// 用于链接属性设置
	Property     map[string]any
	PropertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, router ziface.IMsgHandler) *Connection {
	res := &Connection{
		Server:      server,
		ExitBufChan: make(chan bool),
		Conn:        conn,
		IsClosed:    false,
		ConnID:      connID,
		MsgHandler:  router,
		MsgChan:     make(chan []byte),
		MsgBuffChan: make(chan []byte, util.Globalobject.MaxBufBytes),
		Property: make(map[string]any),
	}
	// 创建连接的时候把该链接添加到ConnManage中
	server.GetConnManage().AddConn(res)
	return res
}

// Implement ziface.IConnection

// 1. 从conn读数据，即read goroutine
// 2. 将拆包后的数据组装成Request, 然后发送到任务队列中/开一个临时的goroutine直接handle
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

		if util.Globalobject.WorkPoolSize > 0 {
			c.MsgHandler.SendTaskToQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandle(req)
		}
	}
}

// 向conn写入业务逻辑处理之后的结果（即write goroutine）。不负责处理业务逻辑
func (c *Connection) startWriter() {
	fmt.Println("Write Goroutine is Running...")
	for {
		select {
		case data := <-c.MsgChan:
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Println("Send data error", err)
				return
			}
		case dataBuf, ok := <-c.MsgBuffChan:
			if ok {
				_, err := c.Conn.Write(dataBuf)
				if err != nil {
					fmt.Println("Send data error", err)
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}

		case <-c.ExitBufChan:
			fmt.Println("Connection has already been closed, Exiting to Write Goroutine...")
			return
		}
	}
}

func (c *Connection) Start() {
	// 读协程负责读取请求数据并将请求数据拆包
	go c.startReader()
	go c.startWriter()

	// 执行Hook
	c.Server.CallStartHook(c)

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

	// 连接退出前执行Hook逻辑(若有)
	c.Server.CallStopHook(c)

	if err := c.Conn.Close(); err != nil {
		fmt.Println("ConnID", c.ConnID, "close error")
		return
	}
	c.Server.GetConnManage().DeleteConn(c)
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

func (c *Connection) SendBufMsg(data []byte, id uint32) error {
	// Message -> []bytes
	msg := NewMessage(data, id)
	dp := NewDataPackage()
	res, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("SendMsg Pack error", err)
		return errors.New("SendMsg Pack error")
	}
	// 获取到TLV数据后再通过c.MsgChan发送给Write goroutine
	c.MsgBuffChan <- res
	return nil
}

func (c *Connection) SetProperty(key string, value any) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	c.Property[key] = value
}

func (c *Connection) GetProperty(key string) (any, error) {
	c.PropertyLock.RLock()
	defer  c.PropertyLock.RUnlock()
	if _, ok := c.Property[key]; !ok {
		return nil, fmt.Errorf("%s Property not Exist", key)
	}else {
		return c.Property[key], nil
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	delete(c.Property, key)
}
