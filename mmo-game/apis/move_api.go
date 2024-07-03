package apis 

import (
    "fmt"
    "google.golang.org/protobuf/proto"
    "my_zinx/ziface"
    "my_zinx/znet"
    "my_demo/mmo-game/pb"
    "my_demo/mmo-game/core"
)

type MoveApi struct {
    znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
    position_proto_msg := &pb.Position{} 
    if err := proto.Unmarshal(request.GetDate(), position_proto_msg); err != nil {
        fmt.Println("Unmarshal Error: ", err)
        return 
    }

    pid, err := request.GetConnection().GetProperty("pid")
    if err != nil {
        fmt.Println("GetProperty Error: ", err)
        return
    }

    core.WorldManagerObj.GetPlayerByPid(pid.(int32)).UpdatePos(
        position_proto_msg.X, 
        position_proto_msg.Y, 
        position_proto_msg.Z, 
        position_proto_msg.V)

}
