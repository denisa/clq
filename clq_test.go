package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangelogShouldSucceed(t *testing.T) {
	executeClq(assert.New(t), "", 0, "", "", "CHANGELOG.md")
}

func TestScenarios(t *testing.T) {
	type Scenario struct {
		Platform  string   `json:"platform,omitempty"`
		Name      string   `json:"name"`
		Result    int      `json:"result"`
		Arguments []string `json:"arguments,omitempty"`
		Input     string   `json:"input,omitempty"`
		Output    string   `json:"output,omitempty"`
		Error     string   `json:"error,omitempty"`
	}
	assert := assert.New(t)
	file, err := os.Open("testdata/scenarios.json")
	require.NoError(t, err)
	defer file.Close()

	dec := json.NewDecoder(bufio.NewReader(file))

	testFiles := allTestFiles(t)
	var scenarios []Scenario
	for {
		if err := dec.Decode(&scenarios); err == io.EOF {
			break
		} else {
			require.NoErrorf(t, err, "Error reading %v", file.Name())
		}
		for _, scenario := range scenarios {
			t.Run(scenario.Name, func(t *testing.T) {
				switch scenario.Platform {
				case "unix":
					if os.Getuid() == -1 {
						t.Skipf("Running on Windows, skipping %v/%v", scenario.Name, scenario.Platform)
						return
					}
				case "windows":
					if os.Getuid() != -1 {
						t.Skipf("Running on Unix, skipping %v/%v", scenario.Name, scenario.Platform)
						return
					}
				}
				delete(testFiles, scenario.Name)

				var args []string
				if scenario.Name == "" {
					args = scenario.Arguments
				} else if scenario.Name == "-" {
					args = append(scenario.Arguments, scenario.Name)
				} else {
					args = append(scenario.Arguments, filepath.Join("testdata", scenario.Name))
				}
				executeClq(assert, scenario.Input, scenario.Result, scenario.Output, scenario.Error, args...)
			})
		}
	}

	if len(testFiles) != 0 {
		var buf strings.Builder
		for val := range testFiles {
			if buf.Len() > 0 {
				buf.WriteString(", ")
			}
			if val, err := json.Marshal(Scenario{Name: val, Result: 999999}); err == nil {
				buf.Write(val)
			}
		}
		require.Failf(t, "Unused test files: %v", buf.String())
	}
}

func executeClq(assert *assert.Assertions, input string, expectedCode int, expectedOutput string, expectedErr string, arguments ...string) {
	var actualOutput strings.Builder
	var actualErr strings.Builder
	clq := &Clq{stdin: strings.NewReader(input), stdout: &actualOutput, stderr: &actualErr}

	var actualCode = clq.entryPoint("clq", arguments...)

	assert.Equal(expectedCode, actualCode)
	assert.Equal(expectedOutput, actualOutput.String())
	assert.Equal(expectedErr, actualErr.String())
}

func allTestFiles(t *testing.T) map[string]bool {
	file, err := os.Open("testdata")
	require.NoError(t, err)
	defer file.Close()

	names, err := file.Readdirnames(0)
	require.NoError(t, err)

	result := make(map[string]bool)
	for _, val := range names {
		if val == "scenarios.json" {
			continue
		}
		result[val] = true
	}
	return result
}
