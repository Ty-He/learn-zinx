package znet

import "my_zinx/ziface"

type BaseRouter struct {}

// apply empty method
func (self *BaseRouter) PreHandle(request ziface.IRequest) {}


func (self *BaseRouter) Handle(request ziface.IRequest) {}


func (self *BaseRouter) PostHandle(request ziface.IRequest) {}
