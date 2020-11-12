package znet

import (
	"Zinx/ziface"
	"fmt"
)

type MsgHandle struct{
	ApiHandles map[uint32] ziface.IRouter

}
func NewMsgHandle() *MsgHandle{
	return &MsgHandle{
		ApiHandles: make(map[uint32] ziface.IRouter),
	}
}
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler,ok := mh.ApiHandles[request.GetMsgId()]
	//检查是否存在处理方法
	if !ok{
		fmt.Println("api msgId = ",request.GetMsgId(),"not found ")
		return
	}
	//执行处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//检查是否已经存在
	if _,ok := mh.ApiHandles[msgId];ok{
		fmt.Println("Warning: msgId:",msgId," is Exist")
		return
	}
	mh.ApiHandles[msgId] = router
	fmt.Println("msgId deal router add success:",msgId)
}
