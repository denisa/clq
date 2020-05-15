package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChangelogShouldSucceed(t *testing.T) {
	executeClq(t, 0, "CHANGELOG.md")
}

func TestScenarios(t *testing.T) {
	type Scenario struct {
		Name      string   `json:"name"`
		Result    int      `json:"result"`
		Arguments []string `json:"arguments,omitempty"`
	}
	require := require.New(t)
	file, err := os.Open("testdata/scenarios.json")
	require.NoError(err)
	defer file.Close()

	dec := json.NewDecoder(bufio.NewReader(file))

	testFiles := allTestFiles(t)
	var scenarios []Scenario
	for {
		if err := dec.Decode(&scenarios); err == io.EOF {
			break
		} else {
			require.NoError(err)
		}
		for _, scenario := range scenarios {
			delete(testFiles, scenario.Name)
			t.Run(scenario.Name, func(t *testing.T) {
				args := append(scenario.Arguments, filepath.Join("testdata", scenario.Name))
				executeClq(t, scenario.Result, args...)
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
		require.Failf("Unused test files: %v", buf.String())
	}
}

func executeClq(t *testing.T, expected int, arguments ...string) {
	var actual = entryPoint("clq", arguments...)
	require.Equal(t, expected, actual)
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
