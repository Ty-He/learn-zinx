package ziface

type IMsghandle interface {
    // handle msg by map
    DoHandle(IRequest)

    // add router to map
    AddRouter(uint32, IRouter)
}
