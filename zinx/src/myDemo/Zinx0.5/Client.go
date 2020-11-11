package main

import (
	"Zinx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client start......")
	time.Sleep(1 * time.Second)
	connect, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("client connect error:", err)
		return
	}
	for {
		dp := znet.NewDataPack()
		// 向服务端发送一个包
		binarySendMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("V0.5 client test... ")))
		if err != nil {
			fmt.Println("client pack msgPackage error:", err)
			return
		}
		if _, err := connect.Write(binarySendMsg); err != nil {
			fmt.Println("client write error:", err)
			return
		}
		// 服务端返回一个id为1的消息
		binaryGetHead := make([]byte,dp.GetHeadLen())
		if _,err := io.ReadFull(connect,binaryGetHead);err != nil{
			fmt.Println("client get messageHead error:",err)
			return
		}
		// 将二进制数据拆包为自定义msg结构中
		msgHead,err := dp.UnPack(binaryGetHead)
		if err != nil {
			fmt.Println("client unPack message error:",err)
			return
		}

		//根据dataLen的值将数据body读取出来
		if msgHead.GetDataLen() > 0 {
			msgBody := make([]byte,msgHead.GetDataLen())
			if _,err := io.ReadFull(connect,msgBody);err != nil{
				fmt.Println("client read message body error:",err)
				return
			}
			fmt.Println("------> Recv message from server message is:",string(msgBody)," ID is:",msgHead.GetMessageId())
		}

		time.Sleep(2 * time.Second)
	}
}
