package ziface

type Server interface {
	Start()
	Stop()
	Serve()
	AddRouter(router IRouter)
}