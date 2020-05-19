package query

import (
	"strings"

	"github.com/denisa/clq/internal/changelog"
	"github.com/yuin/goldmark/util"
)

type Query interface {
	Select(w util.BufWriter, heading changelog.Heading) bool
}

type QueryEngine struct {
	queries []Query
	current int
}

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

func (qe *QueryEngine) Apply(w util.BufWriter, heading changelog.Heading) {
	if qe.current < len(qe.queries) && qe.queries[qe.current].Select(w, heading) {
		qe.current++
	}
}
