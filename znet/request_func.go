package znet

import (
	"tcpserver/ziface"
)


type RequestFunc struct {
	ziface.IRequest
	Conn    ziface.IConnection
	CB      func()
}

func NewRequestFunc(conn ziface.IConnection, cb func()) *RequestFunc {
	return &RequestFunc{
		Conn: conn,
		CB: cb,
	}
}