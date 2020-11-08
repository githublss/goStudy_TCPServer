package znet

import "Zinx/ziface"

type Request struct{
	//已建立的连接
	conn ziface.IConnection
	//请求的数据
	data []byte
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
