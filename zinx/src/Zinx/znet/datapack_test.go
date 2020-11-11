package znet

import (
	"Zinx/ziface"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestDataPack_Pack(t *testing.T) {
	//创建socket TCP server
	listener, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}
	//创建服务器go程，从客户端go程读取粘包数据，解析
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server's accept err:", err)
				continue
			}
			//在go程中参数传递的作用
			go func(conn net.Conn) {
				//创建拆封包对象
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					if _, err := io.ReadFull(conn, headData); err != nil {
						fmt.Println("read head error")
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server Unpack error:", err)
						continue
					}
					if msgHead.GetDataLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())
						if _, err := io.ReadFull(conn, msg.Data); err != nil {
							fmt.Println("server read data from conn error:", err)
							continue
						}
						fmt.Println("====> Recv msg: ID ", msg.Id, "len=", msg.DataLen, "data:", string(msg.Data))
					}
				}
			}(conn)
		}
	}()
	//创建客户端go程，向服务端go程发送数据
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:9999")
		if err != nil {
			fmt.Println("client connection error:", err)
			return
		}
		dp := NewDataPack()
		var msg1 ziface.Imessage = &Message{
			Id:      0,
			Data:    []byte("hello"),
			DataLen: 5,
		}
		sendData1, err := dp.Pack(msg1)
		if err != nil{
			fmt.Println("client packData error:",err)
			return
		}
		var msg2 ziface.Imessage = &Message{
			Id:      1,
			Data:    []byte("hahaha"),
			DataLen: 6,
		}
		sendData2, err := dp.Pack(msg2)
		if err != nil{
			fmt.Println("client packData error:",err)
			return
		}
		sendData1 = append(sendData1, sendData2...)

		if _,err := conn.Write(sendData1);err != nil{
			fmt.Println("clent write error:",err)
			return
		}
	}()
	//阻塞
	select {
	case <-time.After(time.Second*2):
		return
	}
}
