package ziface

type IDatePack interface {
    // get the length of message head
    GetHeadLen() uint32 
    // Message -> []byte 
    Pack(IMessage) ([]byte, error)
    // []byte -> Message(only get head in this functino)
    Unpack([]byte) (IMessage, error)
}
