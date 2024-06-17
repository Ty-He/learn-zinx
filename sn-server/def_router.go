package main 

import (
    "fmt"
    "time"
    "bytes"
    "strings"
    "strconv"
    "database/sql"
    "my_zinx/ziface"
    "my_zinx/znet"
)

// post hook func 
func OnConStop(conn ziface.IConnection) {
    cur := conn.GetConnID()
    sn.rwlock.Lock()
    for uId, connId := range sn.onlineMap {
        if cur == connId {
            delete(sn.onlineMap, uId)
            fmt.Printf("uid[%d] insert into onlineMap\n", uId)
            break
        }
    }
    
    sn.rwlock.Unlock()
} 

type user struct {
    uid string
    name string
    email string
    intro string
}

func (u *user) to_string() string {
    return u.uid + "|" + u.name + "|" + u.email + "|" + u.intro
}

// user authentication
type LoginVerify struct {
    znet.BaseRouter
}

func (lv *LoginVerify) Handle(request ziface.IRequest) {
    strs := strings.Split(string(request.GetDate()), " ")
    u, err := strconv.ParseUint(strs[0], 10, 0)
    if err != nil {
        request.GetConnection().SendMsg(st_err, []byte("information format error"))
        return
    }
    uId, passwd := uint32(u), strs[1]

    query := "select uId from UserInfo where uId = ? and uPasswd = ?;"
    var tmp_id uint32
    // fmt.Printf("uId = %d, uPasswd = %s\n", uId, passwd);
    err = sn.db.QueryRow(query, uId, passwd).Scan(&tmp_id)
    if err != nil {
        // include not found
        request.GetConnection().SendMsg(st_err, []byte("uId or passwd error"))
        return
    }

    // login success
    sn.rwlock.Lock()
    sn.onlineMap[uId] = request.GetConnection().GetConnID()
    sn.rwlock.Unlock()

    fmt.Printf("uid[%d] insert into onlineMap\n", uId)
    request.GetConnection().SendMsg(st_ok, []byte{})
}



// get a user information
type GetUser struct {
    znet.BaseRouter
}

func (gu *GetUser) Handle(request ziface.IRequest) {
    u, err := strconv.ParseUint(string(request.GetDate()), 10, 0)
    if err != nil {
        request.GetConnection().SendMsg(st_err, []byte("uid format error"))
        return
    }
    uId := uint32(u)
    query := "select uId, uName, uEmail, uIntro from UserInfo where uId = ?;"
    var uid uint32
    var name string
    var email string 
    var intro string
    err = sn.db.QueryRow(query, &uId).Scan(&uid, &name, &email, &intro)
    if err != nil {
        request.GetConnection().SendMsg(st_err, []byte("GetUser err"))
        return
    }

    request.GetConnection().SendMsg(st_getUser, []byte(string(request.GetDate()) + "|" + name + "|" + email + "|" + intro))
}



// register a new user
type Register struct {
    znet.BaseRouter 
}

func (r *Register) Handle(request ziface.IRequest) {
    strs := strings.Split(string(request.GetDate()), "|")
    if len(strs) != 5 {
        request.GetConnection().SendMsg(st_err, []byte("bad register information(len)"))
        return
    }

    query := "insert into UserInfo (uId, uName, uEmail, uIntro, uPasswd) values (?, ?, ?, ?, ?);"
    u, err := strconv.ParseUint(strs[0], 10, 0)
    if err != nil {
        request.GetConnection().SendMsg(st_err, []byte("bad register information(id)"))
        return
    }
    _, err = sn.db.Exec(query, uint32(u), strs[1], strs[2], strs[3], strs[4])
    if err != nil {
        request.GetConnection().SendMsg(st_err, []byte("Exec err:" + err.Error()))
        return
    }

    request.GetConnection().SendMsg(st_ok, []byte{})
    fmt.Printf("%s insert to UserInfo\n", request.GetDate())
}


// get all user althougu not online
type GetAllUser struct {
    znet.BaseRouter 
}

func (au *GetAllUser) Handle(request ziface.IRequest) {
    // mysql prepare 
    query := `select uId, uName from UserInfo;`
    rows, err := sn.db.Query(query)
    if err != nil {
        request.GetConnection().SendMsg(st_err, []byte("query error"))
        return
    }

    defer rows.Close()

    var total string
    for rows.Next() {
        var uid string 
        var name string
        err = rows.Scan(&uid, &name)
        if err != nil {
            request.GetConnection().SendMsg(st_err, []byte("query error"))
            return
        }
        total += (uid + "|" + name + "#")
    }

    fmt.Println("total user: ", total)

    request.GetConnection().SendMsg(st_ok, []byte(total))
}




