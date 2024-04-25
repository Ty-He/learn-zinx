package ziface


type IMessage interface {
    // getter
    GetId() uint32 
    GetLen() uint32 
    GetData() []byte

    // setter
    SetId(uint32) 
    SetLen(uint32) 
    SetData([]byte) 
}
