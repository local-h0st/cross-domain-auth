// must run in sgx env!
// 监听端口，收信息，发信息，that's it.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

const serving_port string = ":55555"

// 应该要互斥锁，特别是写操作
var black_list []string

type basicMsg struct {
	method     string
	cipherText string
}
type verifyMsg struct {
	challenge  int
	updateFlag bool
	cipherID   string
	domain     string
	serverID   string // 来自哪个server
	signature  string
}

func handleMsg(msg []byte) {
	basic_msg := basicMsg{}
	err := json.Unmarshal(msg, &basic_msg)
	if err != nil {
		fmt.Println("[handleMsg] basic msg json unmarshal failed.")
		return
	}
	if basic_msg.method == "verifyID" {
		// TODO
	} else {
		// 其他分支还没写
		return
	}
}

func verifyID(cipher string) (string, error) {
	type resultMsg struct {
		challenge int
	}
	msg_text := decryptMsg(cipher)
	verify_msg := verifyMsg{}
	err := json.Unmarshal([]byte(msg_text), &verify_msg)
	if err != nil {
		fmt.Println("[verifyID] verify msg json unmarshal failed.")
		return false, err
	}

}

// TODO 没写，
func decryptMsg(msg string) string {
	return msg
}

// //////////////////////////////////////
func main() {
	fmt.Println("[interactSGX] running in enclave env. listening on", serving_port)
	ln, err := net.Listen("tcp", serving_port)
	if err != nil {
		panic(err)
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
	fmt.Println("[parseMessage] msg received: ", string(msg))
	handleMsg(msg)
}
