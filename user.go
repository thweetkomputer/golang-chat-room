package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// NewUser 创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// 启动监听当前user
	go user.ListenMessage()

	return user
}

// Online 用户的上线业务
func (u *User) Online() {
	// 用户上线，加入OnlineMap
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// 广播当前用户上线消息
	u.server.BroadCast(u, "已上线")
}

// Offline 用户的下线业务
func (u *User) Offline() {
	// 用户下线，移出OnlineMap
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()




	// 广播当前用户上线消息
	u.server.BroadCast(u, "下线")
}

// SendMsg 给客户端发送消息
func (u *User) SendMsg(msg string) {
	//_, err := u.conn.Write([]byte(msg))
	//if err != nil {
	//
	//	return
	//}
	u.C <- msg
}

// DoMessage 用户处理消息的业务
func (u *User) DoMessage(msg string) {
	// 查询在线用户
	if msg == "who" {
		u.server.mapLock.Lock()

		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线。"
			u.SendMsg(onlineMsg)
		}

		u.server.mapLock.Unlock()

		return
	}

	// 修改用户名
	if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式 rename|new name
		newName := strings.Split(msg, "|")[1]

		// 判断newName是否存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("当前用户名已被使用")
			return
		}
		u.server.mapLock.Lock()
		delete(u.server.OnlineMap, u.Name)
		u.server.OnlineMap[newName] = u
		u.server.mapLock.Unlock()

		u.Name = newName
		u.SendMsg("您已更新用户名为 " + u.Name + " 。")
		return
	}

	// 私聊 to|张三｜消息内容
	if len(msg) > 4 && msg[:3] == "to|" {
		params := strings.Split(msg, "|")
		// 1 获取用户名
		remoteName := params[1]
		if remoteName == "" {
			u.SendMsg("消息格式不正确")
			return
		}

		// 2 得到User对象
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMsg("用户不存在")
			return
		}

		// 3 获取消息内容，通过User对象发送消息内容
		content := params[2]
		if content == "" {
			u.SendMsg("消息不能为空")
			return
		}
		remoteUser.SendMsg(u.Name + " 对您说： " + content)
		return
	}

	u.server.BroadCast(u, msg)
}

// ListenMessage 监听当前User channel 的方法， 一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	for msg := range u.C{
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println(u.Name + "\nmsg:" + msg + " \nListenMessage().write " + err.Error())
			return
		}
	}
}
