package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// WFChaincode example simple Chaincode implementation
type WFChaincode struct {
}

func (t *WFChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	_, args := stub.GetFunctionAndParameters()
	var err error
	if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1")
		}
	
	// Initialize the chaincode
	err = stub.PutState("init_wf", []byte(args[0]))
	if err != nil {
		return shim.Error(err.Error())
	}
	
	return shim.Success(nil)
}

func (t *WFChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "write" {
		// calls the write method
		return t.write(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"write\" \"delete\" \"query\"")
}


// Writes an entity to state
func (t *WFChaincode) write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2, Key and Value")
	}
	var Key, Value string
	var err error
	
	Key := args[0]
	Value := args[1]

	// Write the key tothe state in ledger
	err = stub.PutState(Key, []byte(Value))
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to write-put state\",\"Key\":\"" + Key + "\",\"Value\":\"" + string(Value) + "\"}"
		return shim.Error(jsonResp)
	}
	jsonResp := "{\"Key\":\"" + Key + "\",\"Value\":\"" + string(Value) + "\"}"
	fmt.Printf("Write Response:%s\n", jsonResp)
	return shim.Success(nil)
}

// Deletes an entity from state
func (t *WFChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var Key string
	var err error
	
	Key := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(Key)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *WFChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Key string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the Key to query")
	}

	Key = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(Key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + Key + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil value for " + Key + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Key\":\"" + Key + "\",\"Value\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(WFChaincode))
	if err != nil {
		fmt.Printf("Error starting Wellsfargo chaincode: %s", err)
	}
}