// get all private message in uid and tar
type GetPriMsg struct {
    znet.BaseRouter
}

func (gpm *GetPriMsg) Handle(request ziface.IRequest) {
    s := strings.Split(string(request.GetDate()), "|")
    if len(s) != 2 {
        fmt.Println("Split err")
        request.GetConnection().SendMsg(st_err, []byte("uid err"))
        return
    }
    fmt.Println("mSender:", s[0], " mRecver:", s[1])
    
    query := `select mSender, mContent from PrivateMessage 
    where (mSender = ? and mRecver = ?) or (mSender = ? and mRecver = ?)
    order by mTime;`

    rows, err := sn.db.Query(query, s[0], s[1], s[1], s[0])
    if err != nil {
        request.GetConnection().SendMsg(st_err, []byte("query err"))
        fmt.Println("query err: ", err)
        return
    }
    defer rows.Close()

    var buf bytes.Buffer
    for rows.Next() {
        var who, what string 
        rows.Scan(&who, &what)
        buf.WriteString(who)
        buf.WriteString("|")
        buf.WriteString(what)
        buf.WriteString("#")
    }
    fmt.Println(buf.String())
    request.GetConnection().SendMsg(st_ok, buf.Bytes())
}


// someone send a private message
type dealPriMsg struct {
    znet.BaseRouter 
}

func (dpm *dealPriMsg) Handle(request ziface.IRequest) {
    str := strings.Split(string(request.GetDate()), "|")
    if (len(str) != 3) {
        fmt.Println("dealPriMsg Split err")
        request.GetConnection().SendMsg(st_err, []byte{})
        return
    }
    // get maxId then add self
    query := `select max(mId) from PrivateMessage;`
    var pre uint32 
    sn.db.QueryRow(query).Scan(&pre)

    // get current time
    currentTime := time.Now().Format("2006-01-02 15:04:05")

    query = `insert into PrivateMessage (mId, mContent, mTime, mSender, mRecver)
    values (?, ?, ?, ?, ?);`

    _, err := sn.db.Exec(query, pre + 1, str[2], currentTime, str[0], str[1])
    if err != nil {
        fmt.Println("insert pri msg err: ", err)
        request.GetConnection().SendMsg(st_err, []byte{})
        return
    }

    fmt.Println("insert into a new private message")
    // from onlineMap find target, if is online, sned to it
    u, err := strconv.ParseUint(str[1], 10, 0)
    if err != nil {
        fmt.Println("dealPriMsg ParseUint err, ", err)
        return
    }
    target := uint32(u)
    sn.rwlock.RLock()
    defer sn.rwlock.RUnlock()
    for uId, connId := range sn.onlineMap {
        if uId == target {
            conn, err := sn.s.GetConnManager().Get(connId) 
            if err == nil {
                conn.SendMsg(st_recvOtherMsg, request.GetDate())
                fmt.Println("message send ok")
            }
        }
    }
}


// 
type BuildRelationship struct {
    znet.BaseRouter 
}

