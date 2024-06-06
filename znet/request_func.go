package znet

import (
	"tcpserver/ziface"
)

type BaseRequest struct {
	
}

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

func (rf *RequestFunc) CallFunc()  {
	if rf.CB == nil {
		return
	}
	rf.CB()
}