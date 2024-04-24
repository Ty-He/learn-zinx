package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"my_zinx/utils"
	"my_zinx/ziface"
)


type DataPack struct {}

func NewDataPack() *DataPack {
    return &DataPack{}
}

func (this *DataPack) GetHeadLen() uint32 {
    // datalen(4 byte) + messageId(4 byte)
    return 8
}


// Len | Id | Data
func (this *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
    // construct a Buffer
    buf := bytes.NewBuffer([]byte{})
    // write DateLen
    if err := binary.Write(buf, binary.LittleEndian, msg.GetLen()); err != nil {
        return nil, err
    }
    // write Id
    if err := binary.Write(buf, binary.LittleEndian, msg.GetId()); err != nil {
        return nil, err
    }
    // write Date
    if err := binary.Write(buf, binary.LittleEndian, msg.GetData()); err != nil {
        return nil, err 
    }

    return buf.Bytes(), nil
}


func (this *DataPack) Unpack(binData []byte) (ziface.IMessage, error) {
    msg := &Message{}
    // construct ioReader
    reader := bytes.NewReader(binData) 

    // only get the head of msg
    if err := binary.Read(reader, binary.LittleEndian, &msg.Len); err != nil {
        return nil, err 
    }
    // get id 
    if err := binary.Read(reader, binary.LittleEndian, &msg.Id); err != nil {
        return nil, err 
    }

    // judge the size of data package 
    if msg.Len > utils.GlobalObj.MaxPackageSize {
        return nil, errors.New("too large data package")
    }

    // msg.Date == nil 
    return msg, nil 
} 

