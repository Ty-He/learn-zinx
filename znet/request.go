package znet

import "my_zinx/ziface"

type Request struct {
    // socket fd
    conn ziface.IConnection 
    // data
    // data []byte
    // message
    msg ziface.IMessage 
}


func (self *Request) GetConnection() ziface.IConnection {
    return self.conn
}

func (self *Request) GetDate() []byte {
    return self.msg.GetData()
}

func (self *Request) GetMsgId() uint32 {
    return self.msg.GetId()
}
