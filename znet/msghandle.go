package znet

import (
	"fmt"
	"tcpserver/util"
	"tcpserver/ziface"
)

type MsgHandle struct {
	// 存储MsgID与Router的映射关系
	APIs         map[uint32]ziface.IRouter
	// 规范工作goroutine的数量
	WorkPoolSize uint32
	// worker从TaskQueue中消费任务
	TaskQueue    []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		APIs: make(map[uint32]ziface.IRouter),
		WorkPoolSize: util.Globalobject.MaxPackageSize,
		// 为每一个Worker分配一个TaskQueue
		TaskQueue: make([]chan ziface.IRequest, util.Globalobject.WorkPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandle(request ziface.IRequest) {

	router, ok := mh.APIs[request.GetMsgID()]
	if !ok {
		fmt.Println("DoMsgHandle error : Unknow MsgID")
		return
	}

	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.APIs[msgID]; !ok {
		fmt.Printf("AddRouter error, MsgID %d already exist\n", msgID)
	}
	mh.APIs[msgID] = router
}

func (mh *MsgHandle) StartOneWorker(workID int, taskQueue chan ziface.IRequest)  {
	
	// 该worker负责消费该taskQueue
	for {
		select {
		case req := <- taskQueue:
			mh.DoMsgHandle(req)
		}
	}
}

// 
func (mh *MsgHandle) StartWorkerPool()  {
	for i:=0;i<int(util.Globalobject.WorkPoolSize);i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, util.Globalobject.MaxWorkTaskLen)
		// 一次性启动全部goroutine
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}


func (mh *MsgHandle) SendTaskToQueue(request ziface.IRequest)  {
	workID := request.GetConnection().GetConnID() % util.Globalobject.WorkPoolSize
	mh.TaskQueue[workID] <- request
}