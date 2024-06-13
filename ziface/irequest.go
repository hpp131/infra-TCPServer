package ziface

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetMsgID() uint32
	BindRouterSlice(handlers []RouterHandler)
	ExecRouteHandlerNext()
	Abort()
	// 用于requestPoolMode
	Set(key string, value any)
	Get(key string) (value any, exist bool)
	Copy() IRequest
}
