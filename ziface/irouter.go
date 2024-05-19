package ziface


type IRouter interface {
	// 确定不返回一个带有error的返回值？
	PreHandle(request IRequest)
	Handle(request IRequest)
	PostHandle(request IRequest)
}
