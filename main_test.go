package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"github.com/jessevdk/go-flags"
	"log"
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
	args_long := []string{
		"-r test-role",
		"--json",
	}

	args_short := []string{
		"-r test-role",
		"-j",
	}

	args_not_present := []string{
		"-r test-role",
	}
	var options_long options
	args_long, err_long := flags.ParseArgs(&options_long, args_long)
	if err_long != nil {
		log.Fatal(err_long)
		assert.Fail("Something went wrong parsing the long args")
	}

	var options_short options
	args_short, err_short := flags.ParseArgs(&options_short, args_short)
	if err_short != nil {
		log.Fatal(err_short)
		assert.Fail("Something went wrong parsing the short args")
	}

	var options_not_present options
	args_not_present, err_not_present := flags.ParseArgs(&options_not_present, args_not_present)
	if err_not_present != nil {
		log.Fatal(err_not_present)
		assert.Fail("Something went wrong parsing the non presemt args")
	}

	//
	// Execute tests
	//

	// Long flag
	output := defineOutput(options_long)
	assert.True(output.JSON)

	// Short flag
	output = defineOutput(options_short)
	assert.True(output.JSON)

	// Not present
	output = defineOutput(options_not_present)
	assert.False(output.JSON)
}