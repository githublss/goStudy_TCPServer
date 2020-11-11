package ziface

/*
对message进行拆包、封包
面向字节流，为数据添加头部信息，处理粘包粘包问题
 */
type IDataPack interface {
	GetHeadLen() uint32			// 获取包头长度
	Pack(msg Imessage)([]byte,error)	//封包方法
	UnPack([]byte)(Imessage,error)		//拆包问题
}