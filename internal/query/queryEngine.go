package query

import (
	"fmt"
	"strings"

	"github.com/denisa/clq/internal/changelog"
)

// QueryEngine tracks the evaluation of the overall query against the complete changelog.
type QueryEngine struct {
	output  OutputFormat
	queries []Query
	current int
}

// NewQueryEngine parses the query and contructs a new dedicated query engine.
func NewQueryEngine(query string, formatName string) (*QueryEngine, error) {
	outputFormat, err := newOutputFormat(formatName)
	if err != nil {
		return nil, err
	}

	qe := &QueryEngine{output: outputFormat}
	if len(query) > 0 {
		if err := qe.newIntroductionQuery(strings.Split(query, ".")); err != nil {
			return nil, err
		}
	}
	return qe, nil
}

func (qe *QueryEngine) HasQuery() bool { return len(qe.queries) > 0 }

func (qe *QueryEngine) Result() string {
	return qe.output.Result()
}

// Enter lets the query engine evaluates the heading upon entering it.
func (qe *QueryEngine) Enter(heading changelog.Heading) {
	if qe.current == len(qe.queries) {
		// no queries defined...
		return
	}

	for cur := qe.current; ; cur-- {
		if qe.queries[cur].Accept(heading) {
			qe.current = cur
			break
		}
		if cur == 0 {
			return
		}
	}

	ok, project := qe.queries[qe.current].Enter(heading)
	if !ok {
		return
	}
	if project != nil {
		qe.output.Open(heading)
		project(qe.output.(ResultCollector), heading)
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

	for cur := qe.current; ; cur-- {
		if qe.queries[cur].Accept(heading) {
			qe.current = cur
			break
		}
		if cur == 0 {
			return
		}
	}

	ok, project := qe.queries[qe.current].Exit(heading)
	if ok && project != nil {
		qe.output.Open(heading)
		project(qe.output.(ResultCollector), heading)
	}

	qe.output.Close(heading)
}

func parseName(element string) (name, selector string, isList, isRecursive bool, err error) {
	openBracketIndex := strings.Index(element, "[")
	closeBracketIndex := strings.Index(element, "]")
	if openBracketIndex != -1 {
		if closeBracketIndex < openBracketIndex {
			err = fmt.Errorf("Missing closing bracket in %q", element)
			return
		}
		isList = true
	} else if closeBracketIndex != -1 {
		err = fmt.Errorf("Missing opening bracket in %q", element)
		return
	}
	isRecursive = strings.HasSuffix(element, "/")
	if !isList && isRecursive {
		err = fmt.Errorf("Missing closing bracket in %q", element)
		return
	}
	if isList {
		name = element[:openBracketIndex]
		selector = element[openBracketIndex+1 : closeBracketIndex]
	} else {
		name = element
	}
	return
}

func elementIsFinal(name string, isList bool, elements []string) error {
	if isList {
		return fmt.Errorf("%q is a scalar attribue", name)
	}
	if len(elements) != 0 {
		return fmt.Errorf("No further query element allowed after %q", name)
	}
	return nil
}
