package znet

import (
	"errors"
	"fmt"
	"io"
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
    // handleAPI ziface.HandleFunc 
    // router
    Router ziface.IRouter
    // signal of exit
    ExitChan chan bool
}


func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
    obj := &Connection {
        Conn : conn,
        ConnID : connID,
        isClosed : false,
        Router : router,
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
        // buf := make([]byte, utils.GlobalObj.MaxPackageSize)
        // n, err := self.Conn.Read(buf)
        // if err != nil {
        //     if n == 0 {
        //         fmt.Printf("client quit connection, id:%d\n", self.ConnID)
        //         self.Stop()
        //         return
        //     }
        //     fmt.Println("Connection Read Error")
        //     continue
        // }

        // get head (byte) 
        dataPack := NewDataPack()
        head := make([]byte, dataPack.GetHeadLen())
        if _, err := io.ReadFull(self.Conn, head); err != nil {
            if err.Error() == "EOF" {
                fmt.Printf("Client[RemoteAddr-%s] quit connection\n", self.RemoteAddr().String())
            } else {
                fmt.Println("ReadFull Error:", err)
            }
            break
        }

        // Unpack: get head(data)
        msg, err := dataPack.Unpack(head)
        // fmt.Printf("head:%s; Unpack: Id=%d, Len=%d, Data=%s\n", head, msg.GetId(), msg.GetLen(), msg.GetData())
        if err != nil {
            fmt.Println("Unpack Error:", err)
            break
        }

        if msg.GetLen() > 0 {
            // var data []byte
            data := make([]byte, msg.GetLen())
            if _, err := io.ReadFull(self.Conn, data); err != nil {
                fmt.Println("ReadFull Error:", err)
                break
            }
            msg.SetData(data)
        }
        

        // handle the data from client
        // if err := self.handleAPI(self.Conn, buf, n); err != nil {
        //     fmt.Println("handleAPI error", err)
        // }
        
        // constuct request
        req := Request {
            conn : self,
            msg : msg,
        }
        // router handle
        go func(request ziface.IRequest) {
            self.Router.PreHandle(request)
            self.Router.Handle(request)
            self.Router.PostHandle(request)
        }(&req)
    }
}

func (self *Connection) Start() {
    fmt.Printf("Connection Start-> ConnID:%d, RemoteAddr:%s\n", self.ConnID, self.RemoteAddr().String())
    // read goroutine
    go self.start_reader()
    // TODO write goroutine
}


func (self *Connection) Stop() {
    if self.isClosed {
        return 
    }
    fmt.Printf("Connection Stop, ID:%d\n", self.ConnID)
    self.isClosed = true
    // close socket fd
    self.Conn.Close()
    // close channel
    close(self.ExitChan)
}


func (self *Connection) GetTCPConnection() *net.TCPConn{
    return self.Conn
}


func (self *Connection) GetConnID() uint32 {
    return self.ConnID
}


func (self *Connection) RemoteAddr() net.Addr{
    return self.Conn.RemoteAddr()
}

func (self *Connection) SendMsg(msgId uint32, data []byte) error {
    if self.isClosed {
        return errors.New("Connection is already quit")
    }
    // data: byte stream --> binary steram
    dataPack := NewDataPack()
    buf, err := dataPack.Pack(NewMessage(msgId, data))
    if err != nil {
        return err
    }
    
    // send stream
    if _, err := self.Conn.Write(buf); err != nil {
        return err
    }
    return nil
}
