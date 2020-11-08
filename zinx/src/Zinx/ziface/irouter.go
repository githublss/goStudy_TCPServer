package ziface

type IRouter interface {
	PreHandle(request IRequest)		//处理conn业务前的hook方法
	Handle(request IRequest)		//处理方法
	PostHandle(request IRequest)	//处理conn之后的hook方法
}
