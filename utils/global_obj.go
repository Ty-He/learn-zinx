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

    // WorkerPool's count
    WorkerPoolSize uint32
    // limit WorkerPoolSize
    MaxTaskSize uint32 
    
    // app-Version
    AppVersion string
}

// global object 
var GlobalObj *GlobalObject 

// load jsonfile
func (this *GlobalObject) load() {
    data, err := ioutil.ReadFile("/home/ty/go/src/my_demo/zinx_app/conf/zinx.json")
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
        Version : "V0.10",
        Host : "192.168.18.128",
        TcpPort : 8999,
        MaxConn : 10,
        MaxPackageSize : 1024,
        WorkerPoolSize : 10,
        MaxTaskSize : 1024,
    }
    // jsonfile
    GlobalObj.load()
}
