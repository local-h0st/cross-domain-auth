package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/syndtr/goleveldb/leveldb"
)

type SmartContract struct {
	contractapi.Contract
}

var node_id string = os.Getenv("NODE_ID")
var node_pk string // 不知道公钥是不是这个类型，记得初始化

// 不要ID了，直接拿PID当成索引
type PseudoRecord struct {
	IsValid   bool   `json:"IsValid`
	PK        string `json:"PK"` // 代表此pid设备期望以后和某个域通信时用的公钥，感觉不需要公开是哪一个域
	Timestamp string `json:"Timestamp"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	recordJSON, _ := json.Marshal(PseudoRecord{
		false,
		"the record for init.",
		time.Now().Format("2006-01-02 15:04:05"),
	}) // 懒得写err了，这里总不可能出错吧
	ctx.GetStub().PutState("redh3tALWAYS", recordJSON)

	// 初始化leveldb
	_, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		return fmt.Errorf("In func InitLedger(): init level db failed.")
	}
	ctx.GetStub().PutState("db", []byte("./db"))

	return nil
}

func (s *SmartContract) HandleMsgForPseudo(ctx contractapi.TransactionContextInterface, cipher_text string) error {
	// msg 先rsa解密，这里暂时没写，因此测试的时候不应该加密
	msg_text := cipher_text

	type message struct {
		pid               string
		p                 interface{} // // 门限签名技术的那个点，我也不知道用什么格式存储
		tag               string      // node_id
		id_cphier         string
		pk_device2domamin string // 不知道公钥是这哪个类型。另外是建议存到数据库，还是存到内存算了？
	}

	// msg_text被解析到msg结构内
	var msg message
	err := json.Unmarshal([]byte(msg_text), &msg)
	if err != nil {
		return fmt.Errorf("In func HandleMsgForPseudo(): json unmarshal failed.")
	}

	// 判断pid是否已经存在
	assetJSON, err := ctx.GetStub().GetState(msg.pid)
	if err != nil {
		return fmt.Errorf("In func HandleMsgForPseudo(): get state failed.")
	}
	if assetJSON != nil {
		return fmt.Errorf("In func HandleMsgForPseudo(): pseudo already exists.")
	}

	// 非tag node直接存储，tag node需要和sgx交互
	if msg.tag != node_id {
		return storeMsg(msg_text)
	}
	var rcd PseudoRecord
	rcd.IsValid = verifyIDinSGX(msg.id_cphier)
	rcd.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	rcd.PK = msg.pk_device2domamin
	assetJSON, err = json.Marshal(rcd)
	if err != nil {
		return fmt.Errorf("In func HandleMsgForPseudo(): json marshal failed.")
	}
	return ctx.GetStub().PutState(msg.pid, assetJSON)
}

func storeMsg(m string) error {
	// 用于存储pid和点p，以json的格式
	// 这里我选择LevelDB - Fabric自带的内嵌键值存储数据库,可以直接在chaincode中使用。但是数据只保存在一个peer节点上,不可持久化。

	return nil
}

func verifyIDinSGX(cipher string) bool {
	// 在enclave内部解密核验是否在黑名单上，true表示合法
	// 出错了一律直接false
	return true
}
