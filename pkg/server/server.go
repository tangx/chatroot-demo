package server

import (
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

type Server struct {
	Ipaddr string
	Port   uint16

	// 在线用户
	OnlineUsers map[string]*User

	// 广播管道
	BroadCastChan chan string
}

func NewServer(ip string, port uint16) *Server {
	server := &Server{
		Ipaddr:        ip,
		Port:          port,
		OnlineUsers:   make(map[string]*User),
		BroadCastChan: make(chan string),
	}

	return server
}

// Signin 用户登录
func (s *Server) Signin(unique string, user *User) {
	s.OnlineUsers[unique] = user
}

// Handler 用户连接初始化
func (s *Server) Handler(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	logrus.Infof("%s connected\n", addr)

	// 1. 创建用户
	user := NewUser(conn)

	// 2. 用户注册
	name := user.UniqueName()
	s.Signin(name, user)

	// 3. 广播登录信息
	msg := fmt.Sprintf("用户 %s 已上线\n", name)
	s.BroadCast(user, msg)

	// 4. 用户接受消息
	go user.ListenMessage()

	// 5. 服务器接受用户消息并广播
	go func() {
		buf := make([]byte, 4096)
		// buf := []byte{} ????
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.BroadCast(user, "已下线")
				return
			}

			if err != nil && err == io.EOF {
				// 本条消息读取完成
				return
			}

			// 发送用户消息, 不发送 '\n'
			msg := buf[:n-1]
			s.BroadCast(user, string(msg))
		}

	}()

}

func (s *Server) ListenMessager() {
	for {
		// 1. 获取 msg。 无则阻塞
		msg := <-s.BroadCastChan
		// 2. 遍历所有用户， 并发送消息
		for _, user := range s.OnlineUsers {
			user.MsgChan <- msg
		}
	}
}

// BroadCast 广播
func (s *Server) BroadCast(user *User, msg string) {
	msg = fmt.Sprintf("%s: %s\n", user.Name, msg)
	s.BroadCastChan <- msg

}

// Start 启动服务
func (s *Server) Start() {
	// 1. 监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ipaddr, s.Port))
	if err != nil {
		logrus.Fatalf("Server start failed: %v", err)

	}

	// n. defer 关闭
	defer listener.Close()

	// 2. 监听广播
	go s.ListenMessager()

	// 3. 提供服务
	for {
		// 3.1 建立链接
		conn, err := listener.Accept()
		if err != nil {
			if err != nil {
				logrus.Errorf("客户端建立链接错误")
				// 3.1.1 下一次循环继续等待链接
				continue
			}
		}
		// 3.2 开辟 goroutine， 后台处理
		go s.Handler(conn)
		// 3.2.1 下一次循环继续等待链接
	}
}
