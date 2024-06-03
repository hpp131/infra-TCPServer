package async

// import "google.golang.org/api/books/v1"
import (
	"tcpserver/ziface"
)

// 该对象用于操作异步任务执行结果
type AsyncResult struct {
	conn                ziface.IConnection
	resultObj           any    // 异步任务的返回值
	hasReslt            uint32 // 是否已经返回了结果
	callbackFunc        func() // 异步任务执行完成后执行的回调函数
	hasCallbackFunc     uint32 // 是否有回调函数
	hasExecCallbackFunc uint32 // 回调函数是否已经被执行
}

func NewAsyncResult(conn ziface.IConnection) *AsyncResult {
	return &AsyncResult{
		conn: conn,
	}
}

func (ar *AsyncResult) GetAsyncResult() any {
	return ar.resultObj
}

func (ar *AsyncResult) SetAsyncResult(val any) {
	if ar.hasReslt == 1 {
		return
	}
	ar.resultObj = val
}

// 添加回调函数
func (ar *AsyncResult) OnComplete(f func()) {
	if ar.hasCallbackFunc == 1 {
		return
	}
	ar.callbackFunc = f
}

// 执行回调函数
func (ar *AsyncResult) DoComplete() {
	if ar.callbackFunc == nil {
		return
	}
	ar.callbackFunc()
}
