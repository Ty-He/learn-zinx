package znet

type Message struct {
    // the index of this message
    Id uint32 
    // the size of data
    Len uint32
    // data from client
    Data []byte
}

func NewMessage(msgId uint32, data []byte) *Message {
    return &Message {
        Id : msgId,
        Len : uint32(len(data)),
        Data : data,
    }
}


func (this *Message) GetId() uint32 {
    return this.Id
}

func (this *Message) GetLen() uint32 {
    return this.Len
}

func (this *Message) GetData() []byte {
    return this.Data
}


func (this *Message) SetId(id uint32) {
    this.Id = id
}

func (this *Message) SetLen(l uint32) {
    this.Len = l
}

func (this *Message) SetData(d []byte) {
    this.Data = d
}
