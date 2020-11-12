package ziface

type IMsgHandle interface{
	DoMsgHandler(request IRequest)
	AddRouter(msgId uint32, router IRouter)
	StarWorkerPool()
	StartOneWorker(workerId int, taskQueue chan IRequest)
	AddMsgToTaskQueue(request IRequest)
}
