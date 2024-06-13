package znet

import (
	"math"
	"sync"
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
	L sync.RWMutex
	// 存放连接信息
	Conn ziface.IConnection
	// 存放请求数据
	Data ziface.IMessage
	// RouterSlice相关
	Handlers []ziface.RouterHandler
	Index    int8
	// 存放上下文相关信息
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

func (r *Request) BindRouterSlice(handlers []ziface.RouterHandler) {
	r.Handlers = handlers
}

func (r *Request) ExecRouteHandlerNext() {
	// 是否用于责任链模式？
	r.Index++
	for _, f := range r.Handlers {
		f(r)
		r.Index++
	}
}

func (r *Request) Abort() {
	r.Index = int8(len(r.Handlers))
}

func (r *Request) Set(key string, value any) {
	r.L.RLock()
	// 懒加载创建r.Keys
	if r.Keys == nil {
		r.Keys = make(map[string]any)
	}
	r.Keys[key] = value
	r.L.RUnlock()
}

func (r *Request) Get(key string) (value any, exist bool) {
	r.L.RLock()
	defer r.L.RUnlock()
	value, ok := r.Keys[key]
	if ok {
		return value, true
	} else {
		return nil, false
	}
}

// 复制一份request对象（仅复制request中的Keys键值对信息）
func (r *Request) Copy() ziface.IRequest {
	newRequest := &Request{
		Conn:     nil,
		Data:     r.Data,
		Handlers: nil,
		Index:    math.MaxInt8,
		// 存放上下文相关信息
	}
	// copy request.Keys to newRequest
	for key, value := range r.Keys {
		newRequest.Keys[key] = value
	}
	return newRequest
}
