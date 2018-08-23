package main

import (
	"os"
	"github.com/jessevdk/go-flags"
	"github.com/hashicorp/vault/builtin/credential/aws"
	"log"
	"fmt"
	"encoding/json"
	"strings"
)

func buildConcourseFormat(rawLoginData map[string]interface{}) string {
	buffer := make([]string, 0)
	for k, v := range rawLoginData {
		buffer = append(buffer, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	return strings.Join(buffer, ",")
}

func main() {

	// Define flags
	var options struct {
		Role string `short:"r" long:"role" description:"The Vault role to authenticate against" required:"true"`
		Json bool   `short:"j" long:"json" description:"Output data in JSON format"`
		File string `short:"f" long:"file" description:"Write output to file instead of stdout"`
	}

	// Use struct to store file used to
	// write output to
	var output struct {
		file os.File
	}

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

	// Define writer
	if options.File != "" {
		fileWriter, err := os.Create(options.File)
		if err != nil {
			fmt.Println("File does not exists or cannot be created")
			os.Exit(1)
		}
		output.file = *fileWriter
		defer fileWriter.Close()
	} else {
		output.file = *os.Stdout
	}

	// If JSON output is required, use encoding/json to
	// marshall that out
	if options.Json {
		jsonLoginData, _ := json.Marshal(loginData)
		fmt.Fprintln(&output.file, string(jsonLoginData))
		os.Exit(0)
	}

	// Build format required by Concourse and print it out.
	fmt.Fprintln(&output.file, buildConcourseFormat(loginData))

}
