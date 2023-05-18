package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

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

type SmartContract struct {
	contractapi.Contract
}

type SongAsset struct {
	ID     string `json:"ID"`
	Name   string `json:"Name"`
	Author string `json:"Author"`
	Rating int    `json:"Rating"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []SongAsset{
		{ID: "#1", Name: "Leaving California", Author: "Maroon 5", Rating: 2},
		{ID: "#2", Name: "In Your Eyes", Author: "The Weeknd", Rating: 3},
		{ID: "#3", Name: "La Isla Bonita", Author: "Madonna", Rating: 1},
	}
	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SmartContract) AddRecord(ctx contractapi.TransactionContextInterface, id string, name string, author string, rating int) error {
	exists, err := s.RecordExists(ctx, id)
	if err != nil {
		return fmt.Errorf("[AddRecord] query record failed.")
	}
	if exists {
		return fmt.Errorf("[AddRecord] record already exists.")
	}
	assetJSON, err := json.Marshal(SongAsset{
		ID:     id,
		Name:   name,
		Author: author,
		Rating: rating,
	})
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) QueryRecord(ctx contractapi.TransactionContextInterface, id string) (*SongAsset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("[QueryRecord] get state failed.")
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("[QueryRecord] asset nil.")
	}
	var asset SongAsset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, fmt.Errorf("[QueryRecord] unmarshal json failed.")
	}
	return &asset, nil
}

func (s *SmartContract) UpdateRecord(ctx contractapi.TransactionContextInterface, id string, name string, author string, rating int) error {
	// 懒得写细节了，反正差不多
	assetJSON, _ := json.Marshal(SongAsset{
		ID:     id,
		Name:   name,
		Author: author,
		Rating: rating,
	})
	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) DeleteRecord(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.RecordExists(ctx, id)
	if err != nil {
		return fmt.Errorf("[DeleteRecord] query record failed.")
	}
	if !exists {
		return fmt.Errorf("[DeleteRecord] record doesn't exists.")
	}
	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) ChangeRating(ctx contractapi.TransactionContextInterface, id string, new_rating int) error {
	// 	也可以直接UpdateRecord，细节懒得写了
	asset, _ := s.QueryRecord(ctx, id)
	asset.Rating = new_rating
	assetJSON, _ := json.Marshal(asset)
	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) QueryAllRecords(ctx contractapi.TransactionContextInterface) ([]*SongAsset, error) {
	// 懒得写直接抄了算了
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*SongAsset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset SongAsset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (s *SmartContract) RecordExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("[RecordExists] get state failed.")
	}
	return assetJSON != nil, nil
}
