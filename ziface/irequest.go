package ziface

type IRequest interface {
    // socket fd
    GetConnection() IConnection
    // client request data
    GetDate() []byte
}
