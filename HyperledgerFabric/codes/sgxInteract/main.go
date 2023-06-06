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

var domainInfo []msgs.DomainInfoRecord
var serverPubkey map[string]string // serverid-, all kinds of server
var PRVKEY, PUBKEY []byte
var selfID string

const selfAddr string = "localhost:55555"
const servingPort string = ":55555"

func main() {
	fmt.Println("[sgxInteract main] running in enclave env.")
	// PRVKEY, PUBKEY = myrsa.GenRsaKey()	// TODO 为了方便调试暂时先指定了PRVKEY和PUBKEY，实际生成时需要rand
	PUBKEY = []byte("-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCj8hW+keEOHHHLV/7BRO7I0j7a\nXAfxTvkiM8Qyex+aMQ7Ny+cavF4mWlJmdmGo9K3jHFH3LyEd2JuGPh5T0ad/O76C\nor+hX+RvgXkg0HS3MEQIwmzmNg57RSaNxzlJatXEfjpRJJ5Nc+dyA6hpYzaNj9LY\nKex5gvsGFpBMwQZyVwIDAQAB\n-----END PUBLIC KEY-----\n")
	PRVKEY = []byte("-----BEGIN RSA PRIVATE KEY-----\nMIICXgIBAAKBgQCj8hW+keEOHHHLV/7BRO7I0j7aXAfxTvkiM8Qyex+aMQ7Ny+ca\nvF4mWlJmdmGo9K3jHFH3LyEd2JuGPh5T0ad/O76Cor+hX+RvgXkg0HS3MEQIwmzm\nNg57RSaNxzlJatXEfjpRJJ5Nc+dyA6hpYzaNj9LYKex5gvsGFpBMwQZyVwIDAQAB\nAoGBAKK4SZq/Qaf21X8lFIaRO4t5GcczJvL8Fkw7IxWTnNc2r+HU6slfgvcAGN73\nypCeYeSTnEsBrRXpgtun1gQNh/cqvnJU9uCpY/PVuk14vE+lYLhKkX/GAWfsmPs+\n2AUWJZeAVCJuixh9E9jnDSz+X8IWNC77cqZq8CIY/5M+nuCBAkEA1TaANVdqNBPL\n4phbC2dVFddCADWHNyaHbOrVy+bxD0x00CqyDupwvc6QMVbqpLydbbdJplZ3g6mk\nhy1HbpAJlwJBAMTYjTMSk/bLwA6D4SFGw1NyVLMOn9I6bnZzB8ryrbBdq6vbx0vV\nvGIsPNA6bKFTgUJb5DepWRMPisL02qS+dUECQQDSR0Uc1pC0uc18Nlycm5XLy5eZ\nUzF/D+3CWrzus16Ngw81+tXPZiI44E9PifQy8p6lBX6KoX6PiLDubJalkUMTAkBB\n55buwIuVl4YH1hOsBnsjFyZQhNbxleqh8cVsJ3ALmnD9qynAtCDMZa8+sDDqmoCu\nbQGtuR8/iHaW60/A1JuBAkEAgKuNrksiWi0h0KTFnassKgeaBUd2MociEK6hmKwI\nwi7kjPNHeaa1MqJMUQLhhYv33m5xuNFxIip2LTcXeJ+/5g==\n-----END RSA PRIVATE KEY-----\n")
	// fmt.Println("[main] PUBKEY ==> ", string(bytes.Replace(PUBKEY, []byte("\n"), []byte("\\n"), -1)), "\n serving at", selfAddr)
	selfID = os.Getenv("SERVERID")
	serverPubkey = make(map[string]string, 0)
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
		msg, err := parseMessage(conn)
		if err != nil {
			break
		}
		handleMsg(msg)
	}
}
func parseMessage(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 32768)
	n, err := conn.Read(buf)
	if err == io.EOF {
		// conn.Close()
		return nil, err
	} else if err != nil {
		// conn.Close()
		return nil, err
	}
	msg := buf[:n]
	// fmt.Println("[parseMessage] cipher received: ", msg)
	fmt.Println("received one msg.")
	return msg, nil
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
				addServerPubkey(basic_msg.SenderID, basic_msg.Content)
			}
		case "needPubkey":
			// needPubkey(basic_msg.Content)
			fmt.Println("hello joker.")
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
			verifyID(basic_msg.SenderID, basic_msg.Content)
		case "updateServerPubkey":
			updateServerPubkey(basic_msg.SenderID, basic_msg.Content)
		case "updateBlackist":
			fmt.Println("original domaininfo:", domainInfo)
			updateBlackist(basic_msg.Content)
			fmt.Println("up-to-date domaininfo:", domainInfo)
		default:
			fmt.Println("[handleMsg] unknown method.")
			return
		}
	}
}

