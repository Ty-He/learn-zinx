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

// func (self *PingRouter) PreHandle(request ziface.IRequest) {
//     fmt.Println("Call PingRouter.PreHandle")
//     // send message
//     _, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping\n"))
//     if err != nil {
//         fmt.Println("PreHandle.Write Error")
//     }
// }

// func (self *PingRouter) Handle(request ziface.IRequest) {
//     fmt.Println("Call PingRouter.Handle")
//     // send message
//     _, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
//     if err != nil {
//         fmt.Println("Handle.Write Error")
//     }
// }

// func (self *PingRouter) PostHandle(request ziface.IRequest) {
//     fmt.Println("Call PingRouter.PostHandle")
//     // send message
//     _, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
//     if err != nil {
//         fmt.Println("PostHandle.Write Error")
//     }
// }

func (self *PingRouter) Handle(request ziface.IRequest) {
    fmt.Printf("Recv from client: msgId=%d,msgData=%s\n", request.GetMsgId(), request.GetDate())

    // write to client
    if err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping")); err != nil {
        fmt.Println("SendMsg Error:", err)
    }
}



// router-define
type HelloRouter struct {
    znet.BaseRouter
}
func (self *HelloRouter) Handle(request ziface.IRequest) {
    fmt.Printf("Recv from client: msgId=%d,msgData=%s\n", request.GetMsgId(), request.GetDate())

    // write to client
    if err := request.GetConnection().SendMsg(1, []byte("Hello,Zinx!")); err != nil {
        fmt.Println("SendMsg Error:", err)
    }
}


//

func main() {
    // server 
    s := znet.NewServer()
    // add pingrouter
    s.AddRouter(0, &PingRouter{})
    s.AddRouter(1, &HelloRouter{})
    // run serve
    s.Serve()
}
