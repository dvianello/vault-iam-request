package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
)


func TestConcourseFormat(t *testing.T){
	assert := assert.New(t)

	// Setup test
	testLoginData := map[string]interface{}{
		"iam_request_body"			: "test-iam-request-body",
		"role"						: "test",
		"iam_http_request_method" 	: "POST",
		"iam_request_url" 			: "test-iam-request-url",
		"iam_request_headers" 		: "test-iam-request-headers",
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
