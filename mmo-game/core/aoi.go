package core 

import (
    "fmt"
)

type AOIManager struct {
    MinX uint 
    MaxX uint 
    CntX uint 
    MinY uint 
    MaxY uint 
    CntY uint 

    grids map[uint]*Grid
}

func NewAOIManager(minX, maxX, cntX, minY, maxY, cntY uint) *AOIManager {
    m := &AOIManager{
        MinX: minX,
        MaxX: maxX,
        CntX: cntX,
        MinY: minY,
        MaxY: maxY,
        CntY: cntY,
        grids: make(map[uint]*Grid),
    }

    // init every Grid 
    for y := uint(0); y < m.CntY; y++ {
        for x := uint(0); x < m.CntX; x++ {
            gid := y * cntX + x 

            m.grids[gid] = NewGrid(gid, 
                m.MinX + x * m.getGridWid(), m.MinX + (x + 1) * m.getGridWid(),
                m.MinY + y * m.getGridLen(), m.MinY + (y + 1) * m.getGridLen())
        }
    }

    return m
}

// get single grid width
func (m *AOIManager) getGridWid() uint {
    return (m.MaxX - m.MinX) / m.CntX 
}

// get single grid height
func (m *AOIManager) getGridLen() uint {
    return (m.MaxY - m.MinY) / m.CntY
}

func (m *AOIManager) String() string {
    s := fmt.Sprintf("AOIManager:\n MinX:%d, MaxX:%d, CntX:%d, MinY:%d, MaxY:%d, CntY:%d\n Grid:\n",
        m.MinX, m.MaxX, m.CntX, m.MinX, m.MinY, m.CntY)

    for _, g := range m.grids {
        s += fmt.Sprintln(g)
    }

    return s
}

// Get surrounding grid(include self) by ID
func (m* AOIManager) GetSurroundGrid(gid uint) (grids []*Grid) {
    if _, ok := m.grids[gid]; !ok {
        return 
    }

    // init 
    grids = make([]*Grid, 0, 9)
    // used for judge have grid on/under g
    tryGetY := func (g uint) {
        y := g / m.CntX
        if y > 0 {
            grids = append(grids, m.grids[g - m.CntX])
        }

        if y < m.CntY - 1 {
            grids = append(grids, m.grids[g + m.CntX])
        }
    }

    // current grid
    grids = append(grids, m.grids[gid])
    tryGetY(gid)
    x := gid % m.CntX
    // left
    if x > 0 {
        grids = append(grids, m.grids[gid - 1])
        tryGetY(gid - 1)
    }

    // right
    if x < m.CntX - 1 {
        grids = append(grids, m.grids[gid + 1])
        tryGetY(gid + 1)
    }

    return
}


// By position index, get Gid
func (m *AOIManager) GetGid(x, y float32) uint {
    idx, idy := uint(x) / m.getGridWid(), uint(y) / m.getGridLen()
    return idy * m.CntX + idx
}

// By position index, get total playerIds which in same area 
func (m *AOIManager) GetPlayerIds(x, y float32) (playerIds []uint) {
    grids := m.GetSurroundGrid(m.GetGid(x, y))
    
    for _, grid := range grids {
        playerIds = append(playerIds, grid.GetPlayIds()...)
    }
    return
}


// By PlayerId, add player to grid
func (m *AOIManager) AddPidToGrid(pid, gid uint) {
    m.grids[gid].Add(pid)
}

// By PlayerId, remove player to grid
func (m *AOIManager) RemovePidFromGrid(pid, gid uint) {
    m.grids[gid].Remove(pid)
}

// get total player from grid 
func (m *AOIManager) GetPidInGrid(gid uint) []uint {
    return m.grids[gid].GetPlayIds()
}

// By position index, add player to grid 
func (m *AOIManager) AddPosToGrid(pid uint, x, y float32) {
    m.grids[m.GetGid(x, y)].Add(pid)
}

// By position index, remove player to grid 
func (m *AOIManager) RemovePosFromGrid(pid uint, x, y float32) {
    m.grids[m.GetGid(x, y)].Remove(pid)
}
