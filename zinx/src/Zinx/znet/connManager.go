package znet

import (
	"Zinx/ziface"
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	Connections    map[uint32]ziface.IConnection //链接集合
	ConnectionLock sync.RWMutex                  //读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[uint32]ziface.IConnection),
	}
}

//添加链接
func (cm *ConnManager) Add(connection ziface.IConnection) {
	//保护共享资源，加写锁
	cm.ConnectionLock.Lock()
	defer cm.ConnectionLock.Unlock()
	//将conn添加到集合中
	cm.Connections[connection.GetConnID()] = connection

	fmt.Println("connection:", connection.GetConnID(),
		"add to ConnManager successfully;Now count=", cm.Len())

}

//删除链接
func (cm *ConnManager) Remove(connection ziface.IConnection) {
	//保护共享资源，加写锁
	cm.ConnectionLock.Lock()
	defer cm.ConnectionLock.Unlock()
	//如果集合中存在此链接就从集合中删除
	if _, ok := cm.Connections[connection.GetConnID()]; ok == true {
		delete(cm.Connections, connection.GetConnID())
		fmt.Println("connection:", connection.GetConnID(),
			"remove from ConnManager successfully; Now cont=", cm.Len())
	} else {
		fmt.Println(connection.GetConnID(), "is not exit !!!")
	}
}

//利用connId获取链接
func (cm *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	//保护共享资源，加写锁
	cm.ConnectionLock.RLock()
	defer cm.ConnectionLock.RUnlock()
	//如果集合中存在此链接就从集合中删除
	if _, ok := cm.Connections[connId]; ok {
		return cm.Connections[connId], nil
	} else {
		fmt.Println(connId, "is not exit !!!")
		return nil, errors.New("Warning:get connection is not exit!!! ")
	}
}

//获取链接数量
func (cm *ConnManager) Len() int {
	return len(cm.Connections)
}

//删除停止所有链接
func (cm *ConnManager) ClearConn() {
	//保护共享资源，加写锁
	cm.ConnectionLock.Lock()
	defer cm.ConnectionLock.Unlock()

	//停止并删除所有链接
	for id, conn := range cm.Connections {
		conn.Stop()
		delete(cm.Connections, id)
	}

	fmt.Println("Clear all connection successfully, count = ", cm.Len())
}
