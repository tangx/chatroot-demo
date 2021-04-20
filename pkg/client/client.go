package client

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type Client struct {
	ServerAddr string
	ServerPort uint16
	conn       net.Conn
}

func NewClient(ip string, port uint16) *Client {
	c := &Client{
		ServerAddr: ip,
		ServerPort: port,
	}
	return c
}

func (c *Client) Dial() bool {

	address := fmt.Sprintf("%s:%d", c.ServerAddr, c.ServerPort)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		logrus.Errorf("dial server %s failed, err: %v", address, err)
		return false
	}

	c.conn = conn
	return true
}
