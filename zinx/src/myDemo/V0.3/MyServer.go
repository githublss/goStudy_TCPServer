package main

import (
	"Zinx/ziface"
	"Zinx/znet"
	"fmt"
)
type PingRouter struct {
	znet.BaseRouter
}


func (this *PingRouter) Handle(request ziface.IRequest)  {
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("ping...Ping...\n"))
	if err != nil{
		fmt.Println("call back handle error")
	}
}

func main()  {
	s := znet.NewServer("Version 0.3")
	s.AddRouter(&PingRouter{})
	s.Server()
}
