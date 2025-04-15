package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

func executeTemplateFile(input string) (string, error) {
	var (
		err          error
		outputBytes  bytes.Buffer
		outputString string
	)

	if len(ctx) == 0 {
		loadContext()
	}

	if strings.Contains(input, "{{") && strings.Contains(input, "}}") {
		log.Printf("Processing filename as template: %s", input)
		filenameTemplate, err := template.New("filename").Funcs(getFuncMap()).Parse(input)
		if err != nil {
			return "", fmt.Errorf("error parsing filename template: %w", err)
		}

		var processedFilename bytes.Buffer
		if err := filenameTemplate.Execute(&processedFilename, ctx); err != nil {
			return "", fmt.Errorf("error executing filename template: %w", err)
		}

		input = processedFilename.String()
		log.Printf("Processed filename: %s", input)
	}

	inputBytes, err := readInput(input)
	if err != nil {
		return "", err
	}

	tmpl := template.New(input)
	tmpl.Funcs(getFuncMap())

	tmpl, err = tmpl.Parse(string(inputBytes))
	if err != nil {
		return "", err
	}
	if Strict {
		tmpl.Option("missingkey=error")
	}

	err = tmpl.Execute(&outputBytes, ctx)
	if err != nil {
		return "", err
	}

	outputString = outputBytes.String()
	outputString = strings.ReplaceAll(outputString, "<no value>", "")
	return outputString, nil
}

func readInput(input string) ([]byte, error) {
	var (
		err        error
		inputBytes []byte
	)
	if input == "-" {
		inputBytes, err = io.ReadAll(os.Stdin)
	} else {
		inputBytes, err = os.ReadFile(input)
	}
	return inputBytes, err
}
