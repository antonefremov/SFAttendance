package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ============================================================================================================================
// Init Chaincode - Mock stub init wrapper
// ============================================================================================================================
func InitChaincode(test *testing.T) *shim.MockStub {
	stub := shim.NewMockStub("testingStub", new(SmartContract))
	result := stub.MockInit("000", nil)

	if result.Status != shim.OK {
		test.FailNow()
	}
	return stub
}

// ============================================================================================================================
// Invoke wrapper
// ============================================================================================================================
func Invoke(test *testing.T, stub *shim.MockStub, function string, args [][]byte) []byte {
	const transactionID = "000"

	// prepend the function name as the first item
	args = append([][]byte{[]byte(function)}, args...)

	// prepare the parameters for printing
	byteDivider := []byte{','}
	byteArrayToPrint := bytes.Join(args[1:], byteDivider)

	// print information just before the call
	fmt.Println("Call:    ", function, "(", string(byteArrayToPrint), ")")

	// perform the MockInvoke call
	result := stub.MockInvoke(transactionID, args)

	// print the Invoke results
	fmt.Println("RetCode: ", result.Status)
	fmt.Println("RetMsg:  ", result.Message)
	fmt.Println("Payload: ", string(result.Payload))

	if result.Status != shim.OK {
		fmt.Println("Invoke", function, "failed", string(result.Message))
		return nil
	}

	return []byte(result.Payload)
}

// ============================================================================================================================
// Get a mock Attendance
// ============================================================================================================================
func GetAttendaceForTesting() [][]byte {
	return [][]byte{
		// []byte("I300455BC4308531"), // ID
		[]byte("I300455"), // ParticipantID
		[]byte("BC430"),   // SessionID
		[]byte("8531")}    // SessionCode
}

func GetAttendanceForTestingKey() [][]byte {
	return [][]byte{[]byte("I300455BC4308531")}
}

// ============================================================================================================================
// Convert the Attendance passed in as bytes to an Attendance instance presented as bytes
// ============================================================================================================================
func ConvertBytesToAttendanceAsBytes(attendanceAsBytes [][]byte) []byte {
	var attendance Attendance
	attendance.ID = string(attendanceAsBytes[0])
	attendance.ParticipantID = string(attendanceAsBytes[1])
	attendance.SessionID = string(attendanceAsBytes[2])
	attendance.SessionCode = string(attendanceAsBytes[3])
	bagJSON, err := json.Marshal(attendance)
	if err != nil {
		fmt.Println("Error converting an Attendance record to JSON")
		return nil
	}
	return []byte(bagJSON)
}
