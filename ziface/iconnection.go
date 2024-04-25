package ziface

import "net"


type IConnection interface {
    // start a connectino
    Start()
    // stop a connecion
    Stop()
    // get socket obj
    GetTCPConnection() *net.TCPConn
    // get index
    GetConnID() uint32
    // get client information
    RemoteAddr() net.Addr
    // send message
    SendMsg(msgId uint32, data []byte) error 
}

// work functino
type HandleFunc func(*net.TCPConn, []byte, int) error
