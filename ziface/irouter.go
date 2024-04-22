package ziface

type IRouter interface {
    // pre-hook
    PreHandle(request IRequest)
    // main work
    Handle(request IRequest)
    // post-host
    PostHandle(request IRequest)
}
