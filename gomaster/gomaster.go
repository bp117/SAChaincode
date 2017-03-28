package main

import (
	b64 "encoding/base64"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
	"os"
	"strconv"
)

// WFChaincode example simple Chaincode implementation
type WFChaincode struct {
}

var log = logging.MustGetLogger("gomaster")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func (t *WFChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	log.Info("ex02 Init\n")
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	err = stub.PutState("DOCUMENT_INDEX", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *WFChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	log.Infof("invoke is running %s\n", function)

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

	log.Infof("invoke did not find func: %s\n", function)
	return nil, errors.New("Received unknown function invocation: " + function + " expecting init, write, writeDocument, query")
}

// Writes an entity to state
func (t *WFChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3: Key, Value, LogInfo")
	}
	var key, value string
	var err error
	var logData []byte

	key = args[0]
	value = args[1]
	logData, _ = b64.StdEncoding.DecodeString(args[2])

	log.Infof("Running WRITE function :%s\n", string(logData))
	// Write the key to the state in ledger
	err = stub.PutState(key, []byte(value))
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to write-put state\",\"Key\":\"" + key + "\",\"Value\":\"BLOCK DATA\"}"
		return nil, errors.New(jsonResp)
	}
	jsonResp := "{\"Key\":\"" + key + "\",\"Value\":\"BLOCK DATA\"}"
	log.Infof("Write Response:%s\n", jsonResp)
	return nil, nil
}

// Writes an entity to state
func (t *WFChaincode) writeDocument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4: DocKey, DocValue, DocInfo, LogInfo")
	}
	var key, value, docInfo string
	var err error
	var logData, docIndxData []byte
	var docIndx int

	key = args[0]
	value = args[1]
	docInfo = args[2]
	logData, _ = b64.StdEncoding.DecodeString(args[3])

	log.Infof("Running writeDocument function :%s\n", string(logData))
	// Write the key to the state in ledger
	err = stub.PutState(key, []byte(value))
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to writeDocument-put state\",\"Key\":\"" + key + "\",\"Value\":\"BLOCK DATA\"}"
		return nil, errors.New(jsonResp)
	}
	log.Infof("Update DOCUMENT_INDEX\n")
	// read the DOCUMENT_INDEX from the ledger
	docIndxData, err = stub.GetState("DOCUMENT_INDEX")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to read DOCUMENT_INDEX\"}"
		return nil, errors.New(jsonResp)
	}
	docIndx, err = strconv.Atoi(string(docIndxData))
	docIndx++
	err = stub.PutState("DOCUMENT_INDEX", []byte(strconv.Itoa(docIndx)))
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to Update DOCUMENT_INDEX\"}"
		return nil, errors.New(jsonResp)
	}
	// write doc info
	err = stub.PutState("DOCUMENT-"+strconv.Itoa(docIndx), []byte(docInfo))
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to update DOCUMENT-" + strconv.Itoa(docIndx) + " Info \"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Key\":\"" + key + "\", \"DocIndx\":\"DOCUMENT-" + strconv.Itoa(docIndx) + ",\"Value\":\"" + docInfo + "\"}"
	log.Infof("Write Response:%s\n", jsonResp)
	return nil, nil
}

// query callback representing the query of a chaincode
func (t *WFChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	log.Infof("query is running %s\n", function)

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
		log.Infof("Query Response: %s\n", jsonResp)
		return Avalbytes, nil
	} else if function == "readDocuments" {

		Avalbytes, err := t.readDocuments(stub, args)
		if err != nil {
			return nil, err
		}

		if Avalbytes == nil {
			jsonResp := "{\"Error\":\"Nil value for " + args[0] + "\"}"
			return nil, errors.New(jsonResp)
		}

		jsonResp := "{\"Key\":\"" + args[0] + "\",\"Value\":\"BLOCK DATA\"}"
		log.Infof("Query Response: %s\n", jsonResp)
		return Avalbytes, nil
	}
	log.Errorf("query did not find func: %s\n", function)

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
	log.Infof("Running READ function :%s\n", string(logData))
	// reading the state
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *WFChaincode) readDocuments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var docBase = "DOCUMENT-"
	var err error
	var logData, docIndxData []byte
	var pageNum, pageSize, docIndx int

	if len(args) < 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the pageNum, pageSize and LogInfo for query")
	}

	pageNum, err = strconv.Atoi(args[1])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to read pageNum from args 1\"}"
		return nil, errors.New(jsonResp)
	}
	if pageNum > 0 {
		pageSize, err = strconv.Atoi(args[2])
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to read pageSize from args 2\"}"
			return nil, errors.New(jsonResp)
		}
		if pageSize <= 0 {
			pageSize = 15
		}
	}

	logData, _ = b64.StdEncoding.DecodeString(args[1])
	log.Infof("Running readDocuments function :%s\n", string(logData))
	docIndxData, err = stub.GetState("DOCUMENT_INDEX")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to read DOCUMENT_INDEX\"}"
		return nil, errors.New(jsonResp)
	}
	docIndx, err = strconv.Atoi(string(docIndxData))
	var indxStart, indxEnd int
	if pageNum > 0 {
		indxStart = ((pageNum - 1) * pageSize) + 1
		indxEnd = (indxStart + pageSize) - 1
		if indxEnd > docIndx {
			indxEnd = docIndx
		}
	} else {
		indxStart = 1
		indxEnd = docIndx
	}
	var docBaseKey string
	var docData []byte
	jsonResp = "{\"docs\":["
	for x := indxStart; x <= indxEnd; x++ {
		docBaseKey = docBase + strconv.Itoa(x)
		docData, err = stub.GetState(docBaseKey)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get state for " + docBaseKey + "\"}"
			return nil, errors.New(jsonResp)
		}
		jsonResp += "\"" + string(docData) + "\""
		if x < indxEnd {
			jsonResp += ","
		}
	}
	jsonResp += "]}"

	return []byte(jsonResp), nil
}

func main() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)
	err := shim.Start(new(WFChaincode))
	if err != nil {
		log.Errorf("Error starting Wellsfargo chaincode: %s\n", err)
	}
}
