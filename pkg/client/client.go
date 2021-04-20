package client

import (
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

type Client struct {
	ServerAddr string
	ServerPort int
	conn       net.Conn
}

func NewClient(ip string, port int) *Client {
	c := &Client{
		ServerAddr: ip,
		ServerPort: port,
	}
	return c
}

// Dial 链接服务器
func (c *Client) Dial() error {

	address := fmt.Sprintf("%s:%d", c.ServerAddr, c.ServerPort)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

// PublicChat 公聊模式
func (c *Client) PublicChat() {

}

// PrivateChat 私聊模式
func (c *Client) PrivateChat() {}

// UpdateName 更新用户名
func (c *Client) UpdateName() {
	var name string

	fmt.Printf(">> ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		fmt.Printf("输入错误: %v", err)
	}

	msg := fmt.Sprintf("rename|%s", name)

	_, err = c.conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("update name failed: %v \n", err)
		return
	}

	fmt.Printf("update name success")
}

func (c *Client) menu() {
	str := `
菜单:
	0. 退出
	1. 公聊模式
	2. 私聊模式
	3. 更新用户名
`

	fmt.Println(str)
	var choice string
	_, err := fmt.Scanln(&choice)
	if err != nil {
		logrus.Warnf("input failed, err: %v", err)
	}

	switch choice {
	case "0":
		os.Exit(0)
	case "1":
		c.PublicChat()
	case "2":
		c.PrivateChat()
	case "3":
		c.UpdateName()
	}

}
