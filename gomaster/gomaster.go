package main

import (
	"errors"
	"fmt"
	b64 "encoding/base64"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// WFChaincode example simple Chaincode implementation
type WFChaincode struct {
}

func (t *WFChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("ex02 Init")
	var err error
	if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}

	// Initialize the chaincode
	err = stub.PutState("DOCUMENT-LIST", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *WFChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "init" {
		// calls our init method
		return t.Init(stub, "init", args)
	} else if function == "write" {
		// calls the write method
		return t.write(stub, args)
	} else if function == "writeDocument" {
		// calls the write method
		return t.writeDocument(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)
	return nil, errors.New("Received unknown function invocation: " + function + " expecting init, write, writeDocument, query")
}


// Writes an entity to state
func (t *WFChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args ) < 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3: Key, Value, LogInfo")
	}
	var key, value string
	var err error
	var logData []byte

	key = args[0]
	value = args[1]
	logData, _ =  b64.StdEncoding.DecodeString(args[2])

	fmt.Printf("Running WRITE function :%s\n", string(logData))
	// Write the key to the state in ledger
	err = stub.PutState(key, []byte(value))
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to write-put state\",\"Key\":\"" + key + "\",\"Value\":\"BLOCK DATA\"}"
		return nil, errors.New(jsonResp)
	}
	jsonResp := "{\"Key\":\"" + key + "\",\"Value\":\"BLOCK DATA\"}"
	fmt.Printf("Write Response:%s\n", jsonResp)
	return nil, nil
}

// Writes an entity to state
func (t *WFChaincode) writeDocument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args ) < 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4: DocKey, DocValue, DocInfo, LogInfo")
	}
	var key, value, docInfo string
	var err error
	var logData []byte

	key = args[0]
	value = args[1]
	docInfo = args[2]
	logData, _ =  b64.StdEncoding.DecodeString(args[3])

	fmt.Printf("Running writeDocument function :%s\n", string(logData))
	// Write the key to the state in ledger
	err = stub.PutState(key, []byte(value))
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to writeDocument-put state\",\"Key\":\"" + key + "\",\"Value\":\"BLOCK DATA\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Printf("DOCUMENT-LIST")
	// Write the key to the state in ledger
	err = stub.PutState("DOCUMENT-LIST", []byte(docInfo))
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to update DOCUMENT-LIST\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Key\":\"" + key + "\",\"Value\":\"BLOCK DATA\"}"
	fmt.Printf("Write Response:%s\n", jsonResp)
	return nil, nil
}

// query callback representing the query of a chaincode
func (t *WFChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the Key to query")
	}

	// Get the state from the ledger
	if function == "read" {

		Avalbytes, err := t.read(stub, args)
		if err != nil {
			return nil, err
		}

		if Avalbytes == nil {
			jsonResp := "{\"Error\":\"Nil value for " + args[0] + "\"}"
			return nil, errors.New(jsonResp)
		}

		jsonResp := "{\"Key\":\"" + args[0] + "\",\"Value\":\"BLOCK DATA\"}"
		fmt.Printf("Query Response: %s\n", jsonResp)
		return Avalbytes, nil
	}
	fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query: " + function)
}

// read - query function to read key/value pair
func (t *WFChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	var logData []byte

	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key and log data for the query")
	}

	key = args[0]
	logData, _ = b64.StdEncoding.DecodeString(args[1])
	fmt.Printf("Running READ function :%s\n", string(logData))
	// reading the state
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func main() {
	err := shim.Start(new(WFChaincode))
	if err != nil {
		fmt.Printf("Error starting Wellsfargo chaincode: %s", err)
	}
}