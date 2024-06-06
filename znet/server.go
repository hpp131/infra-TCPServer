package znet

import (
	"fmt"
	"net"
	"strconv"
	"tcpserver/util"
	"tcpserver/ziface"
	"time"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	// Router   ziface.IRouter
	// 使用MsgHandler代替Router实现多路由
	MsgHandler ziface.IMsgHandler
	ConnManage  ziface.IConnManager
	StartHook ziface.Hook
	StopHook ziface.Hook
}

func NewServer(name string) *Server {
	return &Server{
		Name:      util.Globalobject.Name,
		IPVersion: "tcp4",
		IP:        util.Globalobject.Host,
		Port:      int(util.Globalobject.TCPPort),
		MsgHandler: NewMsgHandle(),
		ConnManage: NewConnManage(),
	}
}


// 实现 ziface.IServer interface
func (s *Server) Start() {
	fmt.Println("Server version:", util.Globalobject.Version)
	fmt.Println("Listen Port:", util.Globalobject.TCPPort)

	// 开启WorkerPool
	s.MsgHandler.StartWorkerPool()
	// 不阻塞当前goroutine
	go func() {
		addr := fmt.Sprintf("%s:%s", s.IP, strconv.Itoa(s.Port))

		// network: "tcp | tcp6"
		// address: "ip:port"
		tcpAddr, err := net.ResolveTCPAddr(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}
		// 通过net.listenxxx方法获取fd句柄
		listenner, err := net.ListenTCP(s.IPVersion, tcpAddr)
		if err != nil {
			panic(err)
		}
		fmt.Printf("listening successfully... IP is %s\n", addr)
		for {
			// Accept 是阻塞的
			conn, err := listenner.AcceptTCP()
			if err != nil {
				panic(err)
			}
			var cid uint32
			// 创建连接前先判断当前server端的连接数是否达到配置的最大数量
			if s.ConnManage.Len() > util.Globalobject.MaxConn {
				conn.Close()
				continue
			}
			c := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go c.Start()
		}
	}()
}
func (s *Server) Stop() {
	fmt.Println("[STOP] tcp server , name " , s.Name)
	s.ConnManage.ClearAll()	
}

func (s *Server) Serve() {
	s.Start()
	for {
        time.Sleep(10*time.Second)
    }
}

// func (s *Server) AddRouter(MsgID uint32, router ziface.IRouter)  {
// 	s.MsgHandler.AddRouter(MsgID, router)
// }


/*
	以下三个方法为RouterSlice相关:

	AddRouterSlice(msgID uint32, handler ...RouterHandler) IRouterSlice
	Use(msgID uint32, handlers ...RouterHandler) ziface.IRouterSlice
	Group(start, end int, rs IRouterSlice, handlers ...RouterHandler) ziface.IGroupRouterSlice
*/

func (s *Server) AddRouterSlice(msgID uint32, handler ...ziface.RouterHandler) ziface.IRouterSlice {
	return s.MsgHandler.AddRouterSlice(msgID, handler...)
}

func (s *Server) Use(msgID uint32, handlers ...ziface.RouterHandler) ziface.IRouterSlice {
	return s.MsgHandler.Use(msgID, handlers...)
}

func (s *Server) Group(start, end uint32, handlers ...ziface.RouterHandler) ziface.IGroupRouterSlice {
	return s.MsgHandler.Group(start, end, handlers...)
}


func (s *Server) GetConnManage() ziface.IConnManager {
	return s.ConnManage
}

func (s *Server) SetStartHook(in ziface.Hook)  {
	s.StartHook = in
}

func (s *Server) SetStopHook(in ziface.Hook)  {
	s.StopHook = in
}

func (s *Server) CallStartHook(in ziface.IConnection)  {
	s.StartHook(in)
}

func (s *Server) CallStopHook(in ziface.IConnection)  {
	s.StopHook(in)
}