package znet

import (
	"fmt"
	"sync"
	"tcpserver/ziface"
)


type ConnManage struct {
	lock    sync.RWMutex
	connSet map[uint32]ziface.IConnection
}

func NewConnManage() *ConnManage {
	return &ConnManage{
		connSet: make(map[uint32]ziface.IConnection),
	}
}

func(cm *ConnManage) AddConn(conn ziface.IConnection)  {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	id := conn.GetConnID()
	if _, ok := cm.connSet[id]; ok {
		fmt.Printf("Current Connection %d Already Exist, Will Not Execute AddConn", id)
		return
	}
	cm.connSet[id] = conn
}

func (cm *ConnManage) DeleteConn(conn ziface.IConnection)  {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	delete(cm.connSet, conn.GetConnID())
}

func (cm *ConnManage) GetConn(connID uint32) ziface.IConnection {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	if _, ok := cm.connSet[connID]; !ok {
		fmt.Println("Current Connection IS Not Existing")
		return nil
	}
	return cm.connSet[connID]
}

func (cm *ConnManage) Len() int {
	return len(cm.connSet)
}

func (cm *ConnManage) ClearAll()  {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	for id, conn := range cm.connSet {
		conn.Stop()
		delete(cm.connSet, id)
	}
	fmt.Printf("CleaAll Connection Compeltely", cm.Len())
}
