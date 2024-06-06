package ziface

/*
消息管理模块
*/

type IFuncHandle interface {
	CallFunc()
}

type IMsgHandler interface {
	// 添加MsgID与Router的映射
	// AddRouter(msgID uint32, router IRouter)
	// RouterSlice相关
	AddRouterSlice(msgID uint32, handler ...RouterHandler) IRouterSlice
	Use(msgID uint32, handlers ...RouterHandler) IRouterSlice
	Group(start, end uint32, handlers ...RouterHandler) IGroupRouterSlice
	// 执行实际的handle方法
	DoMsgHandle(request IRequest)
	DoFuncHandle(request IFuncHandle)
	StartWorkerPool()
	SendTaskToQueue(request IRequest)
}
