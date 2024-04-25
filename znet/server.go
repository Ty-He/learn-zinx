package znet

import (
	"fmt"
	"my_zinx/utils"
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
    // router for handle
    // Router ziface.IRouter
    
    // msgHandler --> f(msgid)-> handle
    MsgHandler ziface.IMsghandle
}

// handapi, later apply by demo
// func callback_client(conn *net.TCPConn, data []byte, n int) error {
//     // echo 
//     fmt.Println("conn handle callback_client")
//     if _, err := conn.Write(data[:n]); err != nil {
//         fmt.Println("Conn Write Error:", err) 
//         return errors.New("callback_client error")
//     }
//     return nil
// }

func NewServer() ziface.IServer {
    s := &Server {
        Name : utils.GlobalObj.Name,
        IpVersion : "tcp",
        Ip : utils.GlobalObj.Host,
        port : utils.GlobalObj.TcpPort,
        MsgHandler : NewMsgHandle(),
    }
    return s
}

func (self *Server) Start() {
    fmt.Printf("Name:%s; IpVersion:%s, Host: %s, Port:%d, MaxConn:%d, MaxPkgSz:%d\n", 
        utils.GlobalObj.Name,
        utils.GlobalObj.Version,
        utils.GlobalObj.Host,
        utils.GlobalObj.TcpPort,
        utils.GlobalObj.MaxConn,
        utils.GlobalObj.MaxPackageSize)

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
    // fmt.Printf("Server start success: ip:%s, port:%d\n", self.Ip, self.port)

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
        dealConn := NewConnection(conn, cid, self.MsgHandler)
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


func (self *Server) AddRouter(msgId uint32, router ziface.IRouter) {
    fmt.Printf("Add a router, msgId = %d\n", msgId)
    self.MsgHandler.AddRouter(msgId, router)
}
