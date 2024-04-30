package znet

import (
	"errors"
	"fmt"
	"my_zinx/ziface"
	"sync"
)


type ConnManager struct {
    connMap map[uint32] ziface.IConnection 
    mutex sync.RWMutex 
}

func NewConnManager() *ConnManager {
    return &ConnManager {
        connMap : make(map[uint32] ziface.IConnection),
    }
}

func (self *ConnManager) Insert(conn ziface.IConnection) {
    // preserve shared source
    self.mutex.Lock()
    defer self.mutex.Unlock()

    self.connMap[conn.GetConnID()] = conn

    fmt.Println("Insert a conn:", conn.GetConnID(), "Conn total = ", self.Len())
}


func (self *ConnManager) Remove(connId uint32) {
    // preserve shared source
    self.mutex.Lock()
    defer self.mutex.Unlock()

    delete(self.connMap, connId)

    fmt.Println("Remove a conn:", connId, "Conn total = ", self.Len())
}


func (self *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
    // add shared_lock
    self.mutex.RLock()
    defer self.mutex.RUnlock()

    if conn, exist := self.connMap[connId]; exist {
        return conn, nil
    } else {
        return nil, errors.New("Conn not found")
    }
}

func (self *ConnManager) Len() int {
    return len(self.connMap)
}


func (self *ConnManager) Clear() {
    // preserve shared source
    self.mutex.Lock()
    defer self.mutex.Unlock()

    for connId, conn := range self.connMap {
        conn.Stop()

        delete(self.connMap, connId)
    }

    fmt.Println("Clear all conn, conn len is,", self.Len())
}
