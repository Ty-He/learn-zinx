package core 

import (
    "fmt"
    "testing"
)

func TestNewAOIManager(t *testing.T) {
    m := NewAOIManager(0, 250, 5, 0, 250, 5)

    // fmt.Println(m)
    for gid := range m.grids {
        grids := m.GetSurroundGrid(gid)
        fmt.Printf("Get Grid[%d] surroundings Grid, len = %d\n", gid, len(grids))
        ids := make([]uint, 0, len(grids))
        for _, g := range grids {
            ids = append(ids, g.Gid)
        }
        fmt.Println("surroundings grids: ", ids)
    }
}
