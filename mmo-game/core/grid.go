package core 

import (
    "sync"
    "fmt"
)

type Grid struct {
    // ensure a grid
    Gid uint 
    // border
    MinX uint
    MaxX uint 
    MinY uint 
    MaxY uint 
    // all player collection
    playerIds map[uint]struct{}
    pIdLock sync.RWMutex 
}

func NewGrid(gid, minX, maxX, minY, maxY uint) *Grid {
    return &Grid{
        Gid: gid,
        MinX: minX,
        MaxX: maxX,
        MinY: minY,
        MaxY: maxY,
        playerIds: make(map[uint]struct{}),
    }
}

// add a playerId
func (g *Grid) Add(playerId uint) {
    g.pIdLock.Lock()
    defer g.pIdLock.Unlock()

    g.playerIds[playerId] = struct{}{}
}

// Remove a playerId
func (g *Grid) Remove(playerId uint) {
    g.pIdLock.Lock()
    defer g.pIdLock.Unlock()
    
    delete(g.playerIds, playerId)
}

// get all playerIds
func (g *Grid) GetPlayIds() []uint {
    playerIds := make([]uint, len(g.playerIds))
    var idx int 
    g.pIdLock.RLock()
    defer g.pIdLock.RUnlock()

    for playerId := range g.playerIds {
        playerIds[idx] = playerId
        idx ++
    }

    return playerIds
}

// for fmt.Println
func (g *Grid) String() string {
    return fmt.Sprintf("Gid: %d, MinX: %d, MaxX: %d, MinY: %d, MaxY: %d,playerIds: %v\n",
        g.Gid, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIds)
}
