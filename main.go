package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"errors"
	"github.com/hashicorp/vault/builtin/credential/aws"
	"github.com/jessevdk/go-flags"
)

// Define command line options
type options struct {
	Role string `short:"r" long:"role" description:"The Vault role to authenticate against" required:"true"`
	JSON bool   `short:"j" long:"json" description:"Output data in JSON format"`
	File string `short:"f" long:"file" description:"Write output to file instead of stdout"`
}

// STSCall Struct to define output configuration
type STSCall struct {
	Role    string
	File    os.File
	JSON    bool
	Content map[string]interface{}
}

// GenerateLoginData builds the STS call via the Vault awsauth lib and adds
// a 'role' key to it.
func (call *STSCall) GenerateLoginData(role string) {

	// Generate login data via awsauth package
	loginData, err := awsauth.GenerateLoginData("", "", "", "")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Add role to loginData since we need to send it along
	// when authenticating to Vault
	loginData["role"] = role
	call.Content = loginData
}

func (call *STSCall) defineOutput(file string, json bool) (err error) {
	if file != "" {
		fileWriter, err := os.Create(file)
		if err != nil {
			return errors.New("path does not exists or cannot be created")
		}
		call.File = *fileWriter

	} else {
		call.File = *os.Stdout
	}

	if json {
		call.JSON = true
	}

	return
}

// WriteOutput writes the STS call defined output
func (call *STSCall) WriteOutput(file string, JSONOutput bool) (err error) {
	err = call.defineOutput(file, JSONOutput)

	if err != nil {
		return err
	}

	if call.JSON {
		jsonLoginData, _ := json.Marshal(call.Content)
		fmt.Fprintln(&call.File, string(jsonLoginData))

	} else {
		fmt.Fprintln(&call.File, call.buildConcourseFormat())
	}

	call.File.Close()
	return
}

func (call *STSCall) buildConcourseFormat() string {
	buffer := make([]string, 0)
	for k, v := range call.Content {
		buffer = append(buffer, fmt.Sprintf("%s=\"%s\"", k, v))
	}

	return strings.Join(buffer, ",")
}

func main() {

	// Define flags
	var options options

	// Parse command line flags
	_, err := flags.Parse(&options)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var call STSCall
	// Leverage Vault awsauth to generate LoginData
	call.GenerateLoginData(options.Role)
	err = call.WriteOutput(options.File, options.JSON)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
