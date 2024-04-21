package znet

import (
	"errors"
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

// handapi, later apply by demo
func callback_client(conn *net.TCPConn, data []byte, n int) error {
    // echo 
    fmt.Println("conn handle callback_client")
    if _, err := conn.Write(data[:n]); err != nil {
        fmt.Println("Conn Write Error:", err) 
        return errors.New("callback_client error")
    }
    return nil
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
    var cid uint32 
    cid = 0
    for {
        conn, err := listener.AcceptTCP()
        if err != nil {
            fmt.Println("ListenTCP Error")
            continue
        }
        
        // bind conn and Connection
        dealConn := NewConnection(conn, cid, callback_client)
        cid ++

        // in case if obsructive current goroutinue
        go dealConn.Start()

        // Handle
        // go func() {
        //     for {
        //         buf := make([]byte, 512)
        //         n, err := conn.Read(buf)
        //         if err != nil {
        //             if n == 0 {
        //                 fmt.Println("Client Quit Connection")
        //                 return
        //             }
        //             fmt.Println("Read Error")
        //             continue
        //         }
        //         // Write back to client
        //         if _, err := conn.Write(buf[:n]); err != nil {
        //             fmt.Println("Write Error")
        //             continue
        //         }
        //     }
        // }()
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


