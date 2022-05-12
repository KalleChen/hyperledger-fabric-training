package smartcontract

import (
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing archives
type SmartContract struct {
	contractapi.Contract
}

// User Data struct
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Transactions []Transaction `json:"transactions,omitempty" metadata:",optional"`
}

type Transaction struct {
	Hash         string `json:"hash"`
	Amount       string `json:"amount"`
	Currency string `json:"currency"`
	Date    string `json:"date"`
}

type TransactionHashMapUserId struct {
	UserId	string `json:"user_id"`
}

type Bank struct {
	ID	string `json:"id"`
	Name	string `json:"name"`
	TransactionCount	int `json:"transaction_count"`
}

const bankPrefix = "Bank_"

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	var cathayBank Bank = Bank{
		ID: "04231910",
		Name: "國泰世華商業銀行",
		TransactionCount: 0,
	}
	var fubonBank Bank = Bank{
			ID: "03750168",
			Name: "台北富邦商業銀行",
			TransactionCount: 0,
	}
	cathayBankJson, err := json.Marshal(cathayBank)
	if err != nil {
		return err
	}
	fubonBankJson, err := json.Marshal(fubonBank)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(bankPrefix+"04231910", cathayBankJson)
	ctx.GetStub().PutState(bankPrefix+"03750168", fubonBankJson)
	return nil
}

func (s *SmartContract) UserExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	fmt.Println("function UserExists")
	//TODO UserExists
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJSON != nil, nil
}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, id string, name string, email string) error {
	fmt.Println("function CreateUser")
	//TODO CreateUser
	exists, err := s.UserExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the user %s already exists", id)
	}

	user := User{
		ID: id,
		Name: name,
		Email: email,
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(id, userJson)
	return nil
}

func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, id string) (*User, error) {
	fmt.Println("function GetUser")
	//TODO GetUser
	userJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userJson == nil {
		return nil, fmt.Errorf("the user %s does not exist", id)
	}

	var user User
	err = json.Unmarshal(userJson, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, id string, name string, email string) error {
	fmt.Println("function UpdateUser")
	//TODO UpdateUser
	user, err := s.GetUser(ctx, id)
	if err != nil {
		return err
	}
	user.Email = email
	user.Name = name
	userJson, err := json.Marshal(user)

	return ctx.GetStub().PutState(id, userJson)
}

func (s *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, id string) error {
	fmt.Println("function DeleteUser")
	//TODO DeleteUser
	exists, err := s.UserExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the user %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	fmt.Println("function GetAllUsers")
	var users []*User
	//TODO GetAllUsers
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

		var user User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, userId string, hash string, amount string, currency string, date string, bankId string) (bool, error) {
	user, err := s.GetUser(ctx, userId)
	if err != nil {
		return false, err
	}

	var transaction Transaction = Transaction{
		Hash:      hash,
		Amount:    amount,
		Currency:  currency,
		Date:      date,
	}
	user.Transactions = append(user.Transactions, transaction)

	userJson, err := json.Marshal(user)
	if err != nil {
		return false, err
	}

	ctx.GetStub().PutState(userId, userJson)

	var transactionHashMapUserId TransactionHashMapUserId = TransactionHashMapUserId{
		UserId:	userId,
	}

	transactionHashMapJson, err := json.Marshal(transactionHashMapUserId)
	if err != nil {
		return false, err
	}

	ctx.GetStub().PutState(hash, transactionHashMapJson)

	bank, err := s.GetBankByID(ctx, bankId)
	bank.TransactionCount++
	bankJson, err := json.Marshal(bank)
	if err != nil {
		return false, err
	}

	ctx.GetStub().PutState(bankPrefix + bankId, bankJson)

	return true, nil
}

func (s *SmartContract) GetUserByTransactionHash(ctx contractapi.TransactionContextInterface, hash string) (*User, error) {
	fmt.Println("function GetUser")
	//TODO GetUser
	transactionHashMapUserIdJson, err := ctx.GetStub().GetState(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if transactionHashMapUserIdJson == nil {
		return nil, fmt.Errorf("the transaction %s does not exist", hash)
	}
	var transactionHashMapUserId TransactionHashMapUserId
	err = json.Unmarshal(transactionHashMapUserIdJson, &transactionHashMapUserId)
	if err != nil {
		return nil, err
	}
	user, err := s.GetUser(ctx, transactionHashMapUserId.UserId)
	if err != nil {
		return nil, err
	}
	return user, nil	
}

func (s *SmartContract) GetBankByID(ctx contractapi.TransactionContextInterface, bankId string) (*Bank, error) {
	bankJson, err := ctx.GetStub().GetState(bankPrefix + bankId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bankJson == nil {
		return nil, fmt.Errorf("the bank %s does not exist", bankId)
	}

	var bank Bank
	err = json.Unmarshal(bankJson, &bank)
	if err != nil {
		return nil, err
	}
	return &bank, nil
}

func (s *SmartContract) BankExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	//TODO UserExists
	assetJSON, err := ctx.GetStub().GetState(bankPrefix + id)
	fmt.Println("function UserExists")
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJSON != nil, nil
}

func (s *SmartContract) CreateBank(ctx contractapi.TransactionContextInterface, bankId string, name string) error {
	exists, err := s.BankExists(ctx, bankId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the bank %s already exists", bankId)
	}

	bank := Bank{
		ID: bankId,
		Name: name,
		TransactionCount: 0,
	}
	bankJson, err := json.Marshal(bank)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(bankPrefix + bankId, bankJson)
	return nil	
}