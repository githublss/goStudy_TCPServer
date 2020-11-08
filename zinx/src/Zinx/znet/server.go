package znet

import (
	"Zinx/ziface"
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Name string
	IPVersion string
	IP string
	Port int
}
//临时定义处理客户端链接的函数
func CallbackToClient(conn *net.TCPConn, data []byte,count int) error{
	fmt.Println("[Conn handle] CallbackClient")
	if _,err := conn.Write(data[:count]);err!=nil{
		fmt.Println("write back buff error:",err)
		return errors.New("CallbackToClient error")
	}
	return nil
}
func (s *Server) Start()  {
	fmt.Printf("[Start] Server Listenner at IP:%s, Port： %d, is started\n", s.IP, s.Port)
	//
	go func() {
		//解析获取addr
		addr, err := net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err != nil{
			fmt.Println("resolve tcp addr error",err)
			return
		}
		// 监听addr并获取一个listener
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil{
			fmt.Println("listen:",s.IP,"error")
			return
		}
		fmt.Println("start Zinx server success---")
		var connId uint32 = 0
		// 阻塞等待客户端链接，处理客户端链接业务
		for {
			// 接收链接
			conn,err := listener.AcceptTCP()
			if err != nil{
				fmt.Println("Accept err",err)
				continue
			}
			//
			//go func() {
			//	for{
			//		buf := make([]byte,512)
			//		readLen,err := conn.Read(buf)
			//		if err != nil{
			//			fmt.Println("Read err",err)
			//			continue
			//		}
			//		fmt.Printf("from client a message:%s \n",buf)
			//		if _, err := conn.Write(buf[:readLen]); err != nil{
			//			fmt.Println("Write err",err)
			//			continue
			//		}
			//	}
			//}()

			//将新链接conn和方法进行绑定封装
			dealConn := NewConnection(conn,connId,CallbackToClient)
			connId++
			go dealConn.Start()
		}
	}()

}
func (s *Server) Stop(){

}
func (s *Server) Server()  {
	//启动服务
	s.Start()
	//做一些其他事项

	//阻塞
	select {}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name: name,
		IPVersion: "tcp4",
		IP: "0.0.0.0",
		Port: 8999,
	}
	return s
}