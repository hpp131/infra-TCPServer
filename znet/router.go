package znet

import "tcpserver/ziface"


// 基础类，其他类可继承该类并重写部分方法
type BaseRouter struct {}

func (br *BaseRouter) PreHandle(request ziface.IRequest)  {
	// TODO	
}

func (br *BaseRouter) Handle(request ziface.IRequest)  {
	// TODO	
}

func (br *BaseRouter) PostHandle(request ziface.IRequest)  {
	// TODO	
}