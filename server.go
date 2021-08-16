package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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

	user := NewUser(conn, s)

	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接受客户端发来的信息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// 提取用户的消息
			msg := string(buf[:n-1])

			// 用户对msg进行处理
			user.DoMessage(msg)

			// 用户的任意消息，代表当前用户是一个活跃的
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
			// 重置计时器
			// 不做任何事情，为了激活select， 应该重置定时器
		case <-time.After(time.Minute * 5):
			// 已经超时
			// 讲当前的User强制的关闭

			_, err := user.conn.Write([]byte("你被踢了\n"))
			if err != nil {
				fmt.Println(err)
				return
			}

			// 销毁用的资源
			close(user.C)

			// 关闭链接
			err1 := user.conn.Close()
			if err1 != nil {
				fmt.Println("conn.Close err:", err1)
				return
			}

			// 退出当前的Handler
			return // runtime.Exit()
		}
	}

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
