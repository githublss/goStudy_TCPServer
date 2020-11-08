package main

import "Zinx/znet"

func main()  {
	s := znet.NewServer("Version 0.2")
	s.Server()
}
