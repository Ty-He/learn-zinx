package main 

import "my_zinx/znet"

func main() {
    // server 
    s := znet.NewServer("ZinxV0.2")
    // run serve
    s.Serve()
}