func (br *BuildRelationship) Handle(request ziface.IRequest) {
    uid := string(request.GetDate())
    // all select no need transaction

    // search last private message
    query := `select Y.mContent from PrivateMessage Y
        where Y.mSender = ? and not exists (
                select * from PrivateMessage X
                where X.mSender = Y.mSender and X.mTime > Y.mTime
            );`
    var str string 
    sn.db.QueryRow(query, uid).Scan(&str)
    strs := strings.Split(str, "%")
    if len(strs) != 3 || strs[0] != "[rBuild]" {
        // fmt.Println(r.Err().Error())
        request.GetConnection().SendMsg(st_err, []byte("You do not apply!"))
        return
    }
    
    // fmt.Println(str)
    // search target is agree or not
    var s string
    sn.db.QueryRow(query, strs[1]).Scan(&s) // note str ir strs
    tar := strings.Split(s, "%")
    if len(tar) != 3 || tar[0] != "[rBuild]" || tar[1] != uid || tar[2] != strs[2] {
        // fmt.Println(r.Err().Error())
        fmt.Println(s)
        request.GetConnection().SendMsg(st_err, []byte("target do not agree!"))
        return
    }
    
    // is have same relationship
    query = `select rId from Relationship 
        where rWhat = ? and ((rUser0 = ? and rUser1 = ?) or (rUser0 = ? or rUser1 = ?));`
    err := sn.db.QueryRow(query, strs[2], uid, str[1], str[1], uid).Scan()
    
    if err != nil && err.Error() == sql.ErrNoRows.Error() { // empty so insert
        query = `select max(rId) from Relationship;`
        var pre uint32 
        sn.db.QueryRow(query).Scan(&pre)
        pre ++
        query = `insert into Relationship (rId, rWhat, rUser0, rUser1) values
            (?, ?, ?, ?);`
        _, err := sn.db.Exec(query, pre, strs[2], uid, strs[1])
        if err != nil {
            fmt.Println("db Exec err:", err)
            request.GetConnection().SendMsg(st_err, []byte("Exec err"))
            return 
        }
        request.GetConnection().SendMsg(st_ok, []byte("You build relationship successfully."))
        return
    }

    if err != nil {
        fmt.Println("db Exec err:", err)
        request.GetConnection().SendMsg(st_err, []byte("select err"))
        return 
    }

    request.GetConnection().SendMsg(st_err, []byte("You have built this relationship."))
}



type CancelRelationship struct {
    znet.BaseRouter
}

func (cr *CancelRelationship) Handle(request ziface.IRequest) {
    uid := string(request.GetDate())
    // all select no need transaction
    query := `select Y.mContent from PrivateMessage Y
        where Y.mSender = ? and not exists (
                select * from PrivateMessage X
                where X.mSender = Y.mSender and X.mTime > Y.mTime
            );`
    var str string 
    sn.db.QueryRow(query, uid).Scan(&str)
    s := strings.Split(str, "%")
    if len(s) != 3 || s[0] != "[rCancel]" {
        // fmt.Println(r.Err().Error())
        request.GetConnection().SendMsg(st_err, []byte("You do not apply!"))
        return
    }
    
    // uid | str[1]
    query = `delete from Relationship 
        where rWhat = ? and ((rUser0 = ? and rUser1 = ?) or (rUser0 = ? and rUser1 = ?));`

    ret, err := sn.db.Exec(query, s[2], uid, s[1], s[1], uid)

    if err != nil {
        fmt.Println("CancelRelationship Exec err:",err)
        request.GetConnection().SendMsg(st_err, []byte{})
        return
    }
    if t, _ := ret.RowsAffected(); t == 0 {
        request.GetConnection().SendMsg(st_err, []byte("dont hava this relationship"))
        return
    }

    request.GetConnection().SendMsg(st_ok, []byte("CancelRelationship successfully."))
}


type CreateGroup struct {
    znet.BaseRouter 
}

func (cg *CreateGroup) Handle(request ziface.IRequest) {
    // gTopId | gName | gIntro
    s := strings.Split(string(request.GetDate()), "|")
    if len(s) != 3 {
        request.GetConnection().SendMsg(st_err, []byte("format error"))
        return 
    }

    // get current gId
    query := `select max(gId) from Grp;`
    var pre uint32
    sn.db.QueryRow(query).Scan(&pre)
    pre ++

    // Exec insert
    // when insert into grp, ascrip should be update
    // two  insert so begin transaction
    tx, err := sn.db.Begin()
    if err != nil {
        fmt.Println("db Begin err", err)
        request.GetConnection().SendMsg(st_err, []byte("transaction err"))
        return
    }
    // 1 insert
    query = `insert into grp (gId, gName, gTopId, gIntro) values
        (?, ?, ?, ?);`

    _, err = sn.db.Exec(query, pre, s[1], s[0], s[2])

    if err != nil {
        fmt.Println("CreateGroup err: ", err)
        request.GetConnection().SendMsg(st_err, []byte("Exec error"))
        // transaction rollback
        tx.Rollback()
        return 
    }

    // 2 insert
    query = `insert into Ascription (uId, gId, Membership) values
        (?, ?, 'Admin');`
    _, err = sn.db.Exec(query, s[0], pre)
    if err != nil {
        fmt.Println("CreateGroup err: ", err)
        request.GetConnection().SendMsg(st_err, []byte("Exec error"))
        // transaction rollback
        tx.Rollback()
        return 
    }

    // all ok try to commit transaction 
    if err := tx.Commit(); err != nil {
        tx.Rollback() // commit err 
        request.GetConnection().SendMsg(st_err, []byte("commit err"))
        return
    }

    // commit success
    request.GetConnection().SendMsg(st_ok, []byte("crate group successfully"))
}


