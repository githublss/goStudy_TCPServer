package main

import (
	"fmt"
	"net"
	"time"
)

func main()  {
	fmt.Println("client start......")
	time.Sleep(1*time.Second)
	connect,err := net.Dial("tcp","127.0.0.1:9999")
	if err != nil{
		fmt.Println("client connect error:",err)
		return
	}
	for{

		if _,err := connect.Write([]byte("hello"));err != nil{
			fmt.Println("client Write error:",err)
			return
		}
		buff := make([]byte,512)
		readLen,err := connect.Read(buff)
		if err != nil{
			fmt.Println("client Read error",err)
			return
		}
		fmt.Printf("from server message:%s,message len:%d \n",buff,readLen)

		time.Sleep(2*time.Second)
	}
}