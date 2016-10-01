package main

import "fmt"
import "flag"

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
	inputOptions, err := parseFlag()
	if err != nil {
		fmt.Printf("Invalid params.\n")
		fmt.Printf("%s\n", err.message)
		return
	}
	fmt.Printf("Filename: %s\n", inputOptions.fileName)
}
