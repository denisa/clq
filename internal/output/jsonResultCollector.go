package output

import (
	"encoding/json"

	"github.com/denisa/clq/internal/changelog"
)

// a jsonResultCollector produces a json-repesentation of the query result.
type jsonResultCollector struct {
	results    []jsonResult
	collection bool
}

type jsonResult struct {
	value interface{}
	name  string
	kind  changelog.HeadingKind
}

func (rc *jsonResultCollector) Result() string {
	if len(rc.results) == 0 {
		if rc.collection {
			return "[]"
		}
		return ""
	}
	if rc.results[0].value == nil {
		return "{}"
	}
	if result, ok := (rc.results[0].value).(string); ok {
		return result
	}

	jsonString, _ := json.Marshal(rc.results[0].value)
	return string(jsonString)
}

func (rc *jsonResultCollector) Open(heading changelog.Heading) {
	opened := jsonResult{kind: heading.Kind()}
	rc.results = append(rc.results, opened)
}

func (rc *jsonResultCollector) Close(_ changelog.Heading) {
	i := len(rc.results) - 1
	if i == -1 {
		return
	}

	newValue := rc.results[i].value
	if newValue == nil {
		return
	}

	if i == 0 {
		if _, ok := newValue.([]interface{}); !ok && rc.collection {
			rc.results[0] = jsonResult{value: []interface{}{newValue}}
		}
		return
	}
	rc.results = rc.results[:i]
	i--

	if _, ok := (rc.results[i].value).([]interface{}); ok && rc.collection {
		rc.results[i].value = append((rc.results[i].value).([]interface{}), newValue)
		return
	}

	result, _ := (rc.results[i].value).(map[string]interface{})
	collection, _ := result[rc.results[i].name].([]interface{})
	result[rc.results[i].name] = append(collection, newValue)
}

func (rc *jsonResultCollector) SetCollection() {
	rc.collection = true
}

func (rc *jsonResultCollector) Set(value string) {
	rc.results[len(rc.results)-1].value = value
}

func (rc *jsonResultCollector) SetField(name string, value string) {
	i := len(rc.results) - 1
	if rc.results[i].value == nil {
		rc.results[i].value = make(map[string]interface{})
	}
	result, _ := (rc.results[i].value).(map[string]interface{})
	result[name] = value
}

func (rc *jsonResultCollector) Array(name string) {
	i := len(rc.results) - 1
	if rc.results[i].value == nil {
		rc.results[i].value = make(map[string]interface{})
	}
	result, _ := (rc.results[i].value).(map[string]interface{})
	result[name] = make([]interface{}, 0)
	rc.results[i].name = name
}
