package main 

import (
    "my_zinx/ziface"
    "my_zinx/znet"
    "my_demo/mmo-game/core"
    "my_demo/mmo-game/apis"
)


func OnConnStart(conn ziface.IConnection) {
    player := core.NewPlayer(conn)

    player.SyncPid()
    player.BroadCastInitialPos()
    
    conn.SetProperty("pid", player.Pid)
    // add player to WorldManager
    core.WorldManagerObj.AddPlayer(player)

    // sync player
    player.SyncSourrounding()
}

func main() {
    s := znet.NewServer()

    s.SetOnConnStart(OnConnStart)
    s.SetOnConnStop(OnConnStop)

    s.AddRouter(2, &apis.TalkApi{})
    s.AddRouter(3, &apis.MoveApi{})
    s.Serve()
}

func OnConnStop(conn ziface.IConnection) {
    pid, _ := conn.GetProperty("pid")
    core.WorldManagerObj.GetPlayerByPid(pid.(int32)).Offline()
}
