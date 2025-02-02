package test

import (
	"users/smartcontract"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
)

var Stub *shimtest.MockStub
var Scc *contractapi.ContractChaincode
var user1 smartcontract.User = smartcontract.User{
	ID:    "1",
	Name:  "John Lee",
	Email: "john.lee@g.com",
}
var user2 smartcontract.User = smartcontract.User{
	ID:    "2",
	Name:  "Amy Lin",
	Email: "amy.lin@g.com",
}

var transaction1 smartcontract.Transaction = smartcontract.Transaction{
	Hash:      "0x000000001",
	Amount:    "200",
	Currency:  "USD",
	Date:      "2022-04-14",
}

var transaction2 smartcontract.Transaction = smartcontract.Transaction{
	Hash:      "0x000000002",
	Amount:    "500",
	Currency:  "NTD",
	Date:      "2022-04-16",
}

var bank smartcontract.Bank = smartcontract.Bank{
	ID: "04231910",
	Name: "國泰世華商業銀行",
	TransactionCount: 0,
}

var testbank smartcontract.Bank = smartcontract.Bank{
	ID: "123456",
	Name: "Test bank",
	TransactionCount: 0,
}


func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	log.SetOutput(ioutil.Discard)
}

func NewStub() {
	Scc, err := contractapi.NewChaincode(new(smartcontract.SmartContract))
	if err != nil {
		log.Println("NewChaincode failed", err)
		os.Exit(0)
	}
	Stub = shimtest.NewMockStub("main", Scc)
}

func Test_CreateUser(t *testing.T) {
	fmt.Println("Test_CreateUser-----------------")
	NewStub()

	err := MockCreateUser(user1.ID, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}
}

func Test_UserExists(t *testing.T) {
	fmt.Println("Test_UserExists-----------------")
	NewStub()

	err := MockCreateUser(user1.ID, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}

	result, err := MockUserExists(user1.ID)
	if err != nil {
		t.FailNow()
	}
	fmt.Println("result: ", result)
	assert.Equal(t, result, true)
}

func Test_GetUser(t *testing.T) {
	fmt.Println("Test_GetUser-----------------")
	NewStub()

	err := MockCreateUser(user1.ID, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}

	userJson, err := MockGetUser(user1.ID)
	if err != nil {
		fmt.Println("get User error", err)
	}
	fmt.Println("userJson: ", userJson)
	assert.Equal(t, userJson.ID, user1.ID)
	assert.Equal(t, userJson.Name, user1.Name)
	assert.Equal(t, userJson.Email, user1.Email)
}

func Test_UpdateUser(t *testing.T) {
	fmt.Println("Test_UpdateUser-----------------")
	NewStub()

	err := MockCreateUser(user1.ID, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}

	MockUpdateUser(user1.ID, "change name", "change email")

	userJson, err := MockGetUser(user1.ID)
	if err != nil {
		fmt.Println("get User", err)
	}
	fmt.Println("userJson: ", userJson)
	assert.Equal(t, userJson.ID, user1.ID)
	assert.Equal(t, userJson.Name, "change name")
	assert.Equal(t, userJson.Email, "change email")

}

func Test_DeleteUser(t *testing.T) {
	fmt.Println("Test_DeleteUser-----------------")
	NewStub()
	err := MockCreateUser(user1.ID, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}

	MockDeleteUser(user1.ID)

	userJson, err := MockGetUser(user1.ID)
	if err != nil {
		fmt.Println("get User", err)
	}
	fmt.Println(userJson)
	// assert.Equal(t, err, errors.New("GetUser error"))
}

func Test_GetAllUsers(t *testing.T) {
	fmt.Println("MockGetAllUsers-----------------")
	NewStub()

	MockCreateUser(user1.ID, user1.Name, user1.Email)
	MockCreateUser(user2.ID, user2.Name, user2.Email)

	users, err := MockGetAllUsers()
	if err != nil {
		fmt.Println("GetAllUsers error", err)
	}
	fmt.Println("users: ", users)
	// assert.Equal(t, len(users), 2)
}

//新增 Test_CreateTransaction
func Test_CreateTransaction(t *testing.T) {
	fmt.Println("CreateTransaction-----------------")
	NewStub()
	MockInitLedger()
	err := MockCreateUser(user1.ID, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}
	result1, err := MockCreateTransaction(user1.ID, transaction1.Hash, transaction1.Amount, transaction1.Currency, transaction1.Date, bank.ID)
	if err != nil {
		fmt.Println("CreateTransaction User", err)
	}
	fmt.Println("CreateTransaction transaction1", result1)

	result2, err := MockCreateTransaction(user1.ID, transaction2.Hash, transaction2.Amount, transaction2.Currency, transaction2.Date, bank.ID)
	if err != nil {
		fmt.Println("CreateTransaction User", err)
	}
	fmt.Println("CreateTransaction transaction2", result2)

	user, err := MockGetUser(user1.ID)
	if err != nil {
		fmt.Println("get User error", err)
	}
	fmt.Println(user)
	assert.Equal(t, len(user.Transactions), 2)

}

func Test_GetUserByTransactionHash(t *testing.T) {
	fmt.Println("GetUserByTransactionHash-----------------------")
	NewStub()
	MockInitLedger()
	err := MockCreateUser(user1.ID, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}
	result, err := MockCreateTransaction(user1.ID, transaction1.Hash, transaction1.Amount, transaction1.Currency, transaction1.Date, bank.ID)
	if err != nil {
		fmt.Println("CreateTransaction User", err)
	}
	fmt.Println("CreateTransaction transaction", result)

	userJson, err := MockGetUserByTransactionHash(transaction1.Hash)
	if err != nil {
		fmt.Println("Get User By Transaction Hash error", err)
	}
	fmt.Println(userJson)
	assert.Equal(t, userJson.ID, user1.ID)
	assert.Equal(t, userJson.Name, user1.Name)
	assert.Equal(t, userJson.Email, user1.Email)
}

