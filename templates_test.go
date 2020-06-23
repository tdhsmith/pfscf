package main

import (
	"path/filepath"
	"reflect"
	"runtime/debug"
	"testing"
)

var templateTestDir string

func init() {
	templateTestDir = filepath.Join(GetExecutableDir(), "testdata", "templates")
}

func expectEqual(t *testing.T, exp interface{}, got interface{}) {
	if exp == got {
		return
	}
	debug.PrintStack()
	t.Errorf("Expected '%v' (type %v), got '%v' (type %v)", exp, reflect.TypeOf(exp), got, reflect.TypeOf(got))
}

func expectNotEqual(t *testing.T, notExp interface{}, got interface{}) {
	typeNotExp := reflect.TypeOf(notExp)
	typeGot := reflect.TypeOf(got)

	// we always require that both types are identical.
	// Without that, testing can be a real pain
	if typeNotExp != typeGot {
		debug.PrintStack()
		t.Errorf("Types do not match! Expected '%v', got '%v'", typeNotExp, typeGot)
		return
	}

	if notExp == got {
		debug.PrintStack()
		t.Errorf("Expected something different than '%v' (type %v)", notExp, typeNotExp)
	}
}

func expectNil(t *testing.T, got interface{}) {
	// do NOT use with errors! This can lead to strange results
	if !reflect.ValueOf(got).IsNil() {
		debug.PrintStack()
		t.Errorf("Expected nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

func expectNotNil(t *testing.T, got interface{}) {
	// do NOT use with errors! This can lead to strange results
	if reflect.ValueOf(got).IsNil() {
		debug.PrintStack()
		t.Errorf("Expected not nil, got '%v' (Type %v)", got, reflect.TypeOf(got))
	}
}

func expectError(t *testing.T, err error) {
	if err == nil {
		debug.PrintStack()
		t.Error("Expected an error, got nil")
	}
}

func expectNoError(t *testing.T, err error) {
	if err != nil {
		debug.PrintStack()
		t.Errorf("Expected no error, got '%v'", err)
	}
}

func TestGetYamlFile_NonExistantFile(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "nonExistantFile.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func TestGetYamlFile_MalformedFile(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "malformed.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func TestGetYamlFile_ValidFile(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "valid.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	// test values contained in the yaml file
	expectEqual(t, "Helvetica", yFile.Default.Font)
	expectEqual(t, 14.0, yFile.Default.Fontsize)
	expectEqual(t, 2, len(yFile.Content))

	content0 := &(yFile.Content[0])
	expectEqual(t, "foo", content0.ID)
	expectEqual(t, "textCell", content0.Type)

	content1 := &(yFile.Content[1])
	expectEqual(t, "bar", content1.ID)
	expectEqual(t, "textCell", content1.Type)
}

func TestGetYamlFile_EmptyFile(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "empty.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	// empty file => no default section and no content
	expectNil(t, yFile.Default)
	expectEqual(t, 0, len(yFile.Content))
}

func TestGetYamlFile_EmptyContentEntry(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "emptyContentEntry.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNotNil(t, yFile)
	expectNoError(t, err)

	expectEqual(t, 1, len(yFile.Content)) // one empty entry is included
	content0 := &(yFile.Content[0])

	expectEqual(t, "", content0.Desc)
	expectEqual(t, 0.0, content0.X1)
}

func TestGetYamlFile_UnknownFields(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "unknownFields.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}

func TestGetYamlFile_FieldTypeMismatch(t *testing.T) {
	fileToTest := filepath.Join(templateTestDir, "fieldTypeMismatch.yml")
	yFile, err := GetYamlFile(fileToTest)

	expectNil(t, yFile)
	expectError(t, err)
}