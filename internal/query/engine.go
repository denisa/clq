package query

import (
	"fmt"
	"strings"

	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/output"
)

// Engine tracks the evaluation of the overall query against the complete changelog.
type Engine struct {
	output  output.Format
	queries []Query
	current int
}

// NewEngine parses the query and contructs a new dedicated query engine.
// It is not an error for the query to be empty.
func NewEngine(query string, outputFormat output.Format) (*Engine, error) {
	qe := &Engine{output: outputFormat}
	if len(query) > 0 {
		if err := qe.newIntroductionQuery(strings.Split(query, ".")); err != nil {
			return nil, err
		}
		for _, q := range qe.queries {
			if q.isCollection() {
				outputFormat.SetCollection()
				break
			}
		}
	}
	return qe, nil
}

// HasQuery is true the Engine was constructed with a non-empty query.
// a Engine with an empty query is a no-op and an be skipped.
func (qe *Engine) HasQuery() bool { return len(qe.queries) > 0 }

// Result returns the result of the query evaluation.
func (qe *Engine) Result() string {
	return qe.output.Result()
}

// Enter lets the query engine evaluates the heading upon entering it.
func (qe *Engine) Enter(heading changelog.Heading) {
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
		project(qe.output, heading)
	}

	if qe.current+1 < len(qe.queries) {
		qe.current++
	}
}

// Exit lets the query engine evaluates the heading upon leaving it.
func (qe *Engine) Exit(heading changelog.Heading) {
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
		project(qe.output, heading)
	}

	qe.output.Close(heading)
}

func parseName(element string) (name, selector string, isList, isRecursive bool, err error) {
	openBracketIndex := strings.Index(element, "[")
	closeBracketIndex := strings.Index(element, "]")
	if openBracketIndex != -1 {
		if closeBracketIndex < openBracketIndex {
			err = fmt.Errorf("missing closing bracket in %q", element)
			return
		}
		isList = true
	} else if closeBracketIndex != -1 {
		err = fmt.Errorf("missing opening bracket in %q", element)
		return
	}
	isRecursive = strings.HasSuffix(element, "/")
	if !isList && isRecursive {
		err = fmt.Errorf("missing closing bracket in %q", element)
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
		return fmt.Errorf("%q is a scalar attribute", name)
	}
	if len(elements) != 0 {
		return fmt.Errorf("no further query element allowed after %q", name)
	}
	return nil
}

func elementIsCollection(name string, isList bool) error {
	if !isList {
		return fmt.Errorf("%q is a collection attribute", name)
	}
	return nil
}
