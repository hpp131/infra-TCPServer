package ziface

// "tcpserver/znet"

type IRouter interface {
	// 确定不返回一个带有error的返回值？
	PreHandle(request IRequest)
	Handle(request IRequest)
	PostHandle(request IRequest)
}

// RouterHandler用于取代IRouter

type IGroupRouterSlice interface {
	Use(handles ...RouterHandler)
	// 添加用于处理业务逻辑的handler
	AddHandler(msgID uint32, handler ...RouterHandler)
}

type RouterHandler func(req IRequest)

type IRouterSlice interface {
	// 相当于中间件功能，支持一次添加多个处理函数
	Use(handles ...RouterHandler)
	// 添加用于处理业务逻辑的handler
	AddHandler(msgID uint32, handler ...RouterHandler)
	// 路由组
	Group(start, end uint32, handlers ...RouterHandler) IGroupRouterSlice
	// 获取某个msgID所关联的RouterHandler
	GetHandler(msgID uint32) ([]RouterHandler, bool)
}
