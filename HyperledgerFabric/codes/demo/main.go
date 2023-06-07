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
	// 直接拿PID当成索引
	Valid                string
	PubkeyDeviceToDomain string
	Timestamp            string
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	recordJSON, _ := json.Marshal(PseudoRecord{
		Valid:                "validorinvalid",
		PubkeyDeviceToDomain: "pubkeydevicetodomain",
		Timestamp:            time.Now().Format("2006-01-02 15:04:05"),
	})
	ctx.GetStub().PutState("redh3tALWAYS", recordJSON)
	return nil
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
func (s *SmartContract) AddPseudoRecord(ctx contractapi.TransactionContextInterface, pid string, valid string, pubkey_device_to_domain string) error {
	recordJSON, _ := json.Marshal(PseudoRecord{
		Valid:                valid,
		PubkeyDeviceToDomain: pubkey_device_to_domain,
		Timestamp:            time.Now().Format("2006-01-02 15:04:05"),
	})
	ctx.GetStub().PutState(pid, recordJSON)
	return nil
}
func (s *SmartContract) QueryAll(ctx contractapi.TransactionContextInterface) ([]*PseudoRecord, error) {
	// 抄的chaincode
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var assets []*PseudoRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset PseudoRecord
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}
	return assets, nil
}

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
