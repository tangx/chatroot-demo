package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	Ipaddr string
	Port   int

	// 在线用户
	userlock    sync.RWMutex
	OnlineUsers map[string]*User

	// 广播管道
	BroadCastChan chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ipaddr:        ip,
		Port:          port,
		OnlineUsers:   make(map[string]*User),
		BroadCastChan: make(chan string),
	}

	return server
}

// Handler 用户连接初始化
func (s *Server) Handler(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	logrus.Infof("%s connected\n", addr)

	// 1. 创建用户
	user := NewUser(conn, s)

	// 2. 用户注册
	user.Online()

	// 4. 用户接受消息
	go user.ListenMessage()

	// 6. 用户存活检测
	isAlive := make(chan bool)
	// 5. 服务器接受用户消息并广播
	go func() {
		buf := make([]byte, 4096)
		// buf := []byte{} // ????
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err == io.EOF {
				// 本条消息读取完成
				return
			}

			// 发送用户消息, 不发送 '\n'
			msg := buf[:n-1]
			if len(msg) == 0 {
				// 不能发送空消息
				continue
			}
			user.DoMessage(string(msg))

			// 6.1 刷新计时器
			isAlive <- true
		}
	}()

	for {
		select {
		case <-isAlive:
			// 如果这里为空， 则或继续执行随后的 case 。
			// 语法小技巧
		case <-time.After(time.Second * 10 * 60):
			// time.After 是计时器， 返回一个通道
			// 如果执行 time.After 则刷新计时器
			// 当 case 获取到数据，表示计时器超时，执行下线操作； 否则阻塞。
			user.Kickoff()
			return
		}
	}
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
	msg = fmt.Sprintf("[%s]%s: %s\n", user.Addr, user.Name, msg)
	s.BroadCastChan <- msg

}

// Start 启动服务
func (s *Server) Start() {
	// 1. 监听
	listen := fmt.Sprintf("%s:%d", s.Ipaddr, s.Port)
	listener, err := net.Listen("tcp", listen)
	if err != nil {
		logrus.Fatalf("Server start failed: %v", err)
	}
	// n. defer 关闭
	defer listener.Close()

	logrus.Infof("服务启动成功， 监听地址 %s", listen)

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
