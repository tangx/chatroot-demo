package server

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

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
		_, _ = u.conn.Write([]byte(msg))
	}
}

// Online 用户上线 并 广播消息
func (u *User) Online() {
	u.server.userlock.Lock()
	defer u.server.userlock.Unlock()

	u.server.OnlineUsers[u.Name] = u
	u.server.BroadCast(u, "上线")
}

// Offline 用户下线 并 广播消息
func (u *User) Offline() {
	u.server.userlock.Lock()
	defer u.server.userlock.Unlock()

	delete(u.server.OnlineUsers, u.Name)
	u.server.BroadCast(u, "下线")
}

// DoMessage 发送消息
func (u *User) DoMessage(msg string) {
	// 查询在线用户
	if msg == "who" {
		who := ""
		for _, user := range u.server.OnlineUsers {
			who += fmt.Sprintf("[%s]%s : 在线\n", user.Addr, user.Name)
		}
		u.SendMessage(who)
		return
	}

	// 修改用户名
	if ok, _ := regexp.MatchString(`^rename\|`, msg); ok {
		u.Rename(msg)
		return
	}

	// 私聊
	if ok, _ := regexp.MatchString(`^to\|`, msg); ok {
		u.PrivateMessage(msg)
		return
	}

	// 广播发送消息
	u.server.BroadCast(u, msg)
}

func (u *User) PrivateMessage(msg string) {
	parts := strings.Split(msg, "|")
	if len(parts) < 3 {
		u.SendMessage("消息格式错误， ex: to|zhang3|content")
		return
	}

	who := parts[1]
	remoteUser, ok := u.server.OnlineUsers[who]
	if !ok {
		u.SendMessage(fmt.Sprintf("系统消息: 用户 %s 不存在", who))
		return
	}

	content := strings.Join(parts[2:], "|")
	remoteUser.SendMessage(fmt.Sprintf("[%s]私聊你: %s\n", u.Name, content))

}

// SendMessage 给自己发送消息
func (u *User) SendMessage(msg string) {
	_, _ = u.conn.Write([]byte(msg))
}

// Rename 用户指令， 重命名
func (u *User) Rename(msg string) {
	parts := strings.Split(msg, "|")
	if len(parts) < 2 {
		u.SendMessage("rename 指令错误。 ex: rename|zhang31 \n")
		return
	}

	name := parts[1]
	if len(name) == 0 {
		u.SendMessage("rename 指令错误。 用户名不能为空 \n")
		return
	}

	if _, ok := u.server.OnlineUsers[name]; ok {
		u.SendMessage(fmt.Sprintf("系统消息: 用户名 %s 已存在 \n", name))
		return
	}

	u.rename(name)
	u.server.BroadCast(u, "系统消息: 已改名")

}

// rename 修改用户名
func (u *User) rename(name string) {
	u.server.userlock.Lock()
	defer u.server.userlock.Unlock()

	delete(u.server.OnlineUsers, u.Name)
	u.Name = name
	u.server.OnlineUsers[u.Name] = u

}

func (u *User) Kickoff() {
	u.server.userlock.Lock()
	defer u.server.userlock.Unlock()

	// 删除在线记录
	delete(u.server.OnlineUsers, u.Name)
	// 发送通知
	_, _ = u.conn.Write([]byte("系统消息: 长时间不活跃，你已经被系统踢下线"))
	// 关闭链接
	if err := u.conn.Close(); err != nil {
		logrus.Errorf("close %s connect failed: %v\n", u.Addr, err)
	}
}
