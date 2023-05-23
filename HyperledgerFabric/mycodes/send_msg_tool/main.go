package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	// 定义命令行flag
	ip := flag.String("ip", "localhost", "IP address to connect to")
	port := flag.Int("p", 0, "Port to connect to")
	msg := flag.String("m", "", "Message to send")

	// 解析命令行flag
	flag.Parse()

	// 连接到指定IP和端口
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		fmt.Println(err)
		return
	}

	// 发送消息
	n, err := fmt.Fprint(conn, *msg)
	if err != nil {
		fmt.Errorf("send failed.")
	} else {
		fmt.Println("total %d bytes sent.", n)
	}
}
