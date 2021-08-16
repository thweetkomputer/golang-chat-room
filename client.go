package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // 客户端模式
}

// NewClient 创建客户端对象
func NewClient(serverIp string, serverPort int) *Client {
	// 创建对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	// 链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))

	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	// 返回对象
	return client
}

// DealResponse 处理server回应的消息，直接显示到标准输出
func (c *Client) DealResponse() {
	// 一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	_, err := io.Copy(os.Stdout, c.conn)
	if err != nil {
		fmt.Println("stdout err:", err)
		return
	}

}

// Menu 显示菜单
func (c *Client) Menu() bool {
	var f int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	_, err := fmt.Scanln(&f)
	for err != nil {
		_, err = fmt.Scanln(&f)
	}

	if f < 0 || f > 3 {
		fmt.Println("请输入合法范围数字")
		return false
	}

	c.flag = f
	return true
}

// SendMsg 发送信息向服务端
func (c *Client) SendMsg(msg string) bool {
	_, err := c.conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("client send msg err:", err)
		return false
	}
	return true
}

// PublicChat 公聊业务
func (c *Client) PublicChat() {
	// 提示用户输入消息
	var chatMsg string

	fmt.Println("请输入聊天内容，exit退出")
	var err error
	for _, err = fmt.Scanln(&chatMsg); chatMsg != "exit"; _, err = fmt.Scanln(&chatMsg) {
		if len(chatMsg) == 0 || err != nil {
			continue
		}
		c.SendMsg(chatMsg)
	}
	fmt.Println("已退出公聊")
}

// SelectUsers 查询在线用户
func (c *Client) SelectUsers() {
	c.SendMsg("who")
}

// PrivateChat 私聊业务
func (c *Client) PrivateChat() {
	var remoteName, chatMsg string

	c.SelectUsers()
	fmt.Println("请输入聊天对象（用户名）,exit退出.")

	var err error
	for _, err = fmt.Scanln(&remoteName); remoteName != "exit"; _, err = fmt.Scanln(&remoteName) {
		if len(remoteName) == 0 || err != nil {
			continue
		}
		fmt.Println("已经进入与 " + remoteName + " 的私聊 ,exit退出.")
		var err1 error
		for _, err1 = fmt.Scanln(&chatMsg); chatMsg != "exit"; _, err1 = fmt.Scanln(&chatMsg) {
			if len(chatMsg) == 0 || err1 != nil {
				continue
			}
			c.SendMsg("to|" + remoteName + "|" + chatMsg)
		}
	}
}

// UpdateName 修改用户名
func (c *Client) UpdateName() bool {
	fmt.Print("请输入新的用户名(长度1-18位)：")
	_, err1 := fmt.Scanln(&c.Name)
	if len(c.Name) > 18 || len(c.Name) < 1 {
		fmt.Println("新用户名长度不合法")
	}
	for err1 != nil || len(c.Name) > 18 || len(c.Name) < 1{
		_, err1 = fmt.Scanln(&c.Name)
		//fmt.Println(len(c.Name))
		if len(c.Name) > 18 || len(c.Name) < 1 {
			fmt.Println("新用户名长度不合法")
		}
	}

	sendMsg := "rename|" + c.Name

	return c.SendMsg(sendMsg)
}

// Run 运行客户端业务
func (c *Client) Run() {
	for c.flag != 0 {
		// 根据不同模式处理不同的业务
		for !c.Menu() {
		}
		switch c.flag {
		case 1:
			// 公聊模式
			fmt.Println(">>>>>>进入公聊模式<<<<<<")
			c.PublicChat()
		case 2:
			// 私聊模式
			fmt.Println(">>>>>>进入私聊模式<<<<<<")
			c.PrivateChat()
		case 3:
			// 更新用户名
			c.UpdateName()
		}
		//fmt.Println("flag=", c.flag)
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("服务器链接失败。。。")
		return
	}

	fmt.Println(">>>>>>>服务器链接成功<<<<<<<")

	// 单独开启一个goroutine处理server的消息
	go client.DealResponse()

	// 启动业务
	client.Run()
}
