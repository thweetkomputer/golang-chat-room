package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (s *Server) Handler(conn net.Conn) {
	// 当前链接的业务
	fmt.Println("connection build success")
}

// Start 启动服务器的接口
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer func(listen net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("net.Listen.close err:", err)
		}
	}(listener)

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err", err)
			continue
		}

		// do handler
		go s.Handler(conn)
	}

	// accept

	// do handler

	// close listen
}
