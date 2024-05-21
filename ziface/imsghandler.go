package ziface

/*
	消息管理模块
*/
type IMsgHandler interface {
	// 添加MsgID与Router的映射
	AddRouter(msgID uint32,  router IRouter)
	// 执行实际的handle方法
	DoMsgHandle(request IRequest)
}