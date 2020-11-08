package ziface

type IRequest interface {
	GetConnection() IConnection 	//获取请求连接
	GetData() []byte				//获取请求消息数据
}