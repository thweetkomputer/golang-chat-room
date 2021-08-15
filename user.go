package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
}

// NewUser 创建一个用户
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User {
		Name : userAddr,
		Addr: userAddr,
		C : make(chan string),
		conn: conn,
	}

	// 启动监听当前user
	go user.ListenMessage()

	return user
}

// ListenMessage 监听当前User channel 的方法， 一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("ListenMessage().write err")
			return
		}
	}
}
