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
	sharedconfigs "sharedConfigs"
)

// 应该要互斥锁，特别是写操作

var domainInfo []msgs.DomainInfoRecord
var serverPubkey map[string]string // serverid-pubkey, all kinds of server
var PRVKEY, PUBKEY, vsPubkey []byte

func main() {
	fmt.Println("[sgxInteract main] running in enclave env.")
	PUBKEY = []byte(sharedconfigs.EnclavePubkey)
	PRVKEY = []byte(sharedconfigs.EnclavePrvkey)
	vsPubkey = []byte(sharedconfigs.ServerPubkey)
	serverPubkey = make(map[string]string, 0)
	// 监听
	ln, err := net.Listen("tcp", sharedconfigs.EnclavePort)
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
		return nil, err
	} else if err != nil {
		return nil, err
	}
	msg := buf[:n]
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
				fmt.Println("[handleMsg] server pubkey already exists..")
			} else {
				addServerPubkey(basic_msg.SenderID, basic_msg.Content)
			}
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
				SenderID:  sharedconfigs.NodeID,
				Content:   []byte(sharedconfigs.EnclaveAddr),
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
		SenderID:  sharedconfigs.NodeID,
		Content:   nil,
		Signature: nil,
	}
	result_msg := msgs.VerifyResultMsg{
		PID:                  msg.PID,
		Result:               "valid",
		PubkeyDeviceToDomain: msg.PubkeyDeviceToDomain,
	}
	true_id := decryptMsg(msg.CipherID)
	fmt.Println("true id: ", string(true_id))
	for _, domain_record := range domainInfo {
		if domain_record.Domain == msg.Domain {
			fmt.Println("find domain", domain_record.Domain)
			for _, black_name := range domain_record.BlackList {
				fmt.Println("comparing: ", black_name)
				if string(true_id) == black_name {
					result_msg.Result = "invalid"
					break
				}
			}
			break
		}
	}
	// 有可能lack domain info，需要同步信息
	verify_result.Content, _ = json.Marshal(result_msg)
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
}
func addServerPubkey(SenderID string, jsonmsg []byte) {
	fmt.Println("exec addServerPubkey..")
	msg := msgs.AddServerPubkeyMsg{}
	if json.Unmarshal(jsonmsg, &msg) != nil {
		fmt.Println("[addServerPubkey] msg json unmarshal failed.")
		return
	}
	serverPubkey[SenderID] = string(msg.ServerPubkey)
	fmt.Println(serverPubkey) // TODO for test
}
func sendBackPubkey() {

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
