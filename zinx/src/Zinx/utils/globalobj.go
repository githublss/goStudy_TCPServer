package utils

import (
	"Zinx/ziface"
	"encoding/json"
	"io/ioutil"
)

type GlobalObj struct {
	/*
		server
	*/
	TcpServer ziface.IServer //当前zinx的全局server对象
	Host      string
	TcpPort   int
	Name      string //当前服务器名称

	/*
		zinx
	*/
	Version        string //current version
	MaxPacketSize  uint32 //size of max message
	MaxConn        int    //size of max connection num
	WorkerPoolSize uint32 //工作池的数量
	TaskQueueSize  uint32 //每个工作池的等待队列的大小
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("config/zinx.json")
	if err != nil {
		panic(err)
	}
	//将配置文件解析到struct中
	if err := json.Unmarshal(data, g); err != nil {
		panic(err)
	}

}

/*
一个文件的init方法会在第一次加载的时候在main方法之前运行
*/
func init() {
	//默认配置
	GlobalObject = &GlobalObj{
		Name:          "ZinxServer",
		Version:       "V0.1",
		TcpPort:       9999,
		Host:          "0.0.0.0",
		MaxConn:       1024,
		MaxPacketSize: 1024,
		WorkerPoolSize: 10,
		TaskQueueSize: 1024,
	}
	GlobalObject.Reload()
}
