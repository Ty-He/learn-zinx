package main

import (
	"fmt"
	"my_zinx/ziface"
	"my_zinx/znet"
    "my_zinx/utils"
    "os"
    "bufio"
    "strings"
    "bytes"
    "syscall"
)

// when have some error in transfer, sned id == 0 msg
// when have some success in transfer, sned id == 1 msg

// enum of msgID
const (
    ft_err uint32 = iota 
    ft_ok 
    ft_file_send // transfer file content
    ft_list
    ft_change_dir
    ft_make_dir
    ft_file_recv
    ft_fill
    ft_finish
    // TODO 
    ft_rm
)

var homeDir string = "./resource"

type SendFileRouter struct {
    znet.BaseRouter
}

func (sf *SendFileRouter) Handle(request ziface.IRequest) {
    path, _ := request.GetConnection().GetProperty("path")
    target := path.(string) + "/" + string(request.GetDate()[4:len(request.GetDate())])
    file, err := os.Open(string(target))
    if err != nil {
        fmt.Println("Open Error")
        // TODO send error to client
        request.GetConnection().SendMsg(ft_err, []byte("Error file path!!!"))
        return
    }
    defer file.Close()

    // protect shared resource
    syscall.Flock(int(file.Fd()), syscall.LOCK_SH)
    defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

    fi, _ := file.Stat()
    if fi.Size() < int64(utils.GlobalObj.MaxPackageSize) {
        err := totalTransfer(file, request)
        if err != nil {
            fmt.Println("totalTransfer err:", err)
            request.GetConnection().SendMsg(ft_err, []byte("transfer error"))
            return
        }
    } else {
        fmt.Println(target)
        err := blockTransfer(file, request)
        if err != nil {
            fmt.Println("totalTransfer err:", err)
            request.GetConnection().SendMsg(ft_err, []byte("transfer error"))
            return
        }
    }

    // send finish
    if err := request.GetConnection().SendMsg(ft_ok, []byte("server sendFile ok!")); err != nil { 
        fmt.Println("SendMsg Error: ", err)
        return
    }
    // fmt.Println(buf.String())
}

type ListFileRouter struct {
    znet.BaseRouter
}

func (lf *ListFileRouter) Handle(request ziface.IRequest) {
    path, _ := request.GetConnection().GetProperty("path")
    dir, err := os.Open(path.(string))
    if err != nil {
        // server error
        fmt.Println("OpenDir Error:", err)
        return
    }
    defer dir.Close()

    // protect shared resource
    syscall.Flock(int(dir.Fd()), syscall.LOCK_SH)
    defer syscall.Flock(int(dir.Fd()), syscall.LOCK_UN)
    
    reqStr := string(request.GetDate())  
    if reqStr == "ls" {
        // read only file name
        names, err := dir.Readdirnames(0)
        if err != nil {
            fmt.Println("Readdirnames Error:", err)
            return
        }
        total := strings.Join(names, "   ")
        if err := request.GetConnection().SendMsg(ft_ok, []byte(total)); err != nil {
            fmt.Println("SendMsg Error: ", err)
        }
    } else if reqStr == "ll" { 
        entrys, err := dir.ReadDir(0)
        if err != nil {
            fmt.Println("ReadDir Error ", err)
            return
        }
        buf := bytes.NewBuffer([]byte{})
        for i := range entrys {
            info, err := entrys[i].Info()
            if err != nil {
                break
            }
            s := fmt.Sprintf("%s %s %s %d\n", 
                info.Mode().String(),
                info.ModTime().String(),
                info.Name(),
                info.Size(),
                )
            buf.WriteString(s)
        }

        data := buf.Bytes()
        if err := request.GetConnection().SendMsg(ft_ok, data); err != nil {
            fmt.Println("SendMsg Error: ", err)
        }
    } else if reqStr == "pwd" {
        request.GetConnection().SendMsg(ft_ok, []byte(path.(string)))
    } else {
        request.GetConnection().SendMsg(ft_err, []byte("bad command!"))
    }
}