func Test_InitLedger(t *testing.T) {
	fmt.Println("InitLedger-----------------")
	NewStub()
	var err = MockInitLedger()
	assert.Equal(t, err, nil)
}

func Test_GetBankByID(t *testing.T) {
	fmt.Println("Test_GetBankByID-----------------")
	NewStub()
	MockInitLedger()
	err := MockCreateUser(user1.ID, user1.Name, user1.Email)
	if err != nil {
		t.FailNow()
	}
	result, err := MockCreateTransaction(user1.ID, transaction1.Hash, transaction1.Amount, transaction1.Currency, transaction1.Date, bank.ID)
	fmt.Println("CreateTransaction transaction", result)
	if err != nil {
		fmt.Println("CreateTransaction User", err)
	}

	bankJson, err := MockGetBankByID(bank.ID)
	if err != nil {
		fmt.Println("get User error", err)
	}
	fmt.Println("bankJson: ", bankJson)
	assert.Equal(t, bankJson.ID, bank.ID)
	assert.Equal(t, bankJson.Name, bank.Name)
	assert.Equal(t, bankJson.TransactionCount, 1)
}

func Test_CreateBank(t *testing.T) {
	fmt.Println("Test_CreateBank-----------------")
	NewStub()

	err := MockCreateBank(testbank.ID, testbank.Name)
	if err != nil {
		t.FailNow()
	}
}

//
// 
// Mock function

func MockUserExists(id string) (bool, error) {
	res := Stub.MockInvoke("uuid", [][]byte{[]byte("UserExists"), []byte(id)})
	if res.Status != shim.OK {
		return false, errors.New("UserExists error")
	}
	var result bool = false
	json.Unmarshal(res.Payload, &result)
	return result, nil
}

func MockCreateUser(id string, name string, email string) error {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("CreateUser"),
			[]byte(id),
			[]byte(name),
			[]byte(email),
		})

	if res.Status != shim.OK {
		fmt.Println("CreateUser failed", string(res.Message))
		return errors.New("CreateUser error")
	}
	return nil
}

func MockCreateBank(bankId string, name string) error {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("CreateBank"),
			[]byte(bankId),
			[]byte(name),
		})
	
		if res.Status != shim.OK {
			fmt.Println("CreateBank failed", string(res.Message))
			return errors.New("CreateBank error")
		}
		return nil
}

func MockGetUser(id string) (*smartcontract.User, error) {
	var result smartcontract.User
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("GetUser"),
			[]byte(id),
		})
	if res.Status != shim.OK {
		fmt.Println("GetUser failed", string(res.Message))
		return nil, errors.New("GetUser error")
	}
	json.Unmarshal(res.Payload, &result)
	return &result, nil
}

func MockGetUserByTransactionHash(hash string) (*smartcontract.User, error) {
	var result smartcontract.User
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("GetUserByTransactionHash"),
			[]byte(hash),
		})
	if res.Status != shim.OK {
		fmt.Println("GetUserByTransactionHash failed", string(res.Message))
		return nil, errors.New("GetUserByTransactionHash error")
	}
	json.Unmarshal(res.Payload, &result)
	return &result, nil
}

func MockUpdateUser(id string, name string, email string) error {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("UpdateUser"),
			[]byte(id),
			[]byte(name),
			[]byte(email),
		})
	if res.Status != shim.OK {
		fmt.Println("UpdateUser failed", string(res.Message))
		return errors.New("UpdateUser error")
	}
	return nil
}

func MockDeleteUser(id string) error {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("DeleteUser"),
			[]byte(id),
		})
	if res.Status != shim.OK {
		fmt.Println("DeleteUser failed", string(res.Message))
		return errors.New("DeleteUser error")
	}
	return nil
}

func MockGetAllUsers() ([]*smartcontract.User, error) {
	res := Stub.MockInvoke("uuid", [][]byte{[]byte("GetAllUsers")})
	if res.Status != shim.OK {
		fmt.Println("GetAllUsers failed", string(res.Message))
		return nil, errors.New("GetAllUsers error")
	}
	var users []*smartcontract.User
	json.Unmarshal(res.Payload, &users)
	return users, nil
}

// 新增 MockCreateTransaction
func MockCreateTransaction(userId string, hash string, amount string, currency string, date string, bankId string) (bool, error) {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("CreateTransaction"),
			[]byte(userId),
			[]byte(hash),
			[]byte(amount),
			[]byte(currency),
			[]byte(date),
			[]byte(bankId),
		})
	if res.Status != shim.OK {
		fmt.Println("CreateTransaction failed", string(res.Message))
		return false, errors.New("CreateTransaction error")
	}
	var result bool = false
	json.Unmarshal(res.Payload, &result)
	return result, nil
}

func MockGetBankByID(id string) (*smartcontract.Bank, error) {
	var result smartcontract.Bank
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("GetBankByID"),
			[]byte(id),
		})
	if res.Status != shim.OK {
		fmt.Println("GetBankByID failed", string(res.Message))
		return nil, errors.New("GetBankByID error")
	}
	json.Unmarshal(res.Payload, &result)
	return &result, nil
}

func MockInitLedger() (error) {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("InitLedger"),
		})
		if res.Status != shim.OK {
			fmt.Println("MockInitLedger failed", string(res.Message))
			return errors.New("MockInitLedger error")
		}
		return nil
}