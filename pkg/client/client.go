package client

import (
	"fmt"
	"io"
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

	fmt.Printf("输入新用户名 >> ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		fmt.Printf("输入错误: %v", err)
	}

	msg := fmt.Sprintf("rename|%s", name)

	err = c.sendMessage(msg)
	if err != nil {
		fmt.Printf("update name failed: %v \n", err)
		return
	}

	fmt.Printf("update name success")
}

// Recevier 消息接收器
func (c *Client) Recevier() {
	// 以下命令永久阻塞，并监听隧道获取消息
	_, err := io.Copy(os.Stdout, c.conn)
	if err != nil {
		logrus.Fatalf("Recevie Message failed: %v", err)
	}

	/* 以上命令等价于 */
	// for {
	// 	content := make([]byte, 4096)
	// 	_, err := c.conn.Read(content)
	// 	fmt.Println(content)
	// }
}

func (c *Client) sendMessage(msg string) error {
	msg = fmt.Sprintf("%s\n", msg)
	_, err := c.conn.Write([]byte(msg))
	return err

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
	n, err := fmt.Scanln(&choice)
	// 用户输入空行空行
	if n == 0 {
		return
	}
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
