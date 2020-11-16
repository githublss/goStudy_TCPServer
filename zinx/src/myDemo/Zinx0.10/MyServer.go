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
	fmt.Println("--->Recv from [client]", request.GetConnection().RemoteAddr().String(),
		" msg ID is:", request.GetMsgId(), " data is: ", string(request.GetData()))

	if err := request.GetConnection().Send(1, []byte("ding...ding...ding...")); err != nil {
		fmt.Println("connection send message error:", err)
	}

}

type HelloRouter struct {
	znet.BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("call handle ... ")

	//先读取客户端的数据，然后再向客户端写数据
	fmt.Println("--->Recv from [client]", request.GetConnection().RemoteAddr().String(),
		" msg ID is:", request.GetMsgId(), " data is: ", string(request.GetData()))

	if err := request.GetConnection().Send(1, []byte("hello...hello...hello...")); err != nil {
		fmt.Println("connection send message error:", err)
	}

}
func StartHook(connection ziface.IConnection) {
	fmt.Println("||||||call on the startHook")
	//设置属性
	connection.SetProperty("Name","Vector")
	connection.SetProperty("Home","https://github.com/githublss")

	if err := connection.Send(202, []byte("new connection is register,have fun!!! ")); err != nil {
		fmt.Println("connection send message error:", err)
	}
}
func StopHook(connection ziface.IConnection) {
	fmt.Println("||||||call on the stopHook")
	if name,err := connection.GetProperty("Name");err == nil{
		fmt.Println("Name is ",name)
	}
	if name,err := connection.GetProperty("Home");err == nil{
		fmt.Println("Home is ",name)
	}
	if err := connection.Send(1, []byte("see you, good buy!!! ")); err != nil {
		fmt.Println("connection send message error:", err)
	}
}

func main() {
	s := znet.NewServer("Version 0.8")
	//注册添加hook函数
	s.SetOnConnStart(StartHook)
	s.SetOnConnStop(StopHook)

	//添加消息Id对应的路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Server()
}