func updateBlackist(jsonmsg []byte) {
	fmt.Println("exec updateBlackist..")
	record := msgs.DomainInfoRecord{}
	if json.Unmarshal(jsonmsg, &record) != nil {
		fmt.Println("[updateBlacklist] json unmarshal failed.")
		return
	}
	for k := range domainInfo {
		if domainInfo[k].Domain == record.Domain {
			domainInfo[k].PasID = record.PasID
			domainInfo[k].BlackList = nil
			domainInfo[k].BlackList = record.BlackList
			return
		}
	}
	domainInfo = append(domainInfo, record)
	return
}

func verifyID(SenderID string, jsonmsg []byte) {
	fmt.Println("exec verifyID..")
	msg := msgs.VerifyMsg{}
	if json.Unmarshal(jsonmsg, &msg) != nil {
		fmt.Println("[verifyID] verify msg json unmarshal failed.")
		return
	}
	if msg.UpdateFlag {
		fmt.Println("[verifyID] black list out of date, now update black list, please verify later.")
		// 仅向域PAS发送请求
		var pas_id string
		for _, domain_record := range domainInfo {
			if domain_record.Domain == msg.Domain {
				pas_id = domain_record.PasID
				break
			}
		}
		if pas_id == "" {
			fmt.Println("[verifyID] Lack domain info, please update first.")
		} else {
			pas_pubkey := serverPubkey[pas_id]
			request_msg := msgs.BasicMsg{
				Method:    "blacklistNeedUpdate",
				SenderID:  selfID,
				Content:   []byte(selfAddr),
				Signature: nil,
			}
			request_msg.GenSign(PRVKEY)
			request_msg_json, _ := json.Marshal(request_msg)
			cipher := encryptMsg(request_msg_json, []byte(pas_pubkey))
			if sendMsg(msg.DomainPasAddr, string(cipher)) != nil {
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
	server_pubkey := serverPubkey[SenderID]
	if server_pubkey == "" {
		fmt.Println("[verifyID] server pubkey unknown. please update.")
		return
	}

	if sendMsg(msg.SenderAddr, string(encryptMsg(jsonmsg, []byte(serverPubkey[SenderID])))) != nil {
		fmt.Println("[verifyID] verify result send failed..")
	}
	return
}
func addServerPubkey(SenderID string, jsonmsg []byte) {
	fmt.Println("exec addServerPubkey..")
	msg := msgs.AddServerPubkeyMsg{}
	if json.Unmarshal(jsonmsg, &msg) != nil {
		fmt.Println("[addServerPubkey] msg json unmarshal failed.")
		return
	}
	serverPubkey[SenderID] = string(msg.ServerPubkey)
	fmt.Println(serverPubkey)
	return
}
func updateServerPubkey(SenderID string, jsonmsg []byte) {
	fmt.Println("exec updateServerPubkey..")
	msg := msgs.UpdateServerPubkeyMsg{}
	if json.Unmarshal(jsonmsg, &msg) != nil {
		fmt.Println("[updateServerPubkey] update server info msg json unmarshal failed.")
		return
	}
	// if !msg.VerifySign([]byte(serverPubkey[msg.ServerID])) {
	// 	fmt.Println("[updateServerPubkey] signature invalid, reject to update pubkey.")
	// 	return
	// }
	serverPubkey[SenderID] = string(msg.ServerNewPubkey)
	fmt.Println(serverPubkey)
	return
}

/*
func needPubkey(jsonmsg []byte) {
	// 小丑函数，对方没我pubkey怎么给我发的加密消息，连decrypt都过不了
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
		// ServerID:     selfID,
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
	if sendMsg(msg.SenderAddr, string(encryptMsg(jsonmsg, msg.SenderPubkey))) != nil {
		fmt.Println("[needPubkey] pubkey send failed.")
	}
}
*/
// 辅助函数
func sendMsg(addr string, data string) error {
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
	return myrsa.DecryptMsg(cipher, PRVKEY)
}
func encryptMsg(text []byte, pubkey []byte) []byte {
	return myrsa.EncryptMsg(text, pubkey)
}
