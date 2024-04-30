package ziface

// define a server interface 
type IServer interface {
    // Start Server 
    Start()

    // Stop Server 
    Stop()
    
    // run Server 
    Serve()

    // customize router
    AddRouter(uint32, IRouter)

    // Get ConnManager
    GetConnManager() IConnManager

    // register hook after new conn 
    SetOnConnStart(func (IConnection))

    // register hook after new conn 
    SetOnConnStop(func (IConnection))

    // call hook
    CallOnConnStart(IConnection)
    CallOnConnStop(IConnection)
}
