package main

import (
	"fmt"
	"my_zinx/ziface"
	"my_zinx/znet"
)

// router-define
type PingRouter struct {
    znet.BaseRouter
}

func (self *PingRouter) PreHandle(request ziface.IRequest) {
    fmt.Println("Call PingRouter.PreHandle")
    // send message
    _, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping\n"))
    if err != nil {
        fmt.Println("PreHandle.Write Error")
    }
}

func (self *PingRouter) Handle(request ziface.IRequest) {
    fmt.Println("Call PingRouter.Handle")
    // send message
    _, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
    if err != nil {
        fmt.Println("Handle.Write Error")
    }
}

func (self *PingRouter) PostHandle(request ziface.IRequest) {
    fmt.Println("Call PingRouter.PostHandle")
    // send message
    _, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
    if err != nil {
        fmt.Println("PostHandle.Write Error")
    }
}


func main() {
    // server 
    s := znet.NewServer("ZinxV0.3")
    // add pingrouter
    s.AddRouter(&PingRouter{})
    // run serve
    s.Serve()
}
