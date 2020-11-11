package main

import (
	"Zinx/ziface"
	"Zinx/znet"
	"fmt"
)

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call handle ... ")

	//先读取客户端的数据，然后再向客户端写数据
	fmt.Println("--->Recv from [client] msg ID is:", request.GetMsgId(), " data is: ", string(request.GetData()))

	if err := request.GetConnection().Send(1, []byte("ding...ding...ding...")); err != nil {
		fmt.Println("connection send message error:", err)
	}

}

func main() {
	s := znet.NewServer("Version 0.5")

	s.AddRouter(&PingRouter{})
	s.Server()
}
