package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//SmartContract is the data structure which represents this contract and on which various lifecycle functions are attached
type SmartContract struct {
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

// Init is called when the smart contract is instantiated
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke routes invocations to the appropriate function in chaincode
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	// Handle invoke functions
	if function == "init" { // initialise the chaincode state
		return s.Init(stub)
	} else if function == "recordAttendance" { // create an attendance record
		return recordAttendance(stub, args)
	} else if function == "verifyAttendance" { // verify if attendance exists
		return verifyAttendance(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name")
}
