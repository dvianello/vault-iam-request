package main

import (
	"log"
	"strings"
	"testing"

	"encoding/json"
	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"reflect"
)

func rebuildMapFromConcourse(concourseFormat string) map[string]interface{} {
	concourseFormatTokens := strings.Split(concourseFormat, ",")
	testMap := make(map[string]interface{})
	for _, e := range concourseFormatTokens {
		parts := strings.Split(e, "=")
		testMap[parts[0]] = strings.Trim(strings.TrimSuffix(parts[1], "\n"), "\"")
	}

	return testMap
}

func TestConcourseFormat(t *testing.T) {
	assert := assert.New(t)

	var call STSCall

	// Test setup
	testLoginData := map[string]interface{}{
		"iam_request_body":        "test-iam-request-body",
		"role":                    "test",
		"iam_http_request_method": "POST",
		"iam_request_url":         "test-iam-request-url",
		"iam_request_headers":     "test-iam-request-headers",
	}

	call.Content = testLoginData
	concourseFormat := call.buildConcourseFormat()

	// Test that all expected sub-strings are in there

	// Rebuild map out of the returned string
	concourseFormatData := rebuildMapFromConcourse(concourseFormat)

	// Check length is correct
	assert.Len(concourseFormatData, 5)

	// Check map is the same
	assert.Equal(concourseFormatData, testLoginData)

}

func TestOutputConfigurationFormatJSONLong(t *testing.T) {
	assert := assert.New(t)

	//
	// Test setup
	//

	// Define custom args for testing
	var call STSCall

	args := []string{
		"-r test-role",
		"--json",
	}

	var options options
	_, errLong := flags.ParseArgs(&options, args)
	if errLong != nil {
		log.Fatal(errLong)
		assert.Fail("Something went wrong parsing the long args")
	}

	//
	// Execute tests
	//

	// Long flag
	err := call.defineOutput(options.File, options.JSON)
	if err != nil {
		assert.Fail("non-nil error")
	}

	assert.True(call.JSON)

}

func TestOutputConfigurationFormatJSONShort(t *testing.T) {
	assert := assert.New(t)

	//
	// Test setup
	//

	// Define custom args for testing
	var call STSCall

	args := []string{
		"-r test-role",
		"-j",
	}

	var options options
	_, errLong := flags.ParseArgs(&options, args)
	if errLong != nil {
		log.Fatal(errLong)
		assert.Fail("Something went wrong parsing the long args")
	}

	//
	// Execute test
	//

	// Long flag
	err := call.defineOutput(options.File, options.JSON)
	if err != nil {
		assert.Fail("non-nil error")
	}

	assert.True(call.JSON)
}

func TestOutputConfigurationFormatNoJSON(t *testing.T) {
	assert := assert.New(t)

	//
	// Test setup
	//

	// Define custom args for testing
	var call STSCall

	args := []string{
		"-r test-role",
	}

	var options options
	_, errShort := flags.ParseArgs(&options, args)
	if errShort != nil {
		log.Fatal(errShort)
		assert.Fail("Something went wrong parsing the short args")
	}

	//
	// Execute tests
	//

	// Short flag
	err := call.defineOutput(options.File, options.JSON)
	if err != nil {
		assert.Fail("non-nil error")
	}

	assert.False(call.JSON)

}

func TestOutputConfigurationFile(t *testing.T) {
	assert := assert.New(t)
	tmpFile := "/tmp/file"

	var call STSCall

	err := call.defineOutput(tmpFile, false)
	if err != nil {
		assert.Fail("non-nil error")
	}

	assert.Equal(reflect.TypeOf(call.File), reflect.TypeOf(os.File{}))
	assert.FileExists(tmpFile)

	// cleanup
	os.Remove(tmpFile)
}

func TestDefineOutputConfigurationFileError(t *testing.T) {
	assert := assert.New(t)
	tmpFile := "/test"

	var call STSCall

	err := call.defineOutput(tmpFile, false)
	assert.Error(err, "path does not exists or cannot be created")
}

func TestWriteOutputConfigurationFileError(t *testing.T) {
	assert := assert.New(t)
	tmpFile := "/test"

	var call STSCall

	err := call.writeOutput(tmpFile, false)
	assert.Error(err, "path does not exists or cannot be created")
}

func TestOutputConfigurationFileJSON(t *testing.T) {
	assert := assert.New(t)

	tmpFile := "/tmp/json_file"
	var testContent = map[string]interface{}{
		"iam_request_body":        "test-iam-request-body",
		"role":                    "test",
		"iam_http_request_method": "POST",
		"iam_request_url":         "test-iam-request-url",
		"iam_request_headers":     "test-iam-request-headers",
	}

	var call STSCall

	// Load fake data in STSCall
	call.Content = testContent

	// Write out data & check if file exists
	err := call.writeOutput(tmpFile, true)
	if err != nil {
		assert.Fail("non-nil error")
	}

	assert.FileExists(tmpFile)

	// Try to parse JSON back & compare to original
	var data map[string]interface{}
	jsonFileContent, _ := ioutil.ReadFile(tmpFile)
	err = json.Unmarshal(jsonFileContent, &data)
	if err != nil {
		assert.Fail("Failed to unmarshal the JSON file.")
	}

	assert.Equal(testContent, data)

	// cleanup
	os.Remove(tmpFile)

}

func TestOutputConfigurationFileConcourse(t *testing.T) {
	assert := assert.New(t)

	tmpFile := "/tmp/concourse_file"
	var testContent = map[string]interface{}{
		"iam_request_body":        "test-iam-request-body",
		"role":                    "test",
		"iam_http_request_method": "POST",
		"iam_request_url":         "test-iam-request-url",
		"iam_request_headers":     "test-iam-request-headers",
	}

	var call STSCall

	// Load fake data in STSCall
	call.Content = testContent

	// Write out data & check if file exists
	err := call.writeOutput(tmpFile, false)
	if err != nil {
		assert.Fail("non-nil error")
	}

	assert.FileExists(tmpFile)

	// Try to parse Concourse format back & compare to original
	concourseFileContent, _ := ioutil.ReadFile(tmpFile)
	concourseFileData := rebuildMapFromConcourse(string(concourseFileContent))
	assert.Equal(concourseFileData, testContent)

	// cleanup
	os.Remove(tmpFile)

}

func TestFailedAWSCall(t *testing.T) {
	assert := assert.New(t)
	role := "testRole"

	var call STSCall

	err := call.generateLoginData(role)
	if err == nil {
		assert.Fail("Nil error, but it should have been non-nil")
		println(err)
	}

	assert.Error(err, "call to AWS failed. Check your credentials")
}

func TestFailedAWSCallBuild(t *testing.T) {
	assert := assert.New(t)
	role := "testRole"
	file := "/testfile"
	json := false

	var call STSCall

	err := call.BuildCall(role, file, json)
	if err == nil {
		assert.Fail("Nil error, but it should have been non-nil")
		println(err)
	}

	assert.Error(err, "call to AWS failed. Check your credentials")

}

//func TestSuccessfullAWSCall(t *testing.T){
//	assert := assert.New(t)
//	role := "testRole"
//
//
//}
