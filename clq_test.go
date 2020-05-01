package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChangelogShouldSucceed(t *testing.T) {
	executeClq(t, 0, "CHANGELOG.md")
}

func TestScenarios(t *testing.T) {
	type Scenario struct {
		Name      string   `json:"name"`
		Success   int      `json:"success"`
		Arguments []string `json:"arguments,omitempty"`
	}
	file, err := os.Open("testdata/scenarios.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	dec := json.NewDecoder(bufio.NewReader(file))

	testFiles := allTestFiles(t)
	var scenarios []Scenario
	for {
		if err := dec.Decode(&scenarios); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		for _, scenario := range scenarios {
			delete(testFiles, scenario.Name)
			t.Run(scenario.Name, func(t *testing.T) {
				args := append(scenario.Arguments, filepath.Join("testdata", scenario.Name))
				executeClq(t, scenario.Success, args...)
			})
		}
	}

	if len(testFiles) != 0 {
		var buf strings.Builder
		for val, _ := range testFiles {
			if buf.Len() > 0 {
				buf.WriteString(", ")
			}
			if val, err := json.Marshal(Scenario{Name: val, Success: 999999}); err == nil {
				buf.Write(val)
			}
		}
		t.Fatalf("Unused test files: " + buf.String())
	}
}

func executeClq(t *testing.T, expected int, arguments ...string) {
	var actual = entryPoint("clq", arguments...)
	if expected != actual {
		t.Errorf("Expected %d but received %v", expected, actual)
	}
}

func allTestFiles(t *testing.T) map[string]bool {
	file, err := os.Open("testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	names, err := file.Readdirnames(0)
	if err != nil {
		t.Fatal(err)
	}
	result := make(map[string]bool)
	for _, val := range names {
		if val == "scenarios.json" {
			continue
		}
		result[val] = true
	}
	return result
}
