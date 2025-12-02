package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

var inputsDir = "."

func getFuncMap() template.FuncMap {
	f := make(template.FuncMap)

	f["fromInputDir"] = fromInputDir
	f["include"] = include
	f["expandenv"] = expandenv

	f["toBool"] = toBool
	f["toToml"] = toToml
	f["toYaml"] = toYaml
	f["toJson"] = toJson

	f["absPath"] = absPath
	f["fileExists"] = fileExists
	f["fileMode"] = fileMode
	f["fileMtime"] = fileMtime
	f["fileRead"] = fileRead
	f["fileSize"] = fileSize
	f["isDir"] = isDir
	f["isFile"] = isFile
	f["joinPath"] = joinPath
	f["toBackslash"] = toBackslash
	f["toOsPath"] = toOsPath
	f["toSlash"] = toSlash

	f["env"] = getEnv

	f["now"] = currentTime
	f["urlHostname"] = urlHostname
	f["urlPort"] = urlPort

	f["wrapWith"] = wrapWith
	f["uuidv4"] = uuidv4

	return f
}

// --- Includes ------------------------------------------------------------------------------------

func include(input string) (string, error) {
	var err error
	includeDir := inputsDir
	if !filepath.IsAbs(input) {
		includeDir, err = getIncludeDir(input)
		if err != nil {
			return "", err
		}
	}
	input = filepath.Join(includeDir, input)

	outputString, err := executeTemplateFile(input)
	if err != nil {
		return "", err
	}
	return outputString, nil
}

func fromInputDir(input string) (string, error) {
	dir, err := getIncludeDir(input)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, input), nil
}

func getIncludeDir(input string) (string, error) {
	if input == "-" {
		cwd, err := filepath.Abs(".")
		if err != nil {
			return "", err
		}
		return cwd, nil
	} else {
		return inputsDir, nil
	}
}

func expandenv(input string) string {
	return os.ExpandEnv(input)
}

// --- Type conversion -----------------------------------------------------------------------------

func toBool(value string) (bool, error) {
	result, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}
	return result, nil
}

func toToml(v interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	enc := toml.NewEncoder(buf)
	err := enc.Encode(v)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func toYaml(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func toJson(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// --- Paths ---------------------------------------------------------------------------------------

func absPath(file string) (string, error) {
	new_file, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return new_file, nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func fileMode(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return info.Mode().String(), nil
}

func fileMtime(file string) (string, error) {
	info, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return info.ModTime().Format(time.RFC3339), nil
}

func fileRead(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

func fileSize(file string) (int64, error) {
	info, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	return info.Size(), nil
}

func isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func isFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.Mode().IsRegular(), nil
}

func joinPath(segments []any) string {
	string_segments := make([]string, 0, len(segments))
	for _, v := range segments {
		string_segments = append(string_segments, fmt.Sprint(v))
	}
	return filepath.Join(string_segments...)
}

func toBackslash(path string) string {
	return strings.ReplaceAll(path, "/", "\\")
}

func toOsPath(path string) string {
	return strings.ReplaceAll(path, "\\", string(filepath.Separator))
}

func toSlash(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func getEnv(name string) string {
	return os.Getenv(name)
}

// utils

func currentTime() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05Z")
}

func urlHostname(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		// fallback: strip scheme manually
		rawURL = strings.TrimPrefix(rawURL, "http://")
		rawURL = strings.TrimPrefix(rawURL, "https://")
		return strings.SplitN(rawURL, "/", 2)[0]
	}
	return u.Host
}

func urlPort(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	switch u.Scheme {
	case "https":
		return "443"
	case "http":
		return "80"
	default:
		return u.Port()
	}
}

func wrapWith(width int, newline, input string) string {
	if width <= 0 {
		return input // No wrapping
	}

	words := strings.Fields(input)
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	currentLineLength := 0

	for i, word := range words {
		if currentLineLength == 0 {
			// Start of a new line
			result.WriteString(word)
			currentLineLength = len(word)
		} else if currentLineLength+1+len(word) <= width {
			// Word fits on current line
			result.WriteString(" ")
			result.WriteString(word)
			currentLineLength += 1 + len(word)
		} else {
			// Word doesn't fit, start a new line
			result.WriteString(newline)
			result.WriteString(word)
			currentLineLength = len(word)
		}

		if i < len(words)-1 && currentLineLength >= width {
			result.WriteString(newline)
			currentLineLength = 0
		}
	}
	return result.String()
}

func uuidv4() string {
	return uuid.New().String()
}

