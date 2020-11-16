package ziface

/*
链接管理接口
*/

type IConnManager interface {
	Add(connection IConnection)             //添加链接
	Remove(connection IConnection)          //删除链接
	Get(connId uint32) (IConnection, error) //利用connId获取链接
	Len() int                               //获取链接数量
	ClearConn()                             //删除停止所有链接
}
