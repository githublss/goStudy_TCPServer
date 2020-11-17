package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {
	//当前链接的服务
	Server ziface.IServer

	// 当前链接的tcp套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前链接的状态
	IsClosed bool

	//用于读，写两个goroutine之间消息的无缓冲管道
	MsgChan chan []byte
	//告知当前连接已退出的 channel
	ExitChan chan bool
	//消息管理模块，用来绑定MsgId和对应的处理方法
	MsgHandle ziface.IMsgHandle
	//链接属性
	Property map[string]interface{}
	//property保护锁
	PropertyLock sync.Mutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		Server:    server,
		Conn:      conn,
		ConnID:    connID,
		IsClosed:  false,
		MsgHandle: msgHandle,
		ExitChan:  make(chan bool, 1),
		MsgChan:   make(chan []byte),
		Property: map[string]interface{}{},
	}
	//将创建的conn添加到管理器中
	c.Server.GetConnManager().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader] Goroutine is running...")
	defer fmt.Println("connID", c.ConnID, "[Reader] is exit,Remote is -> ", c.RemoteAddr().String())
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

		if utils.GlobalObject.TaskQueueSize > 0 {
			//如果开启了工作池机制，就将任务添加到任务队列中，让worker调用自己的msgHandle方法
			c.MsgHandle.AddMsgToTaskQueue(req)
		} else {
			//没有开启工作池机制的时候，启动一个go程来执行对应的handle方法
			go c.MsgHandle.DoMsgHandler(req)
		}

	}
}
func (c *Connection) StartWrite() {
	fmt.Println("[Writer] goroutine is running...")
	defer fmt.Println("connID", c.ConnID, "[Write] is exit,Remote is -> ", c.RemoteAddr().String())
	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Write data error:", err)
				return
			}
		case <-c.ExitChan:
			fmt.Println("Writer is closed.")
			return
		}
	}

}
func (c *Connection) Start() {
	fmt.Println("Conn start()... connID=", c.ConnID)

	//启动读数据业务
	go c.StartReader()
	//启动写数据业务
	go c.StartWrite()

	//调用用户注册的创建链接时的预处理hook方法
	c.Server.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn stop()...  ConnID=", c.ConnID)

	//调用关闭回调业务
	c.Server.CallOnConnStop(c)

	//当前链接已经关闭
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
	//关闭socket
	if err := c.Conn.Close(); err != nil {
		fmt.Println("Close err", err)
	}

	c.ExitChan <- true

	//将链接从管理器中删除
	c.Server.GetConnManager().Remove(c)

	//关闭reader与writer间的消息管道
	close(c.MsgChan)
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

	c.MsgChan <- msg

	return nil
}
//设置属性
func (c *Connection) SetProperty(name string,value interface{}){
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	c.Property[name] = value
}
//获取属性
func (c *Connection) GetProperty(name string) (interface{},error){
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	if value,ok := c.Property[name];ok{
		return value,nil
	}
	return nil,errors.New("not have the Property")
}
//移除属性
func (c *Connection) RemoveProperty(name string){
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	if _,ok := c.Property[name];ok{
		delete(c.Property,name)
	}
}
