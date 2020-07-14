package query

import (
	"strings"

	"github.com/denisa/clq/internal/changelog"
)

// QueryEngine tracks the evaluation of the overall query against the complete changelog.
type QueryEngine struct {
	queries []Query
	current int
	w       strings.Builder
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
func (qe *QueryEngine) Result() string { return qe.w.String() }

// Enter lets the query engine evaluates the heading upon entering it.
func (qe *QueryEngine) Enter(heading changelog.Heading) {
	if qe.current == len(qe.queries) {
		// no queries defined...
		return
	}

	if qe.queries[qe.current].Enter(&qe.w, heading) {
		if qe.current+1 < len(qe.queries) {
			qe.current++
		}
		return
	}
}

// Exit lets the query engine evaluates the heading upon leaving it.
func (qe *QueryEngine) Exit(heading changelog.Heading) {
	if qe.current == len(qe.queries) {
		// no queries defined...
		return
	}

	for !qe.queries[qe.current].Exit(&qe.w, heading) {
		if qe.current == 0 {
			return
		}
		qe.current--
	}
}
