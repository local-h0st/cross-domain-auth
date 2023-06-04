package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"msgs"
	"myrsa"
	"net"
	"time"
)

var PRVKEY, PUBKEY, OLD_PRVKEY, OLD_PUBKEY []byte

const ServerPubkey string = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCj8hW+keEOHHHLV/7BRO7I0j7a\nXAfxTvkiM8Qyex+aMQ7Ny+cavF4mWlJmdmGo9K3jHFH3LyEd2JuGPh5T0ad/O76C\nor+hX+RvgXkg0HS3MEQIwmzmNg57RSaNxzlJatXEfjpRJJ5Nc+dyA6hpYzaNj9LY\nKex5gvsGFpBMwQZyVwIDAQAB\n-----END PUBLIC KEY-----\n"
const ServerAddr string = "localhost:55555"

func main() {
	// rand.Seed(time.Now().Unix())
	PRVKEY, PUBKEY = myrsa.GenRsaKey()
	sendAddServerPubkeyMsg()
	OLD_PRVKEY, OLD_PUBKEY = PRVKEY, PUBKEY
	PRVKEY, PUBKEY = myrsa.GenRsaKey()
	sendUpdateServerPubkeyMsg()
}

func sendUpdateServerPubkeyMsg() {
	fmt.Println("sendUpdateServerPubkeyMsg()")
	basic_msg := msgs.BasicMsg{
		Method:    "updateServerPubkey",
		SenderID:  "genPayloadServer",
		Content:   nil,
		Signature: nil,
	}
	// 感觉senderID和serverID功能重合了，可以改一下
	basic_msg.Content, _ = json.Marshal(msgs.UpdateServerPubkeyMsg{
		ServerID:        "genPayloadServer",
		ServerNewPubkey: PUBKEY,
	})
	basic_msg.GenSign(OLD_PRVKEY)
	data, _ := json.Marshal(basic_msg)
	cipher := myrsa.EncryptMsg(data, []byte(ServerPubkey))
	sendMsg(ServerAddr, string(cipher))
}

func sendAddServerPubkeyMsg() {
	fmt.Println("sendAddServerPubkeyMsg()")
	basic_msg := msgs.BasicMsg{
		Method:    "addServerPubkey",
		SenderID:  "genPayload",
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.AddServerPubkeyMsg{
		ServerID:     "genPayloadServer",
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
		ServerID:     "test",
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
	fmt.Println(b.ServerID, b.ServerPubkey)
	fmt.Println()
	fmt.Println()
	fmt.Println(pk)
	myrsa.UnitTest()
	myrsa.UnitTest3()
}
