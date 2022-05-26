package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"users/smartcontract"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
)

var Stub *shimtest.MockStub
var Scc *contractapi.ContractChaincode
var file1 smartcontract.File = smartcontract.File{
	ID:   "1",
	Hash: "1234",
	Time: "2022/04/11 12:00:09",
}
var file2 smartcontract.File = smartcontract.File{
	ID:   "2",
	Hash: "5678",
	Time: "2022/04/13 12:00:09",
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

func Test_CreateFile(t *testing.T) {
	fmt.Println("Test_CreateFile-----------------")
	NewStub()

	err := MockCreateFile(file1.ID, file1.Hash, file1.Time)
	if err != nil {
		t.FailNow()
	}
}

func Test_FileExists(t *testing.T) {
	fmt.Println("Test_FileExists-----------------")
	NewStub()

	err := MockCreateFile(file1.ID, file1.Hash, file1.Time)
	if err != nil {
		t.FailNow()
	}

	result, err := MockFileExists(file1.ID)
	if err != nil {
		t.FailNow()
	}
	fmt.Println("result: ", result)
	assert.Equal(t, result, true)
}

func Test_GetFile(t *testing.T) {
	fmt.Println("Test_GetFile-----------------")
	NewStub()

	err := MockCreateFile(file1.ID, file1.Hash, file1.Time)
	if err != nil {
		t.FailNow()
	}

	fileJson, err := MockGetFile(file1.ID)
	if err != nil {
		fmt.Println("get File error", err)
	}
	fmt.Println("fileJson: ", fileJson)
	assert.Equal(t, fileJson.ID, file1.ID)
	assert.Equal(t, fileJson.Hash, file1.Hash)
	assert.Equal(t, fileJson.Time, file1.Time)
}

func Test_GetAllFiles(t *testing.T) {
	fmt.Println("MockGetAllFiles-----------------")
	NewStub()

	MockCreateFile(file1.ID, file1.Hash, file1.Time)
	MockCreateFile(file2.ID, file2.Hash, file2.Time)

	files, err := MockGetAllFiles()
	if err != nil {
		fmt.Println("GetAllFiles error", err)
	}
	fmt.Println("files: ", files)
	assert.Equal(t, len(files), 2)
}

//
//
// Mock function

func MockFileExists(id string) (bool, error) {
	res := Stub.MockInvoke("uuid", [][]byte{[]byte("FileExists"), []byte(id)})
	if res.Status != shim.OK {
		return false, errors.New("FileExists error")
	}
	var result bool = false
	json.Unmarshal(res.Payload, &result)
	return result, nil
}

func MockCreateFile(id string, hash string, time string) error {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("CreateFile"),
			[]byte(id),
			[]byte(hash),
			[]byte(time),
		})

	if res.Status != shim.OK {
		fmt.Println("CreateFile failed", string(res.Message))
		return errors.New("CreateFile error")
	}
	return nil
}

func MockGetFile(id string) (*smartcontract.File, error) {
	var result smartcontract.File
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("GetFile"),
			[]byte(id),
		})
	if res.Status != shim.OK {
		fmt.Println("GetFile failed", string(res.Message))
		return nil, errors.New("GetFile error")
	}
	json.Unmarshal(res.Payload, &result)
	return &result, nil
}

func MockGetAllFiles() ([]*smartcontract.File, error) {
	res := Stub.MockInvoke("uuid", [][]byte{[]byte("GetAllFiles")})
	if res.Status != shim.OK {
		fmt.Println("GetAllFiles failed", string(res.Message))
		return nil, errors.New("GetAllFiles error")
	}
	var files []*smartcontract.File
	json.Unmarshal(res.Payload, &files)
	return files, nil
}

func MockInitLedger() error {
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

