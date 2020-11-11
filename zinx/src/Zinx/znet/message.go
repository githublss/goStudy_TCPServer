package znet

type Message struct{
	Id uint32
	DataLen uint32
	Data []byte
}

//创建一个message包
func NewMsgPackage(id uint32,data []byte)  *Message{
	return &Message{
		Id: id,
		Data: data,
		DataLen: uint32(len(data)),
	}
}

func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

func (msg *Message) GetMessageId() uint32 {
	return msg.Id
}

func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) SetDataLen(u uint32) {
	msg.DataLen = u
}

func (msg *Message) SetData(bytes []byte) {
	msg.Data = bytes
}

func (msg *Message) SetMessageId(u uint32) {
	msg.Id = u
}

