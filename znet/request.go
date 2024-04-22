package znet

import "my_zinx/ziface"

type Request struct {
    // socket fd
    conn ziface.IConnection 
    // data
    data []byte
}


func (self *Request) GetConnection() ziface.IConnection {
    return self.conn
}

func (self *Request) GetDate() []byte {
    return self.data
}
