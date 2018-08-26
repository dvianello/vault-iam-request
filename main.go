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

func buildConcourseFormat(rawLoginData map[string]interface{}) string {
	buffer := make([]string, 0)
	for k, v := range rawLoginData {
		buffer = append(buffer, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	return strings.Join(buffer, ",")
}

// Define command line options
type options struct {
	Role string `short:"r" long:"role" description:"The Vault role to authenticate against" required:"true"`
	JSON bool   `short:"j" long:"json" description:"Output data in JSON format"`
	File string `short:"f" long:"file" description:"Write output to file instead of stdout"`
}

// Struct to define output configuratioj
type output struct {
	File os.File
	JSON bool
}

func defineOutput(options options) output {
	var output output

	if options.File != "" {
		fileWriter, err := os.Create(options.File)
		if err != nil {
			fmt.Println("File does not exists or cannot be created")
			os.Exit(1)
		}
		output.File = *fileWriter
		defer fileWriter.Close()
	} else {
		output.File = *os.Stdout
	}

	if options.JSON {
		output.JSON = true
	}

	return output

}

func main() {

	// Define flags
	var options options

	// Parse command line flags
	_, err := flags.Parse(&options)
	if err != nil {
		os.Exit(1)
	}

	// Leverage Vault awsauth to generate LoginData
	loginData, err := awsauth.GenerateLoginData("", "", "", "")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Add role to login data, as we need to send that along as well to the
	// vault server.
	loginData["role"] = options.Role

	// Define output configuration
	output := defineOutput(options)

	// If JSON output is required, use encoding/json to
	// marshall that out
	if options.JSON {
		jsonLoginData, _ := json.Marshal(loginData)
		fmt.Fprintln(&output.File, string(jsonLoginData))
		os.Exit(0)
	}

	// Build format required by Concourse and print it out.
	fmt.Fprintln(&output.File, buildConcourseFormat(loginData))

}
