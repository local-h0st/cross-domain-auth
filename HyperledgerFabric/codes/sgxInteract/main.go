// must run in sgx env!
// 监听端口，收信息，发信息，that's it.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"msgs"
	"myrsa"
	"net"
	"os"
)

// 应该要互斥锁，特别是写操作
type domainInfoRecord struct {
	Domain    string
	PASPubkey string
	BlackList []string
}

var domainInfo []domainInfoRecord
var serverPubkey map[string]string // serverid-pubkey
var PRVKEY, PUBKEY []byte
var selfID string

const selfAddr string = "localhost:55555"
const servingPort string = ":55555"

func main() {
	fmt.Println("[sgxInteract main] running in enclave env.")
	// 生成公私钥，获取环境变量ID
	PRVKEY, PUBKEY = myrsa.GenRsaKey()
	fmt.Println("[main] PUBKEY: ", string(PUBKEY))
	selfID = os.Getenv("SERVERID")
	// 监听
	ln, err := net.Listen("tcp", servingPort)
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
	fmt.Println("[parseMessage] cipher received: ", string(msg))
	handleMsg(msg)
}
func handleMsg(cipher []byte) {
	basic_msg := msgs.BasicMsg{}
	if json.Unmarshal(decryptMsg(cipher), &basic_msg) != nil {
		fmt.Println("[handleMsg] basic msg json unmarshal failed.")
		return
	}
	// 两类消息，一类应该是不需要签名的，另一类的需要签名的
	if basic_msg.Signature == nil {
		switch basic_msg.Method {
		case "addServerPubkey":
			if serverPubkey[basic_msg.SenderID] != "" {
				fmt.Println("[handleMsg] pubkey already exists..")
			} else {
				addServerPubkey(basic_msg.Content)
			}
		case "needPubkey":
			needPubkey(basic_msg.Content)
		default:
			fmt.Println("[handleMsg] unknown method.")
			return
		}
	} else {
		if !basic_msg.VerifySign([]byte(serverPubkey[basic_msg.SenderID])) {
			fmt.Println("[handleMsg] msg signature invalid.")
		}
		switch basic_msg.Method {
		case "verifyID":
			verifyID(basic_msg.Content)
		case "updateServerPubkey":
			updateServerPubkey(basic_msg.Content)
		default:
			fmt.Println("[handleMsg] unknown method.")
			return
		}
	}

}

func verifyID(jsonmsg []byte) {
	msg := msgs.VerifyMsg{}
	if json.Unmarshal(jsonmsg, &msg) != nil {
		fmt.Println("[verifyID] verify msg json unmarshal failed.")
		return
	}
	if msg.UpdateFlag {
		fmt.Println("[verifyID] black list out of date, now update black list, please verify later.")
		// 仅向域PAS发送请求
		pas_pubkey := serverPubkey[msg.DomainPasID]
		if pas_pubkey == "" {
			fmt.Println("[verifyID] Pas pubkey unknown. please update pas pubkey.")
		} else {
			request_msg := msgs.BasicMsg{
				Method:    "blacklistNeedUpdate",
				SenderID:  selfID,
				Content:   []byte(selfAddr),
				Signature: nil,
			}
			request_msg.GenSign(PRVKEY)
			request_msg_json, _ := json.Marshal(request_msg)
			cipher := encryptMsg(request_msg_json, []byte(pas_pubkey))
			if sendMsg(msg.DomainPasAddr, cipher) != nil {
				fmt.Println("[verifyID] update request send failed.")
			}
		}
		return
	}

	verify_result := msgs.BasicMsg{
		Method:    "verifyResult",
		SenderID:  selfID,
		Content:   []byte("valid"),
		Signature: nil,
	}
	true_id := decryptMsg([]byte(msg.CipherID))
	for _, domain_record := range domainInfo {
		if domain_record.Domain == msg.Domain {
			for _, black_name := range domain_record.BlackList {
				if string(true_id) == black_name {
					verify_result.Content = []byte("invalid")
					break
				}
			}
			break
		}
	}
	verify_result.GenSign(PRVKEY)
	jsonmsg, _ = json.Marshal(verify_result)
	server_pubkey := serverPubkey[msg.ServerID]
	if server_pubkey == "" {
		fmt.Println("[verifyID] server pubkey unknown. please update.")
		return
	}

	if sendMsg(msg.ServerAddr, encryptMsg(jsonmsg, []byte(serverPubkey[msg.ServerID]))) != nil {
		fmt.Println("[verifyID] verify result send failed..")
	}
	return
}
func addServerPubkey(jsonmsg []byte) {
	msg := msgs.AddServerPubkeyMsg{}
	if json.Unmarshal(jsonmsg, &msg) != nil {
		fmt.Println("[addServerPubkey] msg json unmarshal failed.")
		return
	}
	serverPubkey[msg.ServerID] = string(msg.ServerPubkey)
	return
}
func updateServerPubkey(jsonmsg []byte) {
	msg := msgs.UpdateServerPubkeyMsg{}
	if json.Unmarshal(jsonmsg, &msg) != nil {
		fmt.Println("[updateServerPubkey] update server info msg json unmarshal failed.")
		return
	}
	// 先核验旧的pk-sk的签名是否正确，随后再更新pk
	if !msg.VerifySign([]byte(serverPubkey[msg.ServerID])) {
		fmt.Println("[updateServerPubkey] signature invalid, reject to update pubkey.")
		return
	}
	serverPubkey[msg.ServerID] = string(msg.ServerNewPubkey)
	return
}
func needPubkey(jsonmsg []byte) {
	msg := msgs.NeedPubkey{}
	if json.Unmarshal(jsonmsg, &msg) != nil {
		fmt.Println("[needPubkey] msg json unmarshal failed.")
		return
	}
	if serverPubkey[msg.SenderID] == "" {
		serverPubkey[msg.SenderID] = string(msg.SenderPubkey)
	}
	// TODO，对方没有自己的公钥
	add_key_msg := msgs.AddServerPubkeyMsg{
		ServerID:     selfID,
		ServerPubkey: PUBKEY,
	}
	msg_to_send := msgs.BasicMsg{
		Method:    "addServerPubkey",
		SenderID:  selfID,
		Content:   nil,
		Signature: nil,
	}
	msg_to_send.Content, _ = json.Marshal(add_key_msg)
	// add key 消息不需要签名，也不能
	jsonmsg, _ = json.Marshal(msg_to_send)
	if sendMsg(msg.SenderAddr, encryptMsg(jsonmsg, msg.SenderPubkey)) != nil {
		fmt.Println("[needPubkey] pubkey send failed.")
	}
}

// 辅助函数
func sendMsg(addr string, data []byte) error {
	// 连接到指定IP和端口
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	// 发送消息
	n, err := fmt.Fprint(conn, data)
	if err != nil {
		return fmt.Errorf("[sendMsg] send data to %s failed.", addr)
	} else {
		fmt.Println("[sendMsg] total %d bytes sent.", n)
		return nil
	}
}
func decryptMsg(cipher []byte) []byte {
	// 肯定是拿自己的prvkey解密
	return myrsa.RsaDecrypt(cipher, PRVKEY)
}
func encryptMsg(text []byte, pubkey []byte) []byte {
	return myrsa.RsaEncrypt(text, pubkey)
}
