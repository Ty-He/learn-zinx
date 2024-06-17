package main 


import (
    "sync"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "my_zinx/ziface"
)

// handler (tcp and sql) 
type SnHandler struct {
    s ziface.IServer 
    db *sql.DB

    // uId --> connId
    onlineMap map[uint32]uint32
    rwlock sync.RWMutex
}

// global variable
var sn *SnHandler 

func NewSnHandler(s ziface.IServer, db *sql.DB) *SnHandler {
    sn := &SnHandler{
        s : s,
        db : db,
        onlineMap: make(map[uint32]uint32),
    }
    return sn
}

