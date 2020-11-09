package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"fmt"
	"net"
)

type Connection struct{
	// 当前链接的tcp套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前链接的状态
	IsClosed bool

	//告知当前连接已退出的 channel
	ExitChan chan bool


	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:   conn,
		ConnID: connID,
		IsClosed:  false,
		Router: router,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader()  {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID",c.ConnID,"reader is exit,Remote is -> ",c.RemoteAddr().String())
	defer c.Stop()

	for{
		//读取数据
		buf := make([]byte,utils.GlobalObject.MaxPacketSize)
		_,err := c.Conn.Read(buf)
		if err != nil{
			fmt.Println("recv buff err",err)
			continue
		}
		req := Request{
			conn: c,
			data: buf,
		}

		//执行注册路由的的方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn start()... connID=",c.ConnID)
	//启动读数据业务
	go c.StartReader()
	//TODO 启动当前链接的写业务

}

func (c *Connection) Stop() {
	fmt.Println("Conn stop()...  ConnID=",c.ConnID)

	//当前链接已经关闭
	if c.IsClosed == true{
		return
	}
	c.IsClosed = true
	//关闭socket
	if err := c.Conn.Close();err != nil{
		fmt.Println("Close err",err)
	}
	//关闭管道
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	panic("implement me")
}
