package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

type DataPack struct {
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

//拆包封包实例初始化函数
func NewDataPack() *DataPack {
	return &DataPack{}
}

//将message封包，生成二进制发送格式
func (dp *DataPack) Pack(msg ziface.Imessage) ([]byte, error) {
	//创建一个存放bytes的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//写入长度
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//写入ID属性数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMessageId()); err != nil {
		return nil, err
	}
	//写入data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//解包,只读取头部数据
func (dp *DataPack) UnPack(binaryData []byte) (ziface.Imessage, error) {
	readBuff := bytes.NewReader(binaryData)
	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}
	//读取len
	if err := binary.Read(readBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读取ID
	if err := binary.Read(readBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too larger msg received 。。。")
	}

	return msg, nil
}
