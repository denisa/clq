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
	assertions := assert.New(t)
	file, err := os.Open("testdata/scenarios.json")
	require.NoError(t, err)
	defer closeIgnoreError(file)

	dec := json.NewDecoder(bufio.NewReader(file))

	testFiles := allTestFiles(t)
	var scenarios []Scenario
	for {
		if err := dec.Decode(&scenarios); err == io.EOF {
			break
		}
		require.NoErrorf(t, err, "Error reading %v", file.Name())
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
				scenario.executeClq(assertions)
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
		require.Failf(t, "Unused test files", "%v", buf.String())
	}
}

func (s Scenario) executeClq(assertions *assert.Assertions) {
	var actualOutput strings.Builder
	var actualErr strings.Builder
	clq := &Clq{stdin: strings.NewReader(s.Input), stdout: &actualOutput, stderr: &actualErr}

	var actualCode = clq.entryPoint("clq", s.Arguments...)

	assertions.Equal(s.Result, actualCode)
	switch s.OutputFormat {
	case "json":
		assertions.JSONEq(s.Output, actualOutput.String())
	default:
		assertions.Equal(s.Output, actualOutput.String())
	}
	assertions.Equal(s.Error, actualErr.String())
}

func allTestFiles(t *testing.T) map[string]bool {
	file, err := os.Open("testdata")
	require.NoError(t, err)
	defer closeIgnoreError(file)

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

func closeIgnoreError(file *os.File) {
	_ = file.Close()
}
