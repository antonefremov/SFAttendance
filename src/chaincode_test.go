package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestAttendanceCreateRead(test *testing.T) {

	// prepare all the necessary objects and keys
	stub := InitChaincode(test)
	attendanceForTesting := GetAttendaceForTesting()
	attendanceID := GetAttendanceForTestingKey()[0] //attendanceForTesting[0]
	//attendanceForTestingAsBytes := ConvertBytesToAttendanceAsBytes(attendanceForTesting)
	attendanceForTestingKey := GetAttendanceForTestingKey() //[][]byte{attendanceForTesting[0]}

	// invoke the functions
	Invoke(test, stub, "recordAttendance", attendanceForTesting)
	attendanceIDAsBytes := Invoke(test, stub, "verifyAttendance", attendanceForTestingKey)

	// check the results
	if bytes.Compare(attendanceID, attendanceIDAsBytes) != 0 {
		fmt.Println("\n>>> FAILED TEST: verifyAttendance.\n", "\nExpected:\n", string(attendanceID), "\nActual:\n", string(attendanceIDAsBytes), "\n ")
		test.FailNow()
	}
}

func TestAttendanceNotFound(test *testing.T) {
	// prepare all the necessary objects and keys
	stub := InitChaincode(test)
	attendanceForTestingKey := [][]byte{[]byte("1234")}

	// invoke the functions
	attendanceIDAsBytes := Invoke(test, stub, "verifyAttendance", attendanceForTestingKey)

	if len(attendanceIDAsBytes) > 0 {
		test.FailNow()
	}
}
