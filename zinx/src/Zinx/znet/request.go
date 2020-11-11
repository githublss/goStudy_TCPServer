package znet

import "Zinx/ziface"

type Request struct {
	//已建立的连接
	conn ziface.IConnection
	//请求的数据
	msg ziface.Imessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}
//TODO 更好的方式是获取一个message，然后调用者使用返回的message的方法
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMessageId()
}
