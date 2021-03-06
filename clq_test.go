package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Scenario struct {
	Platform     string   `json:"platform,omitempty"`
	Title        string   `json:"title,omitempty"`
	Name         string   `json:"name"`
	Result       int      `json:"result"`
	Arguments    []string `json:"arguments,omitempty"`
	Input        string   `json:"input,omitempty"`
	Output       string   `json:"output,omitempty"`
	OutputFormat string   `json:"output_format,omitempty"`
	Error        string   `json:"error,omitempty"`
}

func TestChangelogShouldSucceed(t *testing.T) {
	Scenario{Arguments: []string{"CHANGELOG.md"}}.executeClq(assert.New(t))
}

func (s Scenario) name() string {
	var buf strings.Builder
	if s.Title != "" {
		buf.WriteString(s.Title)
	} else {
		buf.WriteString(s.Name)
	}
	if s.Platform != "" {
		buf.WriteString(" (")
		buf.WriteString(s.Platform)
		buf.WriteString(")")
	}
	return buf.String()
}

func TestScenarios(t *testing.T) {
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
			t.Run(scenario.name(), func(t *testing.T) {
				switch scenario.Platform {
				case "skip":
					t.SkipNow()
					return
				case "unix":
					if runtime.GOOS == "windows" {
						t.SkipNow()
						return
					}
				case "windows":
					if runtime.GOOS != "windows" {
						t.SkipNow()
						return
					}
				}
				delete(testFiles, scenario.Name)

				if scenario.Name == "-" {
					scenario.Arguments = append(scenario.Arguments, scenario.Name)
				} else if scenario.Name != "" {
					scenario.Arguments = append(scenario.Arguments, filepath.Join("testdata", scenario.Name))
				}
				scenario.executeClq(assert)
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

func (scenario Scenario) executeClq(assert *assert.Assertions) {
	var actualOutput strings.Builder
	var actualErr strings.Builder
	clq := &Clq{stdin: strings.NewReader(scenario.Input), stdout: &actualOutput, stderr: &actualErr}

	var actualCode = clq.entryPoint("clq", scenario.Arguments...)

	assert.Equal(scenario.Result, actualCode)
	switch scenario.OutputFormat {
	case "json":
		assert.JSONEq(scenario.Output, actualOutput.String())
	default:
		assert.Equal(scenario.Output, actualOutput.String())
	}
	assert.Equal(scenario.Error, actualErr.String())
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
