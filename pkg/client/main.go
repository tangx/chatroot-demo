package client

import (
	"flag"

	"github.com/sirupsen/logrus"
)

var client *Client

func init() {
	client = &Client{}

	flag.StringVar(&client.ServerAddr, "server", "127.0.0.1", "服务器地址，默认值为 127.0.0.1")
	flag.IntVar(&client.ServerPort, "port", 8888, "服务器默认端口， 8888")

	flag.Parse()
}
func Execute() {

	if err := client.Dial(); err != nil {
		logrus.Fatalf("dial server failed, err: %v", err)
	}

	// 启动协程 获取消息
	go client.Recevier()

	for {
		client.menu()
	}
}
