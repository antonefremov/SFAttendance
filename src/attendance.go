package main

import (
	"encoding/json"
	"errors"

	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// The Attendance structure
type Attendance struct {
	ID            string `json:"Id"`            // Id
	ParticipantID string `json:"participantId"` // Participant Id
	SessionID     string `json:"sessionId"`     // Session Id
	SessionCode   string `json:"sessionCode"`   // Session Code
	// DateTime	  string ???
	// Confirmation string `json:"confirmation"` // Confirmation response sent as a bar code or a url link ???
}

// ============================================================================================================================
// Record Attendance - create a record of a new Attendance in the chaincode state
// ============================================================================================================================
func recordAttendance(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting exactly 4 arguments")
	}

	var attendance Attendance
	attendance.ParticipantID = args[1]
	attendance.SessionID = args[2]
	attendance.SessionCode = args[3]
	// compositeKey := []string{args[1], args[2], args[3]}
	// attendance.ID, _ = stub.CreateCompositeKey("", compositeKey) //args[0] //attendance.ParticipantID + attendance.SessionID + attendance.SessionCode
	attendance.ID = attendance.ParticipantID + attendance.SessionID + attendance.SessionCode
	fmt.Println("The composite key created is", attendance.ID)
	// attendance.Confirmation = ""

	existingAttendance, _ := getAttendance(stub, attendance.ID)
	if len(existingAttendance.ID) > 0 {
		return shim.Error("Attendance already exists!")
	}

	// if owner.Id == "" && err != nil {
	// 	return shim.Error("Owner does not exist with Id '" + asset.OwnerId + "'")
	// }

	attendanceAsBytes, _ := json.Marshal(attendance)
	err = stub.PutState(attendance.ID, attendanceAsBytes)
	if err != nil {
		// fmt.Println("Could not save an Attendance")
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// Verify Attendance - return verification on Attendance existence
// ============================================================================================================================
func verifyAttendance(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting exactly 1 argument")
	}
	var id = args[0] //+ args[1] + args[2]

	existingAttendance, _ := getAttendance(stub, id)
	if len(existingAttendance.ID) > 0 {
		return shim.Success([]byte(existingAttendance.ID))
	}

	return shim.Error("Attendance does not exist with Id = " + id)
}

// ============================================================================================================================
// Get Attendance - return an instance from the state
// ============================================================================================================================
func getAttendance(stub shim.ChaincodeStubInterface, id string) (Attendance, error) {
	var attendance Attendance
	// fmt.Println("The key inside the getAttendace is", id)
	attendanceAsBytes, err := stub.GetState(id)
	if err != nil {
		return attendance, err
	}

	// fmt.Println("Attendance as bytes is ", attendanceAsBytes)

	json.Unmarshal(attendanceAsBytes, &attendance)
	if len(attendance.ID) == 0 {
		return attendance, errors.New("Attendance was not found")
	}
	return attendance, nil
}
