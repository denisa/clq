package query

import (
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
	if query == "" {
		return qe, nil
	}

	var queryFactory = introductionQueryFactory
	var selector = ""
	var isRecursive = false
	queryElements := strings.Split(query, ".")
	for i := 0; queryFactory != nil; {
		if q, parsedElement, err := queryFactory(selector, isRecursive, queryElements[i:]); err == nil {
			qe.queries = append(qe.queries, q)
			if q.isCollection() {
				outputFormat.SetCollection()
			}
			isRecursive = parsedElement.isRecursive
			queryFactory = parsedElement.queryFactory
			selector = parsedElement.selector
			i = min(i+1, len(queryElements))
		} else {
			return nil, err
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
