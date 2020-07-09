package query

import (
	"encoding/json"
	"strings"

	"github.com/denisa/clq/internal/changelog"
)

// QueryEngine tracks the evaluation of the overall query against the complete changelog.
type QueryEngine struct {
	queries []Query
	current int
	results []result
}

// NewQueryEngine parses the query and contructs a new dedicated query engine.
func NewQueryEngine(query string) (*QueryEngine, error) {
	qe := &QueryEngine{}
	if len(query) > 0 {
		queryElements := strings.Split(query, ".")
		if err := qe.newIntroductionQuery(queryElements[0], queryElements[1:]); err != nil {
			return nil, err
		}
	}
	return qe, nil
}

func (qe *QueryEngine) HasQuery() bool { return len(qe.queries) > 0 }

// Enter lets the query engine evaluates the heading upon entering it.
func (qe *QueryEngine) Enter(heading changelog.Heading) {
	if qe.current == len(qe.queries) {
		// no queries defined...
		return
	}

	ok, project := qe.queries[qe.current].Enter(heading)
	if !ok {
		return
	}

	if project != nil {
		qe.open(heading)
		project(qe, heading)
	}

	if qe.current+1 < len(qe.queries) {
		qe.current++
	}
}

// Exit lets the query engine evaluates the heading upon leaving it.
func (qe *QueryEngine) Exit(heading changelog.Heading) {
	if qe.current == len(qe.queries) {
		// no queries defined...
		return
	}

	for {
		ok, project := qe.queries[qe.current].Exit(heading)
		if !ok {
			return
		}

		if project != nil {
			qe.open(heading)
			project(qe, heading)
		}
		qe.close(heading)

		if qe.current == 0 {
			return
		}
		qe.current--
	}
}

type result struct {
	value interface{}
	name  string
	kind  changelog.HeadingKind
}

func (qe *QueryEngine) Result() string {
	if len(qe.results) == 0 {
		return ""
	}
	if len(qe.results) > 1 {
		panic("results greater than 1")
	}
	if result, ok := (qe.results[0].value).(string); ok {
		return result
	}

	result, _ := (qe.results[0].value).(map[string]interface{})
	if jsonString, err := json.Marshal(result); err != nil {
		panic(err)
	} else {
		return string(jsonString)
	}
}

func (qe *QueryEngine) open(heading changelog.Heading) {
	if i := len(qe.results); i > 0 && qe.results[i-1].kind == heading.Kind() {
		return
	}
	opened := result{kind: heading.Kind()}
	qe.results = append(qe.results, opened)
}

func (qe *QueryEngine) close(heading changelog.Heading) {
	if i := len(qe.results); i < 2 || qe.results[i-1].kind != heading.Kind() {
		return
	}

	i := len(qe.results) - 1
	newValue := qe.results[i].value
	qe.results = qe.results[:i]
	i--
	if result, ok := (qe.results[i].value).(map[string]interface{}); !ok {
		panic("WTF?")
	} else {
		if collection, ok := result[qe.results[i].name].([]interface{}); ok {
			result[qe.results[i].name] = append(collection, newValue)
		} else {
			panic("WTF?")
		}
	}
}

func (qe *QueryEngine) set(value string) {
	qe.results[len(qe.results)-1].value = value
}

func (qe *QueryEngine) setField(name string, value string) {
	i := len(qe.results) - 1
	if qe.results[i].value == nil {
		qe.results[i].value = make(map[string]interface{})
	}
	result, _ := (qe.results[i].value).(map[string]interface{})
	result[name] = value
}

func (qe *QueryEngine) array(name string) {
	i := len(qe.results) - 1
	if qe.results[i].value == nil {
		qe.results[i].value = make(map[string]interface{})
	}
	result, _ := (qe.results[i].value).(map[string]interface{})
	result[name] = make([]interface{}, 0)
	qe.results[i].name = name
}

func selector(element string) string {
	return element[strings.Index(element, "[")+1 : strings.Index(element, "]")]
}
