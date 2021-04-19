package server

import (
	"fmt"
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
	uniquename := user.UniqueName()
	s.Signin(uniquename, user)

	// 3. 广播登录信息
	msg := fmt.Sprintf("%s login\n", uniquename)
	s.BroadCastChan <- msg

	// 4. 用户接受消息
	go user.Receive()

}

func (s *Server) BroadCast() {
	for {
		// 1. 获取 msg。 无则阻塞
		msg := <-s.BroadCastChan
		// 2. 遍历所有用户， 并发送消息
		for _, user := range s.OnlineUsers {
			user.MsgChan <- msg
			// user.conn.Write([]byte(msg))
		}
	}
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
	go s.BroadCast()

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
