package main

import "github.com/tangx/chatroot-demo/pkg/server"

// 启动服务器
func main() {
	s := server.NewServer("127.0.0.1", 8888)
	s.Start()
}
