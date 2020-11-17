package ziface

type IMsgHandle interface{
	//对消息进行处理
	DoMsgHandler(request IRequest)
	//添加路由
	AddRouter(msgId uint32, router IRouter)
	//设置工作worker池模式
	StarWorkerPool()
	//启动一个工作流程
	StartOneWorker(workerId int, taskQueue chan IRequest)
	//添加一个消息到消息队列中
	AddMsgToTaskQueue(request IRequest)
}
