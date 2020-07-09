package query

import (
	"encoding/json"

	"github.com/denisa/clq/internal/changelog"
)

type jsonResultCollector struct {
	results []jsonResult
}

type jsonResult struct {
	value interface{}
	name  string
	kind  changelog.HeadingKind
}

func (rc *jsonResultCollector) Result() string {
	if len(rc.results) == 0 {
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

func (rc *jsonResultCollector) Close(heading changelog.Heading) {
	i := len(rc.results) - 1
	if i < 1 {
		return
	}
	if r, ok := (rc.results[i-1].value).([]interface{}); ok {
		r = append(r, rc.results[i].value)
		rc.results = rc.results[:i-1]
		collection := jsonResult{value: r}
		rc.results = append(rc.results, collection)
		return
	}
	if rc.results[i-1].kind == rc.results[i].kind {
		r := []interface{}{rc.results[i-1].value, rc.results[i].value}
		rc.results = rc.results[:i-1]
		collection := jsonResult{value: r}
		rc.results = append(rc.results, collection)
		return
	}

	newValue := rc.results[i].value
	rc.results = rc.results[:i]
	i--

	if newValue != nil {
		result, _ := (rc.results[i].value).(map[string]interface{})
		collection, _ := result[rc.results[i].name].([]interface{})
		result[rc.results[i].name] = append(collection, newValue)
	}
}

func (rc *jsonResultCollector) set(value string) {
	rc.results[len(rc.results)-1].value = value
}

func (rc *jsonResultCollector) setField(name string, value string) {
	i := len(rc.results) - 1
	if rc.results[i].value == nil {
		rc.results[i].value = make(map[string]interface{})
	}
	result, _ := (rc.results[i].value).(map[string]interface{})
	result[name] = value
}

func (rc *jsonResultCollector) array(name string) {
	i := len(rc.results) - 1
	if rc.results[i].value == nil {
		rc.results[i].value = make(map[string]interface{})
	}
	result, _ := (rc.results[i].value).(map[string]interface{})
	result[name] = make([]interface{}, 0)
	rc.results[i].name = name
}
