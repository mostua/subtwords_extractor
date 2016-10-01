package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/mostua/subtwords_extractor/parser"
)

type inputOptions struct {
	fileName string
}

type parseFlagError struct {
	message string
}

func newParseFlagError(message string) *parseFlagError {
	var result *parseFlagError
	result = &parseFlagError{}
	result.message = message
	return result
}

func parseFlag() (*inputOptions, *parseFlagError) {
	var options *inputOptions
	options = &inputOptions{}
	defaultVal := ""
	flag.StringVar(&options.fileName, "f", defaultVal, "Subtitles filename")
	flag.Parse()
	if options.fileName == defaultVal {
		return nil, newParseFlagError("f parameter is missing")
	}
	return options, nil
}

func main() {
	inputOptions, flagErr := parseFlag()
	if flagErr != nil {
		fmt.Printf("Invalid params.\n")
		fmt.Printf("%s\n", flagErr.message)
		return
	}
	fmt.Printf("Filename: %s\n", inputOptions.fileName)
	bytes, err := ioutil.ReadFile(inputOptions.fileName)
	if err != nil {
		fmt.Printf("Unexpected error occured %v", err)
		return
	}
	subt := string(bytes)
	parser := &parser.SrtParser{}
	_, parseErr := parser.Parse(subt)
	if parseErr != nil {
		fmt.Printf("Parsing error %v", parseErr)
	}

}
