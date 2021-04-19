package server

import "net"

type User struct {
	Name    string
	Addr    string
	MsgChan chan string
	conn    net.Conn
}

func NewUser(conn net.Conn) *User {
	user := &User{
		MsgChan: make(chan string),
		conn:    conn,
	}
	return user
}

// UniqueName 用户唯一ID
func (u *User) UniqueName() string {
	return u.conn.RemoteAddr().String()
}

// Recevie 接受消息
func (u *User) Receive() {

	for {
		select {
		case msg := <-u.MsgChan:
			u.conn.Write([]byte(msg))
		}
	}
}
