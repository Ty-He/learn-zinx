package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
    conn, err := net.Dial("tcp", "192.168.18.128:8999")
    if err != nil {
        fmt.Println("Connect Error")
        return 
    }

    for {
        _, err := conn.Write([]byte("Hello ZinxV0.2 ..."))
        if err != nil {
            fmt.Println("Write Error")
            continue
        }
        buf := make([]byte, 512)
        n, err := conn.Read(buf)
        if n == 0 {
            fmt.Println("Quit Connection")
            break
        }
        if err != nil {
            fmt.Println("Read Error")
            continue
        }
        fmt.Printf("Read buf:%s, len:%d\n", buf, n)
        time.Sleep(1 * time.Second)
    }
}
