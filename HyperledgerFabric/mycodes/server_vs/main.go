package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"serverVS/tools/rsatool"
)

func main() {
	rsatool.UnitTest()
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

func handleMsg(cipher string) {
	// 事实上收到的是加密后的密文，需要用node_sk解密

}
