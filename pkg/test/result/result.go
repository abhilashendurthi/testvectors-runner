/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package result implements functions related to results
*/
package result

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//TestVectorResult stores the testvector name, list of all the test case results and final result of the testvector
type TestVectorResult struct {
	TestVectorName  string
	TestCaseResults []*TestCaseResult
	Result          bool
}

//TestCaseResult stores the test case id, result of the test case execution
type TestCaseResult struct {
	TestCaseName string
	Result       bool
}

var (
	resultDir  = "/tmp/result"
	resultFile = "result"
)

/*type Result interface {
	Print()
}*/

//Print method simulates the go test printing at suite level
/*func (tr *TestVectorResult) Print() {
	if tr.Result {
		fmt.Print("--- PASS: ")
	} else {
		fmt.Printf("--- FAIL: ")
	}
	fmt.Printf(tr.TestVectorName + "\n")
	for _, tcr := range tr.TestCaseResults {
		tcr.Print()
	}
}*/

//Print method simulates the go test printing at test case level
/*func (tcr *TestCaseResult) Print() {
	if tcr.Result {
		fmt.Print("\t--- PASS: ")
	} else {
		fmt.Printf("\t--- FAIL: ")
	}
	fmt.Printf("%s\n", tcr.TestCaseName)
}*/

//NewTestVectorResult returns a TestVectorResult with given test name
func NewTestVectorResult(testName string) *TestVectorResult {
	return &TestVectorResult{TestVectorName: testName,
		TestCaseResults: []*TestCaseResult{},
		Result:          true,
	}
}

//NewTestCaseResult returns a TestCaseResult with given test case name
func NewTestCaseResult(testCaseName string) *TestCaseResult {
	return &TestCaseResult{
		TestCaseName: testCaseName,
	}
}

//PrintStats prints the total number of testvectors executed, passed and failed
func PrintStats(testResults []*TestVectorResult) {
	fmt.Println("\n\n\n" + strings.Repeat("*", 100))
	fmt.Println("Statistics")
	fmt.Println(strings.Repeat("*", 100))
	testVectorCount := len(testResults)
	passedTestVectorCount, failedTestVectorCount := 0, 0

	for _, testResult := range testResults {
		if testResult.Result {
			passedTestVectorCount++
		}
	}
	failedTestVectorCount = testVectorCount - passedTestVectorCount

	fmt.Printf("Total number of test vectors executed: %v\n", testVectorCount)
	fmt.Printf("Total number of test vectors passed: %v\n", passedTestVectorCount)
	fmt.Printf("Total number of test vectors failed: %v\n", failedTestVectorCount)
	fmt.Println(strings.Repeat("*", 100))

}

//WriteResults writes the list of testvector results to csv file
func WriteResults(testVectorResults []*TestVectorResult) {
	file, err := os.Create(resultDir + "/" + resultFile + ".csv")
	if err != nil {
		log.Fatal("Cannot create file ", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	//writer.Write([]string{"Test Name", "Result"})

	for _, testVectorResult := range testVectorResults {
		data := []string{testVectorResult.TestVectorName, strconv.FormatBool(testVectorResult.Result)}
		err := writer.Write(data)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}
}

//SetResultFolder sets the resultDir variable and creates the directory
func SetResultFolder(resultDirName string) {
	resultDir = resultDirName
	err := os.MkdirAll(resultDirName, 0755)
	if err != nil {
		log.Fatal("Cannot create directory ", err)
	}
}

//SetResultFileName sets the resultFile variable
func SetResultFileName(resultFileName string) {
	resultFile = resultFileName
}
