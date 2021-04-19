package server

import "net"

type User struct {
	Name    string
	Addr    string
	MsgChan chan string
	conn    net.Conn
	server  *Server // 绑定 User 与 Server
}

func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		MsgChan: make(chan string),
		Name:    addr,
		Addr:    addr,
		conn:    conn,
		server:  server,
	}
	return user
}

// UniqueName 用户唯一ID
func (u *User) UniqueName() string {
	return u.conn.RemoteAddr().String()
}

// Recevie 接受消息
func (u *User) ListenMessage() {

	for {
		msg := <-u.MsgChan
		u.conn.Write([]byte(msg))
	}
}

// Online 用户上线 并 广播消息
func (u *User) Online() {
	u.server.OnlineUsers[u.Name] = u
	u.server.BroadCast(u, "上线")
}

// Offline 用户下线 并 广播消息
func (u *User) Offline() {
	delete(u.server.OnlineUsers, u.Name)
	u.server.BroadCast(u, "下线")
}

// DoMessage 发送消息
func (u *User) DoMessage(msg string) {
	u.server.BroadCast(u, msg)
}
