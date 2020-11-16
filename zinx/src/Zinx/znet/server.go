package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"fmt"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	// server注册的链接对应的处理业务
	MsgHandle ziface.IMsgHandle
	//链接管理器
	ConnManager ziface.IConnManager
	//链接创建时的hook函数
	OnConnStart func(connection ziface.IConnection)
	//链接断开时的hook函数
	OnConnStop func(connection ziface.IConnection)
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandle:   NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s
}
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP:%s, Port： %d, is started\n", s.IP, s.Port)
	//
	go func() {
		//启动工作池机制
		s.MsgHandle.StarWorkerPool()
		//解析获取addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error", err)
			return
		}
		// 监听addr并获取一个listener
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen:", s.IP, ":", s.Port, "error")
			return
		}
		fmt.Println("start Zinx server success---")
		var connId uint32 = 0
		// 阻塞等待客户端链接，处理客户端链接业务
		for {
			// 接收链接
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			fmt.Println("new connection=====", conn.RemoteAddr().String())
			//将新链接conn和路由进行绑定封装
			dealConn := NewConnection(s, conn, connId, s.MsgHandle)
			connId++
			go dealConn.Start()
		}
	}()

}
func (s *Server) Stop() {

	s.ConnManager.ClearConn()
}
func (s *Server) Server() {
	//启动服务
	s.Start()
	//做一些其他事项

	//阻塞
	select {}
}

// 添加一个路由
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	//检查是否已经存在
	s.MsgHandle.AddRouter(msgId, router)
}

//获取server的链接管理器
func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

//设置链接建立时的hook函数
func (s *Server) SetOnConnStart(hook func(connection ziface.IConnection)) {
	s.OnConnStart = hook
}

// 设置链接断开时的hook函数
func (s *Server) SetOnConnStop(hook func(connection ziface.IConnection)) {
	s.OnConnStop = hook
}

//调用hookStart函数
func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(connection)
	}
}

//调用hookStop函数
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(connection)
	}
}
