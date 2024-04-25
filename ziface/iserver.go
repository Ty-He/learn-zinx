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
}
