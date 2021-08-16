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

func (c *Client) menu() bool {
	var f int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	_, err := fmt.Scanln(&f)
	if err != nil {
		fmt.Println("Scan err:", err)
		return false
	}

	if f < 0 || f > 3 {
		fmt.Println("请输入合法范围数字")
		return false
	}

	c.flag = f
	return true
}

func (c *Client) SendMsg(msg string) bool {
	_, err := c.conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("client send msg err:", err)
		return false
	}
	return true
}

func (c *Client) PublicChat() {
	// 提示用户输入消息
	var chatMsg string

	fmt.Println("请输入聊天内容，exit退出")
	var err error
	for _, err = fmt.Scanln(&chatMsg); chatMsg != "exit"; _, err = fmt.Scanln(&chatMsg){
		if len(chatMsg) == 0 || err != nil{
			continue
		}
		c.SendMsg(chatMsg)
	}
	fmt.Println("exit PC")
}

func (c *Client) UpdateName() bool {
	fmt.Print("请输入新的用户名：")
	_, err1 := fmt.Scanln(&c.Name)
	if err1 != nil {
		fmt.Println("Scan err:", err1)
		return false
	}

	sendMsg := "rename|" + c.Name

	return c.SendMsg(sendMsg)
}

func (c *Client) Run() {
	for c.flag != 0 {
		// 根据不同模式处理不同的业务
		for !c.menu() {
		}
		switch c.flag {
		case 1:
			// 公聊模式
			fmt.Println(">>>>>>进入公聊模式<<<<<<")
			c.PublicChat()
		case 2:
			// 私聊模式
			fmt.Println(">>>>>>进入私聊模式<<<<<<")
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
