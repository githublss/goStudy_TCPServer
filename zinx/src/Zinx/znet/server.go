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
	Router ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP:%s, Port： %d, is started\n", s.IP, s.Port)
	//
	go func() {
		//解析获取addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error", err)
			return
		}
		// 监听addr并获取一个listener
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen:", s.IP,":",s.Port, "error")
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

			//将新链接conn和方法进行绑定封装
			dealConn := NewConnection(conn, connId, s.Router)
			connId++
			go dealConn.Start()
		}
	}()

}
func (s *Server) Stop() {

}
func (s *Server) Server() {
	//启动服务
	s.Start()
	//做一些其他事项

	//阻塞
	select {}
}

// 添加一个路由
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}
