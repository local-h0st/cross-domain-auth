package main

import (
	"crypto/rsa"
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var node_id string = os.Getenv("NODE_ID")
var pk rsa.PublicKey // 不知道公钥是不是这个类型，记得初始化

type SmartContract struct {
	contractapi.Contract
}

// https://youtube.com/shorts/y0cxkflRHto?feature=share

type InfoRecord struct {
	PID       string `json:"PID"`
	Valid     bool
	PK_pid    string `json:"PK_pid`
	Timestamp int
}

// 监听端口，收到消息 就开一个goroutine去处理
func handleReq(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		cipher := string(buf[:n])
	}
	// msg 先rsa解密，这里没写，因此测试的时候不应该加密
	msg := cipher
	// msg 按照 json的格式存储
	type common_msg struct {
		pid string
		p   interface{} // 门限签名技术的那个点，我也不知道用什么格式存储

	}
	type tag_msg struct {
		pid               string
		p                 interface{}
		tag               string // node_id
		id_cphier         string
		pk_device2domamin rsa.PublicKey // 不知道公钥是不是这个类型，建议本地存储？存到数据库吧，还是存到内存算了
	}
	var con_m common_msg
	var tag_m tag_msg
	var msg_type bool
	err := json.Unmarshal(msg, &con_m)
	if err != nil {
		err = json.Unmarshal(msg, &tag_m)
		if err != nil {
			panic("json unmarshal collapsed.")
		} else {
			msg_type = true
		}
	}
	var InfoRecord
	if msg_type { // 是 tag 节点
		if tag_m.tag != node_id {
			panic("tag doesnt match node id.")
		}
		// golang可以直接和Intel SGX交互！！！！！！！！！

	} else {
		storeMsg(msg)
	}

}

func storeMsg(m string) {
	// 用于存储pid和点p，以json的格式
}

func verifyIDinSGX(cipher string) bool {
	// 在enclave内部解密核验是否在黑名单上，true表示合法
	return true
}
