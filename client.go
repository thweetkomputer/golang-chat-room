package main

import (
	"flag"
	"fmt"
	"net"
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

func (client *Client) menu() bool {
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

	client.flag = f
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		// 根据不同模式处理不同的业务
		for !client.menu() {
		}
		switch client.flag {
		case 1:
			// 公聊模式
			fmt.Println(">>>>>>进入公聊模式<<<<<<")
		case 2:
			// 私聊模式
			fmt.Println(">>>>>>进入私聊模式<<<<<<")
		case 3:
			// 更新用户名
			fmt.Println(">>>>>>进入用户名模式<<<<<<")
		}
		fmt.Println("flag=", client.flag)
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

	// 启动业务
	client.Run()
}
