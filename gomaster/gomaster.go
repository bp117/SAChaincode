package main

import (
	"errors"
	"fmt"

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
	err = stub.PutState("init_wf", []byte(args[0]))
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
	} 

	fmt.Println("invoke did not find func: " + function)
	return nil, errors.New("Received unknown function invocation: " + function + "expecting init, write, query")
}


// Writes an entity to state
func (t *WFChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2, Key and Value")
	}
	var Key, Value string
	var err error
	
	Key = args[0]
	Value = args[1]

	// Write the key tothe state in ledger
	err = stub.PutState(Key, []byte(Value))
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to write-put state\",\"Key\":\"" + Key + "\",\"Value\":\"" + string(Value) + "\"}"
		return nil, errors.New(jsonResp)
	}
	jsonResp := "{\"Key\":\"" + Key + "\",\"Value\":\"" + string(Value) + "\"}"
	fmt.Printf("Write Response:%s\n", jsonResp)
	return nil, nil
}

// query callback representing the query of a chaincode
func (t *WFChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the Key to query")
	}

	// Get the state from the ledger
	if function == "read" { 
		//read a variable
		return t.read(stub, args)
	}
	Avalbytes, err := t.read(stub, args)
	if err != nil {
		return nil, err
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil value for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Key\":\"" + args[0] + "\",\"Value\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response: %s\n", jsonResp)
	return Avalbytes, nil
}

// read - query function to read key/value pair
func (t *WFChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
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
