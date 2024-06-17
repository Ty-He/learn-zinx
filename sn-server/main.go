package main 

import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "my_zinx/znet"
)

const (
    st_err uint32 = iota
    st_ok
    st_logVer
    st_getUser
    st_register
    st_allUser
    st_allPriMsg
    st_sendPriMsg
    st_recvOtherMsg
    st_build
    st_cancel
    st_create
    st_join
    st_leave
    st_getRelat
    st_inGroup
    st_allGrpMsg
    st_sendGrpMsg
    st_recvGrpMsg
)


func main() {
    dsn := "root:245869@tcp(192.168.113.112:3306)/social_net"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        fmt.Println("Open sql err:", err)
        return
    }

    err = db.Ping()
    if err != nil {
        fmt.Println("Ping err:", err)
        return
    }

    fmt.Println("Conn sql success")
    
    s := znet.NewServer()
    sn = NewSnHandler(s, db)

    // add router
    s.AddRouter(st_logVer, &LoginVerify{})
    s.AddRouter(st_getUser, &GetUser{})
    s.AddRouter(st_register, &Register{})
    s.AddRouter(st_allUser, &GetAllUser{})
    s.AddRouter(st_allPriMsg, &GetPriMsg{})
    s.AddRouter(st_sendPriMsg, &dealPriMsg{})
    s.AddRouter(st_build, &BuildRelationship{})
    s.AddRouter(st_cancel, &CancelRelationship{})
    s.AddRouter(st_create, &CreateGroup{})
    s.AddRouter(st_join, &JoinGroup{})
    s.AddRouter(st_leave, &LeaveGroup{})
    s.AddRouter(st_getRelat, &getAllRelationship{})
    s.AddRouter(st_inGroup, &InWhichGruop{})
    s.AddRouter(st_allGrpMsg, &AllGroupMessage{})
    s.AddRouter(st_sendGrpMsg, &dealGrpMsg{})

    // set hook func
    s.SetOnConnStart(OnConStop)

    s.Serve()
}
