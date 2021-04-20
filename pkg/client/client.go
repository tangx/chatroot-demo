package client

import (
	"fmt"
	"net"
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

func (c *Client) Dial() error {

	address := fmt.Sprintf("%s:%d", c.ServerAddr, c.ServerPort)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}
