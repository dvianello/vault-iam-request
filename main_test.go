package main

import (
	"log"
	"strings"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
)

func TestConcourseFormat(t *testing.T) {
	assert := assert.New(t)

	// Test setup
	testLoginData := map[string]interface{}{
		"iam_request_body": "test-iam-request-body",
		"role":             "test",
		"iam_http_request_method": "POST",
		"iam_request_url":         "test-iam-request-url",
		"iam_request_headers":     "test-iam-request-headers",
	}

	concourseFormat := buildConcourseFormat(testLoginData)

	// Test that all expected sub-strings are in there

	// Rebuild map out of the returned string
	concourseFormatTokens := strings.Split(concourseFormat, ",")
	testMap := make(map[string]interface{})
	for _, e := range concourseFormatTokens {
		parts := strings.Split(e, "=")
		testMap[parts[0]] = strings.Trim(parts[1], "\"")
	}

	// Check length is correct
	assert.Len(testMap, 5)

	// Check map is the same
	assert.Equal(testMap, testLoginData)

}

func TestOutputConfigurationJson(t *testing.T) {
	assert := assert.New(t)

	//
	// Test setup
	//

	// Define custom args for testing
	argsLong := []string{
		"-r test-role",
		"--json",
	}

	argsShort := []string{
		"-r test-role",
		"-j",
	}

	argsNotPresent := []string{
		"-r test-role",
	}
	var optionsLong options
	_, errLong := flags.ParseArgs(&optionsLong, argsLong)
	if errLong != nil {
		log.Fatal(errLong)
		assert.Fail("Something went wrong parsing the long args")
	}

	var optionsShort options
	_, errShort := flags.ParseArgs(&optionsShort, argsShort)
	if errShort != nil {
		log.Fatal(errShort)
		assert.Fail("Something went wrong parsing the short args")
	}

	var optionsNotPresent options
	_, errNotPresent := flags.ParseArgs(&optionsNotPresent, argsNotPresent)
	if errNotPresent != nil {
		log.Fatal(errNotPresent)
		assert.Fail("Something went wrong parsing the non presemt args")
	}

	//
	// Execute tests
	//

	// Long flag
	output := defineOutput(optionsLong)
	assert.True(output.JSON)

	// Short flag
	output = defineOutput(optionsShort)
	assert.True(output.JSON)

	// Not present
	output = defineOutput(optionsNotPresent)
	assert.False(output.JSON)
}
