package query

import (
	"fmt"

	"github.com/denisa/clq/internal/changelog"
)

func (qe *QueryEngine) newChangeItemQuery(selector string, isRecursive bool, queryElements []string) error {
	if selector != "" {
		return fmt.Errorf("Query change selector %q not yet supported", selector)
	}
	if len(queryElements) > 0 {
		return fmt.Errorf("Query attribute selector %v not yet supported", queryElements)
	}
	queryMe := &changeItemQuery{}
	queryMe.exit = func(rc ResultCollector, h changelog.Heading) {
		if h, ok := h.(changelog.ChangeItem); ok {
			rc.set(h.Title())
		}
	}
	qe.queries = append(qe.queries, queryMe)
	return nil
}

type changeItemQuery struct {
	projections
}

func (q *changeItemQuery) Accept(heading changelog.Heading) bool {
	_, ok := heading.(changelog.ChangeItem)
	return ok
}

func (q *changeItemQuery) Enter(heading changelog.Heading) (bool, Project) {
	if !q.Accept(heading) {
		return false, nil
	}
	return true, q.enter
}

func (q *changeItemQuery) Exit(heading changelog.Heading) (bool, Project) {
	if !q.Accept(heading) {
		return false, nil
	}
	return true, q.exit
}
