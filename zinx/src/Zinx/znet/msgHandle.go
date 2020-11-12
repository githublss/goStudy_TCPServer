package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"fmt"
)

type MsgHandle struct {
	//每个msgId对应的Router
	ApiHandles     map[uint32]ziface.IRouter
	WorkerPoolSize uint32                 //worker池的数量
	TaskQueue      []chan ziface.IRequest //worker负责任务的队列
	//TODO 使用一个任务队列，让worker从队列中依次取出任务
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		ApiHandles:     make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		//一个worker一个chan
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}
func (mh MsgHandle) AddMsgToTaskQueue(request ziface.IRequest)  {
	//根据clientId来轮询的分配消息,使用取余法得到要处理的WorkId
	workId := request.GetConnection().GetConnID() % utils.GlobalObject.WorkerPoolSize
	//将任务添加到对应work的任务队列中
	mh.TaskQueue[workId] <- request
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.ApiHandles[request.GetMsgId()]
	//检查是否存在处理方法
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), "not found ")
		return
	}
	//执行处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//检查是否已经存在
	if _, ok := mh.ApiHandles[msgId]; ok {
		fmt.Println("Warning: msgId:", msgId, " is Exist")
		return
	}
	mh.ApiHandles[msgId] = router
	fmt.Println("msgId deal router add success:", msgId)
}
//启动工作池
func (mh *MsgHandle) StarWorkerPool() {
	//根据worker的数量，生成每个worker的任务队列
	for i:=0;i<int(utils.GlobalObject.WorkerPoolSize);i++{
		//一个worker被启动，给当前worker分配对应的任务队列
		mh.TaskQueue[i] = make(chan ziface.IRequest,utils.GlobalObject.TaskQueueSize)
		//启动当前worker，阻塞的等待对应的任务队列上时候有任务需要处理
		go mh.StartOneWorker(i,mh.TaskQueue[i])
	}
}
//启动一个工作流程
func (mh *MsgHandle) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID=",workerId,"is working...")
	//循环从任务队列中取出任务
	for{
		select {
		//case msg := <- mh.TaskQueue[workerId]:
		case msg := <- taskQueue:
			mh.DoMsgHandler(msg)
		}
	}
}
