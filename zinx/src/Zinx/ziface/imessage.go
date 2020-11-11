package ziface

/*
请求消息封装到message中，定义抽象接口
*/
type Imessage interface {
	GetDataLen() uint32
	GetMessageId() uint32
	GetData() []byte

	SetDataLen(uint32)
	SetData([]byte)
	SetMessageId(uint32)
}
