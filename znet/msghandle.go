package znet

import (
	"fmt"
	"my_zinx/utils"
	"my_zinx/ziface"
)


type MsgHandle struct {
    // msgID --> router
    Apis map[uint32] ziface.IRouter

    // mutil task
    TaskQueue []chan ziface.IRequest
    // 
    WorkerPoolSize uint32
}


func NewMsgHandle() *MsgHandle {
    return &MsgHandle { 
        Apis : make(map[uint32]ziface.IRouter),
        TaskQueue : make([]chan ziface.IRequest, utils.GlobalObj.WorkerPoolSize),
        WorkerPoolSize : utils.GlobalObj.WorkerPoolSize,
    }
}

func (self *MsgHandle) DoHandle(request ziface.IRequest) {
    msgId := request.GetMsgId()
    router, exist := self.Apis[msgId]
    if !exist {
        fmt.Printf("lack router which is for msgId[%d]\n", msgId)
        return
    }

    // do handle 
    router.PreHandle(request)
    router.Handle(request)
    router.PostHandle(request)
    
}

func (self *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
    if _, exist := self.Apis[msgId]; exist {
        panic(fmt.Sprintf("Repeate router's key -> msgId = %d\n", msgId))
    }

    // add router`
    self.Apis[msgId] = router
}

// StartWorkerPool only once.
func (self *MsgHandle) StartWorkerPool() {
    var size int = int(self.WorkerPoolSize)
    for i := 0; i < size; i++ {
        // construct TaskQueue for every worker
        self.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObj.MaxTaskSize)

        // start work 
        go self.work(i)
        
    }
}


func (self *MsgHandle) work(index int) {
    fmt.Printf("Worker[%d] start to work.\n", index)
    for {
        select {
        case request := <-self.TaskQueue[index]:
            self.DoHandle(request)
        }
    }
}

func (self *MsgHandle) PushTask(request ziface.IRequest) {
    // get average value by connId
    connId := request.GetConnection().GetConnID()
    workId := connId % self.WorkerPoolSize 
    fmt.Printf("Conn[%d] add task to worker[%d]\n", connId, workId)

    // push task
    self.TaskQueue[workId] <- request
}
