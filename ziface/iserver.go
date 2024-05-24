package ziface


type Hook func (IConnection)


type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgID uint32, router IRouter)
	GetConnManage() IConnManager
	SetStartHook(Hook)
	SetStopHook(Hook)
	CallStartHook(conn IConnection)
	CallStopHook(conn IConnection)
}