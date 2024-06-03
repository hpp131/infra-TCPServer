package async

import (
	"fmt"
	"sync"
)

var (
	asyncWorkerArray = [100]*AsyncWorker{}
	// 每个AsyncWorker.c中的最大任务数量
	asyncWorkerTask = 100
	asyncWorkerLock = &sync.Mutex{}
)

type AsyncWorker struct {
	// 用于存储需要异步执行的task
	c chan func()
}

func NewAsyncWorker() *AsyncWorker {
	return &AsyncWorker{c: make(chan func(), asyncWorkerTask)}
}

// 添加异步任务
func (aw *AsyncWorker) AddTask(f func()) {
	if f == nil {
		fmt.Println("Asynchronous task is nil")
		return
	}
	aw.c <- func() {
		defer func ()  {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		f()
	}
}

// 执行异步任务
func (aw *AsyncWorker) Process() {
	for {
		task := <-aw.c
		task()
	}
}

// id暂且代表connection id,这里相当于多个connection共用一个异步线程
func GetAsyncWork(id int) *AsyncWorker {
	if id < 0 {
		id = -id
	}
	workID := id % len(asyncWorkerArray)
	work := asyncWorkerArray[workID]
	if work != nil {
		return work
	}

	// 如果没有找到worker，那么进行初始化并开启一个goroutine来后台执行其中的任务
	asyncWorkerLock.Lock()
	defer asyncWorkerLock.Unlock()

	// 加锁后进行double-check
	work = asyncWorkerArray[workID]
	if work != nil {
		return work
	}

	work = &AsyncWorker{
		c: make(chan func(), asyncWorkerTask),
	}
	asyncWorkerArray[workID] = work
	go work.Process()
	return work
}