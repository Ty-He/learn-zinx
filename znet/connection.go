package znet

import (
	"fmt"
	"my_zinx/ziface"
	"net"
)

type Connection struct {
    // tcp-socket
    Conn *net.TCPConn
    // index
    ConnID uint32
    // the status of current connection
    isClosed bool
    // worker-function
    handleAPI ziface.HandleFunc 
    // signal of exit
    ExitChan chan bool
}


func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandleFunc) *Connection {
    obj := &Connection {
        Conn : conn,
        ConnID : connID,
        isClosed : false,
        handleAPI : callback_api,
        ExitChan : make(chan bool, 1),
    }
    return obj
}

func (self *Connection) start_reader() {
    fmt.Println("Reader Goroutinue is runnning ...")
    defer fmt.Printf("Read exit, ConnID:%d, RemoteAddr:%s\n", self.ConnID, self.RemoteAddr().String())
    // if read goroutine is finished, then all conn could stop
    defer self.Stop()

    // read from client 
    for {
        buf := make([]byte, 512)
        n, err := self.Conn.Read(buf)
        if err != nil {
            if n == 0 {
                fmt.Printf("client quit connection, id:%d\n", self.ConnID)
                self.Stop()
                return
            }
            fmt.Println("Connection Read Error")
            continue
        }
        // handle the data from client
        if err := self.handleAPI(self.Conn, buf, n); err != nil {
            fmt.Println("handleAPI error", err)
        }
    }
}

func (self *Connection) Start() {
    fmt.Printf("Connection Start, ID:%d\n", self.ConnID)
    // read goroutine
    go self.start_reader()
    // TODO write goroutine
}


func (self *Connection) Stop() {
    fmt.Printf("Connection Stop, ID:%d\n", self.ConnID)
    if self.isClosed {
        return 
    }
    self.isClosed = true
    // close socket fd
    self.Conn.Close()
    // close channel
    close(self.ExitChan)
}


func (self *Connection) GetTCPConnection() *net.TCPConn{
    return self.Conn
}


func (self *Connection) getConnID() uint32 {
    return self.ConnID
}


func (self *Connection) RemoteAddr() net.Addr{
    return self.Conn.RemoteAddr()
}

func (self *Connection) send(data []byte) error {
    return nil
}
