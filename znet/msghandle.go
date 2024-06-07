package znet

import (
	"fmt"
	"tcpserver/util"
	"tcpserver/ziface"

)

type MsgHandle struct {
	// 存储MsgID与Router的映射关系
	// APIs map[uint32]ziface.IRouter
	// 规范工作goroutine的数量
	WorkPoolSize uint32
	// worker从TaskQueue中消费任务
	TaskQueue   []chan ziface.IRequest
	RouterSlice *RouterSlices
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		// APIs:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: util.Globalobject.MaxPackageSize,
		// 为每一个Worker分配一个TaskQueue
		TaskQueue:   make([]chan ziface.IRequest, util.Globalobject.WorkPoolSize),
		RouterSlice: NewRouterSlices(),
	}
}

// 执行异步任务的回调函数
func (mh *MsgHandle) DoFuncHandle(f ziface.IFuncHandle) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic happended:", err)
		}
	}()
	f.CallFunc()
}


// 执行前置中间件并执行业务路由
func (mh *MsgHandle) DoMsgHandle(request ziface.IRequest) {
	handlers, ok := mh.RouterSlice.GetHandler(request.GetMsgID())
	if !ok {
		fmt.Println("GetHandler failed")
	}
	request.BindRouterSlice(handlers)
	// 执行中间件在内的所有handle函数
	request.ExecRouteHandlerNext()
}


/*
	以下三个方法为RouterSlice相关:

	AddRouterSlice(msgID uint32, handler ...RouterHandler) IRouterSlice
	Use(msgID uint32, handlers ...RouterHandler) ziface.IRouterSlice
	Group(start, end int, rs IRouterSlice, handlers ...RouterHandler) ziface.IGroupRouterSlice
*/

func (mh *MsgHandle) AddRouterSlice(msgID uint32, handler ...ziface.RouterHandler) ziface.IRouterSlice {
	mh.RouterSlice.AddHandler(msgID, handler...)
	return mh.RouterSlice
}

func (mh *MsgHandle) Use(handlers ...ziface.RouterHandler) ziface.IRouterSlice {
	mh.RouterSlice.Use(handlers...)
	return mh.RouterSlice
}

func (mh *MsgHandle) Group(start, end uint32, handlers ...ziface.RouterHandler) ziface.IGroupRouterSlice {
	return mh.RouterSlice.Group(start, end, handlers...)
}

func (mh *MsgHandle) StartOneWorker(workID int, taskQueue chan ziface.IRequest) {

	// 该worker负责消费该taskQueue
	for {
		select {
		case req := <-taskQueue:
			switch typ := req.(type){
			case ziface.IFuncHandle:
				mh.DoFuncHandle(typ)
			case ziface.IRequest:
				mh.DoMsgHandle(typ)
			}

			
		}
	}
}

func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(util.Globalobject.WorkPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, util.Globalobject.MaxWorkTaskLen)
		// 一次性启动全部goroutine，在后台常驻处理Task
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 这里使用hash的方式对Connection做计算得到服务于该链接的Worker编号，由该worker处理该链接上的请求数据
func (mh *MsgHandle) SendTaskToQueue(request ziface.IRequest) {
	workID := request.GetConnection().GetConnID() % util.Globalobject.WorkPoolSize
	mh.TaskQueue[workID] <- request
}
