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
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Before handle"))
	if err != nil {
		fmt.Println("Call Router PreHandle Failed!!!")
	}
}

func (pr *PingRouter) Handle(request ziface.IRequest)  {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Handling"))
	if err != nil {
		fmt.Println("Call Router Handle Failed!!!")
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest)  {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After handling"))
	if err != nil {
		fmt.Println("Call Router PostHandle Failed!!!")
	}
}


type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	Router   ziface.IRouter
}

func NewServer(name string) *Server {
	return &Server{
		Name:      util.Globalobject.Name,
		IPVersion: "tcp4",
		IP:        util.Globalobject.Host,
		Port:      int(util.Globalobject.TCPPort),
		Router: nil,
	}
}

// Connection中的handle成员
func CallBackHandl(conn *net.TCPConn, data []byte, length int) error {
	contentBytes, err := conn.Write(data)
	if err != nil {
		fmt.Println("write error: ", err)
		return err
	}
	fmt.Printf("Back Write data successs, data length is %d\n", contentBytes)
	return nil
}


// 实现 ziface.IServer interface

func (s *Server) Start() {
	fmt.Println("Server version:", util.Globalobject.Version)
	fmt.Println("Listen Port:", util.Globalobject.TCPPort)
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
			c := NewConnection(conn, cid, s.Router)
			cid++
			go c.Start()
		}
	}()
}
func (s *Server) Stop() {
	fmt.Println("[STOP] tcp server , name " , s.Name)
}

func (s *Server) Serve() {
	s.Start()
	for {
        time.Sleep(10*time.Second)
    }
}

func (s *Server) AddRouter(router ziface.IRouter)  {
	// TODO
	s.Router = router
}
