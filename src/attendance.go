package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// The Attendance structure
type Attendance struct {
	ID                           string    `json:"Id"`                           // Id
	Cust_attendance_externalCode string    `json:"cust_attendance_externalCode"` // First part of the key
	ExternalCode                 string    `json:"externalCode"`                 // Second part of the key
	Cust_session_id              string    `json:"cust_session_id"`              // Cust Session Id
	Cust_session_code            string    `json:"cust_session_code"`            // Cust Session Code
	LastModifiedBy               string    `json:"lastModifiedBy"`               // Last Modified By (Participant)
	ExternalName                 string    `json:"externalName"`                 // Participant's Full Name
	LastModifiedDateTime         time.Time `json:"lastModifiedDateTime"`         // Last Modified Time
	// DateTime	  string ???
	// Confirmation string `json:"confirmation"` // Confirmation response sent as a bar code or a url link ???
}

// ============================================================================================================================
// Record Attendance - create a record of a new Attendance in the chaincode state
// ============================================================================================================================
func recordAttendance(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting exactly 6 arguments")
	}

	var attendance Attendance
	attendance.Cust_attendance_externalCode = args[0]
	attendance.ExternalCode = args[1]
	attendance.Cust_session_id = args[2]
	attendance.Cust_session_code = args[3]
	attendance.LastModifiedBy = args[4]
	attendance.ExternalName = args[5]
	attendance.LastModifiedDateTime = time.Now()
	attendance.ID = stub.GetTxID() // attendance.Cust_attendance_externalCode + attendance.ExternalCode

	existingAttendance, _ := getAttendance(stub, attendance.ID)
	if len(existingAttendance.ID) > 0 {
		return shim.Error("Attendance already exists!")
	}

	attendanceAsBytes, _ := json.Marshal(attendance)
	err = stub.PutState(attendance.ID, attendanceAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// Verify Attendance - return Attendance Id as verification of its existence
// ============================================================================================================================
func verifyAttendance(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting exactly 1 argument")
	}
	var id = args[0]

	existingAttendance, _ := getAttendance(stub, id)
	if len(existingAttendance.ID) > 0 {
		return shim.Success([]byte(existingAttendance.ID))
	}

	return shim.Error("Attendance does not exist with TrId = " + id)
}

// ============================================================================================================================
// Read Attendance - return Attendance as sequence of bytes
// ============================================================================================================================
func readAttendance(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting exactly 1 argument")
	}
	var id = args[0]

	attendanceAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get state for Attendance with TrId " + id + "\"}")
	}

	return shim.Success(attendanceAsBytes)
}

// ============================================================================================================================
// Get Attendance - return an instance from the state
// ============================================================================================================================
func getAttendance(stub shim.ChaincodeStubInterface, id string) (Attendance, error) {
	var attendance Attendance

	attendanceAsBytes, err := stub.GetState(id)
	if err != nil {
		return attendance, err
	}

	json.Unmarshal(attendanceAsBytes, &attendance)
	if len(attendance.ID) == 0 {
		return attendance, errors.New("Attendance was not found")
	}
	return attendance, nil
}
