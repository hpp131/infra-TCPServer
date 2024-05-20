package znet

import "tcpserver/ziface"

func NewRequest(conn ziface.IConnection, data []byte) *Request {
	return &Request{
		Conn: conn,
		Data: data,
	}
}

// Implement IRequest interface
type Request struct {
	// 存放连接信息
	Conn ziface.IConnection
	// 存放请求数据
	Data []byte
}




func (r *Request) GetConnection() ziface.IConnection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.Data
}