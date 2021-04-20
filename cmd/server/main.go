package main

import (
	"flag"

	"github.com/tangx/chatroot-demo/pkg/server"
)

var (
	ListenIP   string
	ServerPort int
)

// 启动服务器
func main() {
	s := server.NewServer("127.0.0.1", 8888)
	s.Start()
}

func init() {
	flag.StringVar(&ListenIP, "listen", "127.0.0.1", "服务监听IP， 默认 127.0.0.1")
	flag.IntVar(&ServerPort, "port", 8888, "服务监听端口， 默认 8888")

	flag.Parse()
}
