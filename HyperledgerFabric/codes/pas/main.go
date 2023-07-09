package main

import (
	"encoding/json"
	"fmt"
	"msgs"
	"myrsa"
	"net"
	sharedconfigs "sharedConfigs"
)

const (
	pasPrvkey  = "-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQDoL3urEUspR34r6AWvCdv5z4/MW/Kut1Is7gW2JroFr6u+fRnk\nQBfqrSbG2HPXVvX8YdeX/xaLKUdaxUMrq/zRkt1ggKBx3k2NlS3pG82T9DZszYyJ\nj+GP1fyy40N+OetQqtnPc4Zy1sGmYsbnlxEfFkqTEXqsM5Vpg4oSnFz96wIDAQAB\nAoGBAIBVvoVPibvHSHX8SSf2yx/JGjJaoEjyCvnKll2YCjoaX1Nq0mTXCGEuU8CU\n43KjHlPhwMjCtjM1HbuOTRJWfeZK+EliP8dtrwCYuhBmT2HVtN4hjOiH9frGjmCC\n4LqaJ400w7mW7TogVDP8HE/TcFzTolXLWDhRSJn+Zg420bwBAkEA6LZlsCnWJFqn\n/01cdxr0zc10FXbc4VxGVkpav99Kw2OtmPzkfUDbsCT5RFehOB2MB0LHQKL+axF/\ncGHjvDdSBwJBAP9rlcO796xykdsdxWdpvox3+7oKv4fcG9dgp57cE9agr/x8QEH7\nvttt8jR7Zw0Vm5/Fd9FsAZ3+LlJ6aXhYa/0CQECvmrqKFo1KadJMhbxR0OR4DKF+\nxc0a4i5QQsN85QJE7ddNzJGIesiOrn8xwI2hoO/PvyUXaZMHbR4nB6+kzPcCQF7U\nhsohI5d3AggkSYJXlFN6yI8OJoY+hme0jwdAFm19Q1mul/znhrjZXS93EY+eEiWD\nnzS1sPQDxxcAM+Bmk9ECQQC4JU0VLqQCfNPYxfzJiWZbGOA8igRkBcMrV/wnK14I\nEE0K7mVDmpKwCR/OHpA2Kk+rBs5FBUHqjFr7sGBvcO6m\n-----END RSA PRIVATE KEY-----\n"
	pasPubkey  = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDoL3urEUspR34r6AWvCdv5z4/M\nW/Kut1Is7gW2JroFr6u+fRnkQBfqrSbG2HPXVvX8YdeX/xaLKUdaxUMrq/zRkt1g\ngKBx3k2NlS3pG82T9DZszYyJj+GP1fyy40N+OetQqtnPc4Zy1sGmYsbnlxEfFkqT\nEXqsM5Vpg4oSnFz96wIDAQAB\n-----END PUBLIC KEY-----\n"
	pasID      = "pasCali"
	domainName = "domainCali"
	pasAddr    = "localhost:51234"
)

var input string

func main() {

	// 关于fragment，有多少台server就要send多少次，注意每次密钥都不同
	fmt.Scanln(&input)
	sendFragment()

}

func sendFragment() {
	bm := msgs.BasicMsg{
		Method:   "fragment",
		SenderID: pasID,
	}
	bm.Content, _ = json.Marshal(msgs.FragmentMsg{
		Tag:        sharedconfigs.NodeID,
		PID:        "ppppppppppppppid",
		OrigDomain: "domainCali",
		DestDomain: "domainHang",
		// TODO
	})
	bm.GenSign([]byte(pasPrvkey))
	bm_str, _ := json.Marshal(bm)
	sendMsg(sharedconfigs.ServerAddr, string(myrsa.EncryptMsg(bm_str, []byte(sharedconfigs.ServerPubkey))))
}

func syncBlacklist() {

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
