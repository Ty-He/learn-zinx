package core 

import (
    "fmt"
    "math/rand"
    "sync"
    "google.golang.org/protobuf/proto"
    "my_zinx/ziface"
    "my_demo/mmo-game/pb"
)

type Player struct {
    Pid int32 
    Conn ziface.IConnection
    X float32   // plane x
    Y float32   // height
    Z float32   // plane y 
    V float32   // angle 0-360
}

var PidGen int32 
var mut sync.Mutex 

func NewPlayer(conn ziface.IConnection) *Player {
    mut.Lock()
    PidGen ++
    id := PidGen
    mut.Unlock()

    return &Player{
        Pid : id,
        Conn : conn,
        X : float32(150 + rand.Intn(25)),
        Y : 0,
        Z : float32(150 + rand.Intn(25)),
        V : 0,
    }
}


func (p *Player) SendMsg(msgId uint32, msg proto.Message) {
    // by protobuf to binary data
    buf, err := proto.Marshal(msg)
    if err != nil {
        fmt.Println("Marshal Err", err)
        return
    }

    // send binary data
    if err := p.Conn.SendMsg(msgId, buf); err != nil {
        fmt.Println("SendMsg Err:", err)
        return
    }
}

func (p *Player) SyncPid() {
    // encapsulate protobuf Message
    proto_msg := &pb.SyncPid{
        Pid : p.Pid,
    }

    // send client
    p.SendMsg(1, proto_msg)
}

func (p *Player) BroadCastInitialPos() {
    proto_msg := &pb.BroadCast{
        Pid : p.Pid,
        Tp: 2,
        Data: &pb.BroadCast_P{
            P : &pb.Position{
                X : p.X,
                Y : p.Y,
                Z : p.Z,
                V : p.V,
            },
        },
    }

    p.SendMsg(200, proto_msg)
}

func (p *Player) Talk(content string) {
    // msgId = 200 BroadCast
    proto_msg := &pb.BroadCast{
        Pid: p.Pid,
        Tp: 1, // Talk 
        Data: &pb.BroadCast_Content{
            Content: content,
        },
    }

    players := WorldManagerObj.GetTotalPlayers()

    // BroadCast
    for i := range players {
        players[i].SendMsg(200, proto_msg)
    }
}

// 
func (p *Player) SyncSourrounding() {
    pids := WorldManagerObj.AoiManager.GetPlayerIds(p.X, p.Z)

    players := make([]*Player, 0, len(pids))
    for i := range pids {
        players = append(players, WorldManagerObj.GetPlayerByPid(pids[i]))
    }

    // BroadCast Current Player Position
    BroadCast_proto_msg := &pb.BroadCast{
        Pid : p.Pid,
        Tp : 2,
        Data: &pb.BroadCast_P{
            P : &pb.Position{
                X : p.X,
                Y : p.Y,
                Z : p.Z,
                V : p.Z,
            },
        },
    }

    for i := range players {
        players[i].SendMsg(200, BroadCast_proto_msg)
    }

    // SyncSourrounding Other Player
    ps := make([]*pb.Player, 0, len(pids))
    for _, player := range players {
        p := &pb.Player{
            Pid: player.Pid,
            P: &pb.Position{
                X : player.X,
                Y : player.Y,
                Z : player.Z,
                V : player.V, 
            },
        }
        ps = append(ps, p)
    }
    SyncSourrounding_proto_msg := &pb.SyncPlayers{
        Ps : ps,
    }
    p.SendMsg(202, SyncSourrounding_proto_msg)
} 


func (p *Player) UpdatePos(x, y, z, v float32) {
    // Update Current Position 
    p.X, p.Y, p.Z, p.V = x, y, z, v

    // 
    move_proto_msg := &pb.BroadCast{
        Pid: p.Pid,
        Tp: 4,
        Data: &pb.BroadCast_P{
            P: &pb.Position{
                X: x,
                Y: y,
                Z: z,
                V: v,
            },
        },
    }

    for _, player := range p.GetVisiblePlayers() {
        player.SendMsg(200, move_proto_msg)
    }

}


// get all player which in same area
func (p *Player) GetVisiblePlayers() (players []*Player) {
    pids := WorldManagerObj.AoiManager.GetPlayerIds(p.X, p.Z)
    players = make([]*Player, 0, len(pids))

    for _, pid := range pids {
        players = append(players, WorldManagerObj.GetPlayerByPid(pid))
    }
    return
}

func (p *Player) Offline() {
    offline_proto_msg := &pb.SyncPid{
        Pid: p.Pid,
    }

    for _, player := range p.GetVisiblePlayers() {
        // it's will send to self, by note err
        player.SendMsg(201, offline_proto_msg)
    }

    WorldManagerObj.RemovePlayer(p.Pid)
}
