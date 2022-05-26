package smartcontract

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing archives
type SmartContract struct {
	contractapi.Contract
}

type File struct {
	ID        string `json:"id"`
	FILE_NAME string `json:"file_name"`
	Hash      string `json:"hash"`
	Time      string `json:"time"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) FileExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJSON != nil, nil
}

func (s *SmartContract) CreateFile(ctx contractapi.TransactionContextInterface, id string, file_name string, hash string, time string) error {
	exists, err := s.FileExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the File %s already exists", id)
	}

	file := File{
		ID:   id,
    FILE_NAME: file_name,
		Hash: hash,
		Time: time,
	}
	fileJson, err := json.Marshal(file)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(id, fileJson)
	return nil
}

func (s *SmartContract) GetFile(ctx contractapi.TransactionContextInterface, id string) (*File, error) {
	fileJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if fileJson == nil {
		return nil, fmt.Errorf("the file %s does not exist", id)
	}

	var file File
	err = json.Unmarshal(fileJson, &file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *SmartContract) GetAllFiles(ctx contractapi.TransactionContextInterface) ([]*File, error) {
	var files []*File
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var file File
		err = json.Unmarshal(queryResponse.Value, &file)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}
