package core 

import "sync"

const (
    AOI_MIN_X uint = 50
    AOI_MAX_X uint = 500
    AOI_CNT_X uint = 10
    AOI_MIN_Y uint = 50
    AOI_MAX_Y uint = 500
    AOI_CNT_Y uint = 10
)

type WorldManager struct {
    AoiManager AOIManager 
    Players map[int32] *Player
    rw_lock sync.RWMutex 
}

// global handler
var WorldManagerObj *WorldManager 

func init() {
    WorldManagerObj = &WorldManager{
        AoiManager: *NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNT_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNT_Y),
        Players: make(map[int32] *Player),
    }
}


// AddPlayer By *Players
func (wm *WorldManager) AddPlayer(player *Player) {
    // add to AoiManager 
    wm.AoiManager.AddPosToGrid(player.Pid, player.X, player.Z)
    // add to Players 
    wm.rw_lock.Lock()
    wm.Players[player.Pid] = player
    wm.rw_lock.Unlock()
}

// RemovePlayer By Pid
func (wm *WorldManager) RemovePlayer(pid int32) {
    // 1 remove from AoiManager
    wm.rw_lock.RLock()
    player := wm.Players[pid]
    wm.rw_lock.RUnlock()
    wm.AoiManager.RemovePosFromGrid(pid, player.X, player.Z)

    // 2 remove from Players
    wm.rw_lock.Lock()
    delete(wm.Players, pid)
    wm.rw_lock.Unlock()
}


func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
    wm.rw_lock.RLock()
    defer wm.rw_lock.RUnlock()
    return wm.Players[pid]
}

func (wm *WorldManager) GetTotalPlayers() (players []*Player) {
    wm.rw_lock.RLock()
    defer wm.rw_lock.RUnlock()

    players = make([]*Player, 0, len(wm.Players))
    for _, player := range wm.Players {
        players = append(players, player)
    }
    return
}
