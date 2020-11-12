package znet

import (
	"Zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	// 当前链接的tcp套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前链接的状态
	IsClosed bool

	//告知当前连接已退出的 channel
	ExitChan chan bool
	//消息管理模块，用来绑定MsgId和对应的处理方法
	MsgHandle ziface.IMsgHandle
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		IsClosed:  false,
		MsgHandle: msgHandle,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID", c.ConnID, "reader is exit,Remote is -> ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		dp := NewDataPack()
		//从客户端发送数据读取 msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			fmt.Println("read msg head error ", err)
			return
		}
		//拆包，将包头数据解析放在msg中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("Unpack error ", err)
			return
		}
		var data []byte
		//根据dataLen读取data，放在msg的data中
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.Conn, data); err != nil {
				fmt.Println("from conn read data error ", err)
				return
			}
			msg.SetData(data)
		}
		//得到当前客户端请求的Request数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		//执行注册路由的的方法
		c.MsgHandle.DoMsgHandler(req)
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn start()... connID=", c.ConnID)
	//启动读数据业务
	go c.StartReader()
	//TODO 启动当前链接的写业务

}

func (c *Connection) Stop() {
	fmt.Println("Conn stop()...  ConnID=", c.ConnID)

	//当前链接已经关闭
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
	//关闭socket
	if err := c.Conn.Close(); err != nil {
		fmt.Println("Close err", err)
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

//将data消息进行封包，发送给TCP客户端
func (c *Connection) Send(msgId uint32, data []byte) error {
	// TODO 是否加锁
	if c.IsClosed {
		return errors.New("Connection closed when send msg ")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("when shed pack msgPackage error ", err)
		return errors.New("Pack newMsgPackage error ")
	}

	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Send data write error", err)
		return err
	}
	return nil
}
