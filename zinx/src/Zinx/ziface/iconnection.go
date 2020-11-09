package ziface

import "net"

/*
对connection进行一个封装
 */
type IConnection interface {
	//启动
	Start()
	//停止
	Stop()

	//获取原始TcpConnection
	GetTCPConnection() *net.TCPConn

	//获取当前ID
	GetConnID() uint32

	//获取远端客户端TCP状态 IP port
	RemoteAddr() net.Addr

	//发送数据
	Send(data []byte) error
}

//定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error