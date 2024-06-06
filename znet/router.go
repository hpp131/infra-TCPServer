package znet

import (
	"sync"
	"tcpserver/ziface"
)

// 基础类，其他类可继承该类并重写部分方法
type BaseRouter struct{}

func (br *BaseRouter) PreHandle(request ziface.IRequest) {
	// TODO
}

func (br *BaseRouter) Handle(request ziface.IRequest) {
	// TODO
}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {
	// TODO
}

// Implement ziface.IRouterSlice
type RouterSlices struct {
	APIs map[uint32][]ziface.RouterHandler
	RH   []ziface.RouterHandler
	Mu   sync.RWMutex
}

func NewRouterSlices() *RouterSlices {
	return &RouterSlices{
		APIs: make(map[uint32][]ziface.RouterHandler),
		RH:   make([]ziface.RouterHandler, 0),
	}
}

func (rs *RouterSlices) Use(handler ...ziface.RouterHandler) {
	rs.RH = append(rs.RH, handler...)
}

// 把rs.RH和handler全都append到rs.APIs中
func (rs *RouterSlices) AddHandler(msgID uint32, handler ...ziface.RouterHandler) {
	if _, ok := rs.APIs[msgID]; ok {
		panic("RouterHandler is existed")
	}

	finalSize := len(rs.RH) + len(handler)
	finanSlice := make([]ziface.RouterHandler, finalSize)
	copy(finanSlice, rs.RH)
	copy(finanSlice[len(rs.RH):], handler)
	rs.APIs[msgID] = append(rs.APIs[msgID], handler...)
}

func (rs *RouterSlices) GetHandler(msgID uint32) ([]ziface.RouterHandler, bool) {
	// 为什么需要加锁？
	rs.Mu.RLock()
	defer rs.Mu.RUnlock()
	if res, ok := rs.APIs[msgID]; !ok {
		return nil, false
	} else {
		return res, true
	}
}

func (rs *RouterSlices) Group(start, end uint32, handlers ...ziface.RouterHandler) ziface.IGroupRouterSlice {
	return NewGroupRouter(start, end, rs, handlers...)
}

type GroupRouter struct {
	Start, End uint32
	Handlers   []ziface.RouterHandler
	Router     ziface.IRouterSlice
}

func NewGroupRouter(start, end uint32, router ziface.IRouterSlice, handler ...ziface.RouterHandler) *GroupRouter {

	res :=  &GroupRouter{
		Start: start,
		End:   end,
		Handlers: make([]ziface.RouterHandler, 0, len(handler)),
		Router: router,
	}
	res.Handlers = append(res.Handlers, handler...)
	return res	
}

func (gr *GroupRouter) Use(handles ...ziface.RouterHandler)  {
	gr.Handlers = append(gr.Handlers, handles...)
}

func (gr *GroupRouter) AddHandler(msgID uint32, handler ...ziface.RouterHandler)  {
	if msgID < gr.Start || msgID > gr.End {
		panic("msgID is out of range in current Group")
	}
	size := len(gr.Handlers) + len(handler)
	finalSlice := make([]ziface.RouterHandler, 0, size)
	copy(finalSlice, gr.Handlers)
	copy(finalSlice[len(gr.Handlers):], handler)
	gr.Router.AddHandler(msgID, finalSlice...)
}



