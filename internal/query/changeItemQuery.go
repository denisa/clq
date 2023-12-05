package query

import (
	"fmt"

	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/output"
)

func (qe *QueryEngine) newChangeItemQuery(selector string, _ bool, queryElements []string) error {
	if selector != "" {
		return fmt.Errorf("Query change selector %q not yet supported", selector)
	}
	if len(queryElements) > 0 {
		return fmt.Errorf("Query attribute selector %v not yet supported", queryElements)
	}
	queryMe := &changeItemQuery{}
	queryMe.collection = true
	queryMe.exit = func(of output.OutputFormat, h changelog.Heading) {
		if h, ok := h.(changelog.ChangeItem); ok {
			of.Set(h.DisplayTitle())
		}
	}
	qe.queries = append(qe.queries, queryMe)
	return nil
}

type changeItemQuery struct {
	projections
}

func (q *changeItemQuery) isCollection() bool {
	return q.collection
}

func (q *changeItemQuery) Accept(heading changelog.Heading) bool {
	_, ok := heading.(changelog.ChangeItem)
	return ok
}

func (q *changeItemQuery) Enter(heading changelog.Heading) (bool, project) {
	if !q.Accept(heading) {
		return false, nil
	}
	return true, q.enter
}

func (q *changeItemQuery) Exit(heading changelog.Heading) (bool, project) {
	if !q.Accept(heading) {
		return false, nil
	}
	return true, q.exit
}
