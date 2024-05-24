package znet

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"tcpserver/ziface"
	"tcpserver/util"
)

// 继承BaseRouter
type PingRouter struct {
	BaseRouter
}

// 根据需要重写BaseRouter的方法
func (pr *PingRouter) PreHandle(request ziface.IRequest)  {
	fmt.Printf("PreHandle Request, MsgID %d\n", request.GetMsgID())
	fmt.Println("Call Router PreHandle...")
	err := request.GetConnection().SendMsg([]byte("Before handle"), request.GetMsgID())
	if err != nil {
		fmt.Println("Call Router PreHandle error", err)
	}
}

func (pr *PingRouter) Handle(request ziface.IRequest)  {
	fmt.Printf("Handle Request, MsgID %d\n", request.GetMsgID())
	fmt.Println("Call Router Handle...")
	err := request.GetConnection().SendMsg([]byte("handling"), request.GetMsgID())
	if err != nil {
		fmt.Println("Call Router Handle error", err)
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest)  {
	fmt.Printf("PostHandle Request, MsgID %d\n", request.GetMsgID())
	fmt.Println("Call Router PostHandle...")
	err := request.GetConnection().SendMsg([]byte("PostHandle"), request.GetMsgID())
	if err != nil {
		fmt.Println("Call Router PostHandle error", err)
	}
}


type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	// Router   ziface.IRouter
	// 使用MsgHandler代替Router实现多路由
	MsgHandler ziface.IMsgHandler
	ConnManage  ziface.IConnManager
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

// // Connection中的handle成员
// func CallBackHandl(conn *net.TCPConn, data []byte, length int) error {
// 	contentBytes, err := conn.Write(data)
// 	if err != nil {
// 		fmt.Println("write error: ", err)
// 		return err
// 	}
// 	fmt.Printf("Back Write data successs, data length is %d\n", contentBytes)
// 	return nil
// }


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

func (s *Server) AddRouter(MsgID uint32, router ziface.IRouter)  {
	s.MsgHandler.AddRouter(MsgID, router)
}

func (s *Server) GetConnManage() ziface.IConnManager {
	return s.ConnManage
}