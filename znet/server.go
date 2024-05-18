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

// 实现 ziface.IServer interface

func (s *Server) Start() {
	// 不阻塞当前goroutine
	go func() {
		addr := fmt.Sprintf("%s:%s", s.IP, strconv.Itoa(s.Port))
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

			// One conn, One goroutine
			go func() {
				defer conn.Close()
				for {
					buf := make([]byte, 512)
					content, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						break
					}
					// 使用回显作为tcp返回
					_, err = conn.Write(buf[:content])
					if err != nil {
						fmt.Println("write back buf err ", err)
						break
					}else {
						fmt.Printf("Server端回写成功, content is %s\n", string(buf))
					}
				}
			}()
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
