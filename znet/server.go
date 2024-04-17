package znet

import (
	"fmt"
	"my_zinx/ziface"
	"net"
)

type Server struct {
    // the name of Server
    Name string 
    // the version of Ip
    IpVersion string
    Ip string
    port int
}

func NewServer(name string) ziface.IServer {
    s := &Server {
        Name : name,
        IpVersion : "tcp4",
        Ip : "192.168.18.128",
        port : 8999,
    }
    return s
}

func (self *Server) Start() {
    // 1. get addr
    addr, err := net.ResolveTCPAddr(self.IpVersion, fmt.Sprintf("%s:%d", self.Ip, self.port))
    if err != nil {
        fmt.Println("ResolveTCPAddr Error")
        return
    }

    // 2. listen
    listener, err := net.ListenTCP(self.IpVersion, addr)
    if err != nil {
        fmt.Println("ListenTCP Error")
        return 
    }
    fmt.Printf("Server start success: ip:%s, port:%d\n", self.Ip, self.port)

    // 3. get cilent connections
    for {
        conn, err := listener.AcceptTCP()
        if err != nil {
            fmt.Println("ListenTCP Error")
            continue
        }

        // Handle
        go func() {
            for {
                buf := make([]byte, 512)
                n, err := conn.Read(buf)
                if err != nil {
                    if n == 0 {
                        fmt.Println("Client Quit Connection")
                        return
                    }
                    fmt.Println("Read Error")
                    continue
                }
                // Write back to client
                if _, err := conn.Write(buf[:n]); err != nil {
                    fmt.Println("Write Error")
                    continue
                }
            }
        }()
    }

}

func (self *Server) Stop() {

}
func (self *Server) Serve() {
    go self.Start()

    // TODO do other somethings


    // clog
    select {}


}
