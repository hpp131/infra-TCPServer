package ziface

type Hook func (IConnection)


type IServer interface {
	Start()
	Stop()
	Serve()
	// AddRouter(msgID uint32, router IRouter)
	// RouterSlice相关功能
	AddRouterSlice(msgID uint32, handler ...RouterHandler) IRouterSlice
	Use(handlers ...RouterHandler) IRouterSlice
	Group(start, end uint32, handlers ...RouterHandler) IGroupRouterSlice

	GetConnManage() IConnManager
	SetStartHook(Hook)
	SetStopHook(Hook)
	CallStartHook(conn IConnection)
	CallStopHook(conn IConnection)
}