type ChangeDirRouter struct {
    znet.BaseRouter
} 

func (cd *ChangeDirRouter) Handle(request ziface.IRequest) {
    // The dest is exist or not ?
    path, _ := request.GetConnection().GetProperty("path")
    var target string
    if len(request.GetDate()) > 2 {
        target = string(request.GetDate()[3:])
    }
    if len(target) == 0 {
        // return home dir
        request.GetConnection().SetProperty("path", homeDir)
    } else if target == ".." { // go father dir
        // auto realize relative path, need file-system
        pa := path.(string)
        if pa == homeDir {
            request.GetConnection().SendMsg(ft_err, []byte("bad directory!"))
            return
        } else {
            for i := len(pa) - 1; i >= 0; i-- {
                if pa[i] == '/' {
                    request.GetConnection().SetProperty("path", pa[:i])
                    break
                }
            }
        }
    } else {
        dest := path.(string) + "/" + target
        fi, err := os.Stat(dest)
        if err != nil || !fi.IsDir() {
            request.GetConnection().SendMsg(ft_err, []byte("bad directory!"))
            return
        }
        // it's a exist dir
        // update path
        request.GetConnection().SetProperty("path", dest)
    }
    request.GetConnection().SendMsg(ft_ok, []byte("go this dir ok!"))
}

type MakeDirRouter struct {
    znet.BaseRouter 
}

func (md *MakeDirRouter) Handle(request ziface.IRequest) {
    path, _ := request.GetConnection().GetProperty("path")
    pa := path.(string) 
    
    var dirname string
    if len(request.GetDate()) > 6 {
        dirname = string(request.GetDate()[6:])
    }
    if err := os.Mkdir(pa + "/" + dirname, 0744); err != nil {
        request.GetConnection().SendMsg(ft_err, []byte("bad directory!")) 
    } else {
        request.GetConnection().SendMsg(ft_ok, []byte("mkdir ok!"))
    } 
} 

type RecvFileRouter struct {
    znet.BaseRouter 
}

func (rf *RecvFileRouter) Handle(request ziface.IRequest) {
    path, _ := request.GetConnection().GetProperty("path")
    pa := path.(string) 
    var name string 
    if len(request.GetDate()) > 4 {
        name = string(request.GetDate()[4:])
    }

    // for create or trunc file
    file, err := os.OpenFile(pa + "/" + name, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0644)

    if err != nil {
        request.GetConnection().SendMsg(ft_err, []byte("bad filename! " + pa + "/" + name)) 
    } else if fi, _:= file.Stat(); fi.IsDir() {
        request.GetConnection().SendMsg(ft_err, []byte("try to write directory!")) 
    }  else {
        fmt.Println("creat a file:", name)
        request.GetConnection().SendMsg(ft_fill, []byte{})
        request.GetConnection().SetProperty("filename", pa + "/" + name)
    }
    file.Close()
} 


type FillFileRouter struct {
    znet.BaseRouter
}

func (ff * FillFileRouter) Handle(request ziface.IRequest) {
    fn, err := request.GetConnection().GetProperty("filename")
    if err != nil {
        fmt.Println("GetProperty Error:", err)
    } else {
        filename := fn.(string)
        file, _ := os.OpenFile(filename, os.O_APPEND | os.O_WRONLY, 0644)
        defer file.Close()

        // protect shared resource
        syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
        defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
        
        if _, err := file.Write(request.GetDate()); err != nil {
            request.GetConnection().RemoveProperty("filename")
            request.GetConnection().SendMsg(ft_err, []byte("put file err!"))
        } else {
            // continue to request
            request.GetConnection().SendMsg(ft_fill, []byte{})
        }
    }
}

type FinishFileRouter struct {
    znet.BaseRouter 
} 

func (ff *FinishFileRouter) Handle(request ziface.IRequest) {
    // delete key 
    request.GetConnection().RemoveProperty("filename")
    request.GetConnection().SendMsg(ft_ok, []byte("put file ok."))
}


