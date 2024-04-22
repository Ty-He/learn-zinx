package utils

import (
	"encoding/json"
	"io/ioutil"
	"my_zinx/ziface"
)

type GlobalObject struct {
    // server
    TcpServer ziface.IServer 
    Host string 
    TcpPort int 
    Name string 

    // zinx 
    Version string 
    MaxConn int 
    MaxPackageSize uint32
}

// global object 
var GlobalObj *GlobalObject 

// load jsonfile
func (this *GlobalObject) load() {
    data, err := ioutil.ReadFile("./conf/zinx.json")
    if err != nil {
        panic(err)
    }
    err = json.Unmarshal(data, &GlobalObj)
    if err != nil {
        panic(err)
    }
}

func init() {
    // load default value 
    GlobalObj = &GlobalObject {
        Name : "ZinxServerApp",
        Version : "V0.4",
        Host : "192.168.18.128",
        TcpPort : 8999,
        MaxConn : 10,
        MaxPackageSize : 1024,
    }
    // jsonfile
    GlobalObj.load()
}
