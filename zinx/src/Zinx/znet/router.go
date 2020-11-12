package znet

import "Zinx/ziface"

//实现实际router的时候，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

func (router *BaseRouter) PreHandle(request ziface.IRequest) {}

func (router *BaseRouter) Handle(request ziface.IRequest) {}

func (router *BaseRouter) PostHandle(request ziface.IRequest) {}
