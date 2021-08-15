package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的chan
	Message chan string
}

// NewServer server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// ListenMessage 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线user
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		// 将msg发送给全部的在线user
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// BroadCast 广播消息
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	s.Message <- sendMsg

}

func (s *Server) Handler(conn net.Conn) {
	// 当前链接的业务
	//fmt.Println("connection build success")

	user := NewUser(conn)

	// 用户上线，加入OnlineMap
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	// 广播当前用户上线消息
	s.BroadCast(user, "已上线")

	// 接受客户端发来的信息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.BroadCast(user, "下线")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// 提取用户的消息
			msg := string(buf[:n-1])

			// 广播消息
			s.BroadCast(user, msg)
		}
	}()

	// 当前handler阻塞
	select {}

}

// Start 启动服务器的接口
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen
	defer func(listen net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("net.Listen.close err:", err)
		}
	}(listener)

	// 启动监听Message的goroutine
	go s.ListenMessage()

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

}
