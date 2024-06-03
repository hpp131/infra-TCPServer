package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	// net.TCPConn 是个struct，内部有一个conn成员变量，conn中又有一个fd成员变量，fd是*syscall.FD，即socket句柄。net.TCPConn实现了Conn interface
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	GetRemoteAddr() net.Addr
	GetMsgHandle() IMsgHandler
	// 添加封包方法,发送Message被封包后的数据[]byte
	SendMsg(data []byte, id uint32) error
	SendBufMsg(data []byte, id uint32) error
	SetProperty(key string, value any)
	GetProperty(key string) (any, error)
	RemoveProperty(key string)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
