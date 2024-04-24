package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)


func TestDataPack(t *testing.T) {
    listener, err := net.Listen("tcp", "192.168.18.128:8999")
    if err != nil {
        fmt.Println("Listen Error:", err)
        return 
    }

    // server
    go func() {
        for {
            conn, err := listener.Accept()
            if err != nil {
                fmt.Println("Accept Error:", err)
                break
            }
            go server_func(&conn)
        }    
    }()

    client_func()

    select {}
}

func server_func(conn *net.Conn) {
    dataPack := NewDataPack()
    for {
        // get data head
        head := make([]byte, dataPack.GetHeadLen())
        _, err := io.ReadFull(*conn, head)
        if err != nil {
            fmt.Println("ReadFull Error:", err)
            break
        }
        headMsg, err := dataPack.Unpack(head)
        if err != nil {
            fmt.Println("Unpack Error:", err)
            return
        }
        // get data 
        if headMsg.GetLen() > 0 {
            buf := make([]byte, headMsg.GetLen())
            _, err := io.ReadFull(*conn, buf)
            if err != nil {
                fmt.Println("ReadFull Error:", err)
            }
            headMsg.SetData(buf)
        }
        

        // show message
        fmt.Printf("Message ID=%d,Len=%d,Data:%s\n", headMsg.GetId(), headMsg.GetLen(), headMsg.GetData())
    }
}

func client_func() {
    // connect
    conn, err := net.Dial("tcp", "192.168.18.128:8999")
    if err != nil {
        fmt.Println("Dial Error:", err)
        return 
    }

    // send two message in the same time
    msg1 := Message {
        Id : 1,
        Len : 4,
        Data : []byte("zinx"),
    }
    msg2 := Message {
        Id : 2,
        Len : 11,
        Data : []byte("hello,zinx!"),
    }

    dataPack := NewDataPack()

    buf1, err := dataPack.Pack(&msg1)
    if err != nil {
        fmt.Println("Pack Error:", err)
        return
    }

    buf2, err := dataPack.Pack(&msg2)
    if err != nil {
        fmt.Println("Pack Error:", err)
        return
    }

    buf1 = append(buf1, buf2...)
    
    if _, err := conn.Write(buf1); err != nil {
        fmt.Println("Write Error,", err)
    }
    
}
