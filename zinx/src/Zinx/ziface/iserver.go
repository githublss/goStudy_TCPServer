package ziface

type IServer interface {
	Start()
	Stop()
	Server()
	// 路由添加方法
	AddRouter(msgId uint32, router IRouter)
	//获取server的链接管理器
	GetConnManager() IConnManager
	//设置链接建立时的hook函数
	SetOnConnStart(func(IConnection))
	// 设置链接断开时的hook函数
	SetOnConnStop(hook func(connection IConnection))
	//调用hookStart函数
	CallOnConnStart(connection IConnection)
	//调用hookStop函数
	CallOnConnStop(connection IConnection)
}
