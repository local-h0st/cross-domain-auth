package main

import (
	"encoding/json"
	"fmt"
	"msgs"
	"myrsa"
	"net"
)

const ServerPubkey string = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDQlXmFEiNzbO0iHjdYIUPvbWmq\nPmMJcrGVLjRUrr2HtURh9lcrGsti1r4BesFcuS+QAzBlFZsp50Ytae0snr26jnFL\nOpBGscCDLsyPrL3dlUnWGnQY5SOFjvVpAjsuc16W0TXXzdoaW6yZwX+tKd2yLkgb\ncL0alZeTI1v8lJN9YQIDAQAB\n-----END PUBLIC KEY-----\n"
const ServerAddr string = "localhost:54321"
const EnclavePubkey string = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCj8hW+keEOHHHLV/7BRO7I0j7a\nXAfxTvkiM8Qyex+aMQ7Ny+cavF4mWlJmdmGo9K3jHFH3LyEd2JuGPh5T0ad/O76C\nor+hX+RvgXkg0HS3MEQIwmzmNg57RSaNxzlJatXEfjpRJJ5Nc+dyA6hpYzaNj9LY\nKex5gvsGFpBMwQZyVwIDAQAB\n-----END PUBLIC KEY-----\n"
const EnclaveAddr = "localhost:55555"

// 自己的
const selfID string = "genPayloadForServerVS"
const servingPort string = ":55559"
const selfAddr string = "localhost:55559"

const PAS_ID = "PasDomainB"

var PRVKEY, PUBKEY []byte
var PAS_PRVKEY, PAS_PUBKEY []byte

func main() {
	var input string
	PRVKEY, PUBKEY = myrsa.GenRsaKey()
	PAS_PRVKEY, PAS_PUBKEY = myrsa.GenRsaKey()
	fmt.Scanln(&input)
	addPasPubkeyToEnclave()
	fmt.Scanln(&input)
	sendUpdateBlacklistInfo()
	fmt.Scanln(&input)
	sendFragmentMsg()
	fmt.Scanln(&input)
	sendQueryLedgerMsg()
	fmt.Scanln(&input)
	sendGetFragmentMsg()
	fmt.Scanln(&input)
	sendBadFragmentMsg()
	fmt.Scanln(&input)
	sendQueryLedgerMsg()
	fmt.Scanln(&input)
}

func sendUpdateBlacklistInfo() {
	fmt.Println("sendUpdateBlacklistInfo()")
	basic_msg := msgs.BasicMsg{
		Method:    "updateBlackist",
		SenderID:  PAS_ID,
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.DomainInfoRecord{
		Domain:    "domainB",
		PasID:     PAS_ID,
		BlackList: []string{"theBlackID"},
	})
	basic_msg.GenSign(PAS_PRVKEY)
	data, _ := json.Marshal(basic_msg)
	cipher := myrsa.EncryptMsg(data, []byte(EnclavePubkey))
	sendMsg(EnclaveAddr, string(cipher))
}

func addPasPubkeyToEnclave() {
	fmt.Println("addPasPubkeyToEnclave")
	basic_msg := msgs.BasicMsg{
		Method:    "addServerPubkey",
		SenderID:  PAS_ID,
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.AddServerPubkeyMsg{
		ServerPubkey: PAS_PUBKEY,
	})
	data, _ := json.Marshal(basic_msg)
	cipher := myrsa.EncryptMsg(data, []byte(EnclavePubkey))
	sendMsg(EnclaveAddr, string(cipher))
}

func sendGetFragmentMsg() {
	fmt.Println("sendGetFragmentMsg()")
	basic_msg := msgs.BasicMsg{
		Method:    "getFragment",
		SenderID:  selfID,
		Content:   []byte("IAmInTheDark"),
		Signature: nil,
	}
	data, _ := json.Marshal(basic_msg)
	cipher := myrsa.EncryptMsg(data, []byte(ServerPubkey))
	sendMsg(ServerAddr, string(cipher))
}
func sendBadFragmentMsg() {
	fmt.Println("sendBadFragmentMsg()")
	basic_msg := msgs.BasicMsg{
		Method:    "fragment",
		SenderID:  selfID,
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.FragmentMsg{
		Tag:                  "serverVS001",
		PID:                  "notGoodPID",
		Domain:               "domainB",
		Point:                []byte("thisIsThePoint"),
		CipherID:             myrsa.EncryptMsg([]byte("theBlackID"), []byte(EnclavePubkey)),
		PubkeyDeviceToDomain: []byte("pk4dev2domain"),
	})
	basic_msg.GenSign(PRVKEY)
	data, _ := json.Marshal(basic_msg)
	cipher := myrsa.EncryptMsg(data, []byte(ServerPubkey))
	sendMsg(ServerAddr, string(cipher))
}
func sendFragmentMsg() {
	fmt.Println("sendFragmentMsg()")
	basic_msg := msgs.BasicMsg{
		Method:    "fragment",
		SenderID:  selfID,
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.FragmentMsg{
		Tag:                  "serverVS001",
		PID:                  "IAmInTheDark",
		Domain:               "domainB",
		Point:                []byte("thePointToRecoverTrueID"),
		CipherID:             myrsa.EncryptMsg([]byte("whiteTrueID"), []byte(EnclavePubkey)),
		PubkeyDeviceToDomain: []byte("thePubkeyForDeviceAndDomainToCommunicate"),
	})
	basic_msg.GenSign(PRVKEY)
	data, _ := json.Marshal(basic_msg)
	cipher := myrsa.EncryptMsg(data, []byte(ServerPubkey))
	sendMsg(ServerAddr, string(cipher))
}
func sendQueryLedgerMsg() {
	fmt.Println("sendQueryLedgerMsg()")
	basic_msg := msgs.BasicMsg{
		Method:    "queryLedger",
		SenderID:  selfID,
		Content:   nil,
		Signature: nil,
	}
	basic_msg.GenSign(PRVKEY)
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
