package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// 生成pk和sk
	// 需要公布自己的pk，直接上链
	ln, err := net.Listen("tcp", ":54321")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		go parseMessage(conn)
	}
}

func parseMessage(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err == io.EOF {
		conn.Close()
		return
	} else if err != nil {
		conn.Close()
		return
	}
	msg := buf[:n]
	fmt.Println("msg received: ", string(msg))
	handleMsg(string(msg))
}

var node_id string = os.Getenv("NODE_ID")
var node_sk []byte
var node_pk []byte

func handleMsg(cipher string) error {
	// 事实上收到的是加密后的密文，需要用node_sk解密，这里暂时没写
	msg_text := cipher
	type message struct {
		PID              string
		P                string // // 门限签名技术的那个点，我也不知道用什么格式存储
		Tag              string // node_id
		ID_cipher        string
		PK_device2domain string // 不知道公钥是这哪个类型。另外是建议存到数据库，还是存到内存算了？
	}
	var msg message
	err := json.Unmarshal([]byte(msg_text), &msg)
	if err != nil {
		return fmt.Errorf("In func HandleMsgForPseudo(): json unmarshal failed.")
	}
	// 查询pid是否存在

	return nil
}