type RemoveFileRouter struct {
    znet.BaseRouter
}


func (rm *RemoveFileRouter) Handle(request ziface.IRequest) {
    if len(request.GetDate()) < 4 || request.GetDate()[4] == ' ' {
        request.GetConnection().SendMsg(ft_err, []byte("bad filename!"))
    } else {
        // file name
        name := string(request.GetDate()[3:])
        pa, _ := request.GetConnection().GetProperty("path")
        path := pa.(string)
        if err := os.RemoveAll(path + "/" + name); err != nil {
            request.GetConnection().SendMsg(ft_err, []byte("remove err"))
        } else {
            request.GetConnection().SendMsg(ft_ok, []byte("remove file" + name))
        }
    }
}
// Hook function

func OnConnStart(conn ziface.IConnection) {
    fmt.Println("===>New a connecion.")
    if err := conn.SendMsg(ft_ok, []byte("Success connect server.")); err != nil {
        fmt.Println(err)
    }

    // define property
    conn.SetProperty("name", "Ty")
    conn.SetProperty("github", "github.com/Ty-He")
    conn.SetProperty("location", "Nanchang")

    // current client path
    conn.SetProperty("path", homeDir)
    // the file which be filled
    // conn.SetProperty("filename", "null")
}

func OnConnStop(conn ziface.IConnection) {
    fmt.Printf("===>Connection [%d] is lost.\n", conn.GetConnID())

    // read property
    if name, err := conn.GetProperty("name"); err == nil {
        fmt.Println("name:", name)
    }
    if github, err := conn.GetProperty("github"); err == nil {
        fmt.Println("github:", github)
    }
    if location, err := conn.GetProperty("location"); err == nil {
        fmt.Println("location:", location)
    }
}
//

func main() {
    // server 
    s := znet.NewServer()

    // register hook func
    s.SetOnConnStart(OnConnStart)
    s.SetOnConnStop(OnConnStop)

    // add pingrouter
    s.AddRouter(ft_file_send, &SendFileRouter{})
    s.AddRouter(ft_list, &ListFileRouter{})
    s.AddRouter(ft_change_dir, &ChangeDirRouter{})
    s.AddRouter(ft_make_dir, &MakeDirRouter{})
    s.AddRouter(ft_file_recv, &RecvFileRouter{})
    s.AddRouter(ft_fill, &FillFileRouter{})
    s.AddRouter(ft_finish, &FinishFileRouter{})
    s.AddRouter(ft_rm, &RemoveFileRouter{})
    
    // run serve
    s.Serve()
}

func totalTransfer(file *os.File, request ziface.IRequest) error {
    reader := bufio.NewReader(file)
    buf := bytes.NewBuffer([]byte{})
    for {
        // call once only read 1 line
        line, err := reader.ReadString('\n')
        if err != nil && err.Error() != "EOF" {
            return err
        }
        // fmt.Println("[Readfile:]", line)
        buf.WriteString(line)
        if err != nil && err.Error() == "EOF" {
            break
        }
    }
    if err := request.GetConnection().SendMsg(ft_file_send, buf.Bytes()); err != nil {
        return err
    }
    return nil
}

func blockTransfer(file *os.File, request ziface.IRequest) error {
    reader := bufio.NewReader(file)
    buf := bytes.NewBuffer([]byte{})
    for {
        // call once only read 1 line
        line, err := reader.ReadString('\n')
        if err != nil && err.Error() != "EOF" {
            return err
        }
        // fmt.Println("[Readfile:]", line)
        buf.WriteString(line)
        if buf.Len() > 95 * 1024 {
            if err := request.GetConnection().SendMsg(ft_file_send, buf.Bytes()); err != nil {
                return err
            }
            buf = bytes.NewBuffer([]byte{})
        }
        if err != nil && err.Error() == "EOF" {
            if (buf.Len() > 0) {
                if err := request.GetConnection().SendMsg(ft_file_send, buf.Bytes()); err != nil {
                    return err
                }
            }
            break
        }
    }
    return nil
}
