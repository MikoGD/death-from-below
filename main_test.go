package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	if err := os.Mkdir("test_environment", 0777); err != nil {
		log.Fatalf("Failed to create test_environment directory\nError: %v\n", err)
	}

	dir := os.DirFS("./data/tests/environment")

	if err := os.CopyFS("./test_environment", dir); err != nil {
		log.Fatalf("Failed to copy test environment\nError: %v\n", err)
	}

	code := m.Run()

	if err := os.RemoveAll("test_environment"); err != nil {
		log.Fatalf("Failed to remove test environment\nError: %v\n", err)
	}

	os.Exit(code)
}

type TestFunctionErr struct {
	Name       string
	Parameters []string
  ExpectedError error
}


func TestInvalidArgsLength(t *testing.T) {
	// Not including the first arg passed in by default when calling
	tests := []TestFunctionErr{
		{"No args", []string{"Main"}, InvalidArgsError{}},
		{"One arg", []string{"Main", "src", "pattern"}, InvalidArgsError{}},
		{"Two arg", []string{"Main", "src", "pattern", "from",}, InvalidArgsError{}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			os.Args = test.Parameters
			err := searchAndReplace()

			if !errors.Is(err, InvalidArgsError{}) {
				t.Errorf("Expected InvalidArgsError but received %v\n", err)
			}
		})
	}
}

func TestInvalidArg(t *testing.T) {
	// Not including the first arg passed in by default when calling
	tests := []TestFunctionErr{
    {"Empty source arg", []string{"Main", "", "pattern", "from", "to"}, InvalidSourceArg{}},
    {"Empty glob pattern arg", []string{"Main", "src", "", "from", "to"}, InvalidGlobArg{}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			os.Args = test.Parameters
			err := searchAndReplace()

			if !errors.Is(err, test.ExpectedError) {
				t.Errorf("Expected InvalidSourceArg or InvalidGlobArg but received %v\n", err)
			}
		})
	}
}

func TestReplaceAllOldStringWithNewString(t *testing.T) {
	os.Args = []string{"Main", "test_environment", "*.txt", "<replace line>", "[replaced line]"}
	err := searchAndReplace()
	if err != nil {
		t.Errorf("Expected no errors but received %v\n", err)
	}

	err = filepath.WalkDir("./test_environment/text", func(path string, dir os.DirEntry, err error) error {
		if dir.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		expectedFilePath := strings.Replace(path, "/text/", "/text-replaced/", 1)
		expectedFile, err := os.Open(expectedFilePath)
		if err != nil {
			return err
		}
		defer expectedFile.Close()

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		fileContent, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		expectedFileContent, err := io.ReadAll(expectedFile)
		if err != nil {
			return err
		}

		if string(fileContent) != string(expectedFileContent) {
			t.Errorf("%s does not match %s\nActual:\n%s\n\nExpected:\n%s\n\n", path, expectedFilePath, string(fileContent), string(expectedFileContent))
			return errors.New("Actual file does not match expected file")
		}

		return nil
	})

	if err != nil {
		t.Errorf("Expected no errors after replace but recieved: %s", err.Error())
	}
}
