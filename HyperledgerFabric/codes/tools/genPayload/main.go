package main

import (
	"encoding/json"
	"fmt"
	"io"
	"msgs"
	"myrsa"
	"net"
	"sync"
	"time"
)

// sgxInteract
const ServerPubkey string = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCj8hW+keEOHHHLV/7BRO7I0j7a\nXAfxTvkiM8Qyex+aMQ7Ny+cavF4mWlJmdmGo9K3jHFH3LyEd2JuGPh5T0ad/O76C\nor+hX+RvgXkg0HS3MEQIwmzmNg57RSaNxzlJatXEfjpRJJ5Nc+dyA6hpYzaNj9LY\nKex5gvsGFpBMwQZyVwIDAQAB\n-----END PUBLIC KEY-----\n"
const ServerAddr string = "localhost:55555"

// 自己的
const ServerID string = "genPayloadServer"
const servingPort string = ":55550"
const selfAddr string = "localhost:55550"

var PRVKEY, PUBKEY, OLD_PRVKEY, OLD_PUBKEY []byte
var wg sync.WaitGroup

// ////////////////

func listenForVerifyResult() {
	ln, err := net.Listen("tcp", servingPort)
	if err != nil {
		panic(err)
	}
	conn, err := ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()
	msg, err := func(conn net.Conn) ([]byte, error) {
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
		return msg, nil
	}(conn)
	if err != nil {
		return
	}
	basic_msg := msgs.BasicMsg{}
	if json.Unmarshal(myrsa.DecryptMsg(msg, PRVKEY), &basic_msg) != nil {
		fmt.Println("[Verify Err] basic msg json unmarshal failed.")
		return
	}
	if !basic_msg.VerifySign([]byte(ServerPubkey)) {
		fmt.Println("[Verify Err] invalid sig.")
	} else {
		fmt.Println("[Verify Result]", string(basic_msg.Content))
	}
	wg.Done()
}

// ////////////////
func main() {
	wg.Add(1)
	go listenForVerifyResult()
	PRVKEY, PUBKEY = myrsa.GenRsaKey()
	// 虽然消息有先后顺序，但是sgxInteract那边goroutine并发处理顺序就不一定了，为了有序只能采用sleep
	// sendAddServerPubkeyMsg("testServer")
	sendAddServerPubkeyMsg(ServerID)
	time.Sleep(100 * time.Millisecond)

	OLD_PRVKEY, OLD_PUBKEY = PRVKEY, PUBKEY
	PRVKEY, PUBKEY = myrsa.GenRsaKey()

	sendUpdateServerPubkeyMsg(ServerID)
	time.Sleep(100 * time.Millisecond)

	sendAddServerPubkeyMsg(ServerID)
	time.Sleep(100 * time.Millisecond)

	sendUpdateBlacklistMsg(ServerID)
	time.Sleep(100 * time.Millisecond)

	sendVerifyIDMsg(ServerID)
	time.Sleep(100 * time.Millisecond)

	wg.Wait()
}

func sendVerifyIDMsg(sender_id string) {
	fmt.Println("sendVerifyIDMsg()")
	id_to_verify := []byte("gentleewind")
	basic_msg := msgs.BasicMsg{
		Method:    "verifyID",
		SenderID:  sender_id,
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.VerifyMsg{
		CipherID:      myrsa.EncryptMsg(id_to_verify, []byte(ServerPubkey)),
		Domain:        "domainA",
		SenderAddr:    selfAddr,
		UpdateFlag:    false,
		DomainPasAddr: "",
	})
	basic_msg.GenSign(PRVKEY)
	data, _ := json.Marshal(basic_msg)
	cipher := myrsa.EncryptMsg(data, []byte(ServerPubkey))
	sendMsg(ServerAddr, string(cipher))
}

func sendUpdateBlacklistMsg(sender_id string) {
	fmt.Println("sendUpdateBlacklistMsg()")
	basic_msg := msgs.BasicMsg{
		Method:    "updateBlackist",
		SenderID:  sender_id,
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.DomainInfoRecord{
		Domain:    "domainA",
		PasID:     sender_id,
		BlackList: []string{"gentleewind", "blackname"},
	})
	basic_msg.GenSign(PRVKEY)
	data, _ := json.Marshal(basic_msg)
	cipher := myrsa.EncryptMsg(data, []byte(ServerPubkey))
	sendMsg(ServerAddr, string(cipher))
}

func sendUpdateServerPubkeyMsg(sender_id string) {
	fmt.Println("sendUpdateServerPubkeyMsg()")
	basic_msg := msgs.BasicMsg{
		Method:    "updateServerPubkey",
		SenderID:  sender_id,
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.UpdateServerPubkeyMsg{
		// ServerID:        "genPayloadServer",
		ServerNewPubkey: PUBKEY,
	})
	basic_msg.GenSign(OLD_PRVKEY)
	data, _ := json.Marshal(basic_msg)
	cipher := myrsa.EncryptMsg(data, []byte(ServerPubkey))
	sendMsg(ServerAddr, string(cipher))
}

func sendAddServerPubkeyMsg(sender_id string) {
	fmt.Println("sendAddServerPubkeyMsg()")
	basic_msg := msgs.BasicMsg{
		Method:    "addServerPubkey",
		SenderID:  sender_id,
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.AddServerPubkeyMsg{
		// ServerID:     "genPayloadServer",
		ServerPubkey: PUBKEY,
	})
	data, _ := json.Marshal(basic_msg)
	cipher := myrsa.EncryptMsg(data, []byte(ServerPubkey))
	sendMsg(ServerAddr, string(cipher))
}

func sendMsg(addr string, data string) {
	// 连接到指定IP和端口
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
	}
	// 发送消息
	n, err := fmt.Fprint(conn, data)
	if err != nil {
		fmt.Printf("send to %s failed.", addr)
	} else {
		fmt.Println("msg sent, total", n, "bytes.")
	}
}

func UnitTest() {
	_, PUBKEY := myrsa.GenRsaKey()
	pk := "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDYiWYP6r5aVP7NBQkxyUnr8Pny\n2lZ3NFrWZHF4PQnSvUdlQ43zxzqZ33szAmi+GD2fnAQb2PSgf/+zzFcViqTgHQx0\nHyAHXdf39cVS6REVYtm06Iy/yRIjcWSUNasAg/bD/QKMzuNWmZhTVbHFbahzXXg2\nkKRfsE6+Hr7Ncr0wrQIDAQAB\n-----END PUBLIC KEY-----\n"
	// add key
	basic_msg := msgs.BasicMsg{
		Method:    "addServerPubkey",
		SenderID:  "test",
		Content:   nil,
		Signature: nil,
	}
	add_pubkey_msg := msgs.AddServerPubkeyMsg{
		// ServerID:     "test",
		ServerPubkey: PUBKEY,
	}
	basic_msg.Content, _ = json.Marshal(add_pubkey_msg)
	fmt.Println(string(basic_msg.Content))
	data, _ := json.Marshal(basic_msg)
	fmt.Println(string(data))
	fmt.Println()

	a := msgs.BasicMsg{}
	json.Unmarshal(data, &a)
	if a.Signature == nil {
		fmt.Println("sig nil.")
	} else {
		fmt.Println("bad!")
	}
	b := msgs.AddServerPubkeyMsg{}
	json.Unmarshal(a.Content, &b)
	// fmt.Println(b.ServerID, b.ServerPubkey)
	fmt.Println()
	fmt.Println()
	fmt.Println(pk)
	myrsa.UnitTest()
	myrsa.UnitTest3()
}
