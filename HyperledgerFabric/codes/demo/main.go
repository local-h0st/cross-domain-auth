package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type PseudoRecord struct {
	// 不要ID了，直接拿PID当成索引
	IsValid   bool   `json:"IsValid`
	PK        string `json:"PK"` // 代表此pid设备期望以后和某个域通信时用的公钥，感觉不需要公开是哪一个域
	Timestamp string `json:"Timestamp"`
}

func (s *SmartContract) CheckExistance(ctx contractapi.TransactionContextInterface, pid string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(pid)
	if err != nil {
		return false, fmt.Errorf("[chaincode err in CheckExistance()] get state failed.")
	} else if assetJSON != nil {
		return true, nil
	} else {
		// pid不存在
		return false, nil
	}
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	recordJSON, _ := json.Marshal(PseudoRecord{
		false,
		"the record for init.",
		time.Now().Format("2006-01-02 15:04:05"),
	}) // 懒得写err了，这里总不可能出错吧
	ctx.GetStub().PutState("redh3tALWAYS", recordJSON)
	return nil
}

/*

func (s *SmartContract) HandleMsgForPseudo(ctx contractapi.TransactionContextInterface, cipher_text string) error {
	// msg 先rsa解密，这里暂时没写，因此测试的时候不应该加密
	msg_text := cipher_text

	// 解析msg_text到msg
	type message struct {
		PID              string
		P                string // // 门限签名技术的那个点，我也不知道用什么格式存储
		Tag              string // node_id
		ID_cipher        string
		PK_device2domain string // 不知道公钥是这哪个类型。另外是建议存到数据库，还是存到内存算了？
	}
	var msg message
	err := json.Unmarshal([]byte(msg_text), &msg)
	if err != nil {
		return fmt.Errorf("In func HandleMsgForPseudo(): json unmarshal failed.")
	}

	// 判断pid是否已经存在
	assetJSON, err := ctx.GetStub().GetState(msg.PID)
	if err != nil {
		return fmt.Errorf("In func HandleMsgForPseudo(): get state failed.")
	}
	if assetJSON != nil {
		return fmt.Errorf("In func HandleMsgForPseudo(): pseudo already exists.")
	}

	// 设计中非tag node直接存储，tag node需要和sgx交互
	// 这里每个node都存储
	// 存储不修改账本，因此暂时不会造成共识性问题
	s.storeMsg(ctx, msg.PID, msg_text)

	// 以下这一段是有问题的
	var rcd PseudoRecord
	rcd.IsValid = verifyIDinSGX(msg.ID_cipher)
	rcd.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	rcd.PK = msg.PK_device2domain
	assetJSON, err = json.Marshal(rcd)
	if err != nil {
		return fmt.Errorf("In func HandleMsgForPseudo(): json marshal failed.")
	}
	return ctx.GetStub().PutState(msg.PID, assetJSON)
}

func (s *SmartContract) RequestTrackPID(ctx contractapi.TransactionContextInterface, pid string) (string, bool, error) {
	// 这个接口应该重新设计，传入的参数应该包含PAS提供的判定材料，因此还要判定一次，这里没写
	// 另外按照设计应该是将查询结果分别发给每一个PAS的，我这里选择直接上链
	info_str, err := s.getMsg(ctx, pid)
	if err != nil {
		return "", false, fmt.Errorf("In func RequestTtrackPID(): s.getMsg() failed.")
	}

}
*/

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("[main] create chaincode failed.")
	}
	err = chaincode.Start()
	if err != nil {
		log.Panicf("[main] start chaincode failed.")
	}
}
