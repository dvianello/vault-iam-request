package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/vault/builtin/credential/aws"
	"github.com/jessevdk/go-flags"
)

// Define command line options
type options struct {
	Role string `short:"r" long:"role" description:"The Vault role to authenticate against" required:"true"`
	JSON bool   `short:"j" long:"json" description:"Output data in JSON format"`
	File string `short:"f" long:"file" description:"Write output to file instead of stdout"`
}

// Struct to define output configuration
type STSCall struct {
	Role string
	File os.File
	JSON bool
	Content map[string]interface{}
}

//GenerateLoginData
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

func (call *STSCall) defineOutput(file string, json bool) {
	if file != "" {
		fileWriter, err := os.Create(file)
		if err != nil {
			fmt.Println("File does not exists or cannot be created")
			os.Exit(1)
		}
		call.File = *fileWriter
		defer fileWriter.Close()

	} else {
		call.File = *os.Stdout
	}

	if json {
		call.JSON = true
	}
}

func (call *STSCall) WriteOutput(file string, JSONOutput bool){
	call.defineOutput(file, JSONOutput)

	if call.JSON {
		jsonLoginData, _ := json.Marshal(call.Content)
		fmt.Fprintln(&call.File, string(jsonLoginData))
		os.Exit(0)
	}

	fmt.Fprintln(&call.File, call.buildConcourseFormat())

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
		os.Exit(1)
	}


	var call STSCall
	// Leverage Vault awsauth to generate LoginData
	call.GenerateLoginData(options.Role)
	call.WriteOutput(options.File, options.JSON)
}
