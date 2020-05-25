package query

import (
	"strings"

	"github.com/denisa/clq/internal/changelog"
	"github.com/yuin/goldmark/util"
)

// QueryEngine tracks the evaluation of the overall query against the complete changelog.
type QueryEngine struct {
	queries []Query
	current int
}

// NewQueryEngine parses the query and contructs a new dedicated query engine.
func NewQueryEngine(query string) (*QueryEngine, error) {
	qe := &QueryEngine{}
	if len(query) > 0 {
		queryElements := strings.Split(query, ".")
		if err := qe.newChangelogQuery(queryElements[0], queryElements[1:]); err != nil {
			return nil, err
		}
	}
	return qe, nil
}

// Apply lets the query engine evaluates the heading.
// Apply might write the result of the query to the buffer.
func (qe *QueryEngine) Apply(w util.BufWriter, heading changelog.Heading) {
	if qe.current < len(qe.queries) && qe.queries[qe.current].Select(w, heading) {
		qe.current++
	}
}
