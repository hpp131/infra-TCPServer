package znet

import (
	"fmt"
	"tcpserver/ziface"
)

type MsgHandle struct {
	// 存储MsgID与Router的映射关系
	APIs map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		APIs: make(map[uint32]ziface.IRouter),
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

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) ()  {
	if _, ok := mh.APIs[msgID]; !ok {
		fmt.Printf("AddRouter error, MsgID %d already exist\n", msgID)
	}
	mh.APIs[msgID] = router
}