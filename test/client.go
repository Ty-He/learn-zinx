package main

import (
	"fmt"
	"io"
	"my_zinx/znet"
	"net"
	"time"
)

func test(buf []byte) {
    dataPack := znet.NewDataPack()
    msg, _ := dataPack.Unpack(buf)
    fmt.Printf("Unpack:msgId=%d, dataLen=%d, data=%s\n", msg.GetId(), msg.GetLen(), msg.GetData())
}

func main() {
    conn, err := net.Dial("tcp", "192.168.18.128:8999")
    if err != nil {
        fmt.Println("Connect Error")
        return 
    }

    for {
        dataPack := znet.NewDataPack()
        buf, err := dataPack.Pack(znet.NewMessage(0, []byte("test_message")))
        if err != nil {
            fmt.Println("Pack Error:", err)
            break
        }
        // test(buf)
        if _, err := conn.Write(buf); err != nil {
            fmt.Println("Write Error,", err)
            break
        }

        // recv response from server
        headbuf := make([]byte, dataPack.GetHeadLen())
        if _, err := io.ReadFull(conn, headbuf); err != nil {
            fmt.Println("ReadFull Error:", err)
            break
        }
        pmsg, err := dataPack.Unpack(headbuf)
        if err != nil {
            fmt.Println("Unpack Error")
            break
        }
        
        if pmsg.GetLen() > 0 {
            buf := make([]byte, pmsg.GetLen())
            if _, err := io.ReadFull(conn, buf); err != nil {
                fmt.Println("ReadFull Error:", err)
                break
            }
            pmsg.SetData(buf)
        }
        // show 
        fmt.Printf("recv server respondse:msgId=%d, dataLen=%d, data=%s\n", pmsg.GetId(), pmsg.GetLen(), pmsg.GetData())
        time.Sleep(1 * time.Second)
    }
}

