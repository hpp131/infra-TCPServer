package znet

import (
	"tcpserver/ziface"
)

func NewRequest(conn ziface.IConnection, data ziface.IMessage) *Request {
	return &Request{
		Conn: conn,
		Data: data,
		// 责任链模式？
		Index: -1,
	}
}

// Implement IRequest interface
type Request struct {
	// 存放连接信息
	Conn ziface.IConnection
	// 存放请求数据
	Data ziface.IMessage
	// RouterSlice相关
	Handlers []ziface.RouterHandler
	Index int8
	Keys map[string]any // 懒加载创建
}



func (r *Request) GetConnection() ziface.IConnection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.Data.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.Data.GetMsgID()
}

func (r *Request) BindRouterSlice(handlers []ziface.RouterHandler)  {
	r.Handlers = handlers
}

func (r *Request) ExecRouteHandlerNext()  {
	// 是否用于责任链模式？
	r.Index++
	for _, f := range r.Handlers {
		f(r)
		r.Index++
	}
}


func (r *Request) Abort()  {
	r.Index = int8(len(r.Handlers))
}


