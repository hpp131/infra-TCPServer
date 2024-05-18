package znet

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func NewServer(name string) *Server {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
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
			c := NewConnection(conn, cid, CallBackHandl)
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
