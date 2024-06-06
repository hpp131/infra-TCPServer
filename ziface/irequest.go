package ziface

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetMsgID() uint32
	BindRouterSlice(handlers []RouterHandler)
	ExecRouteHandlerNext()
	Abort()
}