type JoinGroup struct {
    znet.BaseRouter
}

func (jg *JoinGroup) Handle(request ziface.IRequest) {
    // uId | gId
    s := strings.Split(string(request.GetDate()), "|")
    if len(s) != 2 {
        request.GetConnection().SendMsg(st_err, []byte("format error"))
        return 
    }

    query := `insert into Ascription (uId, gId, Membership) values
        (?, ?, 'Member');`

    _, err := sn.db.Exec(query, s[0], s[1])
    if err != nil {
        fmt.Println("JoinGroup err: ", err)
        request.GetConnection().SendMsg(st_err, []byte("Exec err, maybe you have joint it."))
        return
    }

    request.GetConnection().SendMsg(st_ok, []byte("JoinGroup successfully."))
}



type LeaveGroup struct {
    znet.BaseRouter
}

func (lg *LeaveGroup) Handle(request ziface.IRequest) {
    // uid | gid 
    s := strings.Split(string(request.GetDate()), "|")
    // search ascrip 
    query := `select Membership from Ascription where uId = ? and gId = ?;`
    
    var m string 
    err := sn.db.QueryRow(query, s[0], s[1]).Scan(&m)
    if err != nil {
        if err == sql.ErrNoRows {
            request.GetConnection().SendMsg(st_err, []byte("You are not in this group."))
            return
        }
        request.GetConnection().SendMsg(st_err, []byte("LeaveGroup QueryRow err"))
        fmt.Println("LeaveGroup QueryRow err: ", err)
        return
    }

    // uid is Admin or Member  
    if m == "Member" {
        // delete a record
        query = `delete from Ascription where uid = ? and gId = ?;`
        _, err := sn.db.Exec(query, s[0], s[1])
        if err != nil {
            fmt.Println("LeaveGroup Exec err, ", err)
            request.GetConnection().SendMsg(st_err, []byte("Exec err"))
            return
        }
    } else {
        // cascade delete 
        // admin LeaveGroup --> delte this group  --> cascade(PublicMessage Ascription) 
        // 
        tx, err := sn.db.Begin()
        if err != nil {
            fmt.Println("Begin err:", err)
            request.GetConnection().SendMsg(st_err, []byte("Begin err"))
            return 
        }

        query = `delete from PublicMessage where mRecver = ?;`
        if _, err := sn.db.Exec(query, s[1]); err != nil {
            tx.Rollback()
            fmt.Println("Exec 1 err:", err)
            request.GetConnection().SendMsg(st_err, []byte("Exec 1 err"))
            return 
        }

        query = `delete from Ascription where gId = ?;`
        if _, err := sn.db.Exec(query, s[1]); err != nil {
            tx.Rollback()
            fmt.Println("Exec 2 err:", err)
            request.GetConnection().SendMsg(st_err, []byte("Exec 2 err"))
            return 
        }

        query = `delete from grp where gId = ?;`
        if _, err := sn.db.Exec(query, s[1]); err != nil {
            tx.Rollback()
            fmt.Println("Exec 3 err:", err)
            request.GetConnection().SendMsg(st_err, []byte("Exec 3 err"))
            return 
        }

        if err := tx.Commit(); err != nil {
            tx.Rollback()
            fmt.Println("commit err, ", err)
            request.GetConnection().SendMsg(st_err, []byte("commit err"))
            return
        }
    }

    request.GetConnection().SendMsg(st_ok, []byte("You leave this group."))
}


type getAllRelationship struct {
    znet.BaseRouter 
}

func (gar *getAllRelationship) Handle(request ziface.IRequest) {
    uid := string(request.GetDate())

    // can ensure unique relationship, so join directly
    // not with self
    query := `select uId, uName, rWhat from UserInfo
        join Relationship on uId = rUser0 or uId = rUser1
        where uId != ? and (rUser0 = ? or rUser1 = ?);`

    rows, err := sn.db.Query(query, uid, uid, uid)
    if err != nil {
        fmt.Println("getAllRelationship Query err: ", err)
        request.GetConnection().SendMsg(st_err, []byte("query err"))
        return 
    }
    defer rows.Close()

    var buf bytes.Buffer 
    for rows.Next() {
       var s [3]string
       rows.Scan(&s[0], &s[1], &s[2])
       buf.WriteString(s[0] + "|" + s[1] + "|" + s[2] + "#")
    }

    request.GetConnection().SendMsg(st_ok, buf.Bytes())
}


