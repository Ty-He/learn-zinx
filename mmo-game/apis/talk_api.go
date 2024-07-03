package apis 

import (
    "fmt"
    "google.golang.org/protobuf/proto"
    "my_zinx/ziface"
    "my_zinx/znet"
    "my_demo/mmo-game/pb"
    "my_demo/mmo-game/core"
)

type TalkApi struct {
    znet.BaseRouter
}

func (t *TalkApi) Handle(request ziface.IRequest) {
    proto_msg := &pb.Talk{}
    // Unmarshal
    err := proto.Unmarshal(request.GetDate(), proto_msg)
    if err != nil {
        fmt.Println("Unmarshal Error", proto_msg)
        return
    }

    // get pid
    pid, err := request.GetConnection().GetProperty("pid")
    if err != nil {
        fmt.Println("GetProperty Error", err)
        return
    }

    player := core.WorldManagerObj.GetPlayerByPid(pid.(int32))

    player.Talk(proto_msg.Content)
}
