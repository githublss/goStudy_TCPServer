package ziface

type IServer interface {
	Start()
	Stop()
	Server()
	// 路由添加方法
	AddRouter(router IRouter)
}