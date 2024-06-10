package main 

import (
    "net"
    "fmt"
    "io"
    "bufio"
    // "time"
    "os"
    "bytes"
    "strings"
    "my_zinx/znet"
    "my_zinx/ziface"
)

var path string = "./resource"

// control the server todo
const (
    ft_err uint32 = iota 
    ft_ok 
    ft_file_send// transfer file content
    ft_list
    ft_change_dir
    ft_make_dir
    ft_file_recv
    ft_fill
    ft_finish
    ft_rm
)

func getId(s string) uint32 {
    switch  {
    case s == "ls" || s == "ll" || s == "pwd":
        return ft_list
    case strings.HasPrefix(s, "cd"):
        return ft_change_dir 
    case strings.HasPrefix(s, "get"):
        return ft_file_send
    case strings.HasPrefix(s, "mkdir"):
        return ft_make_dir
    case strings.HasPrefix(s, "put"):
        return ft_file_recv
    case strings.HasPrefix(s, "rm"):
        return ft_rm
    default:
        return ft_err
    }
}

func sendMsg(conn net.Conn, id uint32, cmd string) error {
    dp := znet.NewDataPack()

    // pack 
    buf, err := dp.Pack(znet.NewMessage(id, []byte(cmd)))
    if err != nil {
        return err
    }

    // send data
    if _, err := conn.Write(buf); err != nil {
        return err
    }

    return nil
}
func recvMsg(conn net.Conn) (ziface.IMessage, error) {
    dp := znet.NewDataPack()

    // get message head, it is byte seq
    headbuf := make([]byte, dp.GetHeadLen())
    if _, err := io.ReadFull(conn, headbuf); err != nil {
        fmt.Println("ReadFull Error:", err)
        return nil, err
    }

    // Unpack head
    msg, err := dp.Unpack(headbuf)
    if err != nil {
        fmt.Println("Unpack Error:", err)
        return nil, err
    }

    // read date
    if msg.GetLen() > 0 {
        bodybuf := make([]byte, msg.GetLen())
        if _, err := io.ReadFull(conn, bodybuf); err != nil  {
            fmt.Println("Unpack Error:", err)
            return nil, err
        }
        msg.SetData(bodybuf)
    }

    return msg, nil
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    conn, err := net.Dial("tcp", "192.168.18.128:8999")
    if err != nil {
        fmt.Println("Dial Error: ", err)
        return
    }
    defer conn.Close()

    // recv the response of the server 
    if msg, err := recvMsg(conn); err !=nil {
        fmt.Println("recvMsg err:", err)
        return
    } else {
        fmt.Printf("[Note:] recv id [%d], datelen [%d] date:\n%s\n", msg.GetId(), msg.GetLen(), msg.GetData())
    }
    for {
        // input a command
        var cmd string
        if scanner.Scan() {
            cmd = scanner.Text()
            if cmd == "bye" {
                break
            }
        }

        id := getId(cmd)
        if id == ft_err {
            fmt.Println("Not a command!")
            continue
        } 
        
        // snedMsg
        if err := sendMsg(conn, id, cmd); err != nil {
            fmt.Println("sendMsg err:", err)
            return
        }

        // recvMsg
        msg, err := recvMsg(conn)
        if err != nil {
            return
        }


        if msg.GetId() != ft_ok && msg.GetId() != ft_err {
            name := path + "/" +cmd[4:]
            if msg.GetId() == ft_file_send { // download file from server
                os.Remove(name)
            }
            // note: append file
            file, _ := os.OpenFile(name, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0644)
            reader := bufio.NewReader(file)
            defer file.Close()
            for msg != nil && msg.GetId() != ft_ok && msg.GetId() != ft_err {
                // wait all content 
                if msg.GetId() == ft_file_send {
                    file.Write(msg.GetData())
                } else if msg.GetId() == ft_fill { // TODO wite to server
                    buf := bytes.NewBuffer([]byte{})
                    for {
                        line, err := reader.ReadString('\n')
                        if err != nil {
                            if err.Error() == "EOF" {
                                break
                            } else {
                                fmt.Println("ReadString Error:", err)
                                return
                            }
                        }

                        buf.WriteString(line)
                        if buf.Len() > 99 * 1024 { // maxCache
                            break
                        }
                    }
                    if buf.Len() == 0 {
                        fmt.Println("finish")
                        sendMsg(conn, ft_finish, "")
                    } else {
                        fmt.Println("fill")
                        sendMsg(conn, ft_fill, buf.String())
                    }
                }
                msg, err = recvMsg(conn)
                if err != nil {
                    return
                }
            }
        }


        // msg.id in { ok, err }
        fmt.Printf("[Note:] recv id [%d], datelen [%d] date:\n%s\n", msg.GetId(), msg.GetLen(), msg.GetData())
        
        // scan err
        if err := scanner.Err(); err != nil {
            fmt.Println("scanner err:", err)
            return
        }

    }

    fmt.Println("[Note:] disconnect the server.")
}
