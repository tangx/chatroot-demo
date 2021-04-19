package server

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type Server struct {
	Ipaddr string
	Port   uint16
}

func NewServer(ip string, port uint16) *Server {
	server := &Server{
		Ipaddr: ip,
		Port:   port,
	}

	return server
}

func (s *Server) Handler(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	fmt.Printf("%s connected \n", addr)
}

func (s *Server) Start() {
	// 1. 监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ipaddr, s.Port))
	if err != nil {
		logrus.Fatalf("Server start failed: %v", err)

	}

	// 2. 关闭
	defer listener.Close()

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
