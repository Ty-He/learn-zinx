package ziface

// define a server interface 
type IServer interface {
    // Start Server 
    Start()
    // Stop Server 
    Stop()
    // run Server 
    Serve()
}
