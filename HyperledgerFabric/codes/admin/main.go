package main

import (
	"encoding/json"
	"fmt"
	"msgs"
	"myrsa"
	"net"
)

const AdminPrvkey = "-----BEGIN RSA PRIVATE KEY-----\nMIICXgIBAAKBgQDVDTjZaYpCk5R8RNcq86rfrbxn6pss6d+MlodAu/UdgfstJzQK\n6X9gqWceCyAmm30CHdRZbC9E39E2U13uF9NX2oTcpdC0XLRAW+rskd/jCGB0zQSf\ne4ktzKu0a5Cp1pwLWyZtCU7re6A2bhVeOh9JUNADFBiemhYl3JIfetAVDwIDAQAB\nAoGBAKjAjlT3GcJeLvC3fk7RLnl5nZAZ7cuHe8BZwsvtlNtIh3FeagRyqqgfxkOv\nwEmUQ1IX2ojx/gbp2UbUhcP/LzAmKN0CWFlpOF+XodndhPVSMYMOajua6+IWXWIO\nKk70FgFNSTKZ2x+dua6P7UgDtXCBs8EZo/RtTF+ieQ7tBSWxAkEA++TqC+9TWZWo\n6ThnfY+VJF4/xrfZL3xbogZUhMutkpJokQSVZTRsciXLhMvsrq1tFqspRUDE9BlQ\nXjFstizu5QJBANiGOnz5dM9PHVnk4ko41eOUlMHm5bIll6geYms2CDzfSUZCJ6Lu\nXusNuTUduuH8xdoJD+knLgWoLuPHMy1TQOMCQGSpkFaAp6BvTHcXEVR+Iq3L9FSn\nd+WgHsZbHT+MXarrU1pQqJsvHf9n1zMUg1sy9xtN/0orngmmbBWYTsdmoXkCQQDU\n9njScNzmBg99WjUEAZDGLV5+tKaZKIZYkcIFZviFPqyoUOsBQujS0gWm653jJiZH\nhIBEtwd6AuhTmpqIawk3AkEAnpe91sUA0SouCWIzgI2EQu4Z6pEgpqb6AHBuuTfF\nvYkZsveKy8edEi2hGzo+U6nb6S4VYnm/rOQKum1Mx5GxEQ==\n-----END RSA PRIVATE KEY-----\n"
const TargetServerAddr = "localhost:54321"
const TargetServerPubkey = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCU+Q+LZjRCJyWxFMpUeTmdHvzl\nRyNCAsaWxUFI1W9mXywxcItmLmGLEG/tYVzYkjdgjqldvYZfi6q7JNdIYaM7EmXr\n5sgDWbfrH6eeXPe+dTu5dGTSvZ2K3Kvlo17kDKp/H4LcDcG3YpvAPQp/nII9ReWI\n56iq8mohZxb1eMgxPwIDAQAB\n-----END PUBLIC KEY-----\n"

func main() {
	for {
		fmt.Println("=== admin ===")
		fmt.Println("q. syncDomain")
		fmt.Println("w. debugPrintAll")
		fmt.Println("e. queryLedger")

		var input string
		fmt.Scanln(&input)
		switch input {
		case "q":
			syncDomain("domainCali", "pasCali", "localhost:51234", []byte("-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDoL3urEUspR34r6AWvCdv5z4/M\nW/Kut1Is7gW2JroFr6u+fRnkQBfqrSbG2HPXVvX8YdeX/xaLKUdaxUMrq/zRkt1g\ngKBx3k2NlS3pG82T9DZszYyJj+GP1fyy40N+OetQqtnPc4Zy1sGmYsbnlxEfFkqT\nEXqsM5Vpg4oSnFz96wIDAQAB\n-----END PUBLIC KEY-----\n"))
			syncDomain("domainHang", "", "", nil)
		case "w":
			debugPrintAll()
		case "e":
			queryLedger()
		default:
			fmt.Println("^C to quit.")
		}
		fmt.Println()

	}
}

func debugPrintAll() {
	basic_msg := msgs.BasicMsg{
		Method:    "debugPrintAll",
		SenderID:  "admin",
		Content:   nil,
		Signature: nil,
	}
	basic_msg.GenSign([]byte(AdminPrvkey))
	jsonstr, _ := json.Marshal(basic_msg)
	sendMsg(TargetServerAddr, string(myrsa.EncryptMsg(jsonstr, []byte(TargetServerPubkey))))
}

func syncDomain(domain_name, pas_id, pas_addr string, pas_pubkey []byte) {
	basic_msg := msgs.BasicMsg{
		Method:    "syncDomain",
		SenderID:  "admin",
		Content:   nil,
		Signature: nil,
	}
	basic_msg.Content, _ = json.Marshal(msgs.DomainRecord{
		Domain:                       domain_name,
		PasID:                        pas_id,
		PasAddr:                      pas_addr,
		PasPubkey:                    pas_pubkey,
		BlacklistLastUpdateTimestamp: "",
		WaitQ:                        nil,
	})
	basic_msg.GenSign([]byte(AdminPrvkey))
	jsonstr, _ := json.Marshal(basic_msg)
	sendMsg(TargetServerAddr, string(myrsa.EncryptMsg(jsonstr, []byte(TargetServerPubkey))))
}

func queryLedger() {
	basic_msg := msgs.BasicMsg{
		Method:    "queryLedger",
		SenderID:  "admin",
		Content:   nil,
		Signature: nil,
	}
	jsonstr, _ := json.Marshal(basic_msg)
	sendMsg(TargetServerAddr, string(myrsa.EncryptMsg(jsonstr, []byte(TargetServerPubkey))))
}

func sendMsg(addr string, data string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
	}
	n, err := fmt.Fprint(conn, data)
	if err != nil {
		fmt.Printf("send to %s failed.", addr)
	} else {
		fmt.Println(n, "bytes sent.")
	}
}
