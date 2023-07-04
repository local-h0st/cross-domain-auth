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

var blacklists []msgs.BlacklistRecord
var PRVKEY, PUBKEY, vsPubkey []byte

func main() {
	fmt.Println("[sgxInteract main] running in enclave env.")
	PUBKEY = []byte(sharedconfigs.EnclavePubkey)
	PRVKEY = []byte(sharedconfigs.EnclavePrvkey)
	vsPubkey = []byte(sharedconfigs.ServerPubkey)
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
	if json.Unmarshal(myrsa.DecryptMsg(cipher, []byte(sharedconfigs.EnclavePrvkey)), &basic_msg) != nil {
		fmt.Println("[handleMsg] basic msg json unmarshal failed.")
		return
	}
	if !basic_msg.VerifySign(vsPubkey) {
		fmt.Println("[handleMsg] msg signature invalid.")
		return
	}
	switch basic_msg.Method {
	case "verifyID":
		verifyID(basic_msg.Content)
	case "updateBlacklist": // TODO 要改
		updateBlacklist(basic_msg.Content)
		fmt.Println("up-to-date blacklists:", blacklists)
	case "testingConnection":
		testmsg := msgs.BasicMsg{
			Method:   "testingConnection",
			SenderID: sharedconfigs.NodeEnclaveID,
		}
		testmsg.GenSign(PRVKEY)
		teststr, _ := json.Marshal(testmsg)
		sendMsg(sharedconfigs.ServerAddr, string(myrsa.EncryptMsg(teststr, vsPubkey)))
	default:
		fmt.Println(basic_msg.Method, ": unknown method.")
		return
	}
}

func verifyID(jsonmsg []byte) {
	msg := msgs.VerifyMsg{}
	if json.Unmarshal(jsonmsg, &msg) != nil {
		fmt.Println("[verifyID] verify msg json unmarshal failed.")
		return
	}
	verify_result := msgs.BasicMsg{
		Method:    "verifyResult",
		SenderID:  sharedconfigs.NodeEnclaveID,
		Content:   nil,
		Signature: nil,
	}
	result_msg := msgs.VerifyResultMsg{
		PID:        msg.PID,
		DestDomain: msg.DestDomain,
		Result:     "valid",
	}
	true_id := string(myrsa.DecryptMsg(msg.CipherID, PRVKEY))
	for _, domain_record := range blacklists {
		if domain_record.Domain == msg.DestDomain {
			for _, black_name := range domain_record.BlackList {
				if true_id == black_name {
					result_msg.Result = "invalid"
					break
				}
			}
			break
		}
	}
	// TODO 有可能lack domain info，需要同步信息
	verify_result.Content, _ = json.Marshal(result_msg)
	verify_result.GenSign(PRVKEY)
	jsonmsg, _ = json.Marshal(verify_result)
	if sendMsg(sharedconfigs.ServerAddr, string(myrsa.EncryptMsg(jsonmsg, []byte(vsPubkey)))) != nil {
		fmt.Println("[verifyID] verify result send failed..")
	}
}

func updateBlacklist(jsonmsg []byte) {
	record := msgs.BlacklistRecord{}
	if json.Unmarshal(jsonmsg, &record) != nil {
		fmt.Println("[updateBlacklist] json unmarshal failed.")
		return
	}
	index := -1
	for k := range blacklists {
		if blacklists[k].Domain == record.Domain {
			index = k
			break
		}
	}
	if index == -1 {
		blacklists = append(blacklists, record)
	} else {
		blacklists[index] = record
	}
}

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
		fmt.Println("[sendMsg] total ", n, " bytes sent.")
		return nil
	}
}
