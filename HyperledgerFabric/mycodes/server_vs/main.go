package main

import (
	"fmt"
	"io"
	"net"
)

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
	fmt.Println(string(msg))
}

func main() {
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

		// 每接收到一个连接就启动一个goroutine处理
		go handleConn(conn)
	}
}
