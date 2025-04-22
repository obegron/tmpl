package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

var ctx = make(map[string]interface{})

func loadContext() {
	for k, v := range getEnvTemplateVariables() {
		ctx[k] = v
	}

	for _, file := range FilesList {
		for k, v := range getFileVariables(file) {
			ctx[k] = v
		}
	}

	for k, v := range getCliVariables() {
		ctx[k] = v
	}
}

func getEnvTemplateVariables() map[string]interface{} {
	vars := make(map[string]interface{})

	tmplVarsBase64 := os.Getenv("TMPLVARS")
	if tmplVarsBase64 == "" {
		return vars
	}

	tmplData, err := base64.StdEncoding.DecodeString(tmplVarsBase64)
	if err != nil {
		log.Printf("Content of TMPLVARS must be base64 encoded: %v", err)
		return vars
	}

	tmplVars := make(map[string]interface{})
	err = yaml.Unmarshal(tmplData, &tmplVars)
	if err == nil {
		for k, v := range tmplVars {
			vars[k] = v
		}
		log.Printf("Loaded %d variables from TMPLVARS environment variable as YAML", len(tmplVars))
		return vars
	}

	err = json.Unmarshal(tmplData, &tmplVars)
	if err == nil {
		for k, v := range tmplVars {
			vars[k] = v
		}
		log.Printf("Loaded %d variables from TMPLVARS environment variable as JSON", len(tmplVars))
		return vars
	}
	log.Printf("Error parsing TMPLVARS as either YAML or JSON")
	return vars
}

func getFileVariables(file string) map[string]interface{} {
	vars := make(map[string]interface{})

	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Printf("unable to read file\n%v\n", err)
		return vars
	}

	if strings.HasSuffix(file, ".json") {
		err = json.Unmarshal(bytes, &vars)
	} else if strings.HasSuffix(file, ".toml") {
		err = toml.Unmarshal(bytes, &vars)
	} else if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
		err = yaml.Unmarshal(bytes, &vars)
	} else {
		err = fmt.Errorf("bad file type: %s", file)
	}
	if err != nil {
		log.Printf("unable to load data\n%v\n", err)
	}
	return vars
}

func getCliVariables() map[string]string {
	vars := make(map[string]string)
	for _, pair := range VarsList {
		kv := strings.SplitN(pair, "=", 2)

		v := kv[1]
		if strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"") {
			v = v[1 : len(v)-1]
		} else if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			v = v[1 : len(v)-1]
		}

		vars[kv[0]] = v
	}
	return vars
}
