package znet

import (
	"fmt"
	"my_zinx/ziface"
)


type MsgHandle struct {
    Apis map[uint32] ziface.IRouter
}


func NewMsgHandle() *MsgHandle {
    return &MsgHandle { make(map[uint32]ziface.IRouter)  }
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