type InWhichGruop struct {
    znet.BaseRouter
}

func (iw *InWhichGruop) Handle(request ziface.IRequest) {
    uid := string(request.GetDate())

    // get group which uid in 
    query := `select gId, gName, uName, gIntro from Grp 
        join UserInfo on gTopId = uId 
        where exists (
            select * from Ascription A 
            where A.uId = ? and A.gId = Grp.gId
        );`

    rows, err := sn.db.Query(query, uid)
    if err != nil {
        fmt.Println("InWhichGruop err:", err)
        request.GetConnection().SendMsg(st_err, []byte("query err"))
        return
    }
    defer rows.Close()

    var buf bytes.Buffer 
    for rows.Next() {
        var s [4]string 
        rows.Scan(&s[0], &s[1], &s[2], &s[3])
        for i := 0; i < 3; i++ {
            buf.WriteString(s[i] + "|")
        }
        buf.WriteString(s[3] + "#")
    }

    request.GetConnection().SendMsg(st_ok, buf.Bytes())
}


// get all group message
type AllGroupMessage struct {
    znet.BaseRouter
}

func (ag *AllGroupMessage) Handle(request ziface.IRequest) {
    s := strings.Split(string(request.GetDate()), "|")
    if len(s) != 2 {
        fmt.Println("format err: len(2) = ", len(s))
        request.GetConnection().SendMsg(st_err, []byte("format err"))
        return
    }

    // user uid and git to get all public message 
    query := `select mSender, mContent from PublicMessage 
        where mRecver = ?;`
    rows, err := sn.db.Query(query, s[1])
    if err != nil {
        fmt.Println("query err, ", err)
        request.GetConnection().SendMsg(st_err, []byte("query err"))
        return
    }
    defer rows.Close()

    // join 
    var buf bytes.Buffer 
    for rows.Next() {
        var uid, msg string 
        rows.Scan(&uid, &msg)
        buf.WriteString(uid + "|" + msg + "#")
    }

    request.GetConnection().SendMsg(st_ok, buf.Bytes())
}


// the last router, deal gruop message 
type dealGrpMsg struct {
    znet.BaseRouter 
}

func (dg *dealGrpMsg) Handle(request ziface.IRequest) {
    // uid | gid | message
    s := strings.Split(string(request.GetDate()), "|")
    if len(s) != 3 {
        fmt.Println("format err, len(3) = ", len(s))
        // request.GetConnection().SendMsg(st_err, []byte{})
        return
    }

    query := `select max(mId) from PublicMessage;`
    var pre uint32
    sn.db.QueryRow(query).Scan(&pre)
    pre ++
    
    query = `insert into PublicMessage (mId, mContent, mTime, mSender, mRecver) values 
        (?, ?, ?, ?, ?);`
     
    // get current time
    currentTime := time.Now().Format("2006-01-02 15:04:05")

    if _, err := sn.db.Exec(query, pre, s[2], currentTime, s[0], s[1]); err != nil {
        fmt.Println("dealGrpMsg Exec err:", err)
        // request.GetConnection().SendMsg(st_err, []byte{})
        return
    } 

    // update database, then forward message 
    // forward or not depend on uid is in this gruop or not
    // because need many same sql, so start Mysql Prepare 
    query = `select count(*) from Ascription where uId = ? and gId = ?;`
    stmt, err := sn.db.Prepare(query)
    if err != nil {
        fmt.Println("db Prepare err:", err)
        return
    }
    defer stmt.Close()

    u, _ := strconv.ParseUint(s[0], 10, 0)

    sn.rwlock.RLock()
    defer sn.rwlock.RUnlock()
    for uid, connId := range sn.onlineMap {
        if uid == uint32(u) {
            continue
        }
        cnt := 0
        if err := stmt.QueryRow(uid, s[1]).Scan(&cnt); err == nil && cnt > 0 {
            // in this gruop 
            conn, err := sn.s.GetConnManager().Get(connId)
            if err == nil {
                conn.SendMsg(st_recvGrpMsg, request.GetDate())
                fmt.Println("group message forward:", string(request.GetDate()))
            }
            // else outline 
        } 
        // else not in this gruop
    }

}
