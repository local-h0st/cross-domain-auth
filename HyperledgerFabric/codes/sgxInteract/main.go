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
type domainInfoRecord struct {
	Domain    string
	Pubkey    string
	BlackList []string
}

var domainInfo []domainInfoRecord

var serverPubkey map[string]string // serverid-pubkey

var prvKey, pubKey string

type basicMsg struct {
	Method     string
	CipherText string
}

func handleMsg(msg []byte) {
	basic_msg := basicMsg{}
	err := json.Unmarshal(msg, &basic_msg)
	if err != nil {
		fmt.Println("[handleMsg] basic msg json unmarshal failed.")
		return
	}
	switch basic_msg.Method {
	case "verifyID":
		verifyID(basic_msg.CipherText)
	case "updateServerInfo":
		updateServerInfo(basic_msg.CipherText)
	default:
		return
	}

}

type updateServerInfoMsg struct {
	ServerID        string
	serverPubkey    string
	ServerSignature string
}

func updateServerInfo(cipher string) {
	server_info := updateServerInfoMsg{}
	err := json.Unmarshal([]byte(decryptMsg(cipher, prvKey)), &server_info)
	if err != nil {
		fmt.Println("[updateServerInfo] update server info msg json unmarshal failed.")
	} else if serverPubkey[server_info.ServerID] == "" {
		// 初次
		serverPubkey[server_info.ServerID] = server_info.serverPubkey
	} else if server_info.VerifySign(serverPubkey[server_info.ServerID]) {
		serverPubkey[server_info.ServerID] = server_info.serverPubkey
	} else {
		fmt.Println("[updateServerInfo] sign invalid!")
	}

}

type verifyMsg struct {
	UpdateFlag      bool
	DomainPasAddr   string // ip:port形式，例如localhost:6666
	CipherID        string
	Domain          string
	ServerID        string // 来自哪个server
	ServerAddr      string
	ServerSignature string
}

func verifyID(cipher string) {
	verify_msg := verifyMsg{}
	err := json.Unmarshal([]byte(decryptMsg(cipher, prvKey)), &verify_msg)
	if err != nil {
		fmt.Println("[verifyID] verify msg json unmarshal failed.")
	} else if !verify_msg.VerifySign(serverPubkey[verify_msg.ServerID]) {
		fmt.Println("[verifyID] Server signature invalid!")
	} else if verify_msg.UpdateFlag {
		fmt.Println("[verifyID] black list out of date, please verify later.")
		err := updateBlacklist(verify_msg.DomainPasAddr)
		if err != nil {
			fmt.Println("[verifyID] update black list request sent, waiting for reply...")
		}
	} else {
		type resultMsg struct {
			Valid     bool
			Signature string
		}
		result_msg := resultMsg{Valid: true}
		id := decryptMsg(verify_msg.CipherID, prvKey)
		for _, domain_record := range domainInfo {
			if domain_record.Domain == verify_msg.Domain {
				for _, black_name := range domain_record.BlackList {
					if id == black_name {
						result_msg.Valid = false
						break
					}
				}
				break
			}
		}
		jsonstr, _ := json.Marshal(result_msg)
		result_msg.Signature = signMsg(string(jsonstr))
		jsonstr, _ = json.Marshal(result_msg)
		err := sendMsg(verify_msg.ServerAddr, string(jsonstr))
		if err != nil {
			fmt.Printf("[verifyID] send result to server %s failed.", verify_msg.ServerAddr)
		} else {
			fmt.Printf("[verifyID] send result to server %s successfully.", verify_msg.ServerAddr)
		}
	}
}

type updateBlacklistMsg struct {
	Method    string
	Signature string
}

func updateBlacklist(addr string) error {
	// 仅仅发送请求
	msg := updateBlacklistMsg{
		Method:    "updateBlacklist",
		Signature: "",
	}
	jsonstr, _ := json.Marshal(msg)
	msg.Signature = signMsg(string(jsonstr))
	jsonstr, _ = json.Marshal(msg)
	return sendMsg(addr, string(jsonstr))
}

// TODO

func (m updateServerInfoMsg) VerifySign(pubkey string) bool {
	return true
}
func (m verifyMsg) VerifySign(pubkey string) bool {
	return true
}
func decryptMsg(cipher string, prvkey string) string {
	return cipher
}
func signMsg(msg string) string {
	return msg
}

// /////////////////////////////////////
func main() {
	fmt.Println("[sgxInteract] running in enclave env. listening on", serving_port)
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
func sendMsg(addr string, msg string) error {
	// 连接到指定IP和端口
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	// 发送消息
	n, err := fmt.Fprint(conn, msg)
	if err != nil {
		return fmt.Errorf("[sendMsg] send msg to %s failed.", addr)
	} else {
		fmt.Println("[sendMsg] total %d bytes sent.", n)
		return nil
	}
}
