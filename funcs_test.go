package main

import (
	"bytes"
	"testing"
	"text/template"
)

func TestReplace(t *testing.T) {
	tmpl, err := template.New("test").Funcs(getFuncMap()).Parse(`{{ replace "hello:world" ":" "-" }}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	expected := "hello-world"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestReplaceMultiple(t *testing.T) {
	tmpl, err := template.New("test").Funcs(getFuncMap()).Parse(`{{ replace "a:b:c" ":" "-" }}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	expected := "a-b-c"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestReplaceNotFound(t *testing.T) {
	tmpl, err := template.New("test").Funcs(getFuncMap()).Parse(`{{ replace "abc" "x" "-" }}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	expected := "abc"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestReplaceEmpty(t *testing.T) {
	tmpl, err := template.New("test").Funcs(getFuncMap()).Parse(`{{ replace "" "x" "-" }}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	expected := ""
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestReplaceWithEmptyString(t *testing.T) {
	tmpl, err := template.New("test").Funcs(getFuncMap()).Parse(`{{ replace "hello:world" ":" "" }}`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	expected := "helloworld"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}
