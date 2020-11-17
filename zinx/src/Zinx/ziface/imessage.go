package ziface

/*
请求消息封装到message中，定义抽象接口
*/
type Imessage interface {
	//得到消息长度
	GetDataLen() uint32
	//得到消息Id
	GetMessageId() uint32
	//得到消息体
	GetData() []byte
	//存储设置消息长度
	SetDataLen(uint32)
	//设置消息体
	SetData([]byte)
	//设置消息id
	SetMessageId(uint32)
}
