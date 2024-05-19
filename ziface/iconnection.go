package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	// net.TCPConn 是个struct，内部有一个conn成员变量，conn中又有一个fd成员变量，fd是*syscall.FD，即socket句柄。net.TCPConn实现了Conn interface
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	GetRemoteAddr() net.Addr
}

type HandleFunc func(*net.TCPConn, []byte, int) error